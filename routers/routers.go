package routers

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/E-Commerce-Website/handlers"
	"github.com/quyld17/E-Commerce-Website/middlewares"
)

func RegisterAPIHandlers(router *echo.Echo, db *sql.DB) {
	// Authentication
	router.POST("/sign-up", func(c echo.Context) error {
		return handlers.SignUp(c, db)
	})
	router.POST("/sign-in", func(c echo.Context) error {
		return handlers.SignIn(c, db)
	})

	// Users
	router.GET("/users/me", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.GetUserDetails(c, db)
	}))
	router.PUT("/users/password", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.UpdateUserPassword(c, db)
	}))
	router.PUT("/users/me", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.UpdateUserDetails(c, db)
	}))

	// Addresses
	router.GET("/addresses", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.GetAddresses(c, db)
	}))
	router.GET("/default-address", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.GetDefaultAddress(c, db)
	}))
	router.POST("/addresses", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.AddAddress(c, db)
	}))
	router.PUT("/addresses/:addressID", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.UpdateAddress(c, db)
	}))
	router.PUT("/addresses/default/:addressID", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.SetDefaultAddress(c, db)
	}))
	router.DELETE("/addresses/:addressID", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.DeleteAddress(c, db)
	}))

	// Products
	router.GET("/products", func(c echo.Context) error {
		return handlers.GetProductsByPage(c, db)
	})
	router.GET("/products/:productID", func(c echo.Context) error {
		productID := c.Param("productID")
		return handlers.GetProduct(productID, c, db)
	})
	router.GET("/products/search", func(c echo.Context) error {
		return handlers.SearchProducts(c, db)
	})

	// Cart
	router.GET("/cart-products", middlewares.JWTAuthorize(func(c echo.Context) error {
		selected := c.QueryParam("selected")
		return handlers.GetCartProducts(c, db, selected)
	}))
	router.POST("/cart-products", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.AddProductToCart(c, db)
	}))
	router.PUT("/cart-products", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.UpdateCartProducts(c, db)
	}))
	router.DELETE("/cart-products/:cart_product_id", middlewares.JWTAuthorize(func(c echo.Context) error {
		cartProductID := c.Param("cart_product_id")
		return handlers.DeleteCartProduct(cartProductID, c, db)
	}))

	// Orders
	router.GET("/orders/me", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.GetOrders(c, db)
	}))
	router.POST("/orders", middlewares.JWTAuthorize(func(c echo.Context) error {
		return handlers.CreateOrder(c, db)
	}))


	// Admin
	router.GET("/admin/products", middlewares.AdminAuthorize(func(c echo.Context) error {
		return handlers.GetProductsByPage(c, db)
	}))
	router.GET("/admin/orders", middlewares.AdminAuthorize(func(c echo.Context) error {
		return handlers.GetOrdersByPage(c, db)
	}))
	router.GET("/admin/customers", middlewares.AdminAuthorize(func(c echo.Context) error {
		return handlers.GetCustomersByPage(c, db)
	}))
	router.DELETE("/admin/products/:productID", middlewares.AdminAuthorize(func(c echo.Context) error {
		productID := c.Param("productID")
		return handlers.DeleteProduct(productID, c, db)
	}))
	router.PUT("/admin/products", middlewares.AdminAuthorize(func(c echo.Context) error {
		return handlers.UpdateProduct(c, db)
	}))
	router.POST("/admin/products", middlewares.AdminAuthorize(func(c echo.Context) error {
		return handlers.AddProduct(c, db)
	}))
}
