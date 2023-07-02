package user

import (
	"database/sql"
)

type User struct {
	UserId    int     `json:"user_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Balance   float64 `json:"balance"`
}

type Service struct {
	db *sql.DB
}

func MustMakeService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (service Service) GetUsers() ([]User, error) {
	rows, err := service.db.Query("SELECT id, first_name, last_name, balance FROM users")
	if err != nil {
		return nil, err
	}
	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Balance); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (service Service) CreditUserBalance(id int, amount float64) error {
	_, err := service.db.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, id)
	return err
}

func (service Service) CreditUserBalanceWithTx(id int, amount float64, tx *sql.Tx) error {
	_, err := tx.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, id)
	return err
}
