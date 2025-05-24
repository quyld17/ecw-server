package addresses

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Address struct {
	AddressID   int    `json:"address_id"`
	City        string `json:"city"`
	District    string `json:"district"`
	Ward        string `json:"ward"`
	Street      string `json:"street"`
	HouseNumber string `json:"house_number"`
	IsDefault   int    `json:"is_default"`
}

func Add(userID int, city, district, ward, street, houseNumber string, c echo.Context, db *sql.DB) error {
	// Check if user has any addresses
	rows, err := db.Query(`
		SELECT address_id 
		FROM addresses 
		WHERE user_id = ?
		LIMIT 1`, userID)
	if err != nil {
		return fmt.Errorf("Error adding address! Please try again")
	}
	defer rows.Close()

	isDefault := 0
	if !rows.Next() {
		// No existing addresses found, set this as default
		isDefault = 1
	}

	_, err = db.Exec(`
		INSERT INTO addresses (	user_id, 
								city, 
								district, 
								ward, 
								street, 
								house_number,
								is_default)
		VALUES (?, ?, ?, ?, ?, ?, ?);
		`, userID, city, district, ward, street, houseNumber, isDefault)
	if err != nil {
		return fmt.Errorf("Error adding address! Please try again")
	}
	return nil
}

func Get(userID int, addressID int, db *sql.DB) (*Address, error) {
	row := db.QueryRow(`
		SELECT 	address_id, 
				city, 
				district, 
				ward, 
				street, 
				house_number, 
				is_default
		FROM addresses
		WHERE user_id = ? AND address_id = ?;
		`, userID, addressID)

	var address Address
	err := row.Scan(&address.AddressID, &address.City, &address.District, &address.Ward, &address.Street, &address.HouseNumber, &address.IsDefault)
	if err != nil {
		return nil, fmt.Errorf("Error getting address! Please try again")
	}
	return &address, nil
}

func Update(userID int, addressID int, city, district, ward, street, houseNumber string, c echo.Context, db *sql.DB) error {
	_, err := db.Exec(`
		UPDATE addresses
		SET city = ?,
			district = ?,
			ward = ?,
			street = ?,
			house_number = ?
		WHERE user_id = ? AND address_id = ?;
		`, city, district, ward, street, houseNumber, userID)

	if err != nil {
		return fmt.Errorf("Error updating address! Please try again")
	}
	return nil
}

func SetDefault(userID int, addressID int, db *sql.DB) error {
	// Start transaction
	transaction, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	// Set current default address to non-default
	_, err = transaction.Exec(`
		UPDATE addresses 
		SET is_default = 0
		WHERE user_id = ? AND is_default = 1;
		`, userID)
	if err != nil {
		return fmt.Errorf("Error updating address! Please try again")
	}

	// Set the specified address as default
	_, err = transaction.Exec(`
		UPDATE addresses
		SET is_default = 1
		WHERE user_id = ? AND address_id = ?;
		`, userID, addressID)
	if err != nil {
		return fmt.Errorf("Error updating address! Please try again")
	}

	// Commit the transaction
	if err = transaction.Commit(); err != nil {
		return fmt.Errorf("Error updating address! Please try again")
	}

	return nil
}