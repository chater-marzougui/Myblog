package models

import (
	"database/sql"
	"errors"
	"myblog/internal/global"
	"time"
)

type Post struct {
	ID        int
	Title     string
	Content   string
	Image     string
	UserID    int
	CreatedAt time.Time
}

func CreatePost(title, content, image string) error {
	_, err := DB.Exec("INSERT INTO posts (title, content, image ,user_ID) VALUES (?, ?, ?, ?)", title, content, image, global.GetUserID())
	return err
}

func GetPosts() ([]Post, error) {
	rows, err := DB.Query("SELECT id, title, image, content, created_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Image, &post.Content, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPost(id int) (*Post, error) {
	row := DB.QueryRow("SELECT id, title, content, image, created_at FROM posts WHERE id = ?", id)
	var post Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &post, nil
}

func GetUser(id int) (*User, error) {
	row := DB.QueryRow("SELECT ID, Username, Icon, Created_At, email, Password FROM users WHERE id = ?", id)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Icon, &user.CreatedAt, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdatePost(id int, title, content string) error {
	_, err := DB.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", title, content, id)
	return err
}

func DeletePost(id int) error {
	_, err := DB.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}

func UpdateUser(user User) error {
	query := `
        UPDATE users
        SET username = ?, email = ?, password = ?, icon = ?
        WHERE id = ?
    `
	_, err := DB.Exec(query, user.Username, user.Email, user.Password, user.Icon, user.ID)
	if err != nil {
		return errors.New("failed to update user: " + err.Error())
	}
	return nil
}
