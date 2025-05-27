package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	orders "github.com/quyld17/E-Commerce-Website/entities/order"
	products "github.com/quyld17/E-Commerce-Website/entities/product"
	users "github.com/quyld17/E-Commerce-Website/entities/user"
)

func GetOrdersByPage(c echo.Context, db *sql.DB) error {
	orders, err := orders.GetByPageAdmin(c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, orders)
}


func GetCustomersByPage(c echo.Context, db *sql.DB) error {
	customers, err := users.GetByPage(c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, customers)
}

func UpdateProduct(c echo.Context, db *sql.DB) error {
	var req products.UpdateProductData

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if req.Product.ProductID == 0 || req.Product.Name == "" || req.Product.Price <= 0 || req.Product.TotalQuantity < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing or invalid required fields")
	}

	if err := products.Update(req, db); err != nil {
		if err.Error()[:17] == "invalid size name:" {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update product")
	}

	return c.JSON(http.StatusOK, "Product updated successfully")
}
