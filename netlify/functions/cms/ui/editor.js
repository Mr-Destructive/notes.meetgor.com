// API URL is set dynamically in editor.html via window.API_URL
let currentPostId = null;
let tags = [];
let isEditMode = false;

// Post type templates with required fields
const POST_TEMPLATES = {
  article: { 
    fields: ['category'],
    titleRequired: true,
    contentRequired: true,
    description: 'Full-length article'
  },
  review: { 
    fields: ['rating', 'author', 'subject_link'],
    titleRequired: true,
    contentRequired: true,
    description: 'Review of a book, movie, or product'
  },
  thought: {
    fields: [],
    titleRequired: false,
    contentRequired: true,
    description: 'Quick thoughts and reflections'
  },
  link: { 
    fields: ['source_url'],
    titleRequired: false,
    contentRequired: false,
    description: 'Curated link with commentary'
  },
  til: { 
    fields: ['category', 'difficulty'],
    titleRequired: false,
    contentRequired: true,
    description: 'Today I Learned'
  },
  quote: { 
    fields: ['author', 'source'],
    titleRequired: false,
    contentRequired: true,
    description: 'Quote or excerpt'
  },
  list: { 
    fields: ['items', 'list_type'],
    titleRequired: true,
    contentRequired: false,
    description: 'Curated list'
  },
  note: { 
    fields: ['is_private'],
    titleRequired: false,
    contentRequired: true,
    description: 'Quick note'
  },
  snippet: { 
    fields: ['language'],
    titleRequired: false,
    contentRequired: true,
    description: 'Code snippet'
  },
  essay: { 
    fields: ['reading_time'],
    titleRequired: true,
    contentRequired: true,
    description: 'Long-form essay'
  },
  tutorial: { 
    fields: ['difficulty', 'estimated_time', 'tools'],
    titleRequired: true,
    contentRequired: true,
    description: 'Step-by-step guide'
  },
  interview: { 
    fields: ['interviewee_name', 'role', 'company'],
    titleRequired: true,
    contentRequired: true,
    description: 'Q&A interview'
  },
  experiment: { 
    fields: ['tools'],
    titleRequired: true,
    contentRequired: true,
    description: 'Experiment or project'
  }
};

function changePostType() {
  const type = document.getElementById('postType').value;
  const template = POST_TEMPLATES[type] || {};
  const container = document.getElementById('typeSpecificFields');
  
  // Update title requirement
  const titleInput = document.getElementById('title');
  if (template.titleRequired === false) {
    titleInput.removeAttribute('required');
    document.querySelector('label[for="title"]').textContent = 'Title (optional)';
  } else {
    titleInput.setAttribute('required', '');
    document.querySelector('label[for="title"]').textContent = 'Title *';
  }

  // Update content requirement
  const contentInput = document.getElementById('content');
  if (template.contentRequired === false) {
    contentInput.removeAttribute('required');
    document.querySelector('label[for="content"]').textContent = 'Content (optional)';
  } else {
    contentInput.setAttribute('required', '');
    document.querySelector('label[for="content"]').textContent = 'Content *';
  }
  
  container.innerHTML = '';
  
  if (template.fields && template.fields.length > 0) {
    const section = document.createElement('div');
    section.className = 'sidebar-section';
    section.innerHTML = '<h3>Type-Specific Fields</h3>';
    
    template.fields.forEach(field => {
      const group = document.createElement('div');
      group.className = 'form-group';
      
      if (field === 'items') {
        group.innerHTML = `
          <label for="${field}">Items (one per line)</label>
          <textarea id="${field}" placeholder="- Item 1&#10;- Item 2&#10;- Item 3" rows="5"></textarea>
        `;
      } else if (field === 'difficulty' || field === 'list_type') {
        group.innerHTML = `
          <label for="${field}">${formatFieldName(field)}</label>
          <select id="${field}">
            ${field === 'difficulty' 
              ? '<option value="">Select difficulty</option><option value="beginner">Beginner</option><option value="intermediate">Intermediate</option><option value="advanced">Advanced</option>'
              : '<option value="">Select type</option><option value="ordered">Ordered</option><option value="unordered">Unordered</option>'
            }
          </select>
        `;
      } else if (field === 'is_private') {
        group.innerHTML = `
          <label>
            <input type="checkbox" id="${field}">
            Private
          </label>
        `;
      } else if (field === 'rating') {
        group.innerHTML = `
          <label for="${field}">Rating</label>
          <select id="${field}">
            <option value="">Select rating</option>
            <option value="1">⭐ 1</option>
            <option value="2">⭐⭐ 2</option>
            <option value="3">⭐⭐⭐ 3</option>
            <option value="4">⭐⭐⭐⭐ 4</option>
            <option value="5">⭐⭐⭐⭐⭐ 5</option>
          </select>
        `;
      } else if (field === 'source_url') {
        group.innerHTML = `
          <label for="${field}">Source URL *</label>
          <input type="url" id="${field}" placeholder="https://example.com" required onchange="updateTitleFromUrl()">
        `;
      } else {
        group.innerHTML = `
          <label for="${field}">${formatFieldName(field)}</label>
          <input type="text" id="${field}" placeholder="${formatFieldName(field)}...">
        `;
      }
      
      section.appendChild(group);
    });
    
    container.appendChild(section);
  }
}

