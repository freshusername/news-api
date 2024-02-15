package main

import (
	"testing"
)

func TestOpenDBSuccess(t *testing.T) {
	// Arrange
	testDSN := "postgresql://postgres:admin@127.0.0.1:5432/news?sslmode=disable"

	// Act
	db, err := openDB(testDSN)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if db == nil {
		t.Fatal("Expected a database connection, got nil")
	}
	db.Close() // Remember to close the database connection
}

func TestOpenDBFailureInvalidDSN(t *testing.T) {
	// Arrange
	testDSN := "invalid_dsn"

	// Act
	_, err := openDB(testDSN)

	// Assert
	if err == nil {
		t.Fatal("Expected an error for invalid DSN, got none")
	}
}

func TestConnectToDBSuccess(t *testing.T) {
	// Arrange
	testDSN := "postgresql://postgres:admin@127.0.0.1:5432/news?sslmode=disable"
	app := Application{DSN: testDSN}

	// Act
	db, err := app.connectToDB()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var testQuery string = "SELECT 1"
	err = db.QueryRow(testQuery).Scan(&testQuery)

	// Assert
	if err != nil {
		t.Fatalf("Failed to query the database: %v", err)
	}
}

func TestConnectToDBFailureInvalidDSN(t *testing.T) {
	// Arrange
	app := Application{DSN: "invalid_dsn"}

	// Act
	_, err := app.connectToDB()

	// Assert
	if err == nil {
		t.Fatal("Expected an error for invalid DSN, got none")
	}
}
