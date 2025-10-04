package repository

import (
    "database/sql"
    "gonote/internal/domain"
)

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
    query := `
        INSERT INTO users (email, password, name, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    return r.db.QueryRow(query, user.Email, user.Password, user.Name).
        Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
    user := &domain.User{}
    query := `
        SELECT id, email, password, name, created_at, updated_at
        FROM users WHERE email = $1
    `
    err := r.db.QueryRow(query, email).Scan(
        &user.ID, &user.Email, &user.Password, &user.Name,
        &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
    user := &domain.User{}
    query := `
        SELECT id, email, password, name, created_at, updated_at
        FROM users WHERE id = $1
    `
    err := r.db.QueryRow(query, id).Scan(
        &user.ID, &user.Email, &user.Password, &user.Name,
        &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return user, nil
}