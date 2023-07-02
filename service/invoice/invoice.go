package invoice

import (
	"database/sql"
	"errors"
)

type Invoice struct {
	Id     int     `json:"id"`
	UserId int     `json:"user_id"`
	Status Status  `json:"status"`
	Amount float64 `json:"amount"`
	Label  string  `json:"label"`
}

type Status string

const (
	Pending Status = "pending"
	Paid    Status = "paid"
)

type Service struct {
	db *sql.DB
}

var NotFoundError = errors.New("invoice: not found")

func MustMakeService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (service Service) CreateInvoice(userId int, label string, amount float64) error {
	_, err := service.db.Exec(`INSERT INTO "invoices" ("user_id", "label", "amount") VALUES ($1, $2, $3)`, userId, label, amount)
	return err
}

func (service Service) GetInvoiceById(id int) (Invoice, error) {
	row := service.db.QueryRow("SELECT id, user_id, status, label, amount FROM invoices WHERE id = $1", id)
	var invoice Invoice
	err := row.Scan(&invoice.Id, &invoice.UserId, &invoice.Status, &invoice.Label, &invoice.Amount)
	if err == sql.ErrNoRows {
		return invoice, NotFoundError
	}
	return invoice, err
}

// GetInvoiceByIdAndAmount retrieves associated invoice. If not found, returns NotFoundError
func (service Service) GetInvoiceByIdAndAmount(id int, amount float64) (Invoice, error) {
	row := service.db.QueryRow("SELECT id, user_id, status, label, amount FROM invoices WHERE id = $1 AND amount = $2", id, amount)
	var invoice Invoice
	err := row.Scan(&invoice.Id, &invoice.UserId, &invoice.Status, &invoice.Label, &invoice.Amount)
	if err == sql.ErrNoRows {
		return invoice, NotFoundError
	}
	return invoice, err
}

func (service Service) UpdateInvoiceStatus(id int, newStatus Status) error {
	_, err := service.db.Exec(`UPDATE "invoices" SET status = $1 WHERE id = $2`, newStatus, id)
	return err
}

func (service Service) UpdateInvoiceStatusWithTx(id int, newStatus Status, tx *sql.Tx) error {
	_, err := tx.Exec(`UPDATE "invoices" SET status = $1 WHERE id = $2`, newStatus, id)
	return err
}
