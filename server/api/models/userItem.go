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
	UserID uuid.UUID `valid:"required" json:"user_id" db:"user_id"`
	ItemID uuid.UUID `valid:"required" json:"item_id" db:"item_id"`
}

// UserItemQuery represents all of the rows the item can be queried over
type UserItemQuery struct {
	UserID uuid.UUID
	ItemID uuid.UUID
}

// GetByUser gets all items followed by user with userID
func (table *UserItemTable) GetByUser(userID uuid.UUID) (userItems []UserItem, err error) {
	var query string
	var values []interface{}
	query = fmt.Sprintf(`SELECT * FROM %s WHERE user_id=$1;`, UserItemTableName)

	values = append(values, userID)
	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	err = pgxscan.Select(context.Background(), table.connection.Pool, &userItems, query, values...)
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
	_, err = govalidator.ValidateStruct(userItem)
	if err != nil {
		err = errors.Wrap(err, "Missing fields in UserItem")
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
