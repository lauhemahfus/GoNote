package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "gonote/internal/domain"
    "github.com/go-redis/redis/v8"
)

type CacheRepository interface {
    GetNote(id int) (*domain.Note, error)
    SetNote(note *domain.Note) error
    DeleteNote(id int) error
    InvalidateUserNotes(userID int) error
}

type cacheRepository struct {
    client *redis.Client
    ctx    context.Context
}

func NewCacheRepository(client *redis.Client) CacheRepository {
    return &cacheRepository{
        client: client,
        ctx:    context.Background(),
    }
}

func (r *cacheRepository) GetNote(id int) (*domain.Note, error) {
    key := fmt.Sprintf("note:%d", id)
    data, err := r.client.Get(r.ctx, key).Result()
    if err != nil {
        return nil, err
    }
    
    var note domain.Note
    if err := json.Unmarshal([]byte(data), &note); err != nil {
        return nil, err
    }
    return &note, nil
}

func (r *cacheRepository) SetNote(note *domain.Note) error {
    key := fmt.Sprintf("note:%d", note.ID)
    data, err := json.Marshal(note)
    if err != nil {
        return err
    }
    return r.client.Set(r.ctx, key, data, 30*time.Minute).Err()
}

func (r *cacheRepository) DeleteNote(id int) error {
    key := fmt.Sprintf("note:%d", id)
    return r.client.Del(r.ctx, key).Err()
}

func (r *cacheRepository) InvalidateUserNotes(userID int) error {
    pattern := fmt.Sprintf("notes:user:%d:*", userID)
    iter := r.client.Scan(r.ctx, 0, pattern, 0).Iterator()
    for iter.Next(r.ctx) {
        r.client.Del(r.ctx, iter.Val())
    }
    return iter.Err()
}