package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"go-api-test/pkg/domain"
	"log"
	"strings"
)

type UsersStorage struct {
	conn *sql.DB
}

func NewUsersStorage(conn *sql.DB) *UsersStorage {
	return &UsersStorage{conn: conn}
}

func (r *UsersStorage) Create(ctx context.Context, user domain.User) (int64, error) {
	var id int64
	stmt, err := r.conn.PrepareContext(ctx, "INSERT INTO users(name, last_name, surname, gender, status, date_of_birth, created_at) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	if err != nil {
		return 0, err
	}
	err = stmt.QueryRow(user.Name, user.LastName, user.Surname, user.Gender, user.Status, user.DateOfBirth, user.CreatedAt).Scan(&id)
	return id, err
}

func (r *UsersStorage) Delete(ctx context.Context, id int64) error {
	_, err := r.conn.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *UsersStorage) GetByID(ctx context.Context, id int64) (domain.User, error) {
	var user domain.User
	err := r.conn.QueryRowContext(ctx, "SELECT id, name, last_name, surname, gender, status, date_of_birth, created_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.LastName, &user.Surname, &user.Gender, &user.Status, &user.DateOfBirth, &user.CreatedAt)
	return user, err
}

func (r *UsersStorage) Update(ctx context.Context, id int64, user domain.User) error {

	var updateQueryParts []string
	var queryParams []interface{}
	paramCounter := 1

	if user.Name != "" {
		updateQueryParts = append(updateQueryParts, fmt.Sprintf("name = $%d", paramCounter))
		queryParams = append(queryParams, user.Name)
		paramCounter++
	}

	if user.LastName != "" {
		updateQueryParts = append(updateQueryParts, fmt.Sprintf("last_name = $%d", paramCounter))
		queryParams = append(queryParams, user.Surname)
		paramCounter++
	}

	if user.Surname != "" {
		updateQueryParts = append(updateQueryParts, fmt.Sprintf("surname = $%d", paramCounter))
		queryParams = append(queryParams, user.Surname)
		paramCounter++
	}

	if user.Gender != "" {
		updateQueryParts = append(updateQueryParts, fmt.Sprintf("gender = $%d", paramCounter))
		queryParams = append(queryParams, user.Gender)
		paramCounter++
	}

	if user.Status != "" {
		updateQueryParts = append(updateQueryParts, fmt.Sprintf("status = $%d", paramCounter))
		queryParams = append(queryParams, user.Status)
		paramCounter++
	}

	if !user.DateOfBirth.IsZero() {
		updateQueryParts = append(updateQueryParts, fmt.Sprintf("date_of_birth = $%d", paramCounter))
		queryParams = append(queryParams, user.DateOfBirth)
		paramCounter++
	}

	updateQuery := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(updateQueryParts, ", "), paramCounter)
	queryParams = append(queryParams, id)

	if _, err := r.conn.ExecContext(ctx, updateQuery, queryParams...); err != nil {
		return err
	}

	return nil
}

func (r *UsersStorage) GetList(ctx context.Context, params *domain.UsersParam) ([]domain.User, uint64, error) {

	var total uint64

	filters, queryParams, paramCounter := getUsersParams(params)

	if filters != "" {
		filters = " WHERE " + filters
	}

	if err := r.conn.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM users"+filters, queryParams...).Scan(&total); err != nil {
		return nil, 0, err
	}

	var users = make([]domain.User, 0, total)

	query := `
				SELECT
					id,
					name,
					last_name,
					surname,
					gender,
					status,
					date_of_birth,
					created_at
				FROM
					users
				`

	aggregateQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", params.OrderBy, params.OrderDir, paramCounter, paramCounter+1)
	queryParams = append(queryParams, params.Limit, params.Offset)

	rows, err := r.conn.QueryContext(ctx, query+filters+aggregateQuery, queryParams...)
	if err != nil {
		return users, total, err
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Println("GetList error occurred while closing connection err:", err)
		}
	}(rows)

	for rows.Next() {
		var user domain.User
		if err = rows.Scan(&user.ID, &user.Name, &user.LastName, &user.Surname, &user.Gender, &user.Status, &user.DateOfBirth, &user.CreatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return users, total, err
	}

	return users, total, nil
}

func getUsersParams(params *domain.UsersParam) (string, []interface{}, int) {
	var queryFilters []string
	var queryParams []interface{}
	paramCounter := 1

	if params.FullName != "" {
		queryFilters = append(queryFilters, fmt.Sprintf("CONCAT(name, ' ', last_name, ' ', COALESCE(surname, '')) LIKE $%d", paramCounter))
		queryParams = append(queryParams, "%"+params.FullName+"%")
		paramCounter++
	}

	if params.Gender != "" {
		queryFilters = append(queryFilters, fmt.Sprintf("gender = $%d", paramCounter))
		queryParams = append(queryParams, params.Gender)
		paramCounter++
	}

	if params.Status != "" {
		queryFilters = append(queryFilters, fmt.Sprintf("status = $%d", paramCounter))
		queryParams = append(queryParams, params.Status)
		paramCounter++
	}

	filters := strings.Join(queryFilters, " AND ")

	return filters, queryParams, paramCounter
}
