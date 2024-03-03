package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/synthao/orders/internal/domain"
	"time"
)

type oneDTO struct {
	ID        int       `db:"id"`
	Status    int       `db:"status"`
	Sum       float64   `db:"sum"`
	CreatedAt time.Time `db:"created_at"`
}

type listDTO struct {
	ID        int       `db:"id"`
	Status    int       `db:"status"`
	Sum       float64   `db:"sum"`
	CreatedAt time.Time `db:"created_at"`
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) domain.Repository {
	return &repository{db: db}
}

func (r *repository) Create(item *domain.Order) (domain.OrderID, error) {
	args := map[string]interface{}{
		"sum":    item.Sum,
		"status": item.Status,
	}

	exec, err := r.db.NamedExec("INSERT INTO orders(sum, status) VALUES (:name, :text)", args)
	if err != nil {
		return 0, fmt.Errorf("%w, named exec, %w", domain.ErrCreateOrder, err)
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w, get last insert id, %w", domain.ErrCreateOrder, err)
	}

	return domain.OrderID(id), nil
}

func (r *repository) GetOne(id int) (*domain.Order, error) {
	var dest oneDTO

	err := r.db.Get(&dest, "SELECT id, status, sum, created_at FROM orders WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return &domain.Order{
		ID:        dest.ID,
		Status:    dest.Status,
		Sum:       dest.Sum,
		CreatedAt: dest.CreatedAt,
	}, nil
}

func (r *repository) GetList(limit, offset int) ([]domain.Order, error) {
	var dest []listDTO

	err := r.db.Select(&dest, "SELECT id, status, sum, created_at FROM orders")
	if err != nil {
		return nil, err
	}

	return fromListDTOToDomain(dest), nil
}

func (r *repository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM orders WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}

func fromListDTOToDomain(dto []listDTO) []domain.Order {
	res := make([]domain.Order, len(dto))

	for i, item := range dto {
		res[i] = domain.Order{
			ID:        item.ID,
			Status:    item.Status,
			Sum:       item.Sum,
			CreatedAt: item.CreatedAt,
		}
	}

	return res
}
