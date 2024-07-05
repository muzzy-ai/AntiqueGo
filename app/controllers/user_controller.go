package controllers

import (
    "AntiqueGo/app/models"
    // "AntiqueGo/database/seeders"
    // "AntiqueGo/middleware"
	"net/http"
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

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	session,_ := store.Get(r,sessionUser)

	session.Values["id"] = nil
	session.Save(r, w)

	http.Redirect(w,r,"/",http.StatusSeeOther)
}