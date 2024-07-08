package controllers

import (
	"AntiqueGo/app/models"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/unrolled/render"
	"gorm.io/gorm"
)

func GetShoppingCartID(w http.ResponseWriter,r *http.Request) (string) {
	session,_ := store.Get(r,sessionShoppingCart)


	if session.Values["cart-id"] == nil{
		session.Values["cart-id"] = uuid.New().String()
		err := session.Save(r,w)
		if err!= nil {
            http.Error(w,err.Error(),http.StatusInternalServerError)
            return ""
        }
	}
	return fmt.Sprintf("%v",session.Values["cart-id"])

}

func ClearCart(db *gorm.DB, cartID string) error {
	var cart models.Cart

	err := cart.ClearCart(db, cartID)
	if err != nil {
		return err
	}

	return nil
}

func GetShoppingCart(db *gorm.DB, cartID string)(*models.Cart,error){
	var cart models.Cart

	existCart,err := cart.GetCart(db,cartID)
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			existCart,err = cart.CreateCart(db,cartID)
			if err!= nil {
                return nil,err
            }
		}else {
			return nil,err
		}
	}
	_,_=existCart.CalculateCart(db,cartID)
	updatedCart,_:=cart.GetCart(db,cartID)


	totalWeight := 0
	productModel := models.Product{}
	for _,cartItem := range updatedCart.CartItems{
		product,err:= productModel.FindByID(db,cartItem.ProductID)
        if err!= nil {
            return nil,err
        }
        productWeight,_:= product.Weight.Float64()
		ceilWeight := math.Ceil(productWeight)
		itemWeight := cartItem.Qty * int(ceilWeight)
		totalWeight += itemWeight
	}
	updatedCart.TotalWeight = totalWeight
	fmt.Println("total weight ====", updatedCart.TotalWeight)
	return updatedCart,nil
}


func (s *Server) GetCart(w http.ResponseWriter,r *http.Request) {
	render := render.New(render.Options{
		Layout:"layout",
        Extensions: []string{".html", ".tmpl"},
    })
	var cart *models.Cart
	cartID:=GetShoppingCartID(w,r)
	cart,_=GetShoppingCart(s.DB,cartID)
	items,_ := cart.GetItems(s.DB,cartID)
	provinces,err := s.GetProvinces()
	if err!= nil {
        log.Fatal(err)
    }

	fmt.Println("id ===== ",cart.ID)
	fmt.Println("cart item =====", cart.CartItems)

	_= render.HTML(w,http.StatusOK,"cart",map[string]interface{}{
		"cart": cart,
		"items": items,
		"provinces": provinces,
		"success": GetFlash(w,r,"success"),
		"error": GetFlash(w,r,"error"),
		"user": s.CurrentUser(w,r),
	})
}

func (s *Server) AddItemToCart(w http.ResponseWriter,r *http.Request) {
    productID := r.FormValue("product_id")

	qty,_:= strconv.Atoi(r.FormValue("qty"))


	productModel := models.Product{}
	product,err := productModel.FindByID(s.DB,productID)
	if err!= nil {
        http.Redirect(w,r,"/products/"+product.Slug,http.StatusSeeOther)
		return
    }

	if qty > product.Stock{
		SetFlash(w,r,"error","stock tidak mencukupi")
		http.Redirect(w,r,"/products/"+product.Slug,http.StatusFound)
		return
	}


	var cart *models.Cart
	cartID:=GetShoppingCartID(w,r)
	cart,_=GetShoppingCart(s.DB,cartID)

	_,err = cart.AddItem(s.DB, models.CartItem{
		ProductID:       productID,
		Qty:             qty,
	})

	if err!= nil {
        http.Redirect(w,r,"/products/"+product.Slug,http.StatusSeeOther)
    }



	SetFlash(w,r,"success","item berhasil ditambahkan")
	http.Redirect(w,r,"/carts",http.StatusSeeOther)
}

func (s *Server) UpdateCart(w http.ResponseWriter, r *http.Request){
	cartID:=GetShoppingCartID(w,r)
	cart,_:=GetShoppingCart(s.DB,cartID)

	for _,item := range cart.CartItems {
		qty,_ := strconv.Atoi(r.FormValue(item.ID))
		_,err := cart.UpdateItemQty(s.DB,item.ID,qty)
		if err!= nil {
            http.Redirect(w,r,"/carts",http.StatusSeeOther)
        }
	}
	http.Redirect(w,r,"/carts",http.StatusSeeOther)
}

