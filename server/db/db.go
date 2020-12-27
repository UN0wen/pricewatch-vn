// Package db is a db wrapper
package db

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

// AutoParam are parameters that are created automatically
// SkipParam are parameters that are supposed to be skipped over
// TimeParams are parameters of type time.Time
var (
	AutoParam = map[string]bool{
		"timecreated": true,
	}
	SkipParam = map[string]bool{
		"post_results": true,
		"pre_results":  true,
	}
	TimeParam = map[string]bool{
		"timecreated":  true,
		"timeloggedin": true,
		"expiresafter": true,
	}
)

// IsUndeclared uses reflection to see if the value of the field is set or not.
// It takes in an interface to reflect on and returns a boolean if the field is set or not.
func isUndeclared(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

// Db contains the database connection pool and its config
type Db struct {
	Pool *pgxpool.Pool
	cfg  Config
}

// Config is a database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// Setup setups the database
func Setup(cfg Config) (db Db, err error) {
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" ||
		cfg.Password == "" || cfg.Database == "" {
		err = errors.New("Provide all fields for config")
		return
	}
	db.cfg = cfg
	cfgDNS := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port)

	config, err := pgxpool.ParseConfig(cfgDNS)
	if err != nil {
		err = errors.Wrapf(err, "Cannot parse config string")
		return
	}

	// PGXPool configs
	config.MaxConns = 10

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		err = errors.Wrapf(err, "Unable to connect to database")
		return
	}

	db.Pool = pool

	utils.Sugar.Infof("Database created with config string: %s", config.ConnString())
	return db, nil
}

// Close closes the connection pool
func (db *Db) Close() {
	if db.Pool == nil {
		return
	}

	db.Pool.Close()
	return
}

// CreateTable executes the query given
func (db *Db) CreateTable(query string) (err error) {
	utils.Sugar.Infof("SQL Query: %s", query)

	if _, err = db.Pool.Exec(context.Background(), query); err != nil {
		err = errors.Wrapf(err, "Table creation query failed")
	}

	return
}

