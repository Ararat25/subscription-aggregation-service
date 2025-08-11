package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	myError "github.com/Ararat25/subscription-aggregation-service/internal/error"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PGRepo - структура для базы данных
type PGRepo struct {
	conn *pgx.Conn // соединение с бд
}

// ConnectDB производит соединение с бд
func (repo *PGRepo) ConnectDB(ctx context.Context, dbHost string, dbUser string, dbPassword string, dbName string, dbPort int) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	repo.conn = conn

	return nil
}

// CreateSubscription добавляет подписку в бд и возвращает id
func (repo *PGRepo) CreateSubscription(ctx context.Context, s *entity.Subscription) (int64, error) {
	if s == nil {
		return 0, fmt.Errorf("invalid argument error")
	}

	var id int64
	err := repo.conn.QueryRow(ctx,
		`INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
             VALUES ($1, $2, $3, $4, $5)
             RETURNING id`,
		s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// ReadSubscription возвращает подписку по id
func (repo *PGRepo) ReadSubscription(ctx context.Context, id int64) (*entity.Subscription, error) {
	row := repo.conn.QueryRow(ctx,
		`SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`, id)

	var s entity.Subscription
	err := row.Scan(&s.Id, &s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myError.ErrSubscriptionNotFound
		}
		return nil, err
	}

	return &s, nil
}

// UpdateSubscription обновляет данные подписки
func (repo *PGRepo) UpdateSubscription(ctx context.Context, s *entity.Subscription) error {
	if s == nil {
		return fmt.Errorf("invalid argument error: subscription is nil")
	}

	if s.Id == 0 {
		return fmt.Errorf("invalid argument error: missing subscription ID")
	}

	cmdTag, err := repo.conn.Exec(ctx,
		`UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6`,
		s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate, s.Id,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return myError.ErrSubscriptionNotFound
	}

	return nil
}

// DeleteSubscription удаляет подписку
func (repo *PGRepo) DeleteSubscription(ctx context.Context, id int64) error {
	cmdTag, err := repo.conn.Exec(ctx,
		`DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return myError.ErrSubscriptionNotFound
	}

	return nil
}

// ListSubscriptions возвращает список всех подписок
func (repo *PGRepo) ListSubscriptions(ctx context.Context) ([]*entity.Subscription, error) {
	rows, err := repo.conn.Query(ctx,
		`SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*entity.Subscription
	for rows.Next() {
		var s entity.Subscription
		err = rows.Scan(&s.Id, &s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate)
		if err != nil {
			return nil, err
		}
		subs = append(subs, &s)
	}
	return subs, nil
}

// TotalCost возвращает суммарную стоимость подписок за определенный период с фильтрацией по id пользователя и/или названию сервиса
func (repo *PGRepo) TotalCost(ctx context.Context, from, to time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	var total int
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE (($1 BETWEEN start_date AND COALESCE(end_date, CURRENT_DATE))
		OR ($2 BETWEEN start_date AND COALESCE(end_date, CURRENT_DATE))
		OR ($1 = start_date OR $2 = start_date))
	`
	args := []interface{}{from, to}

	if userID != nil {
		query += ` AND user_id = $3`
		args = append(args, *userID)
	}

	if serviceName != nil {
		if userID == nil {
			query += ` AND service_name = $3`
			args = append(args, *serviceName)
		} else {
			query += ` AND service_name = $4`
			args = append(args, *serviceName)
		}
	}

	err := repo.conn.QueryRow(ctx, query, args...).Scan(&total)
	return total, err
}

// Close закрывает соединение с бд
func (repo *PGRepo) Close(ctx context.Context) error {
	err := repo.conn.Close(ctx)
	if err != nil {
		return err
	}

	return nil
}
