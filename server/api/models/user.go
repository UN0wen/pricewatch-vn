package models

import (
	"fmt"
	"time"

	"github.com/UN0wen/pricewatch-vn/server/db"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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
	TimeCreated  time.Time `valid:"-" json:"time_created,omitempty"`
	TimeLoggedIn time.Time `valid:"-" json:"time_logged_in,omitempty"`
}

// UserQuery represents all of the rows the user can be queried over
type UserQuery struct {
	ID           uuid.UUID
	Email        string
	TimeLoggedIn time.Time
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
		err = errors.Wrapf(err, "Could not initialize table %s", UserTableName)
	}
	return
}

// Login accepts a user object and checks that the user's email is in the database
// and that the passwords match
func (table *UserTable) Login(user User) (found User, err error) {
	if !govalidator.IsEmail(user.Email) {
		err = errors.New("Please provide a valid email address")
		return
	} else if len(user.Password) == 0 {
		err = errors.New("Password can't be blank")
		return
	}
	query := UserQuery{Email: user.Email}

	data, err := table.connection.Get(db.GetOptions{Query: query, TableName: UserTableName})

	if err != nil {
		err = errors.Wrapf(err, "Error querying user with email %s", user.Email)
		return
	} else if data == nil {
		err = errors.New("No user with this email address can be found")
		return
	}

	err = mapstructure.Decode(data[0], &found)
	if err != nil {
		return
	}

	// Compare incoming password with db password
	err = bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(user.Password))
	if err != nil {
		err = errors.Wrapf(err, "Provided password does not match")
	}

	// Update time logged in

	timeNow := time.Now()

	updated, err := table.connection.Update(found.ID, UserTableName, UserQuery{TimeLoggedIn: timeNow})
	if err != nil {
		err = errors.Wrapf(err, "Error while updating time logged in")
		return
	}

	err = mapstructure.Decode(updated[0], &found)
	return
}

// Get gets stuffs
func (table *UserTable) Get(userQuery UserQuery) (users []User, err error) {
	allData, err := table.connection.Get(db.GetOptions{Query: userQuery, TableName: UserTableName})
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
func (table *UserTable) Insert(user User) (userID uuid.UUID, err error) {
	userID, _ = uuid.NewUUID()
	user.ID = userID
	err = table.connection.Insert(UserTableName, user)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed for new user: %s", user)
	}
	return
}

// Update will update the user row with an incoming user
func (table *UserTable) Update(id uuid.UUID, newUser User) (updated User, err error) {
	// Unchangable fields
	newUser.Email = ""
	newUser.ID = id

	data, err := table.connection.Update(id, UserTableName, newUser)
	if err != nil {
		return
	}
	err = mapstructure.Decode(data[0], &updated)
	return
}

// Delete permanently removes the user with uuid from table
func (table *UserTable) Delete(id uuid.UUID) (err error) {
	// Delete user
	err = table.connection.Delete(id, UserTableName)
	if err != nil {
		err = errors.Wrapf(err, "Delete query failed for user with id: %s", id)
	}
	return
}
