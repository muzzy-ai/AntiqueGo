package controllers

import (
	// "fmt"
	"AntiqueGo/app/models"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	// "AntiqueGo/app/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/unrolled/render"
	"gorm.io/gorm"
)
type ProductWithImage struct {
	ID          string
	Name        string
	Description string
	Price       decimal.Decimal
	Stock       int
	Slug 		string
	ImagePath   string
}

type Order struct {
    // Fields...
    IsDone bool `gorm:"default:false"`
}



func (s *Server) AdminOrders(w http.ResponseWriter, r *http.Request) {
    render := render.New(render.Options{
        Layout:     "admin_layout",
        Extensions: []string{".html", ".tmpl"},
    })

	user := s.CurrentUser(w, r)
    if user == nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    // Check if the user is an admin
	roleModel := models.Role{}
    hasRole, err := roleModel.HasRole(s.DB, user.ID)
    if err != nil {
        // Handle error (you might want to redirect to an error page or log the error)
        SetFlash(w, r, "error", "Failed to check user role")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    if !hasRole {
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }

	

    orderDetails, err := models.GetAllOrdersWithDetails(s.DB)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Send data to the template
    if err := render.HTML(w, http.StatusOK, "admin_orders", map[string]interface{}{
        "orderDetails": orderDetails, // Ensure this matches the template expectation
		"user": user,
    }); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}


// func (s *Server) MarkOrderAsDone(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     orderID := vars["id"]

//     orderModel := models.Order{}
//     if err := s.DB.Model(&orderModel).Where("id = ?", orderID).Update("is_done", true).Error; err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     http.Redirect(w, r, "/admin/orders", http.StatusSeeOther)
// }





func GetProductsWithImages(db *gorm.DB) ([]ProductWithImage, error) {
	var productsWithImages []ProductWithImage

	// Query untuk mengambil satu gambar per produk
	query := `
		SELECT 
			p.id, 
			p.name, 
			p.description, 
			p.price, 
			p.stock, 
			p.slug,
			COALESCE(pi.path, '') as image_path
		FROM products p
		LEFT JOIN (
			SELECT DISTINCT ON (product_id) product_id, path 
			FROM product_images
			ORDER BY product_id, id
		) pi ON p.id = pi.product_id
		WHERE p.deleted_at IS NULL
	`

	err := db.Raw(query).Scan(&productsWithImages).Error
	if err != nil {
		return nil, err
	}

	return productsWithImages, nil
}



func (s *Server) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:     "admin_layout",
		Extensions: []string{".html", ".tmpl"},
	})

	user := s.CurrentUser(w, r)
    if user == nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    // Check if the user is an admin
	roleModel := models.Role{}
    hasRole, err := roleModel.HasRole(s.DB, user.ID)
    if err != nil {
        // Handle error (you might want to redirect to an error page or log the error)
        SetFlash(w, r, "error", "Failed to check user role")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    if !hasRole {
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }


	products, err := GetProductsWithImages(s.DB)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	_ = render.HTML(w, http.StatusOK, "admin_dashboard", map[string]interface{}{
		// "user":     user,
		"products": products,
	})
}

func (s *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	productModel := models.Product{}
	err := productModel.DeleteProduct(s.DB, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func (s *Server) NewProduct(w http.ResponseWriter, r *http.Request) {
	render := render.New(render.Options{
		Layout:     "admin_layout",
		Extensions: []string{".html", ".tmpl"},
	})
	render.HTML(w, http.StatusOK, "admin_add_product", nil)
}

func (s *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // limit your max input length!
	user := s.CurrentUser(w, r)
    if user == nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    // Check if the user is an admin
	roleModel := models.Role{}
    hasRole, err := roleModel.HasRole(s.DB, user.ID)
    if err != nil {
        // Handle error (you might want to redirect to an error page or log the error)
        SetFlash(w, r, "error", "Failed to check user role")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    if !hasRole {
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }

	name := r.FormValue("name")
	description := r.FormValue("description")
	shortDescription := r.FormValue("shortDescription")
	price, _ := decimal.NewFromString(r.FormValue("price"))
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	imageFile, header, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Unable to read image", http.StatusBadRequest)
        return
    }
    defer imageFile.Close()

    // Ambil nama file asli dan ubah spasi menjadi tanda minus
    originalFilename := header.Filename
    filename := strings.ReplaceAll(originalFilename, " ", "-")

    // Definisikan path untuk menyimpan file
    filePath := filepath.Join("assets", "img", "products", filename)
    
    // Buat file baru di path yang ditentukan
    file, err := os.Create(filePath)
    if err != nil {
        fmt.Println(err) // Log kesalahan
        http.Error(w, "Unable to save image", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    // Decode the uploaded image
    img, _, err := image.Decode(imageFile)
    if err != nil {
        http.Error(w, "Unable to decode image", http.StatusInternalServerError)
        return
    }

    // Save the image as JPEG
    if err := jpeg.Encode(file, img, nil); err != nil {
        http.Error(w, "Unable to save image as JPEG", http.StatusInternalServerError)
        return
    }

	ProductID := uuid.New().String()
	sku := strings.ReplaceAll(name, " ", "-")
	slug := strings.ReplaceAll(name, " ", "-")
	filesave := "img"+ "/" + "products" + "/" + filename

	product := models.Product{
		ID:               ProductID,
		UserID: 		user.ID,
		Name:             name,
		Description:      description,
		ShortDescription: shortDescription,
		Price:            price,
		Stock:            stock,
		Sku:              sku,
		Slug:             slug,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ProductImages:    []models.ProductImage{
			{
				ID:         uuid.New().String(),
				ProductID:  ProductID,
				Path:       filesave,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}

	err = s.DB.Create(&product).Error
	if err != nil {
		http.Error(w, "Unable to create product", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
