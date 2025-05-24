package addresses

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Address struct {
	AddressID int    `json: "address_id"`
	Name      string `json: "name"`
	Address   string `json: "address"`
	IsDefault int    `json: "is_default"`
}

func Add(userID int, name, address string, c echo.Context, db *sql.DB) error {
	row := db.QueryRow(`
		SELECT address_id 
		FROM addresses 
		WHERE 	user_id = ? AND 
				name = ? 
		LIMIT 1
		`, userID, name)
	var existingID int
	if err := row.Scan(&existingID); err == nil {
		return fmt.Errorf("Address with this name already exists!")
	}

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
		isDefault = 1
	}

	_, err = db.Exec(`
		INSERT INTO addresses (	user_id, 
								name,
								address,
								is_default)
		VALUES (?, ?, ?, ?);
		`, userID, name, address, isDefault)
	if err != nil {
		return fmt.Errorf("Error adding address! Please try again")
	}
	return nil
}

func Get(userID int, db *sql.DB) ([]Address, error) {
	rows, err := db.Query(`
		SELECT 	address_id, 
				name,
				address, 
				is_default
		FROM addresses
		WHERE user_id = ?;
		`, userID)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Error getting addresses! Please try again")
	}
	defer rows.Close()

	var addresses []Address
	for rows.Next() {
		var address Address
		err := rows.Scan(&address.AddressID, &address.Name, &address.Address, &address.IsDefault)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("Error getting addresses! Please try again")
		}
		addresses = append(addresses, address)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Error getting addresses! Please try again")
	}
	return addresses, nil
}

func GetDefault(userID int, db *sql.DB) (Address, error) {
	row := db.QueryRow(`
		SELECT 	address_id, 
				name, 
				address, 
				is_default
		FROM addresses
		WHERE user_id = ? AND is_default = 1;
		`, userID)

	var address Address
	if err := row.Scan(&address.AddressID, &address.Name, &address.Address, &address.IsDefault); err != nil {
		return Address{}, fmt.Errorf("Error getting default address! Please try again")
	}
	return address, nil
}

func Update(userID int, addressID int, name, address string, c echo.Context, db *sql.DB) error {
	row := db.QueryRow(`
		SELECT address_id 
		FROM addresses 
		WHERE 	user_id = ? AND 
				name = ? AND 
				address_id != ? 
		LIMIT 1
		`, userID, name, addressID)
	var existingID int
	if err := row.Scan(&existingID); err == nil {
		return fmt.Errorf("Another address with this name already exists!")
	}

	_, err := db.Exec(`
		UPDATE addresses
		SET name = ?,
			address = ?
		WHERE user_id = ? AND address_id = ?;
		`, name, address, userID, addressID)

	if err != nil {
		return fmt.Errorf("Error updating address! Please try again")
	}
	return nil
}

func SetDefault(userID int, addressID int, db *sql.DB) error {
	transaction, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	_, err = transaction.Exec(`
		UPDATE addresses 
		SET is_default = 0
		WHERE user_id = ? AND is_default = 1;
		`, userID)
	if err != nil {
		return fmt.Errorf("Error updating address! Please try again")
	}

	_, err = transaction.Exec(`
		UPDATE addresses
		SET is_default = 1
		WHERE user_id = ? AND address_id = ?;
		`, userID, addressID)
	if err != nil {
		return fmt.Errorf("Error updating address! Please try again")
	}

	if err = transaction.Commit(); err != nil {
		return fmt.Errorf("Error updating address! Please try again")
	}

	return nil
}

func Delete(userID int, addressID int, db *sql.DB) error {
	row := db.QueryRow(`
		SELECT is_default
		FROM addresses 
		WHERE user_id = ? AND address_id = ?;
		`, userID, addressID)

	var isDefault int
	if err := row.Scan(&isDefault); err != nil {
		return fmt.Errorf("Error deleting address! Please try again")
	}

	if isDefault == 1 {
		return fmt.Errorf("Can not delete default address")
	}

	_, err := db.Exec(`
		DELETE FROM addresses
		WHERE user_id = ? AND address_id = ?;
		`, userID, addressID)
	if err != nil {
		return fmt.Errorf("Error deleting address! Please try again")
	}
	return nil
}
