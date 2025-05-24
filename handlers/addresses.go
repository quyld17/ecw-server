package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	addresses "github.com/quyld17/E-Commerce-Website/entities/address"
	users "github.com/quyld17/E-Commerce-Website/entities/user"
)

func AddAddress(c echo.Context, db *sql.DB) error {
	userID, err := users.GetID(c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var address addresses.Address
	if err := c.Bind(&address); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := addresses.Add(userID, address.City, address.District, address.Ward, address.Street, address.HouseNumber, c, db); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "Address added successfully")
}

func GetAddresses(c echo.Context, db *sql.DB) error {
	userID, err := users.GetID(c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	addresses, err := addresses.Get(userID, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, addresses)
}

func UpdateAddress(c echo.Context, db *sql.DB) error {
	userID, err := users.GetID(c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	addressID, err := strconv.Atoi(c.Param("addressID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var address addresses.Address
	if err := c.Bind(&address); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := addresses.Update(userID, addressID, address.City, address.District, address.Ward, address.Street, address.HouseNumber, c, db); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "Address updated successfully")
}

func SetDefaultAddress(c echo.Context, db *sql.DB) error {
	userID, err := users.GetID(c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	addressID, err := strconv.Atoi(c.Param("addressID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := addresses.SetDefault(userID, addressID, db); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "Set default address successfully")
}

func DeleteAddress(c echo.Context, db *sql.DB) error {
	userID, err := users.GetID(c, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	addressID, err := strconv.Atoi(c.Param("addressID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := addresses.Delete(userID, addressID, db); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "Address deleted successfully")
}	