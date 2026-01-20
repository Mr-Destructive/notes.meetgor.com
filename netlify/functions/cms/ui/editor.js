// API URL is set dynamically in editor.html via window.API_URL
let currentPostId = null;
let tags = [];
let isEditMode = false;
let autoSaveInterval = null;
let isDirty = false;
let isSaving = false;

// Simple markdown to HTML converter
function markdownToHtml(markdown) {
  if (!markdown) return '';
  
  let html = markdown
    // Headers
    .replace(/^### (.*?)$/gm, '<h3>$1</h3>')
    .replace(/^## (.*?)$/gm, '<h2>$1</h2>')
    .replace(/^# (.*?)$/gm, '<h1>$1</h1>')
    // Bold
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/__（.*?)__/g, '<strong>$1</strong>')
    // Italic
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    .replace(/_(.*?)_/g, '<em>$1</em>')
    // Links
    .replace(/\[(.*?)\]\((.*?)\)/g, '<a href="$2">$1</a>')
    // Code blocks
    .replace(/```(.*?)```/gs, '<pre><code>$1</code></pre>')
    // Inline code
    .replace(/`(.*?)`/g, '<code>$1</code>')
    // Line breaks
    .replace(/\n\n/g, '</p><p>')
    .replace(/\n/g, '<br>');
  
  // Wrap in paragraphs
  if (html && !html.startsWith('<h') && !html.startsWith('<pre')) {
    html = '<p>' + html + '</p>';
  }
  
  return html;
}

// Update preview as user types
function updatePreview() {
  const content = document.getElementById('content').value;
  const preview = document.getElementById('preview');
  preview.innerHTML = markdownToHtml(content);
}

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
    fields: ['source_url'],
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

// Fetch metadata from URL
async function fetchUrlMetadata(url) {
  try {
    const response = await fetch(`${window.API_URL}/metadata?url=${encodeURIComponent(url)}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
      }
    });
    if (response.ok) {
      return await response.json();
    }
  } catch (e) {
    console.warn('Could not fetch metadata:', e);
  }
  return null;
}

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
          <label for="${field}">Source URL ${(document.getElementById('postType').value === 'link' || document.getElementById('postType').value === 'thought') ? '*' : ''}</label>
          <input type="url" id="${field}" placeholder="https://example.com" onchange="updateFromUrl()" onblur="updateFromUrl()">
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
  
  markDirty();
}

function formatFieldName(field) {
  return field
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
}

async function updateFromUrl() {
  const type = document.getElementById('postType').value;
  const sourceUrl = document.getElementById('source_url')?.value;
  
  if (!sourceUrl) return;
  
  try {
    const url = new URL(sourceUrl);
    const title = document.getElementById('title');
    const excerpt = document.getElementById('excerpt');
    
    // Try to fetch metadata from URL
    showStatus('Fetching metadata...', 'info');
    const metadata = await fetchUrlMetadata(sourceUrl);
    
    if (metadata) {
      // Auto-fill title if empty
      if (!title.value && metadata.title) {
        title.value = metadata.title;
      }
      
      // Auto-fill excerpt if empty
      if (!excerpt.value && metadata.description) {
        excerpt.value = metadata.description;
      }
      
      // Store author info in metadata
      if (metadata.author && !document.getElementById('author')?.value) {
        const authorField = document.getElementById('author');
        if (authorField) {
          authorField.value = metadata.author;
        }
      }
      
      showStatus('Metadata loaded', 'success');
    } else {
      // Fallback: extract title from hostname
      const hostname = url.hostname.replace('www.', '');
      if (!title.value) {
        title.value = hostname.charAt(0).toUpperCase() + hostname.slice(1);
      }
      showStatus('Using domain name as title', 'info');
    }
    
    // Auto-generate slug from URL if not set
    if (!document.getElementById('slug').value) {
      document.getElementById('slug').value = generateSlugFromUrl(url);
    }
  } catch (e) {
    console.warn('Error parsing URL:', e);
  }
}

function generateSlugFromUrl(url) {
  try {
    const pathname = url.pathname.split('/').filter(p => p).pop() || url.hostname.replace('www.', '');
    return pathname
      .toLowerCase()
      .replace(/[^\w\s-]/g, '')
      .replace(/\s+/g, '-')
      .replace(/-+/g, '-')
      .trim();
  } catch (e) {
    return '';
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
  markDirty();
}

function removeTag(tag) {
  tags = tags.filter(t => t !== tag);
  renderTags();
  markDirty();
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
          const items = input.value
            .split('\n')
            .map(i => i.trim())
            .filter(i => i);
          if (items.length > 0) {
            metadata[field] = items;
          }
        } else if (field === 'tools') {
          const tools = input.value
            .split(',')
            .map(t => t.trim())
            .filter(t => t);
          if (tools.length > 0) {
            metadata[field] = tools;
          }
        } else if (field === 'is_private') {
          metadata[field] = input.checked;
        } else if (input.value) {
          metadata[field] = input.value;
        }
      }
    });
  }
  
  return metadata;
}

function generateUniqueId() {
  return `draft-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

async function saveDraft() {
  const title = document.getElementById('title').value.trim();
  const content = document.getElementById('content').value.trim();
  const type = document.getElementById('postType').value;
  
  if (!content && !title) {
    showAlert('Add at least some content or a title', 'error');
    return;
  }
  
  await performSave(false);
}

async function publishPost() {
  const type = document.getElementById('postType').value;
  const title = document.getElementById('title').value.trim();
  const content = document.getElementById('content').value.trim();
  const template = POST_TEMPLATES[type];
  
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
  
  await performSave(true);
}

async function performSave(isPublish = false) {
  if (isSaving) return;
  
  isSaving = true;
  updateStatus('saving');
  
  const type = document.getElementById('postType').value;
  const title = document.getElementById('title').value.trim();
  const slug = document.getElementById('slug').value.trim();
  
  // Auto-generate slug if missing
  let finalSlug = slug;
  if (!finalSlug) {
    if (title) {
      finalSlug = generateSlug(title);
    } else {
      // Use unique ID for untitled drafts
      finalSlug = generateUniqueId();
    }
  }
  
  // Update title with unique ID if it's a new draft without title
  let finalTitle = title;
  if (!finalTitle && !isPublish && !isEditMode) {
    finalTitle = `Draft ${generateUniqueId().slice(6, 13)}`;
  }
  
  const postData = {
    type_id: type,
    title: finalTitle || null,
    slug: finalSlug,
    excerpt: document.getElementById('excerpt').value.trim() || null,
    content: document.getElementById('content').value.trim() || null,
    status: isPublish ? 'published' : 'draft',
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
      throw new Error(error.error || 'Save failed');
    }
    
    const data = await response.json();
    currentPostId = data.id;
    isEditMode = true;
    
    // Update slug and title in form
    document.getElementById('slug').value = data.slug;
    document.getElementById('title').value = data.title;
    
    isDirty = false;
    updateStatus('saved');
    
    if (isPublish) {
      showAlert('Post published successfully', 'success');
    } else {
      showAlert('Draft saved', 'success');
    }
  } catch (error) {
    updateStatus('error');
    showAlert(error.message || 'Save failed', 'error');
  } finally {
    isSaving = false;
  }
}

function markDirty() {
  isDirty = true;
  updateStatus('unsaved');
}

function updateStatus(status) {
  const dot = document.getElementById('statusDot');
  const text = document.getElementById('statusText');
  
  if (!dot || !text) return;
  
  dot.className = 'status-dot';
  
  switch (status) {
    case 'saved':
      dot.classList.add('saved');
      text.textContent = 'Saved';
      break;
    case 'saving':
      dot.classList.add('saving');
      text.textContent = 'Saving...';
      break;
    case 'unsaved':
      text.textContent = 'Unsaved changes';
      break;
    case 'error':
      dot.style.background = '#e74c3c';
      text.textContent = 'Error saving';
      break;
    default:
      dot.classList.add('saved');
      text.textContent = 'Ready';
  }
}

function showAlert(message, type) {
  const alert = document.getElementById('alert');
  
  const prefixes = {
    success: '[OK]',
    error: '[ERROR]',
    info: '[INFO]'
  };
  
  alert.className = `alert ${type}`;
  alert.innerHTML = `<span class="alert-icon">${prefixes[type] || ''}</span> <span>${escapeHtml(message)}</span>`;
  
  setTimeout(() => {
    alert.textContent = '';
    alert.className = '';
  }, 4000);
}

function showStatus(message, type) {
  showAlert(message, type || 'info');
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
    
    isDirty = false;
    updateStatus('ready');
    showAlert('Form cleared', 'success');
    updatePreview();
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

function formatDate(dateString) {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
}

function logout() {
  localStorage.removeItem('auth_token');
  window.location.href = '/login';
}

// Mark form as dirty on any input change
document.addEventListener('input', (e) => {
  if (!isSaving) {
    markDirty();
  }
  // Update preview when content changes
  if (e.target.id === 'content') {
    updatePreview();
  }
});

document.addEventListener('change', () => {
  if (!isSaving) {
    markDirty();
  }
});

// Auto-save drafts every 5 minutes
function startAutoSave() {
  autoSaveInterval = setInterval(() => {
    if (isDirty && !isSaving) {
      saveDraft();
    }
  }, 5 * 60 * 1000); // 5 minutes
}

// Stop auto-save on unload
window.addEventListener('beforeunload', () => {
  if (autoSaveInterval) clearInterval(autoSaveInterval);
});

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
  
  // Start auto-save
  startAutoSave();
  updateStatus('ready');
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
        
        if (typeof post.metadata === 'string') {
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
                const items = Array.isArray(metadata[field]) ? metadata[field] : [metadata[field]];
                input.value = items.join('\n');
              } else if (field === 'tools') {
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

    isDirty = false;
    updateStatus('saved');
    
    // Show post info
    if (post.created_at) {
      document.getElementById('createdDate').textContent = formatDate(post.created_at);
    }
    if (post.updated_at) {
      document.getElementById('updatedDate').textContent = formatDate(post.updated_at);
    }
    
    // Update preview
    updatePreview();
    
    showAlert('Editing: ' + (post.title || 'Untitled'), 'success');
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
