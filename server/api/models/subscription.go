package models

// TODO: subscriptions table
// UserID -> email to send to

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

// SubscriptionTableName is the name of the subscription table in the db
const (
	SubscriptionTableName = "subscriptions"
)

// SubscriptionTable represents the connection to the db instance
type SubscriptionTable struct {
	connection *db.Db
}

// Subscription represents a single row in the UserItemTable
type Subscription struct {
	UserID      uuid.UUID `valid:"required" json:"user_id" db:"user_id"`
	ItemID      uuid.UUID `valid:"required" json:"item_id" db:"json_id"`
	Email       string    `valid:"required" json:"email"`
	TargetPrice int64     `valid:"required" json:"target_price" db:"target_price"`
}

// SubscriptionQuery represents all of the rows the item can be queried over
type SubscriptionQuery struct {
	UserID uuid.UUID
	ItemID uuid.UUID
}

// GetByUser gets all items followed by user with userID
func (table *SubscriptionTable) GetByUser(userID uuid.UUID) (subscriptions []Subscription, err error) {
	var query string
	var values []interface{}

	query = fmt.Sprintf(`SELECT * FROM %s WHERE user_id=$1;`, SubscriptionTableName)

	values = append(values, userID)
	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	err = pgxscan.Select(context.Background(), table.connection.Pool, &subscriptions, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}
	return
}

// GetByItem gets all users who follows an item
func (table *SubscriptionTable) GetByItem(itemID uuid.UUID) (subscriptions []Subscription, err error) {
	var query string
	var values []interface{}

	query = fmt.Sprintf(`SELECT * FROM %s WHERE user_id=$1;`, SubscriptionTableName)

	values = append(values, itemID)
	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	err = pgxscan.Select(context.Background(), table.connection.Pool, &subscriptions, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}
	return
}

// Insert adds a new item into the table.
func (table *SubscriptionTable) Insert(subscription Subscription) (returnedSubscription Subscription, err error) {
	var query string
	var values []interface{}
	_, err = govalidator.ValidateStruct(subscription)
	if err != nil {
		err = errors.Wrap(err, "Missing fields in Subscription")
		return
	}

	if subscription.ItemID == uuid.Nil || subscription.UserID == uuid.Nil {
		err = errors.Wrap(err, "Missing ItemID/UserID in Subscription")
		return
	}

	values = append(values, subscription.UserID, subscription.ItemID, subscription.Email, subscription.TargetPrice)
	query = fmt.Sprintf(`INSERT INTO "%s" (user_id, item_id, email, target_price) VALUES ($1, $2, $3, $4) RETURNING *;`, UserTableName)

	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	returnedSubscription = Subscription{}
	err = pgxscan.Get(context.Background(), table.connection.Pool, &returnedSubscription, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed to execute")
	}

	return
}

// Update will update the item row with an incoming item
func (table *SubscriptionTable) Update(id uuid.UUID, newSubscription Subscription) (updated Subscription, err error) {
	data, err := table.connection.Update(id, SubscriptionTableName, newSubscription)
	if err != nil {
		return
	}
	err = mapstructure.Decode(data[0], &updated)
	return
}

// Delete permanently removes the subscription with uuid from table
func (table *SubscriptionTable) Delete(userID, itemID uuid.UUID) (err error) {
	options := db.SearchOptions{
		Query: SubscriptionQuery{
			UserID: userID,
			ItemID: itemID,
		},
		Op:        "AND",
		TableName: SubscriptionTableName,
	}
	err = table.connection.Delete(options)
	return
}
