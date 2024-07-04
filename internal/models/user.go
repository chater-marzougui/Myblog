package models

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	Icon      string
}

func CreateUser(username, email, password, icon string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	_, err = DB.Exec("INSERT INTO users (username, email, password, icon) VALUES (?, ?, ?, ?)", username, email, hashedPassword, icon)
	return err
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password must not be empty")
	}
	if len(password) < 8 {
		return "", fmt.Errorf("password must have at least 8 characters")
	}
	hashedPassword := ""
	for v := range password {
		hashedPassword += string(int(password[v]) + 1)
	}
	return string(hashedPassword), nil
}

func AuthenticateUser(usernameOrEmail, password string) (ID int, err error) {
	hashedPassword, err := HashPassword(password)
	var storedHashedPassword string
	var userID int
	if err != nil {
		return 0, err
	}
	typeToCheck := "username"
	for i := range usernameOrEmail {
		if usernameOrEmail[i] == '@' {
			typeToCheck = "email"
			break
		}
	}
	query := fmt.Sprintf("SELECT ID, password FROM users WHERE %s = ?", typeToCheck)
	err = DB.QueryRow(query, usernameOrEmail).Scan(&userID, &storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no user found with the given %s", typeToCheck)
		}
		return 0, fmt.Errorf("error querying user: %v", err)
	}
	if storedHashedPassword != hashedPassword {
		return 0, fmt.Errorf("invalid password")
	}
	return 0, nil
}

func DeleteUser(db *sql.DB, id int) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with the given ID")
	}

	return nil
}
