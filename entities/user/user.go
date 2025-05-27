package users

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserId             int       `json:"user_id"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	NewPassword        string    `json:"new_password"`
	FullName           string    `json:"full_name"`
	DateOfBirth        time.Time `json:"date_of_birth"`
	DateOfBirthDisplay string    `json:"date_of_birth_display"`
	PhoneNumber        string    `json:"phone_number"`
	Gender             int       `json:"gender"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedAtDisplay   string    `json:"created_at_display"`
}

func Authenticate(account User, db *sql.DB) error {
	var hashedPassword []byte
	err := db.QueryRow(`	
		SELECT password 
		FROM users
		WHERE email = ?
		`, account.Email).Scan(&hashedPassword)

	if err == sql.ErrNoRows {
		return fmt.Errorf("Invalid email or password! Please try again")
	}
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(account.Password))
	if err != nil {
		return fmt.Errorf("Invalid email or password! Please try again")
	}

	return nil
}

func Create(newUser User, db *sql.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error processing password")
	}

	_, err = db.Exec(`	
		INSERT INTO users (email, password) 
		VALUES (?, ?)
		`, newUser.Email, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func GetDetails(userID int, db *sql.DB) (*User, error) {
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
		return nil, err
	}

	var user User
	if row.Next() {
		var nullFullName, nullPhoneNumber sql.NullString
		var nullGender sql.NullInt64
		var nullDateOfBirth sql.NullTime
		var email string

		err := row.Scan(&email, &nullFullName, &nullPhoneNumber, &nullGender, &nullDateOfBirth)
		if err != nil {
			return nil, err
		}

		user.Email = email
		user.FullName = nullFullName.String
		user.PhoneNumber = nullPhoneNumber.String
		user.Gender = int(nullGender.Int64)

		if nullDateOfBirth.Valid {
			user.DateOfBirth = nullDateOfBirth.Time
			user.DateOfBirthDisplay = user.DateOfBirth.Format("2006-01-02")
		}
	}
	defer row.Close()

	return &user, nil
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
	var hashedPassword string
	err := db.QueryRow(`
		SELECT password
		FROM users 
		WHERE user_id = ?`, userID).Scan(&hashedPassword)
	if err != nil {
		return fmt.Errorf("Error while changing password! Please try again")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("Wrong password! Please try again")
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error while changing password! Please try again")
	}

	_, err = db.Exec(`
		UPDATE users 
		SET password = ?
		WHERE user_id = ?`, string(hashedNewPassword), userID)
	if err != nil {
		return fmt.Errorf("Error while changing password! Please try again")
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

func GetRole(email string, db *sql.DB) (string, error) {
	var role string
	err := db.QueryRow(`
		SELECT role_name
		FROM roles
		WHERE role_id = (
			SELECT role_id
			FROM users
			WHERE email = ?
		)
		`, email).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

func GetByPage(offset, limit int, search string, db *sql.DB) ([]User, error) {
	var query string
	var rows *sql.Rows
	var err error

	if search != "" {
		query = `
			SELECT 
				user_id,
				email,
				full_name,
				date_of_birth,
				phone_number,
				gender,
				created_at
			FROM users
			WHERE 
				role_id = 1 AND
				full_name LIKE ? OR 
				email LIKE ? OR 
				phone_number LIKE ? 
			LIMIT ? OFFSET ?;
		`
		rows, err = db.Query(query, "%"+search+"%", "%"+search+"%", "%"+search+"%", limit, offset)
	} else {
		query = `
			SELECT 
				user_id,
				email,
				full_name,
				date_of_birth,
				phone_number,
				gender,
				created_at
			FROM users
			WHERE role_id = 1
			LIMIT ? OFFSET ?;
		`
		rows, err = db.Query(query, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserId, &user.Email, &user.FullName, &user.DateOfBirth, &user.PhoneNumber, &user.Gender, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		user.CreatedAtDisplay = user.CreatedAt.Format("2006-01-02 15:04:05")
		user.DateOfBirthDisplay = user.DateOfBirth.Format("2006-01-02") 
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}