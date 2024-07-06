package controllers

import (
	"AntiqueGo/app/models"
	// "AntiqueGo/database/seeders"
	// "AntiqueGo/middleware"
	"net/http"

	"github.com/google/uuid"
	"github.com/unrolled/render"
)

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:  "layout",
        Extensions: []string{".html", ".tmpl"},
	})

	_ = render.HTML(w,http.StatusOK,"login", map[string]interface{}{
		"error": GetFlash(w,r,"error"),
	})
}

func (s *Server) DoLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	userModel := models.User{}
	user,err := userModel.FindByEmail(s.DB,email)
	if err!= nil {
        SetFlash(w,r,"error","Invalid email or password")
        http.Redirect(w,r,"/login",http.StatusSeeOther)
        return
    }

	if !ComparePassword(password,user.Password){
		SetFlash(w,r,"error","Invalid email or password")
        http.Redirect(w,r,"/login",http.StatusSeeOther)
        return
    }

	session, _ := store.Get(r,sessionUser)
	session.Values["id"] = user.ID
	session.Save(r, w)

	http.Redirect(w,r,"/",http.StatusSeeOther)

}


func (s *Server) Register (w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:  "layout",
        Extensions: []string{".html", ".tmpl"},
	})

	_ = render.HTML(w,http.StatusOK,"register", map[string]interface{}{
		"error": GetFlash(w,r,"error"),
	})
}

func (s *Server) DoRegister(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if firstName == "" || lastName == "" || email == "" || password == ""{
		SetFlash(w,r,"error","All fields are required")
        http.Redirect(w,r,"/register",http.StatusSeeOther)
        return
    }

	userModel := models.User{}
	exisUser,_ := userModel.FindByEmail(s.DB, email)
	if exisUser!= nil{
        SetFlash(w,r,"error","Sorry, Email already exists")
        http.Redirect(w,r,"/register",http.StatusSeeOther)
        return
    }

	hashedPassword, _ := MakePassword(password)
	params := &models.User{
		ID : 		uuid.New().String(),
		FirstName:     firstName,
        LastName:    lastName,
        Email:         email,
        Password:      hashedPassword,
	}
	user,err:= userModel.CreateUser(s.DB, params)
	if err!= nil {
        SetFlash(w,r,"error","Failed to create user")
        http.Redirect(w,r,"/register",http.StatusSeeOther)
        return
    }
	session, _ := store.Get(r,sessionUser)
	session.Values["id"] = user.ID
	session.Save(r, w)

	http.Redirect(w,r,"/",http.StatusSeeOther)
}

func (s *Server) Logout (w http.ResponseWriter, r *http.Request) {
	session,_ := store.Get(r,sessionUser)

	session.Values["id"] = nil
	session.Save(r, w)

	http.Redirect(w,r,"/",http.StatusSeeOther)
}
