package controllers

import (
	// "database/sql"
	// "errors"
	"fmt"
	"net/http"
	// "os"
	// "time"

	// "github.com/gieart87/gotoko/app/core/session/auth"
	// "github.com/gieart87/gotoko/app/core/session/flash"

	// "github.com/gorilla/mux"
	"github.com/unrolled/render"

	// "github.com/google/uuid"
	// "github.com/midtrans/midtrans-go"
	// "github.com/midtrans/midtrans-go/snap"

	// "AntiqueGo/app/consts"

	"AntiqueGo/app/models"
	// "github.com/shopspring/decimal"
)

type CheckoutRequest struct {
	Cart            *models.Cart
	ShippingFee     *ShippingFee
	ShippingAddress *ShippingAddress
}

type ShippingFee struct {
	Courier     string
	PackageName string
	Fee         float64
}

type ShippingAddress struct {
	FirstName  string
	LastName   string
	CityID     string
	ProvinceID string
	Address1   string
	Address2   string
	Phone      string
	Email      string
	PostCode   string
}

func (s *Server) Checkout(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:"layout",
        Extensions: []string{".html", ".tmpl"},
    })

	if !IsLoggedIn(r){
		SetFlash(w,r,"error","anda perlu login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	user := s.CurrentUser(w,r)
	fmt.Println(user)

	_ = render.HTML(w,http.StatusOK, "checkout",map[string]interface{}{
		"user": user,
	})
}

// func (server *Server) ShowOrder(w http.ResponseWriter, r *http.Request) {
// 	render := render.New(render.Options{
// 		Layout:     "layout",
// 		Extensions: []string{".html", ".tmpl"},
// 	})

// 	vars := mux.Vars(r)

// 	if vars["id"] == "" {
// 		http.Redirect(w, r, "/products", http.StatusSeeOther)
// 		return
// 	}

// 	orderModel := models.Order{}
// 	order, err := orderModel.FindByID(server.DB, vars["id"])
// 	if err != nil {
// 		http.Redirect(w, r, "/products", http.StatusSeeOther)
// 		return
// 	}

// 	_ = render.HTML(w, http.StatusOK, "show_order", map[string]interface{}{
// 		"order":   order,
// 		"success": flash.GetFlash(w, r, "success"),
// 		"user":    auth.CurrentUser(server.DB, w, r),
// 	})
// }

// func (server *Server) getSelectedShippingCost(w http.ResponseWriter, r *http.Request) (float64, error) {
// 	origin := os.Getenv("API_ONGKIR_ORIGIN")
// 	destination := r.FormValue("city_id")
// 	courier := r.FormValue("courier")
// 	shippingFeeSelected := r.FormValue("shipping_fee")

// 	cartID := GetShoppingCartID(w, r)
// 	cart, _ := GetShoppingCart(server.DB, cartID)

// 	if destination == "" {
// 		return 0, errors.New("invalid destination")
// 	}

// 	shippingFeeOptions, err := server.CalculateShippingFee(models.ShippingFeeParams{
// 		Origin:      origin,
// 		Destination: destination,
// 		Weight:      cart.TotalWeight,
// 		Courier:     courier,
// 	})

// 	if err != nil {
// 		return 0, errors.New("failed shipping calculation")
// 	}

// 	var shippingCost float64
// 	for _, shippingFeeOption := range shippingFeeOptions {
// 		if shippingFeeOption.Service == shippingFeeSelected {
// 			shippingCost = float64(shippingFeeOption.Fee)
// 		}
// 	}

// 	return shippingCost, nil
// }

// func (server *Server) SaveOrder(user *models.User, r *CheckoutRequest) (*models.Order, error) {
// 	var orderItems []models.OrderItem

// 	orderID := uuid.New().String()

// 	paymentURL, err := server.createPaymentURL(user, r, orderID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(r.Cart.CartItems) > 0 {
// 		for _, cartItem := range r.Cart.CartItems {
// 			orderItems = append(orderItems, models.OrderItem{
// 				ProductID:       cartItem.ProductID,
// 				Qty:             cartItem.Qty,
// 				BasePrice:       cartItem.BasePrice,
// 				BaseTotal:       cartItem.BaseTotal,
// 				TaxAmount:       cartItem.TaxAmount,
// 				TaxPercent:      cartItem.TaxPercent,
// 				DiscountAmount:  cartItem.DiscountAmount,
// 				DiscountPercent: cartItem.DiscountPercent,
// 				SubTotal:        cartItem.SubTotal,
// 				Sku:             cartItem.Product.Sku,
// 				Name:            cartItem.Product.Name,
// 				Weight:          cartItem.Product.Weight,
// 			})
// 		}
// 	}

// 	orderCustomer := &models.OrderCustomer{
// 		UserID:     user.ID,
// 		FirstName:  r.ShippingAddress.FirstName,
// 		LastName:   r.ShippingAddress.LastName,
// 		CityID:     r.ShippingAddress.CityID,
// 		ProvinceID: r.ShippingAddress.ProvinceID,
// 		Address1:   r.ShippingAddress.Address1,
// 		Address2:   r.ShippingAddress.Address2,
// 		Phone:      r.ShippingAddress.Phone,
// 		Email:      r.ShippingAddress.Email,
// 		PostCode:   r.ShippingAddress.PostCode,
// 	}

// 	orderData := &models.Order{
// 		ID:                  orderID,
// 		UserID:              user.ID,
// 		OrderItems:          orderItems,
// 		OrderCustomer:       orderCustomer,
// 		Status:              0,
// 		OrderDate:           time.Now(),
// 		PaymentDue:          time.Now().AddDate(0, 0, 7),
// 		PaymentStatus:       consts.OrderPaymentStatusUnpaid,
// 		BaseTotalPrice:      r.Cart.BaseTotalPrice,
// 		TaxAmount:           r.Cart.TaxAmount,
// 		TaxPercent:          r.Cart.TaxPercent,
// 		DiscountAmount:      r.Cart.DiscountAmount,
// 		DiscountPercent:     r.Cart.DiscountPercent,
// 		ShippingCost:        decimal.NewFromFloat(r.ShippingFee.Fee),
// 		GrandTotal:          r.Cart.GrandTotal,
// 		ShippingCourier:     r.ShippingFee.Courier,
// 		ShippingServiceName: r.ShippingFee.PackageName,
// 		PaymentToken:        sql.NullString{String: paymentURL, Valid: true},
// 	}

// 	orderModel := models.Order{}
// 	order, err := orderModel.CreateOrder(server.DB, orderData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return order, nil
// }

// func (server *Server) createPaymentURL(user *models.User, r *CheckoutRequest, orderID string) (string, error) {
// 	midtransServerKey := os.Getenv("API_MIDTRANS_SERVER_KEY")

// 	midtrans.ServerKey = midtransServerKey

// 	var enabledPaymentTypes []snap.SnapPaymentType

// 	enabledPaymentTypes = append(enabledPaymentTypes, snap.AllSnapPaymentType...)

// 	snapRequest := &snap.Request{
// 		TransactionDetails: midtrans.TransactionDetails{
// 			OrderID:  orderID,
// 			GrossAmt: r.Cart.GrandTotal.IntPart(),
// 		},
// 		CustomerDetail: &midtrans.CustomerDetails{
// 			FName: user.FirstName,
// 			LName: user.LastName,
// 			Email: user.Email,
// 		},
// 		EnabledPayments: enabledPaymentTypes,
// 	}

// 	snapResponse, err := snap.CreateTransaction(snapRequest)
// 	if err != nil {
// 		return "", err
// 	}

// 	return snapResponse.RedirectURL, nil
// }
