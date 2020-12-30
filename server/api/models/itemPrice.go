package models

import (
	"fmt"
	"time"

	"github.com/UN0wen/pricewatch-vn/server/db"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// ItemPriceTableName is the name of the item's price table in the db
const (
	ItemPriceTableName = "itemprices"
)

// ItemPriceTable represents the connection to the db instance
type ItemPriceTable struct {
	connection *db.Db
}

// ItemPrice represents a single row in the ItemPriceTable
type ItemPrice struct {
	ID        uuid.UUID `valid:"required" json:"id"`
	Time      time.Time `valid:"-" json:"time"`
	Price     int64     `valid:"required" json:"price"`
	Available bool      `valid:"required" json:"available"`
}

// ItemPriceQuery represents all of the rows the item can be queried over
type ItemPriceQuery struct {
	ID        uuid.UUID
	Time      time.Time
	Price     int64
	Available bool
}

// NewItemPriceTable creates a new table in the database for items.
// It takes a reference to an open db connection and returns the constructed table
func NewItemPriceTable(db *db.Db) (itemPriceTable ItemPriceTable, err error) {
	// Ensure connection is alive
	if db == nil {
		err = errors.New("Invalid database connection")
		return
	}
	itemPriceTable.connection = db
	// Construct query
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
			id uuid NOT NULL REFERENCES %s(id) ON DELETE CASCADE, 
			time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			price INT,
			available BOOLEAN DEFAULT TRUE,
			PRIMARY KEY (id, time)
		)`, ItemPriceTableName, ItemTableName)
	// Create the actual table
	if err = itemPriceTable.connection.CreateTable(query); err != nil {
		err = errors.Wrapf(err, "Could not initialize table %s", ItemPriceTableName)
	}
	return
}

// Get gets stuffs
func (table *ItemPriceTable) Get(itemPriceQuery ItemPriceQuery, orderBy string, limit int64) (itemPrices []ItemPrice, err error) {
	options := db.SearchOptions{
		Query:      itemPriceQuery,
		TableName:  ItemPriceTableName,
		OrderQuery: "time",
		Order:      orderBy,
		Limit:      limit,
	}

	allData, err := table.connection.Get(options)
	if err != nil {
		return
	}
	for _, data := range allData {
		itemPrice := ItemPrice{}
		err = mapstructure.Decode(data, &itemPrice)
		if err != nil {
			return
		}
		itemPrices = append(itemPrices, itemPrice)
	}
	return
}

// Insert adds a new item into the table.
func (table *ItemPriceTable) Insert(itemPrice ItemPrice) (err error) {
	itemPrice.Time = time.Now()
	err = table.connection.Insert(ItemPriceTableName, itemPrice)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed for new itemprice: %s", itemPrice)
	}
	return
}

// Update will update the item row with an incoming item
func (table *ItemPriceTable) Update(id uuid.UUID, newItemPrice ItemPrice) (updated ItemPrice, err error) {
	data, err := table.connection.Update(id, ItemPriceTableName, newItemPrice)
	if err != nil {
		return
	}
	err = mapstructure.Decode(data[0], &updated)
	return
}

// DeleteByID permanently removes the item with uuid from table
func (table *ItemPriceTable) DeleteByID(id uuid.UUID) (err error) {
	err = table.connection.DeleteByID(id, ItemPriceTableName)
	if err != nil {
		err = errors.Wrapf(err, "Delete query failed for itemprice with id: %s", id)
	}
	return
}
