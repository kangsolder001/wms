package repository

import (
	"context"
	"database/sql"
	"fmt"

	"wms/domain/entity"
	"wms/pkg/logger"
)

type postgresUserRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewPostgresUserRepository(db *sql.DB, log logger.Logger) *postgresUserRepository {
	return &postgresUserRepository{db: db, log: log}
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, password_hash, full_name, role, is_active, created_at, updated_at 
		 FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (r *postgresUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, password_hash, full_name, role, is_active, created_at, updated_at 
		 FROM users WHERE username = $1`, username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, password_hash, full_name, role, is_active, created_at, updated_at 
		 FROM users WHERE email = $1`, email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (r *postgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (username, email, password_hash, full_name, role, is_active, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		user.Username, user.Email, user.PasswordHash, user.FullName, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET username=$1, email=$2, full_name=$3, role=$4, is_active=$5, updated_at=$6 WHERE id=$7`,
		user.Username, user.Email, user.FullName, user.Role, user.IsActive, user.UpdatedAt, user.ID,
	)
	return err
}

func (r *postgresUserRepository) List(ctx context.Context, page, limit int) ([]*entity.User, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, username, email, full_name, role, is_active, created_at, updated_at 
		 FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		u := &entity.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			continue
		}
		users = append(users, u)
	}

	return users, total, nil
}
