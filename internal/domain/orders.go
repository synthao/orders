package domain

import (
	"errors"
	"time"
)

var ErrCreateOrder = errors.New("failed to create a order")
var ErrUpdateOrder = errors.New("failed to update a order")
var ErrGetOrder = errors.New("failed to get a order")

type Service interface {
	Create(item *Order) (OrderID, error)
	GetOne(id int) (*Order, error)
	GetList(limit, offset int) ([]Order, error)
	Delete(id int) error
	UpdateStatus(id, status int) error
}

type Repository interface {
	Create(item *Order) (OrderID, error)
	GetOne(id int) (*Order, error)
	GetList(limit, offset int) ([]Order, error)
	Delete(id int) error
	Update(order *Order) error
}

type OrderID int

const (
	StatusNew = iota + 1
	StatusPacking
	StatusDelivery
	StatusCanceled
	StatusDone
)

type Order struct {
	ID        int
	Status    int
	Sum       float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOrder(sum float64) *Order {
	return &Order{
		Status:    StatusNew,
		Sum:       sum,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