function formatFieldName(field) {
  return field
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
}

function updateTitleFromUrl() {
  const type = document.getElementById('postType').value;
  if (type === 'link') {
    const sourceUrl = document.getElementById('source_url')?.value;
    if (sourceUrl && !document.getElementById('title').value) {
      try {
        const url = new URL(sourceUrl);
        // Try to extract a meaningful title from the URL
        const hostname = url.hostname.replace('www.', '');
        const pathname = url.pathname.split('/').filter(p => p).pop() || hostname;
        document.getElementById('title').value = pathname
          .replace(/[-_]/g, ' ')
          .replace(/\.[^.]+$/, '')
          .replace(/\b\w/g, l => l.toUpperCase());
      } catch (e) {
        // Invalid URL, skip auto-titling
      }
    }
  }
}

function addTag() {
  const input = document.getElementById('tagInput');
  const tag = input.value.trim().toLowerCase();
  
  if (!tag) return;
  if (tags.includes(tag)) {
    showAlert('Tag already added', 'error');
    return;
  }
  
  tags.push(tag);
  input.value = '';
  renderTags();
}

function removeTag(tag) {
  tags = tags.filter(t => t !== tag);
  renderTags();
}

function renderTags() {
  const container = document.getElementById('tags');
  container.innerHTML = tags.map(tag => `
    <div class="tag">
      ${escapeHtml(tag)}
      <button type="button" onclick="removeTag('${tag}')">×</button>
    </div>
  `).join('');
}

function getMetadata() {
  const type = document.getElementById('postType').value;
  const template = POST_TEMPLATES[type];
  const metadata = {};
  
  if (template?.fields) {
    template.fields.forEach(field => {
      const input = document.getElementById(field);
      if (input) {
        if (field === 'items') {
          // Split items by newline and filter empty
          const items = input.value
            .split('\n')
            .map(i => i.trim())
            .filter(i => i);
          if (items.length > 0) {
            metadata.items = items;
          }
        } else if (field === 'is_private') {
          metadata.is_private = input.checked;
        } else if (field === 'tools') {
          // Handle comma-separated tools
          const tools = input.value
            .split(',')
            .map(t => t.trim())
            .filter(t => t);
          if (tools.length > 0) {
            metadata.tools = tools;
          }
        } else if (input.value) {
          // Only add non-empty fields
          metadata[field] = input.value;
        }
      }
    });
  }
  
  return metadata;
}

