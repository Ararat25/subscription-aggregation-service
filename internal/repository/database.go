package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PGRepo - структура для базы данных
type PGRepo struct {
	Conn *pgx.Conn
}

// ConnectDB производит соединение с бд
func (repo *PGRepo) ConnectDB(dbHost string, dbUser string, dbPassword string, dbName string, dbPort int) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	repo.Conn = conn

	return nil
}

func (repo *PGRepo) CreateSubscription(s *entity.Subscription) (int64, error) {
	if s == nil {
		return 0, fmt.Errorf("invalid argument error")
	}

	var id int64
	err := repo.Conn.QueryRow(context.Background(),
		`INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
             VALUES ($1, $2, $3, $4, $5)
             RETURNING id`,
		s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *PGRepo) ReadSubscription(id int64) (*entity.Subscription, error) {
	row := repo.Conn.QueryRow(context.Background(),
		`SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`, id)

	var s entity.Subscription
	err := row.Scan(&s.Id, &s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("subscription with id %d not found", id)
		}
		return nil, err
	}

	return &s, nil
}

func (repo *PGRepo) UpdateSubscription(s *entity.Subscription) error {
	if s == nil {
		return fmt.Errorf("invalid argument error: subscription is nil")
	}

	if s.Id == 0 {
		return fmt.Errorf("invalid argument error: missing subscription ID")
	}

	cmdTag, err := repo.Conn.Exec(context.Background(),
		`UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6`,
		s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate, s.Id,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no subscription found with id %d", s.Id)
	}

	return nil
}

func (repo *PGRepo) DeleteSubscription(id int64) error {
	cmdTag, err := repo.Conn.Exec(context.Background(),
		`DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("subscription with id %d not found", id)
	}

	return nil
}

func (repo *PGRepo) ListSubscriptions() ([]*entity.Subscription, error) {
	rows, err := repo.Conn.Query(context.Background(),
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

func (repo *PGRepo) TotalCost(from, to time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	var total int
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE (($1 BETWEEN start_date AND end_date)
		OR ($2 BETWEEN start_date AND end_date)
		OR ($1 = start_date OR $2 = start_date))
	`
	args := []interface{}{from, to}

	if userID != nil {
		query += ` AND user_id = $3`
		args = append(args, *userID)
	}

	if serviceName != nil {
		query += ` AND service_name = $4`
		args = append(args, *serviceName)
	}

	fmt.Println(args)

	err := repo.Conn.QueryRow(context.Background(), query, args...).Scan(&total)
	return total, err
}
