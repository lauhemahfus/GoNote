
let currentPage = 1;
let totalNotes = 0;
let currentNoteId = null;
let deleteNoteId = null;

document.addEventListener('DOMContentLoaded', () => {
    console.log('Dashboard loaded');
    
    if (!localStorage.getItem('token')) {
        console.log('No token found, redirecting to login');
        window.location.href = 'login.html';
        return;
    }
    
    console.log('Token found, initializing dashboard');
    initializeDashboard();
});

function initializeDashboard() {
    document.getElementById('logoutBtn').addEventListener('click', handleLogout);
    document.getElementById('newNoteBtn').addEventListener('click', () => openNoteModal());
    document.getElementById('noteForm').addEventListener('submit', handleSaveNote);
    document.getElementById('loadMoreBtn').addEventListener('click', loadMoreNotes);
    document.getElementById('getSummaryBtn').addEventListener('click', handleGetSummary);
    document.getElementById('confirmDeleteBtn').addEventListener('click', handleConfirmDelete);
    
    document.querySelectorAll('.modal-close, .modal-cancel').forEach(btn => {
        btn.addEventListener('click', closeModals);
    });
    
    window.addEventListener('click', (e) => {
        const noteModal = document.getElementById('noteModal');
        const deleteModal = document.getElementById('deleteModal');
        if (e.target === noteModal || e.target === deleteModal) {
            closeModals();
        }
    });
    
    loadNotes();
}

async function loadNotes(page = 1) {
    console.log(`Loading notes for page ${page}`);
    
    try {
        const data = await api.notes.getAll(page, 10);
        console.log('Notes data received:', data);
        
        totalNotes = data.total;
        currentPage = page;
        
        if (page === 1) {
            document.getElementById('notesGrid').innerHTML = '';
        }
        
        if (data.notes && data.notes.length > 0) {
            console.log(`Rendering ${data.notes.length} notes`);
            data.notes.forEach(note => renderNote(note));
            updateLoadMoreButton();
        } else if (page === 1) {
            console.log('No notes found, showing empty state');
            showEmptyState();
        }
    } catch (error) {
        console.error('Error loading notes:', error);
        alert('Failed to load notes: ' + error.message);
    }
}

function renderNote(note) {
    console.log('Rendering note:', note);
    
    const notesGrid = document.getElementById('notesGrid');
    const noteElement = document.createElement('div');
    noteElement.className = 'note-item';
    noteElement.setAttribute('data-note-id', note.id);
    
    noteElement.innerHTML = `
        <h3>${escapeHtml(note.title)}</h3>
        <p class="note-preview">${escapeHtml(note.content || 'No content')}</p>
        <div class="note-meta">
            <div class="note-time">
                <span>Created: ${formatDate(note.created_at)}</span>
                <span>Updated: ${formatDate(note.updated_at)}</span>
            </div>
            <div class="note-actions">
                <button class="edit-btn" data-note-id="${note.id}" title="Edit">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="delete-btn" data-note-id="${note.id}" title="Delete">
                    <i class="fas fa-trash"></i>
                </button>
            </div>
        </div>
    `;
    
    noteElement.addEventListener('click', (e) => {
        if (!e.target.closest('.note-actions')) {
            openNoteModal(note.id);
        }
    });
    
    const editBtn = noteElement.querySelector('.edit-btn');
    editBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        openNoteModal(parseInt(e.currentTarget.getAttribute('data-note-id')));
    });
    
    const deleteBtn = noteElement.querySelector('.delete-btn');
    deleteBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        openDeleteModal(parseInt(e.currentTarget.getAttribute('data-note-id')));
    });
    
    notesGrid.appendChild(noteElement);
}

function showEmptyState() {
    const notesGrid = document.getElementById('notesGrid');
    notesGrid.innerHTML = `
        <div class="empty-state">
            <h2>No Notes Yet</h2>
            <p>Start creating your first note to get organized!</p>
        </div>
    `;
    updateLoadMoreButton();
}

function updateLoadMoreButton() {
    const loadedNotes = document.querySelectorAll('.note-item').length;
    const loadMoreContainer = document.getElementById('loadMoreContainer');
    
    console.log(`Loaded notes: ${loadedNotes}, Total notes: ${totalNotes}`);
    
    if (loadedNotes < totalNotes && loadedNotes > 0) {
        loadMoreContainer.classList.add('show');
    } else {
        loadMoreContainer.classList.remove('show');
    }
}

async function loadMoreNotes() {
    await loadNotes(currentPage + 1);
}

