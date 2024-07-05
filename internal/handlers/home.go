package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"myblog/internal/global"
	"myblog/internal/models"
	"net/http"
	"strconv"
	"sync"
)

const mainPage = "/posts"
const notAuthenticated = "User not authenticated"

type PageData struct {
	Posts []models.Post
	User  models.User
}

type PostPageData struct {
	Post models.Post
	User models.User
}

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

			err := models.AuthenticateUser(usernameOrEmail, password)
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
			http.Error(w, notAuthenticated, http.StatusUnauthorized)
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
		data, _ := models.GetUser(global.GetUserID())
		tmpl.Execute(w, data)
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

func ViewProfile(w http.ResponseWriter, r *http.Request) {
	if !global.IsAuthenticated() {
		http.Error(w, notAuthenticated, http.StatusUnauthorized)
		return
	}

	userID := global.GetUserID()

	user, err := models.GetUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts, err := models.GetUserPosts(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		User  models.User
		Posts []models.Post
	}{
		User:  *user,
		Posts: posts,
	}

	tmpl, err := template.ParseFiles("internal/templates/view_profile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ListPosts(w http.ResponseWriter, r *http.Request) {
	if !global.IsAuthenticated() {
		http.Error(w, notAuthenticated, http.StatusUnauthorized)
		return
	}
	posts, err := models.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := models.GetUser(global.GetUserID())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Posts: posts,
		User:  *user,
	}

	tmpl, err := template.ParseFiles("internal/templates/list_posts.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ViewPost(w http.ResponseWriter, r *http.Request) {
	if !global.IsAuthenticated() {
		http.Error(w, notAuthenticated, http.StatusUnauthorized)
		return
	}
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
	user, err := models.GetUser(global.GetUserID())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PostPageData{
		Post: *post,
		User: *user,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
		user, _ := models.GetUser(global.GetUserID())
		data := PostPageData{
			Post: *post,
			User: *user,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

func EditProfile(w http.ResponseWriter, r *http.Request) {
	if !global.IsAuthenticated() {
		http.Error(w, notAuthenticated, http.StatusUnauthorized)
		return
	}

	userID := global.GetUserID()

	if r.Method == http.MethodGet {
		user, err := models.GetUser(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("internal/templates/edit_profile.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, user)
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password, _ := models.HashPassword(r.FormValue("password"))
		icon := r.FormValue("icon")

		user := models.User{
			ID:       userID,
			Username: username,
			Email:    email,
			Password: password,
			Icon:     icon,
		}

		err := models.UpdateUser(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}