// CreateIndex creates an index on tablename
func (db *Db) CreateIndex(tableName, indexType, indexColumn string) (err error) {
	var query bytes.Buffer
	query.WriteString(fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s_%s_idx ON %s USING %s (%s)", tableName, indexColumn, tableName, indexType, indexColumn))

	utils.Sugar.Infof("SQL Query: %s", query.String())
	_, err = db.Pool.Exec(context.Background(), query.String())
	if err != nil {
		err = errors.Wrapf(err, "Index creation query failed")
	}
	return
}

// CreateView creates an view on tablename with query
func (db *Db) CreateView(viewName, viewQuery string) (err error) {
	var query bytes.Buffer

	query.WriteString(fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s", viewName, viewQuery))

	utils.Sugar.Infof("SQL Query: %s", query.String())

	_, err = db.Pool.Exec(context.Background(), query.String())
	if err != nil {
		err = errors.Wrapf(err, "View creation query failed")
	}
	return
}

// Simple ORM to wrap database calls

// Get attempts to provide a generalized search through the specified table based on the provided queries.
// It takes a query for the queryable fields, and an operator such as "AND" or "OR" to define the context of the search. It takes in a table name to act on.
// It returns all the data for all found objects and an error if one exists.
func (db *Db) Get(mQuery interface{}, op, compareOp, table string) (objects []map[string]interface{}, err error) {
	var query bytes.Buffer
	if compareOp == "" {
		compareOp = "="
	}

	query.WriteString(fmt.Sprintf("SELECT * FROM %s", table))

	// Use reflection to analyze object fields
	fields := reflect.ValueOf(mQuery)
	first := true
	var values []interface{}
	vIdx := 1
	for i := 0; i < fields.NumField(); i++ {
		v := fields.Field(i).Interface()
		// Skip fields that are not set to query on
		if !isUndeclared(v) {
			if first {
				query.WriteString(" WHERE ")
				first = false
			} else {
				if op != "" {
					query.WriteString(fmt.Sprintf(" %s ", op))
				}
			}
			k := strings.ToLower(fields.Type().Field(i).Name)
			v = fmt.Sprintf("%v", v)
			values = append(values, v)
			query.WriteString(fmt.Sprintf("%s%s$%d", k, compareOp, vIdx))
			vIdx++
		}
	}
	query.WriteString(";")

	utils.Sugar.Infof("SQL Query: %s", query.String())
	utils.Sugar.Infof("Values: %v", values)

	rows, err := db.Pool.Query(context.Background(), query.String(), values...)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}

	objects, err = mapRows(rows)
	return
}

// GetByID is GET() but for an ID.
// It returns the data for the found object, or an error if one exists.
func (db *Db) GetByID(id uuid.UUID, table string) (data map[string]interface{}, err error) {
	query := idQuery{ID: id}
	objs, err := db.Get(query, "", "=", table)
	if err != nil {
		err = errors.Wrapf(err, "Search error")
		return
	} else if objs == nil {
		err = errors.New(fmt.Sprintf("Failed to find object with id: %s", id.String()))
		return
	} else if len(objs) != 1 {
		err = errors.New(fmt.Sprintf("Found duplicate objects with id: %s", id.String()))
		return
	}
	data = objs[0]
	return
}

// Insert will put a new model row within the specified table in the DB, verifying all fields are valid.
// It takes in the object to insert and the table name.
// It returns id if it is successful.
func (db *Db) Insert(table, model interface{}) (err error) {
	modelName := reflect.TypeOf(model).String()

	_, err = govalidator.ValidateStruct(model)
	if err != nil {
		err = errors.Wrapf(err, "Invalid field in model: %s", modelName)
		return
	}

	var values []interface{}
	var vStr, kStr bytes.Buffer
	vIdx := 1
	fields := reflect.ValueOf(model)
	if fields.NumField() < 1 {
		err = errors.New("Invalid number of fields given")
		return
	}
	first := true
	for i := 0; i < fields.NumField(); i++ {
		k := strings.ToLower(fields.Type().Field(i).Name)
		v := fields.Field(i).Interface()
		// Skip auto params and skippable params
		if AutoParam[k] || SkipParam[k] {
			continue
		}
		if first {
			first = false
		} else {
			vStr.WriteString(", ")
			kStr.WriteString(", ")
		}
		kStr.WriteString(k)
		vStr.WriteString(fmt.Sprintf("$%d", vIdx))
		if k == "password" {
			hash, e := hashPassword(v.(string))
			if e != nil {
				err = errors.Wrapf(e, "Password hash failed")
				return
			}
			values = append(values, hash)
		} else {
			values = append(values, v)
		}

		vIdx++
	}
	var query bytes.Buffer
	query.WriteString(fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s);`, table, kStr.String(), vStr.String()))

	utils.Sugar.Infof("SQL Query: %s", query.String())
	utils.Sugar.Infof("Values: %v", values)

	_, err = db.Pool.Exec(context.Background(), query.String(), values...)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed to execute")
	}

	return
}

type idQuery struct {
	ID uuid.UUID
}

// Update updates the model row(s) in the table based on the incoming object.
// It takes in an id to identify the object in the DB, a string representing the table name and the fileds to update the object on.
// It returns the data representing an updated model.
func (db *Db) Update(id uuid.UUID, table string, updates interface{}) (data []map[string]interface{}, err error) {
	_, err = db.GetByID(id, table)
	if err != nil {
		return
	}
	var query bytes.Buffer
	query.WriteString(fmt.Sprintf("UPDATE %s SET", table))
	var values []interface{}
	vIdx := 1
	fields := reflect.ValueOf(updates)
	if fields.NumField() < 1 {
		err = errors.New("Invalid number of query fields")
		return
	}
	first := true
	for i := 0; i < fields.NumField(); i++ {
		k := strings.ToLower(fields.Type().Field(i).Name)
		v := fields.Field(i).Interface()
		// Skip auto params or unset fields on the incoming User
		if AutoParam[k] || SkipParam[k] || isUndeclared(fields.Field(i).Interface()) {
			continue
		} else if TimeParam[k] { // convert time types to String
			if t, ok := v.(time.Time); ok {
				v = t.Format(time.RFC3339)
			}
		}
		if first {
			query.WriteString(" ")
			first = false
		} else {
			query.WriteString(", ")
		}

		if k == "password" {
			hash, e := hashPassword(v.(string))
			if e != nil {
				err = errors.Wrapf(e, "Password hash failed")
				return
			}
			values = append(values, hash)
		} else {
			values = append(values, v)
		}

		query.WriteString(fmt.Sprintf("%v=$%d", k, vIdx))
		vIdx++
	}
	query.WriteString(fmt.Sprintf(" WHERE id='%s' RETURNING *;", id))

	utils.Sugar.Infof("SQL Query: %s", query.String())
	utils.Sugar.Infof("Values: %v", values)

	rows, err := db.Pool.Query(context.Background(), query.String(), values...)
	if err != nil {
		err = errors.Wrapf(err, "Update query failed")
		return
	}

	// Get updated model
	data, err = mapRows(rows)
	return
}

// Delete removes one row from table where id=id
func (db *Db) Delete(id uuid.UUID, table string) (err error) {
	_, err = db.GetByID(id, table)
	if err != nil {
		return
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", table)

	utils.Sugar.Infof("SQL Query: %s", query)

	_, err = db.Pool.Exec(context.Background(), query, id)
	if err != nil {
		err = errors.Wrapf(err, "Delete query failed")
		return
	}
	return
}

// DeleteByCol removes the one row from table where column=columnValue
func (db *Db) DeleteByCol(column, columnValue, table string) (err error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s=$1", table, column)

	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Value: %s", columnValue)

	_, err = db.Pool.Exec(context.Background(), query, columnValue)
	if err != nil {
		err = errors.Wrapf(err, "Delete query preparation failed")
		return
	}
	return
}

// DeleteAll permanently removes all objects from the table.
// It takes in a string representing the table name.
// It returns an error if one exists.
func (db *Db) DeleteAll(table string) (err error) {
	query := fmt.Sprintf("DELETE FROM %s;", table)
	utils.Sugar.Infof("SQL Query: %s", query)

	_, err = db.Pool.Exec(context.Background(), query)
	if err != nil {
		err = errors.Wrapf(err, "Delete all query preparation failed")
		return
	}
	return
}