async function savePost() {
  const type = document.getElementById('postType').value;
  const template = POST_TEMPLATES[type];
  const title = document.getElementById('title').value.trim();
  const content = document.getElementById('content').value.trim();
  
  // Validate required fields
  if (template.titleRequired && !title) {
    showAlert('Title is required for this post type', 'error');
    return;
  }
  
  if (template.contentRequired && !content) {
    showAlert('Content is required for this post type', 'error');
    return;
  }

  // For link posts, require source_url
  if (type === 'link') {
    const sourceUrl = document.getElementById('source_url')?.value;
    if (!sourceUrl) {
      showAlert('Source URL is required for link posts', 'error');
      return;
    }
  }
  
  const postData = {
    type_id: type,
    title: title || null,
    slug: document.getElementById('slug').value || generateSlug(title || 'untitled'),
    excerpt: document.getElementById('excerpt').value.trim() || null,
    content: content || null,
    status: 'draft',
    tags: tags,
    metadata: getMetadata(),
    is_featured: false
  };
  
  try {
    const token = localStorage.getItem('auth_token');
    const url = currentPostId 
      ? `${window.API_URL}/posts/${currentPostId}`
      : `${window.API_URL}/posts`;
    
    const response = await fetch(url, {
      method: currentPostId ? 'PUT' : 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(postData)
    });
    
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Save failed');
    }
    
    const data = await response.json();
    currentPostId = data.id;
    isEditMode = true;
    showAlert('✓ Post saved as draft', 'success');
  } catch (error) {
    showAlert(error.message || 'Save failed', 'error');
  }
}

async function publishPost() {
  const type = document.getElementById('postType').value;
  const template = POST_TEMPLATES[type];
  const title = document.getElementById('title').value.trim();
  const content = document.getElementById('content').value.trim();
  
  // Validate required fields
  if (template.titleRequired && !title) {
    showAlert('Title is required to publish', 'error');
    return;
  }
  
  if (template.contentRequired && !content) {
    showAlert('Content is required to publish', 'error');
    return;
  }

  // For link posts, require source_url
  if (type === 'link') {
    const sourceUrl = document.getElementById('source_url')?.value;
    if (!sourceUrl) {
      showAlert('Source URL is required for link posts', 'error');
      return;
    }
  }
  
  const postData = {
    type_id: type,
    title: title || null,
    slug: document.getElementById('slug').value || generateSlug(title || 'untitled'),
    excerpt: document.getElementById('excerpt').value.trim() || null,
    content: content || null,
    status: 'published',
    tags: tags,
    metadata: getMetadata(),
    is_featured: document.getElementById('isFeatured').checked
  };
  
  try {
    const token = localStorage.getItem('auth_token');
    const url = currentPostId 
      ? `${window.API_URL}/posts/${currentPostId}`
      : `${window.API_URL}/posts`;
    
    const response = await fetch(url, {
      method: currentPostId ? 'PUT' : 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(postData)
    });
    
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Publish failed');
    }
    
    const data = await response.json();
    currentPostId = data.id;
    isEditMode = true;
    showAlert('🚀 Post published!', 'success');
  } catch (error) {
    showAlert(error.message || 'Publish failed', 'error');
  }
}

function previewPost() {
  alert('Preview functionality coming soon!');
}

function resetForm() {
  if (confirm('Reset form? This will clear all fields.')) {
    document.getElementById('title').value = '';
    document.getElementById('slug').value = '';
    document.getElementById('excerpt').value = '';
    document.getElementById('content').value = '';
    document.getElementById('status').value = 'draft';
    document.getElementById('isFeatured').checked = false;
    tags = [];
    renderTags();
    currentPostId = null;
    isEditMode = false;
    
    // Clear type-specific fields
    document.getElementById('typeSpecificFields').innerHTML = '';
    
    showAlert('Form cleared', 'success');
  }
}

