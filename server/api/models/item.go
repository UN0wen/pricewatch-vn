package models

import (
	"context"
	"fmt"

	"github.com/UN0wen/pricewatch-vn/server/db"
	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/asaskevich/govalidator"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// ItemTableName is the name of the item table in the db
// ItemLatestView is the name of the item + price view
const (
	ItemTableName = "items"
)

// ItemTable represents the connection to the db instance
type ItemTable struct {
	connection *db.Db
}

// Item represents a single row in the ItemTable
type Item struct {
	ID          uuid.UUID `valid:"-" json:"id"`
	Name        string    `valid:"required" json:"name"`
	Description string    `valid:"required" json:"description"`
	ImageURL    string    `valid:"required" json:"image_url" db:"image_url"`
	URL         string    `valid:"required" json:"url"`
	Currency    string    `valid:"required" json:"currency"`
}

// ItemQuery represents all of the rows the item can be queried over
type ItemQuery struct {
	ID   uuid.UUID
	Name string
}

// GetAll gets all items from the table
func (table *ItemTable) GetAll() (items []Item, err error) {
	var query string

	query = fmt.Sprintf(`SELECT * FROM %s;`, ItemTableName)

	utils.Sugar.Infof("SQL Query: %s", query)

	err = pgxscan.Select(context.Background(), table.connection.Pool, &items, query)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}
	return
}

// GetByID finds an item by id
func (table *ItemTable) GetByID(id uuid.UUID) (item Item, err error) {
	var query string
	var values []interface{}
	query = fmt.Sprintf(`SELECT * FROM %s WHERE id=$1;`, ItemTableName)

	values = append(values, id)
	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	err = pgxscan.Get(context.Background(), table.connection.Pool, &item, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}

	return
}

// Insert adds a new item into the table.
func (table *ItemTable) Insert(item Item) (returnedItem Item, err error) {
	var query string
	var values []interface{}
	_, err = govalidator.ValidateStruct(item)
	if err != nil {
		err = errors.Wrap(err, "Missing fields in Item")
		return
	}

	values = append(values, item.Name, item.Description, item.ImageURL, item.URL, item.Currency)
	query = fmt.Sprintf(`INSERT INTO "%s" (name, description, image_url, url, currency) VALUES ($1, $2, $3, $4, $5) RETURNING *;`, ItemTableName)

	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	returnedItem = Item{}
	err = pgxscan.Get(context.Background(), table.connection.Pool, &returnedItem, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed to execute")
	}

	return
}

// Update will update the item row with an incoming item
func (table *ItemTable) Update(id uuid.UUID, newItem Item) (updated Item, err error) {
	data, err := table.connection.Update(id, ItemTableName, newItem)
	if err != nil {
		return
	}
	err = mapstructure.Decode(data[0], &updated)
	return
}

// DeleteByID permanently removes the item with uuid from table
func (table *ItemTable) DeleteByID(id uuid.UUID) (err error) {
	// cascade
	err = table.connection.DeleteByID(id, ItemTableName)
	if err != nil {
		err = errors.Wrapf(err, "Delete query failed for item with id: %s", id)
	}
	return
}
