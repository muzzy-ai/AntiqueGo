package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func (s *Server) InitializeRoutes() {
	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/", s.Home).Methods("GET")

	s.Router.HandleFunc("/login", s.Login).Methods("GET")
	s.Router.HandleFunc("/login", s.DoLogin).Methods("POST")
	s.Router.HandleFunc("/register", s.Register).Methods("GET")
	s.Router.HandleFunc("/register", s.DoRegister).Methods("POST")
	s.Router.HandleFunc("/logout", s.Logout).Methods("GET")

	s.Router.HandleFunc("/products", s.Products).Methods("GET")
	s.Router.HandleFunc("/products/{slug}", s.GetProductBySlug).Methods("GET")

	s.Router.HandleFunc("/carts",s.GetCart).Methods("GET")
	s.Router.HandleFunc("/carts",s.AddItemToCart).Methods("POST")
	s.Router.HandleFunc("/carts/update",s.UpdateCart).Methods("POST")
	s.Router.HandleFunc("/carts/cities",s.GetCitiesByProvince).Methods("GET")
	s.Router.HandleFunc("/carts/calculate-shipping",s.CalculateShipping).Methods("POST")
	s.Router.HandleFunc("/carts/apply-shipping",s.ApplyShipping).Methods("POST")
	s.Router.HandleFunc("/carts/remove/{id}",s.RemoveItemByID).Methods("GET")

	s.Router.HandleFunc("/orders/checkout", s.Checkout).Methods("POST")
	s.Router.HandleFunc("/orders/{id}", s.ShowOrder).Methods("GET")

	s.Router.HandleFunc("/payments/midtrans", s.Midtrans).Methods("POST")

	s.Router.HandleFunc("/carts/my-order", s.ShowMyOrder).Methods("GET")

	s.Router.HandleFunc("/admin/dashboard", s.AdminDashboard).Methods("GET")
	s.Router.HandleFunc("/products/delete/{id}", s.DeleteProduct).Methods("POST")

	s.Router.HandleFunc("/admin/orders", s.AdminOrders).Methods("GET")
    // s.Router.HandleFunc("/admin/orders/{id}/done", s.MarkOrderAsDone).Methods("POST")

	s.Router.HandleFunc("/admin/products/new", s.NewProduct).Methods("GET")
	s.Router.HandleFunc("/admin/products/create", s.CreateProduct).Methods("POST")


	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/public/", http.FileServer(staticFileDirectory))
	
	// Custom handler to set the correct MIME type for CSS files
	s.Router.PathPrefix("/public/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if filepath.Ext(r.URL.Path) == ".css" {
			w.Header().Set("Content-Type", "text/css")
		}
		staticFileHandler.ServeHTTP(w, r)
	})).Methods("GET")
}