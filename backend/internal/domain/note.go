package domain

import "time"

type Note struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type NoteRepository interface {
    Create(note *Note) error
    GetByID(id, userID int) (*Note, error)
    GetByUserID(userID, limit, offset int) ([]*Note, error)
    Update(note *Note) error
    Delete(id, userID int) error
    Count(userID int) (int, error)
}