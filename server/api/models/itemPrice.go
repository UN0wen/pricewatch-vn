package models

import (
	"context"
	"fmt"
	"time"

	"github.com/UN0wen/pricewatch-vn/server/db"
	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/asaskevich/govalidator"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// ItemPriceTableName is the name of the item's price table in the db
const (
	ItemPriceTableName = "item_prices"
)

// ItemPriceTable represents the connection to the db instance
type ItemPriceTable struct {
	connection *db.Db
}

// ItemPrice represents a single row in the ItemPriceTable
type ItemPrice struct {
	ItemID    uuid.UUID `valid:"required" json:"item_id" db:"item_id"`
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

// GetAllPrices gets all prices for a certain item
func (table *ItemPriceTable) GetAllPrices(itemID uuid.UUID) (itemPrices []ItemPrice, err error) {
	var query string
	var values []interface{}

	query = fmt.Sprintf(`SELECT * FROM %s WHERE item_id=$1 ORDER BY time ASC;`, ItemPriceTableName)

	values = append(values, itemID)
	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	err = pgxscan.Select(context.Background(), table.connection.Pool, &itemPrices, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}
	return
}

// GetPrice gets the most current price of an item
func (table *ItemPriceTable) GetPrice(itemID uuid.UUID) (itemPrice ItemPrice, err error) {
	var query string
	var values []interface{}
	query = fmt.Sprintf(`SELECT * FROM %s WHERE item_id=$1 ORDER BY time DESC LIMIT 1;`, ItemPriceTableName)

	values = append(values, itemID)
	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	err = pgxscan.Get(context.Background(), table.connection.Pool, &itemPrice, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}

	return
}

// Insert adds a new item into the table.
func (table *ItemPriceTable) Insert(itemPrice ItemPrice) (returnedItemPrice ItemPrice, err error) {
	var query string
	var values []interface{}
	_, err = govalidator.ValidateStruct(itemPrice)
	if err != nil {
		err = errors.Wrap(err, "Missing fields in ItemPrice")
		return
	}

	values = append(values, itemPrice.ItemID, time.Now().Format(time.RFC3339), itemPrice.Price, itemPrice.Available)
	query = fmt.Sprintf(`INSERT INTO "%s" (item_id, time, price, available) VALUES ($1, $2, $3, $4) RETURNING *;`, ItemPriceTableName)

	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	returnedItemPrice = ItemPrice{}
	err = pgxscan.Get(context.Background(), table.connection.Pool, &returnedItemPrice, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed to execute")
	}

	return
}
