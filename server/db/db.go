package db

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
)

func CreateAccount(db *sql.DB, username string, password string, email string, uuid uuid.UUID) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to panic:", p)
		} else if err != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	// Prepare the query using the transaction
	query := "INSERT INTO users (username, email, password, user_id) VALUES (?, ?, ?,?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement with user inputs
	_, err = stmt.Exec(username, email, password, uuid.String())
	if err != nil {
		log.Println("Error executing statement:", err)
		return err
	}

	return nil
}

func Login(db *sql.DB, username string) (string, string, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return "", "", err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to panic:", p)
		} else if err != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	query := "SELECT user_id, password FROM users WHERE username = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	defer stmt.Close()

	// Execute the query
	result, err := stmt.Query(username)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	defer result.Close()

	if !result.Next() {
		log.Printf("No user with username: %s found\n", username)
		return "", "", nil
	}

	var userID, password string
	err = result.Scan(&userID, &password) // Scan both user_id and password
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	return userID, password, nil
}
