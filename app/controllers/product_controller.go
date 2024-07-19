package controllers

import (
	"net/http"
	"strconv"

	"AntiqueGo/app/models"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func (s *Server) Products(w http.ResponseWriter,r *http.Request) {
	render:= render.New(render.Options{
		Layout:"layout",
		Extensions: []string{".html", ".tmpl"},
	})

	q:=r.URL.Query()
	page,_:=strconv.Atoi(q.Get("page"))
	if page <= 0 {
		page=1
	}

	perPage := 9

	searchQuery := q.Get("search")

	if searchQuery != "" {
        productModel := models.Product{}
        products, totalRows, err := productModel.SearchProducts(s.DB, searchQuery, perPage, page)
		if err!= nil {
            return 
        }
		

		pagination,_:=GetPaginationLinks(s.AppConfig, PaginationParams{
			Path:	"products",
			TotalRows: int64(totalRows),
			PerPage: int64(perPage),
			CurrentPage: int64(page),
		})
	
		cartID := GetShoppingCartID(w, r)
		cart, _ := GetShoppingCart(s.DB, cartID)
		itemCount := len(cart.CartItems)
	
		_ = render.HTML(w,http.StatusOK, "products",map[string]interface{}{
			"products": products,
			"pagination":pagination,
			"user": s.CurrentUser(w,r),
			"itemCount": itemCount,
		})
		return

    } else {
		productModel :=models.Product{}
		products,totalRows,err := productModel.GetProducts(s.DB,perPage,page)
		if err!= nil {
            return 
        }

		
		
		
		pagination,_:=GetPaginationLinks(s.AppConfig, PaginationParams{
			Path:	"products",
			TotalRows: int64(totalRows),
			PerPage: int64(perPage),
			CurrentPage: int64(page),
		})
	
		cartID := GetShoppingCartID(w, r)
		cart, _ := GetShoppingCart(s.DB, cartID)
		itemCount := len(cart.CartItems)
	
		_ = render.HTML(w,http.StatusOK, "products",map[string]interface{}{
			"products": products,
			"pagination":pagination,
			"user": s.CurrentUser(w,r),
			"itemCount": itemCount,
		})
		return

	}


}

func (s *Server) GetProductBySlug(w http.ResponseWriter, r *http.Request){
	render:= render.New(render.Options{
        Layout:"layout",
		Extensions: []string{".html", ".tmpl"},
    })

	vars:= mux.Vars(r)

	if vars["slug"]==""{
		return 
	}

	productModel:= models.Product{}
	product, err := productModel.FindBySlug(s.DB,vars["slug"])
	if err!= nil {
        return 
    }

	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(s.DB, cartID)
	itemCount := len(cart.CartItems)

	_=render.HTML(w,http.StatusOK,"product",map[string]interface{}{
		"product": product,
		"success": GetFlash(w,r,"success"),
		"error": GetFlash(w,r,"error"),
		"user": s.CurrentUser(w,r),
		"itemCount": itemCount,
	})
}