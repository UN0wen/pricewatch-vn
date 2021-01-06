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
	ID       uuid.UUID `valid:"-" json:"id"`
	Username string    `valid:"required" json:"username"`
	Email    string    `valid:"required" json:"email"`
	Password string    `valid:"required" json:"password"`
	Created  time.Time `valid:"-" json:"created" db:"created"`
	LoggedIn time.Time `valid:"-" json:"logged_in" db:"logged_in"`
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

	found, err = table.GetByEmail(user.Email)

	if err != nil {
		err = errors.Wrapf(err, "Error querying user with email %s", user.Email)
		return
	} else if found == (User{}) {
		err = errors.New("No user with this email address can be found")
		return
	}

	// Compare incoming password with db password
	err = bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(user.Password))
	if err != nil {
		err = errors.Wrapf(err, "Provided password does not match")
		return
	}

	// Update time logged in

	timeNow := time.Now()

	updated, err := table.connection.Update(found.ID, UserTableName, User{LoggedIn: timeNow})
	if err != nil {
		err = errors.Wrapf(err, "Error while updating time logged in")
		return
	}

	err = mapstructure.Decode(updated[0], &found)
	return
}

// GetByEmail gets an user by email
func (table *UserTable) GetByEmail(email string) (user User, err error) {
	var query string
	var values []interface{}
	query = fmt.Sprintf(`SELECT * FROM %s WHERE email=$1;`, UserTableName)

	values = append(values, email)
	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	user = User{}
	err = pgxscan.Get(context.Background(), table.connection.Pool, &user, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}

	return
}

// GetByID finds a user with id
func (table *UserTable) GetByID(id uuid.UUID) (user User, err error) {
	var query string
	var values []interface{}
	query = fmt.Sprintf(`SELECT * FROM %s WHERE id=$1;`, UserTableName)

	values = append(values, id)
	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	user = User{}
	err = pgxscan.Get(context.Background(), table.connection.Pool, &user, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Get query failed to execute")
		return
	}

	return
}

// Insert adds a new user into the table.
func (table *UserTable) Insert(user User) (returnedUser User, err error) {
	var query string
	var values []interface{}
	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		err = errors.Wrap(err, "Missing fields in User")
		return
	}

	hash, e := utils.HashPassword(user.Password)
	if e != nil {
		err = errors.Wrapf(e, "Password hash failed")
		return
	}
	values = append(values, user.Email, user.Username, hash)
	query = fmt.Sprintf(`INSERT INTO "%s" (email, username, password) VALUES ($1, $2, $3) RETURNING *;`, UserTableName)

	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %s", values)

	returnedUser = User{}
	err = pgxscan.Get(context.Background(), table.connection.Pool, &returnedUser, query, values...)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed to execute")
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

// DeleteByID permanently removes the user with uuid from table
func (table *UserTable) DeleteByID(id uuid.UUID) (err error) {
	// Delete user
	err = table.connection.DeleteByID(id, UserTableName)
	if err != nil {
		err = errors.Wrapf(err, "Delete query failed for user with id: %s", id)
	}
	return
}
