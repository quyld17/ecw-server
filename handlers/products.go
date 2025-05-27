package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	products "github.com/quyld17/E-Commerce-Website/entities/product"
	"github.com/quyld17/E-Commerce-Website/middlewares"
)

func GetProductsByPage(c echo.Context, db *sql.DB) error {
	itemsPerPage := 10
	offset, err := middlewares.Pagination(c, itemsPerPage)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	sortParam := c.QueryParam("sort")
	var orderBy string
	switch sortParam {
		case "price_desc":
			orderBy = "products.price DESC"
		case "price_asc":
			orderBy = "products.price ASC"
		case "name_desc":
			orderBy = "products.product_name DESC"
		case "name_asc":
			orderBy = "products.product_name ASC"
		default:
			orderBy = "products.product_id DESC"
	}

	searchParam := c.QueryParam("search")

	products, numOfProds, err := products.GetByPage(c, db, itemsPerPage, offset, orderBy, searchParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to retrieve products at the moment. Please try again")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"products":     products,
		"num_of_prods": numOfProds,
	})
}

func GetProduct(productID string, c echo.Context, db *sql.DB) error {
	id, err := strconv.Atoi(productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	productDetail, productImages, productSizes, err := products.GetProductDetails(id, c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve product's details")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"product_detail": productDetail,
		"product_images": productImages,
		"product_sizes":  productSizes,
	})
}

func SearchProducts(c echo.Context, db *sql.DB) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusOK, []products.Product{})
	}

	products, err := products.Search(query, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to search products")
	}
	return c.JSON(http.StatusOK, products)
}

func DeleteProduct(productID string, c echo.Context, db *sql.DB) error {
	id, err := strconv.Atoi(productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = products.Delete(id, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "Product deleted successfully")
}

func AddProduct(c echo.Context, db *sql.DB) error {
	var req products.UpdateProductData

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if req.Product.Name == "" || req.Product.Price <= 0 || req.Product.TotalQuantity < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing or invalid required fields")
	}
	
	err := products.Add(req, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add product")
	}

	return c.JSON(http.StatusOK, "Product added successfully")
}
