package controllers

import (
	// "fmt"
	"net/http"

	// "AntiqueGo/app/utils"


	"github.com/unrolled/render"
)

func (s *Server) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:     "admin_layout",
		Extensions: []string{".html", ".tmpl"},
	})

	user := s.CurrentUser( w, r)

	// fmt.Println("user ===>", utils.PrintJSON(user))

	_ = render.HTML(w, http.StatusOK, "admin_dashboard", map[string]interface{}{
		"user": user,
	})
}
