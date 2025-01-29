package database

import (
	"database/sql"
	"github.com/group-project/aitu-fan/aitu-fan/models"
	_ "github.com/lib/pq"
	"log"
)

func GetAllUsers() ([]models.User, error) {
	var users []models.User

	rows, err := db.Query("SELECT id, username, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func DeleteUserByID(userID string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := db.Exec(query, userID)
	return err
}
