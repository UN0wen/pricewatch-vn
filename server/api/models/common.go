package models

import (
	"sync"

	"github.com/UN0wen/pricewatch-vn/server/db"
	"github.com/UN0wen/pricewatch-vn/server/utils"
)

// Represents the layer for the model by exposing the
// different models' tables.
type layer struct {
	User         *UserTable
	Item         *ItemTable
	UserItem     *UserItemTable
	ItemPrice    *ItemPriceTable
	Session      *SessionTable
	Subscription *SubscriptionTable
}

// Singleton reference to the model layer.
var instance *layer

// Lock for running only once.
var once sync.Once

// LayerInstance gets the static singleton reference
// using double check synchronization.
// It returns the reference to the layer.
func LayerInstance() *layer {
	once.Do(func() {
		// Create DB only once
		db, err := db.Setup(db.Config{
			Host:     utils.DBHost,
			Port:     utils.DBPort,
			User:     utils.DBUser,
			Password: utils.DBPassword,
			Database: utils.DBName,
			URL:      utils.DBUrl,
		})
		utils.CheckError(err)

		// Create the layer only once
		instance = &layer{
			User:         &UserTable{connection: &db},
			Item:         &ItemTable{connection: &db},
			UserItem:     &UserItemTable{connection: &db},
			ItemPrice:    &ItemPriceTable{connection: &db},
			Session:      &SessionTable{connection: &db},
			Subscription: &SubscriptionTable{connection: &db},
		}
	})
	return instance
}