function generateSlug(title) {
  if (!title) return '';
  return title
    .toLowerCase()
    .replace(/[^\w\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .trim();
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

function escapeHtml(text) {
  const map = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  };
  return String(text).replace(/[&<>"']/g, m => map[m]);
}

function logout() {
  localStorage.removeItem('auth_token');
  window.location.href = '/login';
}

// Check auth on page load
window.addEventListener('load', async () => {
  const token = localStorage.getItem('auth_token');
  if (!token) {
    window.location.href = '/login';
    return;
  }

  // Check if editing existing post
  const urlParams = new URLSearchParams(window.location.search);
  const postId = urlParams.get('id') || sessionStorage.getItem('editPostId');
  
  if (postId) {
    await loadPostForEditing(postId);
    sessionStorage.removeItem('editPostId');
  } else {
    changePostType(); // Initialize type-specific fields for new post
  }
});

async function loadPostForEditing(postId) {
  try {
    const token = localStorage.getItem('auth_token');
    const response = await fetch(`${window.API_URL}/posts/${postId}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    if (!response.ok) {
      showAlert('Post not found', 'error');
      return;
    }

    const post = await response.json();
    currentPostId = post.id;
    isEditMode = true;

    // Set post type first so type-specific fields load
    document.getElementById('postType').value = post.type_id || 'article';
    changePostType(); // Initialize type-specific fields

    // Populate form fields
    document.getElementById('title').value = post.title || '';
    document.getElementById('slug').value = post.slug || '';
    document.getElementById('excerpt').value = post.excerpt || '';
    document.getElementById('content').value = post.content || '';
    document.getElementById('status').value = post.status || 'draft';
    document.getElementById('isFeatured').checked = post.is_featured || false;

    // Load tags
    if (post.tags) {
      try {
        tags = JSON.parse(post.tags);
        if (!Array.isArray(tags)) tags = [];
        renderTags();
      } catch (e) {
        tags = [];
      }
    }

    // Load metadata
    if (post.metadata) {
      try {
        let metadata = {};
        
        // Handle both string and object metadata
        if (typeof post.metadata === 'string') {
          // Try to parse if it's a string
          if (post.metadata && post.metadata !== '{}' && !post.metadata.includes('[object Object]')) {
            try {
              metadata = JSON.parse(post.metadata);
            } catch (parseError) {
              console.warn('Could not parse metadata JSON:', post.metadata);
              metadata = {};
            }
          }
        } else if (typeof post.metadata === 'object') {
          metadata = post.metadata;
        }
        
        const type = document.getElementById('postType').value;
        const template = POST_TEMPLATES[type] || {};
        
        if (template.fields && Object.keys(metadata).length > 0) {
          template.fields.forEach(field => {
            const input = document.getElementById(field);
            if (input && metadata[field] !== undefined && metadata[field] !== null) {
              if (field === 'items') {
                // items should be an array
                const items = Array.isArray(metadata[field]) ? metadata[field] : [metadata[field]];
                input.value = items.join('\n');
              } else if (field === 'tools') {
                // tools should be an array or comma-separated
                const tools = Array.isArray(metadata[field]) ? metadata[field] : [metadata[field]];
                input.value = tools.join(', ');
              } else if (field === 'is_private') {
                input.checked = Boolean(metadata[field]);
              } else {
                input.value = metadata[field];
              }
            }
          });
        }
      } catch (e) {
        console.warn('Error loading metadata (non-fatal):', e);
      }
    }

    showAlert(`✓ Editing: ${post.title || 'Untitled'}`, 'success');
  } catch (error) {
    showAlert(error.message || 'Failed to load post', 'error');
  }
}

// Allow Enter to add tag
document.addEventListener('keypress', (e) => {
  if (e.key === 'Enter' && e.target.id === 'tagInput') {
    e.preventDefault();
    addTag();
  }
});

// Validate form before submit
document.addEventListener('submit', (e) => {
  if (e.target.tagName === 'FORM') {
    e.preventDefault();
  }
});
