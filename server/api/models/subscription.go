package models

// TODO: subscriptions table
// UserID -> email to send to

import (
	"fmt"

	"github.com/UN0wen/pricewatch-vn/server/db"
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
	UserID      uuid.UUID `valid:"required" json:"userid"`
	ItemID      uuid.UUID `valid:"required" json:"itemid"`
	Email       string    `valid:"required" json:"email"`
	TargetPrice int64     `valid:"required" json:"targetprice"`
}

// SubscriptionQuery represents all of the rows the item can be queried over
type SubscriptionQuery struct {
	UserID uuid.UUID
	ItemID uuid.UUID
}

// NewSubscriptionTable creates a new table in the database for items.
// It takes a reference to an open db connection and returns the constructed table
func NewSubscriptionTable(db *db.Db) (subscriptionTable SubscriptionTable, err error) {
	// Ensure connection is alive
	if db == nil {
		err = errors.New("Invalid database connection")
		return
	}
	subscriptionTable.connection = db
	// Construct query
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
			userid uuid NOT NULL REFERENCES %s (id) ON DELETE CASCADE,
			itemid uuid NOT NULL REFERENCES %s (id) ON DELETE CASCADE,
			email TEXT NOT NULL, 
			targetprice INT,
			PRIMARY KEY (userid, itemid)
			
		)`, SubscriptionTableName, UserTableName, ItemTableName)
	// Create the actual table
	if err = subscriptionTable.connection.CreateTable(query); err != nil {
		err = errors.Wrapf(err, "Could not initialize table %s", ItemTableName)
	}
	return
}

// Get gets stuffs
func (table *SubscriptionTable) Get(subscriptionQuery SubscriptionQuery) (subscriptions []Subscription, err error) {
	allData, err := table.connection.Get(db.SearchOptions{Query: subscriptionQuery, TableName: SubscriptionTableName, Op: "AND"})
	if err != nil {
		return
	}
	for _, data := range allData {
		subscription := Subscription{}
		err = mapstructure.Decode(data, &subscription)
		if err != nil {
			return
		}
		subscriptions = append(subscriptions, subscription)
	}
	return
}

// Insert adds a new item into the table.
func (table *SubscriptionTable) Insert(subscription Subscription) (err error) {
	err = table.connection.Insert(SubscriptionTableName, subscription)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed for new subscription: %s", subscription)
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
