package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	orders "github.com/quyld17/E-Commerce-Website/entities/order"
	products "github.com/quyld17/E-Commerce-Website/entities/product"
	users "github.com/quyld17/E-Commerce-Website/entities/user"
	"github.com/quyld17/E-Commerce-Website/middlewares"
)

func GetOrdersByPage(c echo.Context, db *sql.DB) error {
	itemsPerPage := 10
	offset, err := middlewares.Pagination(c, itemsPerPage)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	sortParam := c.QueryParam("sort")
	searchParam := c.QueryParam("search")

	orders, err := orders.GetByPageAdmin(offset, itemsPerPage, sortParam, searchParam, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, orders)
}


func GetCustomersByPage(c echo.Context, db *sql.DB) error {
	itemsPerPage := 10
	offset, err := middlewares.Pagination(c, itemsPerPage)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	searchParam := c.QueryParam("search")

	customers, err := users.GetByPage(offset, itemsPerPage, searchParam, db)
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


func UpdateOrder(c echo.Context, db *sql.DB) error {
	var order orders.Order
	if err := c.Bind(&order); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if order.OrderID == 0 || order.Status == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing or invalid required fields")
	}

	err := orders.Update(order.OrderID, order.Status, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update order")
	}

	return c.JSON(http.StatusOK, "Order updated successfully")
}

func GetCustomerOrders(customerID string, c echo.Context, db *sql.DB) error {
	id, err := strconv.Atoi(customerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid customer ID")
	}

	orders, err := orders.GetByPage(id, c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get customer orders")
	}
	return c.JSON(http.StatusOK, orders)
}

