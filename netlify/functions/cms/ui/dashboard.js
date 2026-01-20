// API URL - always use /api which is redirected to Netlify function
const API_URL = `${window.location.origin}/api`;

let allPosts = [];
let currentFilter = 'all';

// Check auth on page load
window.addEventListener('load', async () => {
    const token = localStorage.getItem('auth_token');
    if (!token) {
        window.location.href = '/login';
        return;
    }

    // Load posts (this will fail if token is invalid, redirecting to login)
    loadPosts();
});

async function loadPosts() {
    try {
        const token = localStorage.getItem('auth_token');
        const response = await fetch(`${API_URL}/posts`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        // If unauthorized, redirect to login
        if (response.status === 401) {
            localStorage.removeItem('auth_token');
            window.location.href = '/login';
            return;
        }

        if (!response.ok) {
            throw new Error('Failed to load posts');
        }

        allPosts = await response.json() || [];
        renderPosts();
        updateStats();
        loadDraftPosts();
    } catch (error) {
        showAlert(error.message || 'Failed to load posts', 'error');
        renderEmptyState();
    }
}

function renderPosts() {
    const container = document.getElementById('postsContainer');

    // Filter posts
    let filtered = allPosts;
    if (currentFilter !== 'all') {
        filtered = allPosts.filter(post => post.status === currentFilter);
    }

    if (filtered.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <h3>No posts yet</h3>
                <p>Create your first post by clicking "✏️ New Post" in the sidebar</p>
            </div>
        `;
        return;
    }

    container.innerHTML = `
        <div class="posts-list">
            ${filtered.map(post => `
                <div class="post-item ${post.status || 'draft'}">
                    <div class="post-item-header">
                        <div>
                            <h3>
                                <span class="post-type-badge ${post.type_id || 'article'}">${post.type_id || 'article'}</span>
                                ${escapeHtml(post.title || 'Untitled')}
                            </h3>
                            <div class="post-meta">
                                <span class="post-status ${post.status || 'draft'}">
                                    ${post.status || 'draft'}
                                </span>
                                <span>${formatDate(post.created_at)}</span>
                                ${post.updated_at && post.updated_at !== post.created_at ? `
                                    <span>Updated: ${formatDate(post.updated_at)}</span>
                                ` : ''}
                            </div>
                        </div>
                        <div class="post-item-actions" onclick="event.stopPropagation()">
                            <button class="action-btn" onclick="viewPost('${post.id}')">View</button>
                            <button class="action-btn" onclick="editPost('${post.id}')">Edit</button>
                            <button class="action-btn delete" onclick="deletePost('${post.id}')">Delete</button>
                        </div>
                    </div>
                    ${post.excerpt || post.content ? `
                    <div class="post-item-preview">
                        ${escapeHtml((post.excerpt || post.content || '').substring(0, 200))}...
                    </div>
                    ` : ''}
                </div>
            `).join('')}
        </div>
    `;
}

function renderEmptyState() {
    document.getElementById('postsContainer').innerHTML = `
        <div class="empty-state">
            <h3>Unable to load posts</h3>
            <p>Please check your connection and try again</p>
        </div>
    `;
}

function updateStats() {
    const total = allPosts.length;
    const published = allPosts.filter(p => p.status === 'published').length;
    const drafts = allPosts.filter(p => p.status === 'draft').length;
    const archived = allPosts.filter(p => p.status === 'archived').length;

    document.getElementById('totalCount').textContent = total;
    document.getElementById('publishedCount').textContent = published;
    document.getElementById('draftStatsCount').textContent = drafts;
    document.getElementById('archivedCount').textContent = archived;
}

function loadDraftPosts() {
    const drafts = allPosts.filter(p => p.status === 'draft');
    const draftList = document.getElementById('draftList');
    const draftCount = document.getElementById('draftCount');

    draftCount.textContent = `(${drafts.length})`;

    if (drafts.length === 0) {
        draftList.innerHTML = '<div style="padding: 10px; font-size: 12px; color: #bdc3c7;">No drafts</div>';
        return;
    }

    draftList.innerHTML = drafts.map(draft => `
        <div class="draft-item" onclick="editPost('${draft.id}')" title="${escapeHtml(draft.title)}">
            ${escapeHtml(draft.title || 'Untitled')}
        </div>
    `).join('');
}

function filterPosts(status) {
    currentFilter = status;

    // Update active filter button
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');

    renderPosts();
}

async function viewPost(postId) {
    navigateTo(`/view?id=${postId}`);
}

async function editPost(postId) {
    // Store the post ID in sessionStorage for the editor to retrieve
    sessionStorage.setItem('editPostId', postId);
    navigateTo(`/editor?id=${postId}`);
}

async function deletePost(postId) {
    if (!confirm('Are you sure you want to delete this post?')) {
        return;
    }

    try {
        const token = localStorage.getItem('auth_token');
        const response = await fetch(`${API_URL}/posts/${postId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to delete post');
        }

        showAlert('Post deleted successfully', 'success');
        loadPosts();
    } catch (error) {
        showAlert(error.message || 'Failed to delete post', 'error');
    }
}

function toggleDrafts() {
    const draftList = document.getElementById('draftList');
    const icon = document.getElementById('draftIcon');
    draftList.classList.toggle('open');
    icon.classList.toggle('collapsed');
}

function navigateTo(path) {
    window.location.href = path;
}

function logout() {
    localStorage.removeItem('auth_token');
    window.location.href = '/login';
}

function showAlert(message, type) {
    const alert = document.getElementById('alert');
    alert.className = `alert ${type}`;
    alert.textContent = message;
    setTimeout(() => {
        alert.textContent = '';
        alert.className = '';
    }, 5000);
}

function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
}

function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return text.replace(/[&<>"']/g, m => map[m]);
}
