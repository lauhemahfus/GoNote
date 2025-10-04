package service

import (
    "gonote/internal/domain"
    "gonote/internal/repository"
)

type NoteService interface {
    CreateNote(userID int, title, content string) (*domain.Note, error)
    GetNote(id, userID int) (*domain.Note, error)
    GetNotes(userID, page, limit int) ([]*domain.Note, int, error)
    UpdateNote(id, userID int, title, content string) (*domain.Note, error)
    DeleteNote(id, userID int) error
    GetSummary(id, userID int) (string, error)
}

type noteService struct {
    noteRepo  domain.NoteRepository
    cacheRepo repository.CacheRepository
    aiService AIService
}

func NewNoteService(noteRepo domain.NoteRepository, cacheRepo repository.CacheRepository, aiService AIService) NoteService {
    return &noteService{
        noteRepo:  noteRepo,
        cacheRepo: cacheRepo,
        aiService: aiService,
    }
}

func (s *noteService) CreateNote(userID int, title, content string) (*domain.Note, error) {
    note := &domain.Note{
        UserID:  userID,
        Title:   title,
        Content: content,
    }
    
    if err := s.noteRepo.Create(note); err != nil {
        return nil, err
    }
    
    s.cacheRepo.SetNote(note)
    return note, nil
}

func (s *noteService) GetNote(id, userID int) (*domain.Note, error) {
    note, err := s.cacheRepo.GetNote(id)
    if err == nil {
        if note.UserID == userID {
            return note, nil
        }
    }
    
    note, err = s.noteRepo.GetByID(id, userID)
    if err != nil {
        return nil, err
    }
    
    s.cacheRepo.SetNote(note)
    return note, nil
}

func (s *noteService) GetNotes(userID, page, limit int) ([]*domain.Note, int, error) {
    offset := (page - 1) * limit
    notes, err := s.noteRepo.GetByUserID(userID, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    
    total, err := s.noteRepo.Count(userID)
    if err != nil {
        return nil, 0, err
    }
    
    return notes, total, nil
}

func (s *noteService) UpdateNote(id, userID int, title, content string) (*domain.Note, error) {
    note, err := s.noteRepo.GetByID(id, userID)
    if err != nil {
        return nil, err
    }
    
    note.Title = title
    note.Content = content
    
    if err := s.noteRepo.Update(note); err != nil {
        return nil, err
    }
    
    s.cacheRepo.SetNote(note)
    return note, nil
}

func (s *noteService) DeleteNote(id, userID int) error {
    if err := s.noteRepo.Delete(id, userID); err != nil {
        return err
    }
    
    s.cacheRepo.DeleteNote(id)
    return nil
}

func (s *noteService) GetSummary(id, userID int) (string, error) {
    note, err := s.GetNote(id, userID)
    if err != nil {
        return "", err
    }
    
    return s.aiService.GenerateSummary(note.Content)
}