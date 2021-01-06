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
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

// AutoParam are parameters that are created automatically
// SkipParam are parameters that are supposed to be skipped over
// TimeParams are parameters of type time.Time
var (
	AutoParam = map[string]bool{
		"created": true,
	}
	SkipParam = map[string]bool{
		"post_results": true,
		"pre_results":  true,
	}
	TimeParam = map[string]bool{
		"created":       true,
		"logged_in":     true,
		"expires_after": true,
		"time":          true,
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
	URL      string
}

// Setup setups the database
func Setup(cfg Config) (Db, error) {
	var db Db
	if cfg.URL == "" {
		if cfg.Host == "" || cfg.Port == "" || cfg.User == "" ||
			cfg.Password == "" || cfg.Database == "" {
			err := errors.New("Provide all fields for config")
			return db, err
		}

		db.cfg = cfg
		cfgDNS := fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port)

		config, err := pgxpool.ParseConfig(cfgDNS)
		if err != nil {
			err = errors.Wrapf(err, "Cannot parse config string")
			return db, err
		}

		// PGXPool configs
		config.MaxConns = 10

		pool, err := pgxpool.ConnectConfig(context.Background(), config)
		if err != nil {
			err = errors.Wrapf(err, "Unable to connect to database")
			return db, err
		}

		db.Pool = pool
		utils.Sugar.Infof("Database created with config string: %s", config.ConnString())

	} else {
		pool, err := pgxpool.Connect(context.Background(), cfg.URL)
		if err != nil {
			err = errors.Wrapf(err, "Unable to connect to database")
			return db, err
		}
		db.Pool = pool
		utils.Sugar.Infof("Database created with URL: %s", cfg.URL)

	}
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

// Simple ORM to wrap database calls

// SearchOptions is a list of options to call db.Get with
type SearchOptions struct {
	Query      interface{} // The query to use
	Op         string      // AND / OR
	CompareOp  string      // The comparison operator to use (default =)
	TableName  string      // Name of table to query from
	Order      string      // ASC or DESC
	OrderQuery string      // Columns to order by
	Limit      int64       // Number of rows to select
}

// Update updates the model row(s) in the table based on the incoming object.
// It takes in an id to identify the object in the DB, a string representing the table name and the fileds to update the object on.
// It returns the data representing an updated model.
func (db *Db) Update(id uuid.UUID, table string, updates interface{}) (data []map[string]interface{}, err error) {
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
		// Check if there's a tag for db row name
		k := fields.Type().Field(i).Tag.Get("db")

		if k == "" {
			k = strings.ToLower(fields.Type().Field(i).Name)
		}

		v := fields.Field(i).Interface()
		// Skip auto params or unset fields on the incoming User
		// Also skip the ID field since we dont ever want to update it
		if AutoParam[k] || SkipParam[k] || isUndeclared(v) || k == "id" {
			continue
		} else if TimeParam[k] { // convert time types to String
			if t, ok := v.(time.Time); ok {
				utils.Sugar.Infof("%s", v)
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
			hash, e := utils.HashPassword(v.(string))
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

// Delete attempts to provide a generalized search through the specified table based on the provided queries.
// It takes a query for the queryable fields, and an operator such as "AND" or "OR" to define the context of the search. It takes in a table name to act on.
// It returns all the data for all found objects and an error if one exists.
func (db *Db) Delete(options SearchOptions) (err error) {
	var query bytes.Buffer
	if options.CompareOp == "" {
		options.CompareOp = "="
	}

	query.WriteString(fmt.Sprintf("DELETE FROM %s", options.TableName))

	// Use reflection to analyze object fields
	fields := reflect.ValueOf(options.Query)
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
				if options.Op != "" {
					query.WriteString(fmt.Sprintf(" %s ", options.Op))
				}
			}
			k := strings.ToLower(fields.Type().Field(i).Name)
			v = fmt.Sprintf("%v", v)
			values = append(values, v)
			query.WriteString(fmt.Sprintf("%s%s$%d", k, options.CompareOp, vIdx))
			vIdx++
		}
	}

	query.WriteString(";")

	utils.Sugar.Infof("SQL Query: %s", query.String())
	utils.Sugar.Infof("Values: %v", values)

	_, err = db.Pool.Exec(context.Background(), query.String(), values...)
	if err != nil {
		err = errors.Wrapf(err, "Delete query failed to execute")
		return
	}
	return
}

// DeleteByID removes one row from table where id=id
func (db *Db) DeleteByID(id uuid.UUID, table string) (err error) {
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
