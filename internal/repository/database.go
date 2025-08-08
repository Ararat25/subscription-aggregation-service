package repository

import (
	"context"
	"fmt"
	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
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

func (repo *PGRepo) CreateSubscription(s *entity.Subscription) error {
	_, err := repo.Conn.Exec(context.Background(),
		`INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		 VALUES ($1, $2, $3, $4, $5)`,
		s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate,
	)
	return err
}

func (repo *PGRepo) ReadSubscription(id int64) (*entity.Subscription, error) {
	row := repo.Conn.QueryRow(context.Background(),
		`SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`, id)

	var s entity.Subscription
	err := row.Scan(&s.Id, &s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (repo *PGRepo) UpdateSubscription(s *entity.Subscription) error {
	_, err := repo.Conn.Exec(context.Background(),
		`UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6`,
		s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate, s.Id,
	)
	return err
}

func (repo *PGRepo) DeleteSubscription(id int64) error {
	_, err := repo.Conn.Exec(context.Background(),
		`DELETE FROM subscriptions WHERE id = $1`, id)
	return err
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
		var s *entity.Subscription
		err = rows.Scan(&s.Id, &s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
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
