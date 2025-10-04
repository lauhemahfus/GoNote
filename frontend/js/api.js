const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = {
    async request(endpoint, options = {}) {
        const token = localStorage.getItem('token');
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers,
        };
        
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        
        try {
            const response = await fetch(`${API_BASE_URL}${endpoint}`, {
                ...options,
                headers,
            });
            
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || 'An error occurred');
            }
            
            return data;
        } catch (error) {
            throw error;
        }
    },
    
    auth: {
        signup(name, email, password) {
            return api.request('/auth/signup', {
                method: 'POST',
                body: JSON.stringify({ name, email, password }),
            });
        },
        
        login(email, password) {
            return api.request('/auth/login', {
                method: 'POST',
                body: JSON.stringify({ email, password }),
            });
        },
    },
    
    notes: {
        create(title, content) {
            return api.request('/notes', {
                method: 'POST',
                body: JSON.stringify({ title, content }),
            });
        },
        
        getAll(page = 1, limit = 10) {
            return api.request(`/notes?page=${page}&limit=${limit}`);
        },
        
        getOne(id) {
            return api.request(`/notes/${id}`);
        },
        
        update(id, title, content) {
            return api.request(`/notes/${id}`, {
                method: 'PUT',
                body: JSON.stringify({ title, content }),
            });
        },
        
        delete(id) {
            return api.request(`/notes/${id}`, {
                method: 'DELETE',
            });
        },
        
        getSummary(id) {
            return api.request(`/notes/${id}/summary`, {
                method: 'POST',
            });
        },
    },
};