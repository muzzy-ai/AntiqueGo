package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	ID               string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	ParentID         string `gorm:"size:36;index"`
	User             User
	UserID           string `gorm:"size:36;index"`
	ProductImages    []ProductImage
	Category         []Category      `gorm:"many2many:product_categories;"`
	Sku              string          `gorm:"size:100;index"`
	Name             string          `gorm:"size:255"`
	Slug             string          `gorm:"size:255"`
	Price            decimal.Decimal `gorm:"type:decimal(16,2);"`
	Stock            int
	Weight           decimal.Decimal `gorm:"type:decimal(10,2);"`
	ShortDescription string          `gorm:"type:text"`
	Description      string          `gorm:"type:text"`
	Status           int             `gorm:"default:0"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}
type ProductWithImage struct {
	ID          string
	Name        string
	Description string
	Price       decimal.Decimal
	Stock       int
	Slug 		string
	ImagePath   string
}


func (p *Product) GetProducts(db *gorm.DB, perPage int, page int) ([]ProductWithImage, int64, error) {
	var count int64
	var productsWithImages []ProductWithImage

	// Calculate offset
	offset := (page - 1) * perPage

	// Query to get products with images, excluding those with non-null deletedAt
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
		LIMIT ? OFFSET ?
	`

	err := db.Raw(query, perPage, offset).Scan(&productsWithImages).Error
	if err != nil {
		return nil, count, err
	}

	// Query to get the total count of products excluding those with non-null deletedAt
	countQuery := `
		SELECT COUNT(*)
		FROM products
		WHERE deleted_at IS NULL
	`

	err = db.Raw(countQuery).Scan(&count).Error
	if err != nil {
		return nil, count, err
	}

	return productsWithImages, count, nil
}

func (p *Product) FindBySlug(db *gorm.DB, slug string) (*Product, error) {
	var err error
	var product Product

	err = db.Debug().Preload("ProductImages").Model(&Product{}).Where("slug=?", slug).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, err
}

func (p *Product) FindByID(db *gorm.DB, productID string) (*Product, error) {
	var err error
	var product Product

	err = db.Debug().Preload("ProductImages").Model(&Product{}).Where("id=?", productID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, err
}

func (p *Product) SearchProducts(db *gorm.DB, searchQuery string, perPage, page int) ([]Product, int64, error) {
	var products []Product
	var totalRows int64

	offset := (page - 1) * perPage
	query := db.Model(&Product{}).
		Where("deleted_at IS NULL").
		Where("name LIKE ?", "%"+searchQuery+"%").
		Count(&totalRows).
		Offset(offset).
		Limit(perPage).
		Find(&products)

	return products, totalRows, query.Error
}
func (p *Product) GetAllProducts(db *gorm.DB) (*[]Product, error) {
	var err error
	var products []Product

	err = db.Debug().Preload("ProductImages").Order("created_at desc").Find(&products).Error
	if err != nil {
		return nil, err
	}


	return &products, err
}

func (p *Product) DeleteProduct(db *gorm.DB, productID string) error {
    // Enable debug mode to see SQL queries
    db = db.Debug()

    var product Product
    if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("product with ID %s not found: %w", productID, err)
        }
        return fmt.Errorf("error querying product with ID %s: %w", productID, err)
    }

    if err := db.Delete(&product).Error; err != nil {
        return fmt.Errorf("error deleting product: %w", err)
    }

    return nil
}

