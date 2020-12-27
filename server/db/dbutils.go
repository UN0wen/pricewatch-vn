package db

import (
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func mapRows(rows pgx.Rows) (objects []map[string]interface{}, err error) {
	// Get all the rows that matched the query
	cols := rows.FieldDescriptions()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPtrs := make([]interface{}, len(cols))
		for i := range columns {
			columnPtrs[i] = &columns[i]
		}
		if err = rows.Scan(columnPtrs...); err != nil {
			err = errors.Wrapf(err, "Get query failed to execute")
			return
		}
		m := make(map[string]interface{})
		for i, col := range cols {
			val := columnPtrs[i].(*interface{})
			m[strings.Title(string(col.Name))] = *val
		}
		objects = append(objects, m)
	}

	return
}

// Hashes a password string
func hashPassword(pass string) (hash []byte, err error) {
	hash, err = bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		err = errors.Wrapf(err, "Password hash failed")
		return
	}
	return
}
