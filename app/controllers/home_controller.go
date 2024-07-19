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

	products, err := GetProductsWithImages(s.DB)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	_ = render.HTML(w,http.StatusOK, "home",map[string]interface{}{
		"user": user,
		"itemCount": itemCount,
		"products": products,
	})
}