package models

import (
	"context"
	"fmt"
	"time"

	"github.com/UN0wen/pricewatch-vn/server/db"
	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// SessionTableName is the name of the session table in the db
// SessionView is the name of the view where only current sessions are shown
const (
	SessionTableName = "sessions"
	SessionView      = "current_sessions"
)

// SessionTable represents the connection to the db instance
type SessionTable struct {
	connection *db.Db
}

// Session represents a single row in the SessionTable
type Session struct {
	ID           uuid.UUID `valid:"required" json:"id"`
	UserID       uuid.UUID `valid:"required" json:"user_id" db:"user_id"`
	ExpiresAfter time.Time `valid:"required" json:"expires_after" db:"expires_after"`
	JWT          string    `valid:"required" json:"jwt"`
}

// SessionQuery represents all of the rows the item can be queried over
type SessionQuery struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

// GetUser returns a single session for user
func (table *SessionTable) GetUser(userID uuid.UUID) (session Session, err error) {
	var query string
	var values []interface{}
	query = fmt.Sprintf(`SELECT * FROM %s WHERE user_id=$1 ORDER BY expires_after DESC LIMIT 1;`, SessionView)

	values = append(values, userID)
	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	session = Session{}
	err = pgxscan.Get(context.Background(), table.connection.Pool, &session, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}

	return
}

// GetSession finds a session by id
func (table *SessionTable) GetSession(sessionID uuid.UUID) (session Session, err error) {
	var query string
	var values []interface{}
	query = fmt.Sprintf(`SELECT * FROM %s WHERE id=$1 ORDER BY expires_after DESC LIMIT 1;`, SessionView)

	values = append(values, sessionID)
	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	session = Session{}
	err = pgxscan.Get(context.Background(), table.connection.Pool, &session, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}

	return
}

// Insert adds a new session into the table and return the new session.
func (table *SessionTable) Insert(userID uuid.UUID) (returnedSession Session, err error) {
	expiresAfter := time.Now().Add(time.Hour * 24 * 7)
	sessionID, _ := uuid.NewRandom()

	jwt := utils.GenerateJWT(sessionID)
	var query string
	var values []interface{}
	if err != nil {
		err = errors.Wrap(err, "Missing fields in ItemPrice")
		return
	}

	values = append(values, sessionID, userID, expiresAfter, jwt)
	query = fmt.Sprintf(`INSERT INTO "%s" (id, user_id, expires_after, jwt) VALUES ($1, $2, $3, $4) RETURNING *;`, SessionTableName)

	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	returnedSession = Session{}
	err = pgxscan.Get(context.Background(), table.connection.Pool, &returnedSession, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed to execute")
	}

	return
}

// DeleteByID permanently removes the session with uuid from table
func (table *SessionTable) DeleteByID(id uuid.UUID) (err error) {
	// cascade
	err = table.connection.DeleteByID(id, SessionTableName)
	if err != nil {
		err = errors.Wrapf(err, "Delete query failed for session with id: %s", id)
	}
	return
}