func (s *Server) RemoveItemByID(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)

	if vars["id"] == ""{
		http.Redirect(w,r,"/carts",http.StatusSeeOther)
	}

	cartID:=GetShoppingCartID(w,r)
	cart,_:=GetShoppingCart(s.DB,cartID)

	err := cart.RemoveItemByID(s.DB,vars["id"])
	if err!= nil {
        http.Redirect(w,r,"/carts",http.StatusSeeOther)
    }
	http.Redirect(w,r,"/carts",http.StatusSeeOther)


}

func (s *Server) GetCitiesByProvince(w http.ResponseWriter, r *http.Request){
	provinceID := r.URL.Query().Get("province_id")

	cities, err := s.GetCitiesByProvinceID(provinceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := Result{Code: 200, Data: cities, Message: "Successfully"}
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func (s *Server) CalculateShipping(w http.ResponseWriter, r *http.Request){
	origin := os.Getenv("API_ONGKIR_ORIGIN")
	destination := r.FormValue("city_id")
	courier := r.FormValue("courier")

	if destination == ""{
		http.Error(w, "Invalid destionation", http.StatusInternalServerError)

	}

	cartID := GetShoppingCartID(w,r)
	cart,_:= GetShoppingCart(s.DB,cartID)

	shippingFeeOptions,err:= s.CalculateShippingFee(models.ShippingFeeParams{
		Origin:           origin,
        Destination:      destination,
        Weight:           cart.TotalWeight,
        Courier:          courier,
	})
	if err!= nil {
        http.Error(w,err.Error(), http.StatusInternalServerError)
    }
	res := Result{Code: 200, Data: shippingFeeOptions, Message: "Successfully"}
	result, err := json.Marshal(res)
	if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

	w.Header().Set("Content-Type", "application/json")	
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func (s *Server) ApplyShipping(w http.ResponseWriter, r *http.Request){
	origin := os.Getenv("API_ONGKIR_ORIGIN")
    destination := r.FormValue("city_id")
    courier := r.FormValue("courier")
    shippingPackage := r.FormValue("shipping_package")

    cartID := GetShoppingCartID(w, r)
    cart, _ := GetShoppingCart(s.DB, cartID)

    if destination == "" {
        http.Error(w, "invalid destination", http.StatusInternalServerError)
        return
    }

    shippingFeeOptions, err := s.CalculateShippingFee(models.ShippingFeeParams{
        Origin:      origin,
        Destination: destination,
        Weight:      cart.TotalWeight,
        Courier:     courier,
    })

    if err != nil {
        http.Error(w, "invalid shipping calculation", http.StatusInternalServerError)
        return
    }

    var selectedShipping models.ShippingFeeOption
    found := false
    for _, shippingOption := range shippingFeeOptions {
        if shippingOption.Service == shippingPackage {
            selectedShipping = shippingOption
            found = true
            continue
        }
    }

    if !found {
        http.Error(w, "invalid shipping package", http.StatusInternalServerError)
        return
    }

    type ApplyShippingResponse struct {
        TotalOrder  decimal.Decimal `json:"total_order"`
        ShippingFee decimal.Decimal `json:"shipping_fee"`
        GrandTotal  decimal.Decimal `json:"grand_total"`
        TotalWeight decimal.Decimal `json:"total_weight"`
    }

	var grandTotal float64
    cartGrandTotal, _ := cart.GrandTotal.Float64()
    shippingFee := float64(selectedShipping.Fee)
    grandTotal = cartGrandTotal + shippingFee

    applyShippingResponse := ApplyShippingResponse{
        TotalOrder:  cart.GrandTotal,
        ShippingFee: decimal.NewFromInt(selectedShipping.Fee),
        GrandTotal:  decimal.NewFromFloat(grandTotal),
        TotalWeight: decimal.NewFromInt(int64(cart.TotalWeight)),
    }

    res := Result{Code: 200, Data: applyShippingResponse, Message: "Success"}
    result, err := json.Marshal(res)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(result)
}