async function openNoteModal(noteId = null) {
    const modal = document.getElementById('noteModal');
    const modalTitle = document.getElementById('modalTitle');
    const noteTitle = document.getElementById('noteTitle');
    const noteContent = document.getElementById('noteContent');
    const summarySection = document.getElementById('summarySection');
    const getSummaryBtn = document.getElementById('getSummaryBtn');
    
    currentNoteId = noteId;
    summarySection.style.display = 'none';
    
    if (noteId) {
        modalTitle.textContent = 'Edit Note';
        getSummaryBtn.style.display = 'block';
        try {
            const note = await api.notes.getOne(noteId);
            noteTitle.value = note.title;
            noteContent.value = note.content;
        } catch (error) {
            showError('Failed to load note');
        }
    } else {
        modalTitle.textContent = 'Create Note';
        getSummaryBtn.style.display = 'none';
        noteTitle.value = '';
        noteContent.value = '';
    }
    
    modal.classList.add('show');
}

async function handleSaveNote(e) {
    e.preventDefault();
    
    const title = document.getElementById('noteTitle').value;
    const content = document.getElementById('noteContent').value;
    const errorDiv = document.getElementById('noteError');
    
    console.log('Saving note:', { title, content, currentNoteId });
    
    try {
        if (currentNoteId) {
            await api.notes.update(currentNoteId, title, content);
        } else {
            await api.notes.create(title, content);
        }
        
        closeModals();
        currentPage = 1;
        await loadNotes(1);
    } catch (error) {
        console.error('Error saving note:', error);
        errorDiv.textContent = error.message;
        errorDiv.classList.add('show');
    }
}

async function handleGetSummary() {
    const summarySection = document.getElementById('summarySection');
    const summaryContent = document.getElementById('summaryContent');
    const getSummaryBtn = document.getElementById('getSummaryBtn');
    
    if (!currentNoteId) return;
    
    getSummaryBtn.textContent = 'Generating...';
    getSummaryBtn.disabled = true;
    
    try {
        const data = await api.notes.getSummary(currentNoteId);
        summaryContent.textContent = data.summary;
        summarySection.style.display = 'block';
    } catch (error) {
        showError('Failed to generate summary: ' + error.message);
    } finally {
        getSummaryBtn.textContent = 'Get AI Summary';
        getSummaryBtn.disabled = false;
    }
}

function openDeleteModal(noteId) {
    console.log('Opening delete modal for note:', noteId);
    deleteNoteId = noteId;
    document.getElementById('deleteModal').classList.add('show');
}

async function handleConfirmDelete() {
    if (!deleteNoteId) return;
    
    console.log('Deleting note:', deleteNoteId);
    
    try {
        await api.notes.delete(deleteNoteId);
        closeModals();
        currentPage = 1;
        await loadNotes(1);
    } catch (error) {
        console.error('Error deleting note:', error);
        showError('Failed to delete note');
    }
}

function closeModals() {
    document.getElementById('noteModal').classList.remove('show');
    document.getElementById('deleteModal').classList.remove('show');
    document.getElementById('noteError').classList.remove('show');
    currentNoteId = null;
    deleteNoteId = null;
}

function handleLogout() {
    console.log('Logging out');
    localStorage.removeItem('token');
    window.location.href = '/';
}

function formatDate(dateString) {
    const date = new Date(dateString);
    
    // Convert UTC to GMT+6 (Bangladesh Time)
    const utcTime = date.getTime();
    const gmt6Offset = 6 * 60 * 60 * 1000;
    const localDate = new Date(utcTime + gmt6Offset);
    
    const now = new Date();
    const nowGMT6 = new Date(now.getTime() + gmt6Offset);
    
    const diff = nowGMT6 - localDate;
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    const hours = Math.floor(diff / (1000 * 60 * 60));
    const minutes = Math.floor(diff / (1000 * 60));
    
    if (minutes < 1) return 'Just now';
    if (minutes < 60) return `${minutes} min${minutes > 1 ? 's' : ''} ago`;
    if (hours < 24) return `${hours} hr${hours > 1 ? 's' : ''} ago`;
    if (days === 0) return 'Today';
    if (days === 1) return 'Yesterday';
    if (days < 7) return `${days} day${days > 1 ? 's' : ''} ago`;
    if (days < 30) {
        const weeks = Math.floor(days / 7);
        return `${weeks} wk${weeks > 1 ? 's' : ''} ago`;
    }
    if (days < 365) {
        const months = Math.floor(days / 30);
        return `${months} mo${months > 1 ? 's' : ''} ago`;
    }
    
    return localDate.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showError(message) {
    const errorDiv = document.getElementById('noteError');
    errorDiv.textContent = message;
    errorDiv.classList.add('show');
    setTimeout(() => errorDiv.classList.remove('show'), 5000);
}
