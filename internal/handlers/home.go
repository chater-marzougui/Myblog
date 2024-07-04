package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"myblog/internal/models"
	"net/http"
	"strconv"
	"sync"
)

const mainPage = "/posts"

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("internal/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.ParseFiles("internal/templates/login.html"))
			err := tmpl.Execute(w, nil)
			if err != nil {
				log.Println("Error executing template:", err)
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
			}
		} else if r.Method == http.MethodPost {
			usernameOrEmail := r.FormValue("usernameOrEmail")
			password := r.FormValue("password")

			_, err := models.AuthenticateUser(usernameOrEmail, password)
			if err != nil {
				log.Println("Error authenticating user:", err)
				http.Error(w, "Invalid username/email or password", http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, mainPage, http.StatusSeeOther)
		}
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("internal/templates/register.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		userName := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		pdp := r.FormValue("pdp")
		if confirmPassword != password {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}
		err := models.CreateUser(userName, email, password, pdp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, mainPage, http.StatusSeeOther)
	}
}

func DeleteAccountHandler(db *sql.DB, authenticatedUsers *struct {
	sync.RWMutex
	m map[string]int
}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usernameOrEmail := r.FormValue("usernameOrEmail")

		authenticatedUsers.RLock()
		userID, authenticated := authenticatedUsers.m[usernameOrEmail]
		authenticatedUsers.RUnlock()

		if !authenticated {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		err := models.DeleteUser(db, userID)
		if err != nil {
			log.Println("Error deleting user:", err)
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}

		authenticatedUsers.Lock()
		delete(authenticatedUsers.m, usernameOrEmail)
		authenticatedUsers.Unlock()

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func NewPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("internal/templates/new_post.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		image := r.FormValue("imageLink")
		err := models.CreatePost(title, content, image)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, mainPage, http.StatusSeeOther)
	}
}

func ListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("internal/templates/list_posts.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, posts)
}

func ViewPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/post/"):])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	post, err := models.GetPost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("internal/templates/view_post.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, post)
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/post/edit/"):])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		post, err := models.GetPost(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if post == nil {
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.ParseFiles("internal/templates/edit_post.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, post)
	} else if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		err := models.UpdatePost(id, title, content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/post/"+strconv.Itoa(id), http.StatusSeeOther)
	}
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/post/delete/"):])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = models.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, mainPage, http.StatusSeeOther)
}
