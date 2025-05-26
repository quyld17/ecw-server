package products

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

type Product struct {
	ProductID     int    `json:"product_id"`
	CartProductID int    `json:"cart_product_id"`
	CategoryID    int    `json:"category_id"`
	ProductName   string `json:"product_name"`
	Price         int    `json:"price"`
	ImageURL      string `json:"image_url"`
	TotalQuantity int    `json:"total_quantity"`
	Quantity      int    `json:"quantity"`
	Selected      bool   `json:"selected"`
	SizeID        int    `json:"size_id"`
	SizeName      string `json:"size_name"`
	SizeQuantity  int    `json:"size_quantity"`
}

type ProductImage struct {
	ProductID   int    `json:"product_id"`
	ImageURL    string `json:"image_url"`
	IsThumbnail int    `json:"is_thumbnail"`
}

type ProductSize struct {
	SizeID   int    `json:"size_id"`
	SizeName string `json:"size_name"`
	Quantity int    `json:"quantity"`
}

func GetByPage(c echo.Context, db *sql.DB, limit, offset int) ([]Product, int, error) {
	rows, err := db.Query(`
        SELECT 
			products.product_id, 
			products.category_id, 
			products.product_name, 
			products.price,
			product_images.image_url,
			products.total_quantity
		FROM products
		JOIN product_images ON products.product_id = product_images.product_id
		WHERE 
			product_images.is_thumbnail = 1 AND
			products.total_quantity > 0
		LIMIT ? 
		OFFSET ?;
		`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var numOfProds int
	row := db.QueryRow(`
		SELECT COUNT(*) 
		FROM products;
		`)
	if err := row.Scan(&numOfProds); err != nil {
		return nil, 0, err
	}

	productDetails := []Product{}
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ProductID, &product.CategoryID, &product.ProductName, &product.Price, &product.ImageURL, &product.TotalQuantity)
		if err != nil {
			return nil, 0, err
		}
		productDetails = append(productDetails, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}

	return productDetails, numOfProds, nil
}

func GetProductDetails(productID int, c echo.Context, db *sql.DB) (*Product, []ProductImage, []ProductSize, error) {
	rows, err := db.Query(`
		SELECT 
			products.product_id,
			products.product_name,
			products.price, 
			product_images.image_url, 
			product_images.is_thumbnail 
		FROM products 
		JOIN product_images 
		ON products.product_id = product_images.product_id 
		WHERE products.product_id = ?;
		`, productID)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	productDetail := Product{}
	productImages := []ProductImage{}

	for rows.Next() {
		var product Product
		var productImage ProductImage

		err := rows.Scan(&product.ProductID, &product.ProductName, &product.Price, &productImage.ImageURL, &productImage.IsThumbnail)
		if err != nil {
			return nil, nil, nil, err
		}

		productDetail = product
		productImages = append(productImages, productImage)
	}

	err = rows.Err()
	if err != nil {
		return nil, nil, nil, err
	}

	sizeRows, err := db.Query(`
		SELECT 
			sizes.size_id,
			sizes.size_name,
			size_quantity.quantity
		FROM 
			sizes
		JOIN 
			size_quantity ON sizes.size_id = size_quantity.size_id
		WHERE 
			size_quantity.product_id = ? AND
			size_quantity.quantity > 0;
		`, productID)
	if err != nil {
		return nil, nil, nil, err
	}
	defer sizeRows.Close()

	productSizes := []ProductSize{}
	for sizeRows.Next() {
		var size ProductSize
		err := sizeRows.Scan(&size.SizeID, &size.SizeName, &size.Quantity)
		if err != nil {
			return nil, nil, nil, err
		}
		productSizes = append(productSizes, size)
	}

	err = sizeRows.Err()
	if err != nil {
		return nil, nil, nil, err
	}

	return &productDetail, productImages, productSizes, nil
}

func Search(query string, db *sql.DB) ([]Product, error) {
	rows, err := db.Query(`
		SELECT 
			products.product_id,
			products.product_name,
			products.price,
			product_images.image_url
		FROM 
			products
		JOIN 
			product_images 
		ON 
			products.product_id = product_images.product_id
		WHERE
			products.product_name LIKE ? AND
			product_images.is_thumbnail = 1 AND 
			
		LIMIT 5;
		`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ProductID, &product.ProductName, &product.Price, &product.ImageURL)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
