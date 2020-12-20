package models

import (
	"fmt"

	"github.com/UN0wen/pricewatch-vn/server/db"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// ItemTableName is the name of the user table in the db
const (
	ItemTableName = "items"
)

// ItemTable represents the connection to the db instance
type ItemTable struct {
	connection *db.Db
}

// Item represents a single row in the ItemTable
type Item struct {
	ID           uuid.UUID `valid:"required" json:"id"`
	Name         string    `valid:"-" json:"name"`
	Description  string    `valid:"required" json:"description"`
	ImageURL     string    `valid:"required" json:"imageurl"`
	URL          string    `valid:"required" json:"url"`
	Currency     string    `valid:"required" json:"currency"`
	Availability bool      `valid:"-" json:"availability"`
}

// ItemQuery represents all of the rows the item can be queried over
type ItemQuery struct {
	ID   uuid.UUID
	Name string
}

// NewItemTable creates a new table in the database for items.
// It takes a reference to an open db connection and returns the constructed table
func NewItemTable(db *db.Db) (itemTable ItemTable, err error) {
	// Ensure connection is alive
	if db == nil {
		err = errors.New("Invalid database connection")
		return
	}
	itemTable.connection = db
	// Construct query
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
			id uuid NOT NULL, 
			name TEXT NOT NULL,
			description TEXT,
			imageurl TEXT NOT NULL, 
			url TEXT NOT NULL,
			currency TEXT NOT NULL,
			availability BOOLEAN DEFAULT True,
			PRIMARY KEY (id)
		)`, ItemTableName)
	// Create the actual table
	if err = itemTable.connection.CreateTable(query); err != nil {
		err = errors.Wrapf(err, "Could not initialize table: %s", ItemTableName)
	}
	return
}

// Get gets stuffs
func (table *ItemTable) Get(itemQuery ItemQuery, op string) (items []Item, err error) {
	allData, err := table.connection.Get(itemQuery, op, UserTableName)
	if err != nil {
		return
	}
	for _, data := range allData {
		item := Item{}
		err = mapstructure.Decode(data, &item)
		if err != nil {
			return
		}
		items = append(items, item)
	}
	return
}

// GetByID finds an item by id
func (table *ItemTable) GetByID(id uuid.UUID) (item Item, err error) {
	data, err := table.connection.GetByID(id, ItemTableName)
	if err != nil {
		return
	}
	err = mapstructure.Decode(data, &item)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed for user with id: %s", id)
	}
	return
}

// Insert adds a new item into the table.
func (table *ItemTable) Insert(item Item) (err error) {
	err = table.connection.Insert(ItemTableName, item)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed for new user: %s", item)
	}
	return
}

// Update will update the item row with an incoming item
func (table *ItemTable) Update(id uuid.UUID, newItem Item) (updated Item, err error) {
	data, err := table.connection.Update(id, UserTableName, newItem)
	if err != nil {
		return
	}
	err = mapstructure.Decode(data[0], &updated)
	return
}

// Delete permanently removes the item with uuid from table
// TODO: finish
func (table *ItemTable) Delete(id uuid.UUID) (err error) {
	// TODO: delete all from user-item table
	// Delete user
	err = table.connection.Delete(id, ItemTableName)
	return
}
