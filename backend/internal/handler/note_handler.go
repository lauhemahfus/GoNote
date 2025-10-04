package handler

import (
    "net/http"
    "strconv"
    "gonote/internal/service"
    "github.com/gin-gonic/gin"
)

type NoteHandler struct {
    noteService service.NoteService
}

func NewNoteHandler(noteService service.NoteService) *NoteHandler {
    return &NoteHandler{noteService: noteService}
}

type CreateNoteRequest struct {
    Title   string `json:"title" binding:"required"`
    Content string `json:"content"`
}

type UpdateNoteRequest struct {
    Title   string `json:"title" binding:"required"`
    Content string `json:"content"`
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
    userID := c.GetInt("user_id")
    
    var req CreateNoteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    note, err := h.noteService.CreateNote(userID, req.Title, req.Content)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
        return
    }
    
    c.JSON(http.StatusCreated, note)
}

func (h *NoteHandler) GetNotes(c *gin.Context) {
    userID := c.GetInt("user_id")
    
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    notes, total, err := h.noteService.GetNotes(userID, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "notes": notes,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

func (h *NoteHandler) GetNote(c *gin.Context) {
    userID := c.GetInt("user_id")
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }
    
    note, err := h.noteService.GetNote(id, userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
        return
    }
    
    c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
    userID := c.GetInt("user_id")
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }
    
    var req UpdateNoteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    note, err := h.noteService.UpdateNote(id, userID, req.Title, req.Content)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
        return
    }
    
    c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
    userID := c.GetInt("user_id")
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }
    
    if err := h.noteService.DeleteNote(id, userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func (h *NoteHandler) GetSummary(c *gin.Context) {
    userID := c.GetInt("user_id")
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }
    
    summary, err := h.noteService.GetSummary(id, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate summary"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"summary": summary})
}