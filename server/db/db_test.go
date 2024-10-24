package db

import (
	"database/sql"
	_ "database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test case 1: Successfully creating a new user account
func TestCreateAccount_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sampleUUID := uuid.New()

	// Mock transaction start
	mock.ExpectBegin()

	// Mock SQL insert execution
	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WithArgs("testuser", "test@example.com", "hashedpassword", sampleUUID.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock transaction commit
	mock.ExpectCommit()

	// Call the function
	err = CreateAccount(db, "testuser", "hashedpassword", "test@example.com", sampleUUID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test case 2: Username already exists
func TestCreateAccount_UsernameExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sampleUUID := uuid.New()

	// Mock transaction start
	mock.ExpectBegin()

	// Mock SQL insert execution failure (username already exists)
	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WithArgs("testuser", "test@example.com", "hashedpassword", sampleUUID.String()).
		WillReturnError(sql.ErrNoRows) // Simulate unique constraint failure for username

	// Mock transaction rollback
	mock.ExpectRollback()

	// Call the function
	err = CreateAccount(db, "testuser", "hashedpassword", "test@example.com", sampleUUID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test case 3: Email already exists
func TestCreateAccount_EmailExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sampleUUID := uuid.New()

	// Mock transaction start
	mock.ExpectBegin()

	// Mock SQL insert execution failure (email already exists)
	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WithArgs("testuser", "test@example.com", "hashedpassword", sampleUUID.String()).
		WillReturnError(sql.ErrNoRows) // Simulate unique constraint failure for email

	// Mock transaction rollback
	mock.ExpectRollback()

	// Call the function
	err = CreateAccount(db, "testuser", "hashedpassword", "test@example.com", sampleUUID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test case 4: Transaction begin failure
func TestCreateAccount_TransactionBeginFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Simulate failure to begin a transaction
	mock.ExpectBegin().WillReturnError(sql.ErrTxDone)

	// Call the function
	err = CreateAccount(db, "testuser", "hashedpassword", "test@example.com", uuid.New())

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test case 5: SQL query execution failure
func TestCreateAccount_SQLExecutionFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sampleUUID := uuid.New()

	// Mock transaction start
	mock.ExpectBegin()

	// Mock SQL insert execution failure
	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WithArgs("testuser", "test@example.com", "hashedpassword", sampleUUID.String()).
		WillReturnError(sql.ErrConnDone) // Simulate an execution failure

	// Mock transaction rollback
	mock.ExpectRollback()

	// Call the function
	err = CreateAccount(db, "testuser", "hashedpassword", "test@example.com", sampleUUID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock transaction start
	mock.ExpectBegin()

	// Mock SQL query execution
	mock.ExpectPrepare("SELECT user_id, password FROM users WHERE username = ?").
		ExpectQuery().
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "password"}).
			AddRow("sample-uuid", "hashedpassword"))

	// Mock transaction commit
	mock.ExpectCommit()

	// Call the function
	userID, password, err := Login(db, "testuser")

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, "sample-uuid", userID)
	assert.Equal(t, "hashedpassword", password)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test case 2: Username not found
func TestLogin_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock transaction start
	mock.ExpectBegin()

	// Mock SQL query that returns no rows
	mock.ExpectPrepare("SELECT user_id, password FROM users WHERE username = ?").
		ExpectQuery().
		WithArgs("nonexistentuser").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "password"}))

	// Mock transaction commit
	mock.ExpectCommit()

	// Call the function
	userID, password, err := Login(db, "nonexistentuser")

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, "", userID)
	assert.Equal(t, "", password)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test case 3: SQL execution failure
func TestLogin_SQLExecutionFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock transaction start
	mock.ExpectBegin()

	// Mock SQL query execution failure
	mock.ExpectPrepare("SELECT user_id, password FROM users WHERE username = ?").
		ExpectQuery().
		WithArgs("testuser").
		WillReturnError(sql.ErrConnDone) // Simulate SQL error

	// Mock transaction rollback
	mock.ExpectRollback()

	// Call the function
	userID, password, err := Login(db, "testuser")

	// Assert results
	assert.Error(t, err)
	assert.Equal(t, "", userID)
	assert.Equal(t, "", password)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test case 4: Transaction begin failure
func TestLogin_TransactionBeginFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Simulate failure to begin a transaction
	mock.ExpectBegin().WillReturnError(sql.ErrTxDone)

	// Call the function
	userID, password, err := Login(db, "testuser")

	// Assert results
	assert.Error(t, err)
	assert.Equal(t, "", userID)
	assert.Equal(t, "", password)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test case 5: Scan failure during login
func TestLogin_ScanFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock transaction start
	mock.ExpectBegin()

	// Mock SQL query with a result set that fails to scan
	mock.ExpectPrepare("SELECT user_id, password FROM users WHERE username = ?").
		ExpectQuery().
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "password"}).
			AddRow(nil, "hashedpassword")) // Invalid data for user_id

	// Mock transaction rollback
	mock.ExpectRollback()

	// Call the function
	userID, password, err := Login(db, "testuser")

	// Assert results
	assert.Error(t, err)
	assert.Equal(t, "", userID)
	assert.Equal(t, "", password)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test for valid user ID with groups
func TestGetGroup_ValidUserWithGroups(t *testing.T) {
	// Mock the database connection and response
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"group_id"}).
		AddRow("group1").
		AddRow("group2").
		AddRow("group3")

	// Mock the query execution
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT group_id FROM user_subscriptions where user_id = ?").
		ExpectQuery().
		WithArgs("valid_user").
		WillReturnRows(rows)
	mock.ExpectCommit()

	// Call the function
	groups, err := GetGroup(db, "valid_user")

	// Check the results
	assert.NoError(t, err)
	assert.Equal(t, []string{"group1", "group2", "group3"}, groups)

	// Ensure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test for valid user ID with no groups
func TestGetGroup_ValidUserNoGroups(t *testing.T) {
	// Mock the database connection and response
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock an empty result set (user has no groups)
	rows := sqlmock.NewRows([]string{"group_id"})

	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT group_id FROM user_subscriptions where user_id = ?").
		ExpectQuery().
		WithArgs("valid_user").
		WillReturnRows(rows)
	mock.ExpectCommit()

	// Call the function
	groups, err := GetGroup(db, "valid_user")

	// Check the results
	assert.NoError(t, err)
	assert.Empty(t, groups)

	// Ensure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test for invalid user ID (error scenario)
func TestGetGroup_InvalidUserID(t *testing.T) {
	// Mock the database connection and response
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Simulate an error when the query is executed
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT group_id FROM user_subscriptions where user_id = ?").
		ExpectQuery().
		WithArgs("invalid_user").
		WillReturnError(errors.New("user not found"))
	mock.ExpectRollback()

	// Call the function
	groups, err := GetGroup(db, "invalid_user")

	// Check the results
	assert.Error(t, err)
	assert.Nil(t, groups)

	// Ensure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())
}
