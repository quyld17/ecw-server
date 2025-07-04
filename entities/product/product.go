package products

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
)

type Product struct {
	ProductID     int    `json:"product_id"`
	CartProductID int    `json:"cart_product_id"`
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
	ProductID int   `json:"product_id"`
	SizeID   int    `json:"size_id"`
	SizeName string `json:"size_name"`
	Quantity int    `json:"quantity"`
}

type UpdateProductData struct {
	Product struct {
		ProductID     int     `json:"product_id"`
		Name          string  `json:"name"`
		Price         float64 `json:"price"`
		TotalQuantity int     `json:"total_quantity"`
	} `json:"product"`
	Sizes []struct {
		SizeName string `json:"size_name"`
		Quantity int    `json:"quantity"`
	} `json:"sizes"`
	ImageURLs []string `json:"image_urls"`
}

func GetByPage(c echo.Context, db *sql.DB, limit, offset int, orderBy, search string) ([]Product, int, error) {
	var query string
	var rows *sql.Rows
	var err error
	var countQuery string
	var count int

	if search == "" {
		query = `
			SELECT 
				products.product_id,
				products.product_name, 
				products.price,
				product_images.image_url,
				products.total_quantity
			FROM products
			JOIN product_images ON products.product_id = product_images.product_id
			WHERE 
				product_images.is_thumbnail = 1 AND
				products.total_quantity > 0
			ORDER BY ` + orderBy + `
			LIMIT ? 
			OFFSET ?
			;`
		rows, err = db.Query(query, limit, offset)

		// Count query
		countQuery = `
			SELECT COUNT(*)
			FROM products
			JOIN product_images ON products.product_id = product_images.product_id
			WHERE 
				product_images.is_thumbnail = 1 AND
				products.total_quantity > 0
		`
		err2 := db.QueryRow(countQuery).Scan(&count)
		if err2 != nil {
			return nil, 0, err2
		}
	} else {
		query = `
			SELECT 
				products.product_id,
				products.product_name, 
				products.price,
				product_images.image_url,
				products.total_quantity
			FROM products
			JOIN product_images ON products.product_id = product_images.product_id
			WHERE 
				product_images.is_thumbnail = 1 AND
				products.product_name LIKE ? AND
				products.total_quantity > 0
			ORDER BY ` + orderBy + `
			LIMIT ? 
			OFFSET ?
			;`
		rows, err = db.Query(query, "%"+search+"%", limit, offset)

		// Count query
		countQuery = `
			SELECT COUNT(*)
			FROM products
			JOIN product_images ON products.product_id = product_images.product_id
			WHERE 
				product_images.is_thumbnail = 1 AND
				products.product_name LIKE ? AND
				products.total_quantity > 0
		`
		err2 := db.QueryRow(countQuery, "%"+search+"%").Scan(&count)
		if err2 != nil {
			return nil, 0, err2
		}
	}

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	productDetails := []Product{}
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ProductID, &product.ProductName, &product.Price, &product.ImageURL, &product.TotalQuantity)
		if err != nil {
			return nil, 0, err
		}
		productDetails = append(productDetails, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}

	return productDetails, count, nil
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
			size_id,
			size_name,
			product_id,
			quantity
		FROM 
			sizes
		WHERE 
			product_id = ? AND
			quantity > 0;
		`, productID)
	if err != nil {
		return nil, nil, nil, err
	}
	defer sizeRows.Close()

	productSizes := []ProductSize{}
	for sizeRows.Next() {
		var size ProductSize
		err := sizeRows.Scan(&size.SizeID, &size.SizeName, &size.ProductID, &size.Quantity)
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
			products.total_quantity > 0
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

func Delete(productID int, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		DELETE FROM product_images
		WHERE product_id = ?;
		`, productID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM sizes
		WHERE product_id = ?;
		`, productID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM products
		WHERE product_id = ?;
		`, productID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func Update(data UpdateProductData, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE products 
		SET 
			product_name = ?,
			price = ?,
			total_quantity = ?
		WHERE product_id = ?`,
		data.Product.Name, data.Product.Price, data.Product.TotalQuantity, data.Product.ProductID)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`	
		DELETE FROM product_images 
		WHERE product_id = ?`, 
		data.Product.ProductID); err != nil {
		return err
	}

	for i, imageURL := range data.ImageURLs {
		isThumbnail := 0
		if i == 0 {
			isThumbnail = 1
		}
		if _, err := tx.Exec(
			`INSERT INTO product_images (product_id, image_url, is_thumbnail) 
			VALUES (?, ?, ?)`,
			data.Product.ProductID, imageURL, isThumbnail,
		); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(`	
		DELETE FROM sizes 
		WHERE product_id = ?`, 
		data.Product.ProductID); err != nil {
		return err
	}

	for _, size := range data.Sizes {
		if _, err := tx.Exec(`
			INSERT INTO sizes (size_name, product_id, quantity) 
			VALUES (?, ?, ?)`,
			size.SizeName, data.Product.ProductID, size.Quantity); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func Add(data UpdateProductData, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	result, err := tx.Exec(`
		INSERT INTO products (product_name, price, total_quantity) 
		VALUES (?, ?, ?)`,
		data.Product.Name, data.Product.Price, data.Product.TotalQuantity)
	if err != nil {
		return err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	for i, imageURL := range data.ImageURLs {
		isThumbnail := 0
		if i == 0 {
			isThumbnail = 1
		}
		if _, err := tx.Exec(`
			INSERT INTO product_images (product_id, image_url, is_thumbnail) 
			VALUES (?, ?, ?)`,
			productID, imageURL, isThumbnail); err != nil {
			return err
		}
	}

	for _, size := range data.Sizes {
		_, err := tx.Exec(`
			INSERT INTO sizes (size_name, product_id, quantity) 
			VALUES (?, ?, ?)`,
			size.SizeName, productID, size.Quantity)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func CheckProductExists(productID int, db *sql.DB) error {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM products
			WHERE product_id = ?
		)
		`, productID).Scan(&exists)

	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Product not found")
	}

	return nil
}
