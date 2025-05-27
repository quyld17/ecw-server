package cart

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
	products "github.com/quyld17/E-Commerce-Website/entities/product"
)

func GetProducts(selected string, userID int, c echo.Context, db *sql.DB) ([]products.Product, error) {
	var args []interface{}

	query := `
		SELECT 
			cp.id,
			cp.product_id, 
			cp.quantity, 
			cp.selected, 
			cp.size_id,
			p.product_name, 
			p.price, 
			pi.image_url,
			s.size_name,
			s.quantity
		FROM 
			cart_products cp
		JOIN 
			products p ON cp.product_id = p.product_id
		JOIN 
			product_images pi ON cp.product_id = pi.product_id
		JOIN 
			sizes s ON cp.size_id = s.size_id
		WHERE 
			cp.user_id = ? AND 
			pi.is_thumbnail = 1
		`
	args = append(args, userID)

	if selected == "true" {
		query += "AND cp.selected = ?"
		args = append(args, 1)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cartProducts := []products.Product{}
	for rows.Next() {
		var product products.Product
		err := rows.Scan(&product.CartProductID, 
			&product.ProductID, 
			&product.Quantity, 
			&product.Selected, 
			&product.SizeID, 
			&product.ProductName, 
			&product.Price, 
			&product.ImageURL, 
			&product.SizeName, 
			&product.SizeQuantity)
		if err != nil {
			return nil, err
		}
		cartProducts = append(cartProducts, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return cartProducts, nil
}

func UpSertProduct(userID int, productID int, quantity int, sizeID int, c echo.Context, db *sql.DB) error {
	var availableQuantity int
	err := db.QueryRow(`
		SELECT quantity 
		FROM sizes 
		WHERE product_id = ? AND size_id = ?
	`, productID, sizeID).Scan(&availableQuantity)
	if err != nil {
		return fmt.Errorf("Failed to check product availability")
	}

	var existingQuantity int
	err = db.QueryRow(`
		SELECT quantity 
		FROM cart_products 
		WHERE user_id = ? AND product_id = ? AND size_id = ?
	`, userID, productID, sizeID).Scan(&existingQuantity)

	if err == sql.ErrNoRows {
		if quantity > availableQuantity {
			quantity = availableQuantity
		}
		_, err = db.Exec(`
			INSERT INTO cart_products (user_id, product_id, quantity, size_id, selected)
			VALUES (?, ?, ?, ?, 0)
		`, userID, productID, quantity, sizeID)
	} else if err != nil {
		return fmt.Errorf("Failed to check cart")
	} else {
		newQuantity := existingQuantity + quantity
		if newQuantity > availableQuantity {
			newQuantity = availableQuantity
		}
		_, err = db.Exec(`
			UPDATE cart_products 
			SET quantity = ?
			WHERE user_id = ? AND product_id = ? AND size_id = ?
		`, newQuantity, userID, productID, sizeID)
	}

	if err != nil {
		return fmt.Errorf("Failed to add product to cart! Please try again")
	}

	return nil
}

func Update(userID, cartProductID, quantity int, selected bool, c echo.Context, db *sql.DB) error {
	row, err := db.Query(`	
		SELECT * 
		FROM cart_products
		WHERE 
			user_id = ? AND 
			id = ?;
		`, userID, cartProductID)
	if err != nil {
		return err
	}
	defer row.Close()

	if row.Next() {
		if quantity <= 0 {
			_, err := db.Exec(`	
				DELETE FROM cart_products 
				WHERE 
					user_id = ? AND
					id = ?;
				`, userID, cartProductID)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec(`
				UPDATE cart_products
				SET
					quantity = ?,
					selected = ?
				WHERE
					user_id = ? AND
					id = ?;
				`, quantity, selected, userID, cartProductID)
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("Product not in cart. Please try again")
	}

	return nil
}

func DeleteProduct(userID, cartProductID int, c echo.Context, db *sql.DB) error {
	row, err := db.Query(`	
		SELECT * 
		FROM cart_products
		WHERE 
			user_id = ? AND 
			id = ?;
		`, userID, cartProductID)
	if err != nil {
		return err
	}
	defer row.Close()

	if row.Next() {
		_, err := db.Exec(`	
			DELETE FROM cart_products
			WHERE 
				user_id = ? AND 
				id = ?;
			`, userID, cartProductID)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Product not in cart. Please try again")
	}

	return nil
}
