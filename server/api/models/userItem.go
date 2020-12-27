package models

import (
	"fmt"

	"github.com/UN0wen/pricewatch-vn/server/db"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// UserItemTableName is the name of the user table in the db
const (
	UserItemTableName = "useritems"
)

// UserItemTable represents the connection to the db instance
type UserItemTable struct {
	connection *db.Db
}

// UserItem represents a single row in the UserItemTable
type UserItem struct {
	UserID uuid.UUID `valid:"required" json:"userid"`
	ItemID uuid.UUID `valid:"required" json:"itemid"`
}

// UserItemQuery represents all of the rows the item can be queried over
type UserItemQuery struct {
	UserID uuid.UUID
	ItemID uuid.UUID
}

// NewUserItemTable creates a new table in the database for items.
// It takes a reference to an open db connection and returns the constructed table
func NewUserItemTable(db *db.Db) (userItemTable UserItemTable, err error) {
	// Ensure connection is alive
	if db == nil {
		err = errors.New("Invalid database connection")
		return
	}
	userItemTable.connection = db
	// Construct query
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
			userid uuid NOT NULL REFERENCES %s (id) ON DELETE CASCADE,
			itemid uuid NOT NULL REFERENCES %s (id) ON DELETE CASCADE,
			PRIMARY KEY (userid, itemid)
			
		)`, UserItemTableName, UserTableName, ItemTableName)
	// Create the actual table
	if err = userItemTable.connection.CreateTable(query); err != nil {
		err = errors.Wrapf(err, "Could not initialize table %s", UserItemTableName)
	}
	return
}

// Get gets stuffs
func (table *UserItemTable) Get(userItemQuery UserItemQuery, op, compareOp string) (userItems []UserItem, err error) {
	allData, err := table.connection.Get(userItemQuery, op, compareOp, UserItemTableName)
	if err != nil {
		return
	}
	for _, data := range allData {
		userItem := UserItem{}
		err = mapstructure.Decode(data, &userItem)
		if err != nil {
			return
		}
		userItems = append(userItems, userItem)
	}
	return
}

// GetByID finds an item by id
func (table *UserItemTable) GetByID(id uuid.UUID) (userItem UserItem, err error) {
	data, err := table.connection.GetByID(id, UserItemTableName)
	if err != nil {
		return
	}
	err = mapstructure.Decode(data, &userItem)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed for useritem with id: %s", id)
	}
	return
}

// Insert adds a new item into the table.
func (table *UserItemTable) Insert(userItem UserItem) (err error) {
	err = table.connection.Insert(UserItemTableName, userItem)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed for new useritem: %s", userItem)
	}
	return
}

// Update will update the item row with an incoming item
func (table *UserItemTable) Update(id uuid.UUID, newUserItem UserItem) (updated UserItem, err error) {
	data, err := table.connection.Update(id, UserItemTableName, newUserItem)
	if err != nil {
		return
	}
	err = mapstructure.Decode(data[0], &updated)
	return
}

// Delete permanently removes the item with uuid from table
func (table *UserItemTable) Delete(id uuid.UUID) (err error) {
	err = table.connection.Delete(id, UserItemTableName)
	return
}
