package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/UN0wen/pricewatch-vn/server/db"
	"github.com/UN0wen/pricewatch-vn/server/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// UserTableName is the name of the user table in the db
const (
	UserTableName = "users"
)

// UserTable represents the connection to the db instance
type UserTable struct {
	connection *db.Db
}

// User represents a single row in the UserTable
type User struct {
	ID           uuid.UUID `valid:"required" json:"id"`
	Username     string    `valid:"required" json:"username"`
	Email        string    `valid:"required" json:"email"`
	Password     string    `valid:"required" json:"password"`
	TimeCreated  time.Time `valid:"-" json:"time_created"`
	TimeLoggedIn time.Time `valid:"-" json:"time_logged_in"`
}

// UserQuery represents all of the rows the user can be queried over
type UserQuery struct {
	ID uuid.UUID
}

// NewUser is used to create a new user object from an incoming HTTP request.
// It takes in the HTTP request in JSON format.
// It returns the user constructed and an error if one exists.
func NewUser(r *http.Request) (user User, err error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		err = errors.Wrapf(err, "Couldn't read request body")
	}
	// Converts JSON to user
	json.Unmarshal(b, &user)
	return user, err
}

// GenerateJWT creates a JSON Web Token for a user based on the id,
// with an expiration time of 1 day
// It returns the token string
func (user *User) GenerateJWT() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(utils.GPTokenSecret)
	utils.CheckError(err)
	return tokenString
}

// NewUserTable creates a new table in the database for users.
// It takes a reference to an open db connection and returns the constructed table
func NewUserTable(db *db.Db) (userTable UserTable, err error) {
	// Ensure connection is alive
	if db == nil {
		err = errors.New("Invalid database connection")
		return
	}
	userTable.connection = db
	// Construct query
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
			id uuid NOT NULL, 
			username TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE, 
			password TEXT NOT NULL,
			timeCreated TIMESTAMPTZ NOT NULL DEFAULT now(), 
			timeLoggedIn TIMESTAMPTZ NOT NULL DEFAULT now(), 
			PRIMARY KEY (id)
		)`, UserTableName)
	// Create the actual table
	if err = userTable.connection.CreateTable(query); err != nil {
		err = errors.Wrapf(err, "Could not initialize table: %s", UserTableName)
	}
	return
}

// Get gets stuffs
func (table *UserTable) Get(userQuery UserQuery, op, compareOp string) (users []User, err error) {
	allData, err := table.connection.Get(userQuery, op, compareOp, UserTableName)
	if err != nil {
		return
	}
	for _, data := range allData {
		user := User{}
		err = mapstructure.Decode(data, &user)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	return
}

// GetByID finds a user with id
func (table *UserTable) GetByID(id uuid.UUID) (user User, err error) {
	data, err := table.connection.GetByID(id, UserTableName)
	if err != nil {
		return
	}
	err = mapstructure.Decode(data, &user)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed for user with id: %s", id)
	}
	return
}

// Insert adds a new user into the table.
func (table *UserTable) Insert(user User) (err error) {
	err = table.connection.Insert(UserTableName, user)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed for new user: %s", user)
	}
	return
}

// Update will update the user row with an incoming user
func (table *UserTable) Update(id uuid.UUID, newUser User) (updated User, err error) {
	data, err := table.connection.Update(id, UserTableName, newUser)
	if err != nil {
		return
	}
	err = mapstructure.Decode(data[0], &updated)
	return
}

// Delete permanently removes the user with uuid from table
// TODO: finish
func (table *UserTable) Delete(id uuid.UUID) (err error) {
	// TODO: delete all from user-product table
	// Delete user
	err = table.connection.Delete(id, UserTableName)
	return
}
