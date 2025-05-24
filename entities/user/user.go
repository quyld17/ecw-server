package users

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	address "github.com/quyld17/E-Commerce-Website/entities/address"
)

type User struct {
	UserId            int       `json:"user_id"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	NewPassword       string    `json:"new_password"`
	FullName          string    `json:"full_name"`
	DateOfBirth       time.Time `json:"date_of_birth"`
	DateOfBirthString string    `json:"date_of_birth_string"`
	PhoneNumber       string    `json:"phone_number"`
	Gender            int       `json:"gender"`
}

func Authenticate(account User, db *sql.DB) error {
	rows, err := db.Query(`	
		SELECT 
			email, 
			password 
		FROM users
		WHERE 
			email = ? AND 
			password = ?
		`, account.Email, account.Password)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return nil
	}
	return fmt.Errorf("Invalid email or password! Please try again")
}

func Create(newUser User, db *sql.DB) error {
	_, err := db.Exec(`	
		INSERT INTO users (email, password) 
		VALUES (?, ?)
		`, newUser.Email, newUser.Password)
	if err != nil {
		return err
	}
	return nil
}

func GetDetails(userID int, db *sql.DB) (*User, *address.Address, error) {
	row, err := db.Query(`
		SELECT
			email,
			full_name,
			phone_number,
			gender,
			date_of_birth
		FROM users
		WHERE user_id = ?;
		`, userID)
	if err != nil {
		return nil, nil, err
	}

	var user User
	if row.Next() {
		var nullFullName, nullPhoneNumber sql.NullString
		var nullGender sql.NullInt64
		var nullDateOfBirth sql.NullTime
		var email string

		err := row.Scan(&email, &nullFullName, &nullPhoneNumber, &nullGender, &nullDateOfBirth)
		if err != nil {
			return nil, nil, err
		}

		user.Email = email
		user.FullName = nullFullName.String
		user.PhoneNumber = nullPhoneNumber.String
		user.Gender = int(nullGender.Int64)

		if nullDateOfBirth.Valid {
			user.DateOfBirth = nullDateOfBirth.Time
			user.DateOfBirthString = user.DateOfBirth.Format("2006-01-02")
		}
	}
	defer row.Close()

	row, err = db.Query(`
		SELECT
			city,
			district,
			ward,
			street,
			house_number
		FROM addresses
		WHERE
			user_id = ? AND
			is_default = 1;
		`, userID)
	if err != nil {
		return nil, nil, err
	}
	defer row.Close()

	var address address.Address
	if row.Next() {
		var nullCity, nullDistrict, nullWard, nullStreet, nullHouseNumber sql.NullString
		err := row.Scan(&nullCity, &nullDistrict, &nullWard, &nullStreet, &nullHouseNumber)
		if err != nil {
			return nil, nil, err
		}

		address.City = nullCity.String
		address.District = nullDistrict.String
		address.Ward = nullWard.String
		address.Street = nullStreet.String
		address.HouseNumber = nullHouseNumber.String
	}

	return &user, &address, nil
}

func GetID(c echo.Context, db *sql.DB) (int, error) {
	email := c.Get("email").(string)
	row := db.QueryRow(`
		SELECT user_id 
		FROM users
		WHERE email = ?;
		`, email)
	var userID int
	if err := row.Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}

func ChangePassword(userID int, password, newPassword string, c echo.Context, db *sql.DB) error {
	row, err := db.Query(`
		SELECT password
		FROM users
		WHERE
			user_id = ? AND
			password = ?;
		`, userID, password)
	if err != nil {
		return fmt.Errorf("Error while changing password! Please try again")
	}
	defer row.Close()
	if row.Next() {
		_, err := db.Exec(`	
			UPDATE users
			SET password = ? 
			WHERE user_id = ?;
			`, newPassword, userID)
		if err != nil {
			return fmt.Errorf("Error while changing password! Please try again")
		}
	} else {
		return fmt.Errorf("Wrong password! Plase try again")
	}

	return nil
}

func UpdateDetails(userID int, fullName, phoneNumber string, gender int, dateOfBirth time.Time, c echo.Context, db *sql.DB) error {
	_, err := db.Exec(`
		UPDATE users
		SET full_name = ?, 
			phone_number = ?, 
			gender = ?, 
			date_of_birth = ?
		WHERE user_id = ?;
		`, fullName, phoneNumber, gender, dateOfBirth, userID)
	if err != nil {
		return fmt.Errorf("Error updating profile! Please try again")
	}
	return nil
}
