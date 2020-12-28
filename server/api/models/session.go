package models

import (
	"fmt"
	"time"

	"github.com/UN0wen/pricewatch-vn/server/db"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
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
	UserID       uuid.UUID `valid:"required" json:"userid"`
	ExpiresAfter time.Time `valid:"required" json:"expiresafter"`
}

// SessionQuery represents all of the rows the item can be queried over
type SessionQuery struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

// NewSessionTable creates a new table in the database for sessions.
// It takes a reference to an open db connection and returns the constructed table
func NewSessionTable(db *db.Db) (sessionTable SessionTable, err error) {
	// Ensure connection is alive
	if db == nil {
		err = errors.New("Invalid database connection")
		return
	}
	sessionTable.connection = db
	// Construct query
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
			id uuid NOT NULL,
			userid uuid REFERENCES users(id) ON DELETE CASCADE,
			expiresafter timestamptz NOT NULL,
			PRIMARY KEY (id, userid)
		)`, SessionTableName)
	// Create the actual table
	if err = sessionTable.connection.CreateTable(query); err != nil {
		err = errors.Wrapf(err, "Could not initialize table %s", SessionTableName)
	}

	// Create the view
	viewQuery := fmt.Sprintf(`
	SELECT * FROM %s
	WHERE expiresafter > now()
	`, SessionTableName)

	// Create the view
	if err = sessionTable.connection.CreateView(SessionView, viewQuery); err != nil {
		err = errors.Wrapf(err, "Could not initialize view %s", SessionView)
	}
	return
}

// Get gets stuffs
func (table *SessionTable) Get(sessionQuery SessionQuery) (sessions []Session, err error) {
	allData, err := table.connection.Get(db.GetOptions{Query: sessionQuery, TableName: SessionTableName})
	if err != nil {
		return
	}
	for _, data := range allData {
		session := Session{}
		err = mapstructure.Decode(data, &session)
		if err != nil {
			return
		}
		sessions = append(sessions, session)
	}
	return
}

// GetByID finds a session by id
func (table *SessionTable) GetByID(id uuid.UUID) (session Session, err error) {
	data, err := table.connection.GetByID(id, SessionView)
	if err != nil {
		return
	}
	err = mapstructure.Decode(data, &session)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed for session with id: %s", id)
	}
	return
}

// Insert adds a new session into the table and return the session uuid.
func (table *SessionTable) Insert(session Session) (sessionID uuid.UUID, err error) {
	session.ExpiresAfter = time.Now().Add(time.Hour * 24)
	sessionID, _ = uuid.NewUUID()
	session.ID = sessionID
	err = table.connection.Insert(SessionTableName, session)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed for new session: %s", session)
	}
	return
}

// Delete permanently removes the session with uuid from table
func (table *SessionTable) Delete(id uuid.UUID) (err error) {
	// cascade
	err = table.connection.Delete(id, SessionTableName)
	if err != nil {
		err = errors.Wrapf(err, "Delete query failed for session with id: %s", id)
	}
	return
}
