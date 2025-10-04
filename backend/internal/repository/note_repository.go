package repository

import (
    "database/sql"
    "gonote/internal/domain"
)

type noteRepository struct {
    db *sql.DB
}

func NewNoteRepository(db *sql.DB) domain.NoteRepository {
    return &noteRepository{db: db}
}

func (r *noteRepository) Create(note *domain.Note) error {
    query := `
        INSERT INTO notes (user_id, title, content, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    return r.db.QueryRow(query, note.UserID, note.Title, note.Content).
        Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
}

func (r *noteRepository) GetByID(id, userID int) (*domain.Note, error) {
    note := &domain.Note{}
    query := `
        SELECT id, user_id, title, content, created_at, updated_at
        FROM notes WHERE id = $1 AND user_id = $2
    `
    err := r.db.QueryRow(query, id, userID).Scan(
        &note.ID, &note.UserID, &note.Title, &note.Content,
        &note.CreatedAt, &note.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return note, nil
}

func (r *noteRepository) GetByUserID(userID, limit, offset int) ([]*domain.Note, error) {
    query := `
        SELECT id, user_id, title, content, created_at, updated_at
        FROM notes WHERE user_id = $1
        ORDER BY updated_at DESC
        LIMIT $2 OFFSET $3
    `
    rows, err := r.db.Query(query, userID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var notes []*domain.Note
    for rows.Next() {
        note := &domain.Note{}
        err := rows.Scan(
            &note.ID, &note.UserID, &note.Title, &note.Content,
            &note.CreatedAt, &note.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        notes = append(notes, note)
    }
    return notes, nil
}

func (r *noteRepository) Update(note *domain.Note) error {
    query := `
        UPDATE notes
        SET title = $1, content = $2, updated_at = NOW()
        WHERE id = $3 AND user_id = $4
        RETURNING updated_at
    `
    return r.db.QueryRow(query, note.Title, note.Content, note.ID, note.UserID).
        Scan(&note.UpdatedAt)
}

func (r *noteRepository) Delete(id, userID int) error {
    query := `DELETE FROM notes WHERE id = $1 AND user_id = $2`
    _, err := r.db.Exec(query, id, userID)
    return err
}

func (r *noteRepository) Count(userID int) (int, error) {
    var count int
    query := `SELECT COUNT(*) FROM notes WHERE user_id = $1`
    err := r.db.QueryRow(query, userID).Scan(&count)
    return count, err
}