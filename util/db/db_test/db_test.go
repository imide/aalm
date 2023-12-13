package db_test

import (
	"github.com/imide/aafl/util/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

var uri string = "mongodb+srv://imide:E4wbWnfxcCBY6L48Eqe9@aafl-db.a0ahcv2.mongodb.net/?retryWrites=true&w=majority"

func TestUpdateMiles(t *testing.T) {
	err := db.UpdateMiles("test_id", 100, db.Add)
	assert.NoError(t, err)
}

func TestUpdateMilesWithInvalidOperation(t *testing.T) {
	err := db.UpdateMiles("test_id", 100, 10) // 10 is not a valid operation
	assert.Error(t, err)
}

func TestExistsWithExistingUser(t *testing.T) {
	exists, err := db.Exists("existing_user_id")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestExistsWithNonExistingUser(t *testing.T) {
	exists, err := db.Exists("non_existing_user_id")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestAddUser(t *testing.T) {
	err := db.AddUser("new_user_id")
	assert.NoError(t, err)
}

func TestGetMiles(t *testing.T) {
	miles, err := db.GetMiles("existing_user_id")
	assert.NoError(t, err)
	assert.Equal(t, int64(100), miles)
}

func TestGetMilesWithNonExistingUser(t *testing.T) {
	_, err := db.GetMiles("non_existing_user_id")
	assert.Error(t, err)
}
