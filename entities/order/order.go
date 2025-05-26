package orders

import (
	"database/sql"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/E-Commerce-Website/entities/cart"
	products "github.com/quyld17/E-Commerce-Website/entities/product"
)

type Order struct {
	OrderID          int            `json:"order_id"`
	UserID           int            `json:"user_id"`
	TotalPrice       int            `json:"total_price"`
	PaymentMethod    string         `json:"payment_method"`
	Address          string         `json:"address"`
	AddressID        int            `json:"address_id"`
	Status           string         `json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	CreatedAtDisplay string         `json:"created_at_display"`
	Products         []OrderProduct `json:"products"`
}

type OrderProduct struct {
	ID          int    `json:"id"`
	OrderID     int    `json:"order_id"`
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	ImageURL    string `json:"image_url"`
	SizeID      int    `json:"size_id"`
	SizeName    string `json:"size_name"`
}

func Create(orderedProducts []products.Product, userID, totalPrice int, paymenMethod, address string, c echo.Context, db *sql.DB) error {
	transaction, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	result, err := transaction.Exec(`
		INSERT INTO`+"`orders`"+`
			(user_id, 
			total_price, 
			payment_method,
			address,
			status) 
		VALUES (?, ?, ?, ?, ?)
		`, userID, totalPrice, paymenMethod, address, "Delivering")
	if err != nil {
		return err
	}
	orderID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	orderProduct, err := transaction.Prepare(`	
		INSERT INTO order_products
			(order_id, 
			product_id, 
			product_name, 
			quantity, 
			price, 
			image_url,
			size_id)
		VALUES (?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer orderProduct.Close()

	adjustQuantity, err := transaction.Prepare(`
		UPDATE size_quantity
		SET quantity = quantity - ?
		WHERE size_id = ? AND product_id = ?;`)
	if err != nil {
		return err
	}
	defer adjustQuantity.Close()

	for _, product := range orderedProducts {
		_, err := orderProduct.Exec(orderID, product.ProductID, product.ProductName, product.Quantity, product.Price, product.ImageURL, product.SizeID)
		if err != nil {
			return err
		}
		_, err = adjustQuantity.Exec(product.Quantity, product.SizeID, product.ProductID)
		if err != nil {
			return err
		}
		if err = cart.DeleteProduct(userID, product.CartProductID, c, db); err != nil {
			return err
		}
	}

	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}

func GetByPage(userID int, c echo.Context, db *sql.DB) ([]Order, error) {
	rows, err := db.Query(`
		SELECT 
			order_id,
			total_price,
			status,
			address,
			created_at,
			payment_method
		FROM `+"`orders`"+`
		WHERE user_id = ?
		ORDER BY created_at DESC;
		`, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []Order{}
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.OrderID, &order.TotalPrice, &order.Status, &order.Address, &order.CreatedAt, &order.PaymentMethod)
		if err != nil {
			return nil, err
		}
		order.CreatedAtDisplay = order.CreatedAt.Format("2006-01-02 15:04:05")

		productRows, err := db.Query(`
			SELECT 
				order_products.id,
				order_products.product_id,
				order_products.product_name,
				order_products.quantity,
				order_products.price,
				order_products.image_url,
				order_products.size_id,
				sizes.size_name
			FROM order_products
			JOIN sizes ON order_products.size_id = sizes.size_id
			WHERE order_id = ?;
			`, order.OrderID)
		if err != nil {
			return nil, err
		}
		defer productRows.Close()
		for productRows.Next() {
			var product OrderProduct
			err := productRows.Scan(&product.ID, &product.ProductID, &product.ProductName, &product.Quantity, &product.Price, &product.ImageURL, &product.SizeID, &product.SizeName)
			if err != nil {
				return nil, err
			}
			order.Products = append(order.Products, product)
		}
		err = productRows.Err()
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}	

	return orders, nil
}
