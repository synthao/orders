package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/synthao/orders/internal/domain"
	"go.uber.org/zap"
)

var ErrUpdateStatus = errors.New("failed to update orders status")

type service struct {
	logger     *zap.Logger
	repository domain.Repository
	kw         *kafka.Writer
}

func NewService(logger *zap.Logger, repository domain.Repository, kw *kafka.Writer) domain.Service {
	return &service{logger: logger, repository: repository, kw: kw}
}

func (s *service) Create(item *domain.Order) (domain.OrderID, error) {
	id, err := s.repository.Create(item)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("payload", item))
	}

	return id, nil
}

func (s *service) GetOne(id int) (*domain.Order, error) {
	return s.repository.GetOne(id)
}

func (s *service) GetList(limit, offset int) ([]domain.Order, error) {
	return s.repository.GetList(limit, offset)
}

func (s *service) UpdateStatus(id, status int) error {
	one, err := s.repository.GetOne(id)
	if err != nil {
		s.logger.Error(err.Error())
		return fmt.Errorf("%w, GetOne, %w", ErrUpdateStatus, err)
	}

	// TODO check if status is valid

	prevStatus := one.Status

	one.Status = status

	err = s.repository.Update(one)
	if err != nil {
		s.logger.Error(err.Error())
		return fmt.Errorf("%w, %w", ErrUpdateStatus, err)
	}

	var data struct {
		OrderID        int `json:"order_id"`
		PreviousStatus int `json:"previous_status"`
		Status         int `json:"status"`
	}

	data.OrderID = one.ID
	data.PreviousStatus = prevStatus
	data.Status = status

	jsonData, err := json.Marshal(data)
	if err != nil {
		s.logger.Error(err.Error())
		return fmt.Errorf("%w, json marshal, %w", ErrUpdateStatus, err)
	}

	msg := kafka.Message{
		Value: jsonData,
	}

	err = s.kw.WriteMessages(context.Background(), msg)
	if err != nil {
		s.logger.Error(err.Error())
		return fmt.Errorf("%w, write message, %w", ErrUpdateStatus, err)
	}

	return nil
}

func (s *service) Delete(id int) error {
	if _, err := s.repository.GetOne(id); err != nil {
		return err
	}

	if err := s.repository.Delete(id); err != nil {
		s.logger.Error(err.Error(), zap.Int("id", id))
		return err
	}

	return nil
}
