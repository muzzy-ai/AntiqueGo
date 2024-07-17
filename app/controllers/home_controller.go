package controllers

import (
	// "fmt"
	"net/http"
	

	"github.com/unrolled/render"
)

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
    render:= render.New(render.Options{
		Layout:"layout",
		Extensions: []string{".html", ".tmpl"},
	})

	user := s.CurrentUser(w,r)
	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(s.DB, cartID)
	itemCount := len(cart.CartItems)

	_ = render.HTML(w,http.StatusOK, "home",map[string]interface{}{
		"user": user,
		"itemCount": itemCount,
	})
}