package models

import (
	"context"
	"fmt"

	"github.com/UN0wen/pricewatch-vn/server/db"
	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// UserItemTableName is the name of the user table in the db
const (
	UserItemTableName = "user_items"
)

// UserItemTable represents the connection to the db instance
type UserItemTable struct {
	connection *db.Db
}

// UserItem represents a single row in the UserItemTable
type UserItem struct {
	UserID uuid.UUID `valid:"-" json:"user_id" db:"user_id"`
	ItemID uuid.UUID `valid:"-" json:"item_id" db:"item_id"`
}

// UserItemQuery represents all of the rows the item can be queried over
type UserItemQuery struct {
	UserID uuid.UUID
	ItemID uuid.UUID
}

// GetByUser gets all items followed by user with userID
func (table *UserItemTable) GetByUser(userID uuid.UUID) (items []ItemWithPrice, err error) {
	var query string
	var values []interface{}
	query = fmt.Sprintf(`SELECT i.* FROM %s inner join %s i on user_items.item_id = i.id WHERE user_id=$1;`, UserItemTableName, ItemLatestView)

	values = append(values, userID)
	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	err = pgxscan.Select(context.Background(), table.connection.Pool, &items, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}
	return
}

// GetByUserItem gets an UserItem by both UserID and ItemID
func (table *UserItemTable) GetByUserItem(userID uuid.UUID, itemID uuid.UUID) (returnedUserItem UserItem, err error) {
	var query string
	var values []interface{}
	query = fmt.Sprintf(`SELECT * FROM %s WHERE user_id=$1,item_id=$2;`, UserItemTableName)

	values = append(values, userID, itemID)
	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	returnedUserItem = UserItem{}
	err = pgxscan.Get(context.Background(), table.connection.Pool, &returnedUserItem, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}
	return
}

// Insert adds a new item into the table.
func (table *UserItemTable) Insert(userItem UserItem) (returnedUserItem UserItem, err error) {
	var query string
	var values []interface{}

	if userItem.ItemID == uuid.Nil || userItem.UserID == uuid.Nil {
		err = errors.Wrap(err, "Missing ItemID/UserID in UserItem")
		return
	}

	values = append(values, userItem.UserID, userItem.ItemID)
	query = fmt.Sprintf(`INSERT INTO "%s" (user_id, item_id) VALUES ($1, $2) RETURNING *;`, UserItemTableName)

	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	returnedUserItem = UserItem{}
	err = pgxscan.Get(context.Background(), table.connection.Pool, &returnedUserItem, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed to execute")
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

// DeleteByID permanently removes the item with uuid from table
func (table *UserItemTable) DeleteByID(id uuid.UUID) (err error) {
	err = table.connection.DeleteByID(id, UserItemTableName)
	return
}
