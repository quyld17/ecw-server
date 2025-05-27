package orders

import (
	"database/sql"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/E-Commerce-Website/entities/cart"
	products "github.com/quyld17/E-Commerce-Website/entities/product"
	users "github.com/quyld17/E-Commerce-Website/entities/user"
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
	User             users.User     `json:"user"`
}

type OrderProduct struct {
	ID          int    `json:"id"`
	OrderID     int    `json:"order_id"`
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	ImageURL    string `json:"image_url"`
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
			size_name)
		VALUES (?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer orderProduct.Close()

	adjustQuantity, err := transaction.Prepare(`
		UPDATE sizes
		SET quantity = quantity - ?
		WHERE size_id = ? AND product_id = ?;`)
	if err != nil {
		return err
	}
	defer adjustQuantity.Close()

	updateTotalQuantity, err := transaction.Prepare(`
		UPDATE products 
		SET total_quantity = total_quantity - ?
		WHERE product_id = ?;`)
	if err != nil {
		return err
	}
	defer updateTotalQuantity.Close()

	for _, product := range orderedProducts {
		_, err := orderProduct.Exec(orderID, product.ProductID, product.ProductName, product.Quantity, product.Price, product.ImageURL, product.SizeName)
		if err != nil {
			return err
		}
		_, err = adjustQuantity.Exec(product.Quantity, product.SizeID, product.ProductID)
		if err != nil {
			return err
		}
		_, err = updateTotalQuantity.Exec(product.Quantity, product.ProductID)
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
			SELECT *
			FROM order_products
			WHERE order_id = ?;
			`, order.OrderID)
		if err != nil {
			return nil, err
		}
		defer productRows.Close()
		for productRows.Next() {
			var product OrderProduct
			err := productRows.Scan(
				&product.ID, 
				&product.OrderID, 
				&product.ProductID, 
				&product.ProductName, 
				&product.Quantity,
				&product.Price, 
				&product.ImageURL, 
				&product.SizeName)
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

func GetByPageAdmin(offset, limit int, sortParam, search string, db *sql.DB) ([]Order, error) {
	var orderBy string
	switch sortParam {
	case "date_desc":
		orderBy = "o.created_at DESC"
	case "date_asc":
		orderBy = "o.created_at ASC"
	case "amount_desc":
		orderBy = "o.total_price DESC"
	case "amount_asc":
		orderBy = "o.total_price ASC"
	default:
		orderBy = "o.created_at DESC"
	}

	var query string
	var rows *sql.Rows
	var err error

	if search != "" {
		query = `
			SELECT 
				o.order_id,
				o.user_id,
				o.total_price,
				o.status,
				o.address,
				o.created_at,
				o.payment_method,
				u.email,
				u.phone_number,
				u.full_name
			FROM ` + "`orders`" + ` AS o
			JOIN users AS u ON o.user_id = u.user_id
			WHERE o.order_id LIKE ? OR u.email LIKE ? OR u.phone_number LIKE ? OR u.full_name LIKE ?
			ORDER BY ` + orderBy + ` 
			LIMIT ? OFFSET ?;
		`
		rows, err = db.Query(query, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", limit, offset)
	} else {
		query = `
			SELECT 
				o.order_id,
				o.user_id,
				o.total_price,	
				o.status,
				o.address,
				o.created_at,
				o.payment_method,
				u.email,
				u.phone_number,
				u.full_name
			FROM ` + "`orders`" + ` AS o
			JOIN users AS u ON o.user_id = u.user_id
			ORDER BY ` + orderBy + `
			LIMIT ? OFFSET ?;
		`
		rows, err = db.Query(query, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []Order{}
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.OrderID, 
			&order.UserID,
			&order.TotalPrice, 
			&order.Status, 
			&order.Address, 
			&order.CreatedAt,
			&order.PaymentMethod, 
			&order.User.Email, 
			&order.User.PhoneNumber, 
			&order.User.FullName)
		if err != nil {
			return nil, err
		}	
		
		order.CreatedAtDisplay = order.CreatedAt.Format("2006-01-02 15:04:05")

		productRows, err := db.Query(`
			SELECT *
			FROM order_products
			WHERE order_id = ?;
			`, order.OrderID)
		if err != nil {
			return nil, err
		}
		defer productRows.Close()
		
		for productRows.Next() {
			var product OrderProduct
			err := productRows.Scan(
				&product.ID, 
				&product.OrderID, 
				&product.ProductID, 
				&product.ProductName, 
				&product.Quantity,
				&product.Price, 
				&product.ImageURL, 
				&product.SizeName)
			if err != nil {
				return nil, err
			}
			order.Products = append(order.Products, product)
		}
		if err = productRows.Err(); err != nil {
			return nil, err
		}
		
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return orders, nil
}


func Update(orderID int, status string, db *sql.DB) error {
	_, err := db.Exec(`
		UPDATE orders
		SET status = ?
		WHERE order_id = ?;
		`, status, orderID)
	if err != nil {
		return err
	}
	return nil
}