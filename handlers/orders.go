package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/E-Commerce-Website/entities/cart"
	orders "github.com/quyld17/E-Commerce-Website/entities/order"
	users "github.com/quyld17/E-Commerce-Website/entities/user"
)

func CreateOrder(c echo.Context, db *sql.DB) error {
	userID, err := users.GetID(c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var order orders.Order
	if err := c.Bind(&order); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	orderedProducts, err := cart.GetProducts("true", userID, c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	totalPrice := 0
	for _, product := range orderedProducts {
		if product.Selected {
			totalPrice += product.Quantity * product.Price
		}
	}

	if err := orders.Create(orderedProducts, userID, totalPrice, order.PaymentMethod, order.Address, c, db); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Place order successfully!")
}

func GetOrders(c echo.Context, db *sql.DB) error {
	userID, err := users.GetID(c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	orders, err := orders.GetByPage(userID, c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, orders)
}
