package models

import (
	"database/sql"
	"errors"
	"fmt"
	"myblog/internal/global"
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
	mailTest := false
	for i := range email {
		if email[i] == '@' {
			mailTest = true
			break
		}
	}
	if !mailTest {
		return fmt.Errorf("invalid email format")
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

func AuthenticateUser(usernameOrEmail, password string) error {
	hashedPassword, err := HashPassword(password)
	var storedHashedPassword string
	var userID int
	if err != nil {
		return err
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
			return fmt.Errorf("no user found with the given %s", typeToCheck)
		}
		return fmt.Errorf("error querying user: %v", err)
	}
	if storedHashedPassword != hashedPassword {
		return fmt.Errorf("invalid password")
	}
	global.ModifyUser(userID)
	global.SetAuthenticated()
	return nil
}

func DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := DB.Exec(query, id)
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

func GetUserPosts(userID int) ([]Post, error) {
	query := "SELECT id, title, content, image, user_id FROM posts WHERE user_id = ?"
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, errors.New("failed to get posts: " + err.Error())
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.UserID)
		if err != nil {
			return nil, errors.New("failed to scan post: " + err.Error())
		}
		posts = append(posts, post)
	}
	return posts, nil
}
