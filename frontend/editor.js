const API_URL = 'https://notes-meetgor-com.vercel.app';
let currentPostId = null;
let tags = [];

// Initialize TinyMCE editor
tinymce.init({
  selector: 'textarea#content',
  height: 400,
  plugins: 'anchor autolink charmap codesample emoticons image link lists media searchreplace table visualblocks wordcount',
  toolbar: 'undo redo | formatselect | bold italic underline strikethrough | link image media | alignleft aligncenter alignright | outdent indent | lists | charmap emoticons',
  menubar: 'file edit view insert format tools table help',
  skin: 'oxide',
  content_css: 'default'
});

// Post type templates
const POST_TEMPLATES = {
  article: { fields: ['category'] },
  review: { fields: ['type', 'rating', 'author', 'subject_link'] },
  thought: {},
  link: { fields: ['source_url'] },
  til: { fields: ['category', 'difficulty'] },
  quote: { fields: ['author', 'source', 'context'] },
  list: { fields: ['items', 'list_type'] },
  note: { fields: ['is_private'] },
  snippet: { fields: ['language'] },
  essay: { fields: ['reading_time'] },
  tutorial: { fields: ['difficulty', 'estimated_time', 'tools'] },
  interview: { fields: ['interviewee_name', 'role', 'company'] },
  experiment: { fields: ['tools'] }
};

function changePostType() {
  const type = document.getElementById('postType').value;
  const template = POST_TEMPLATES[type] || {};
  const container = document.getElementById('typeSpecificFields');
  
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
          <label for="${field}">${field.charAt(0).toUpperCase() + field.slice(1)}</label>
          <textarea id="${field}" placeholder="One item per line..." rows="5"></textarea>
        `;
      } else if (field === 'difficulty' || field === 'list_type') {
        group.innerHTML = `
          <label for="${field}">${field.charAt(0).toUpperCase() + field.slice(1)}</label>
          <select id="${field}">
            ${field === 'difficulty' 
              ? '<option>beginner</option><option>intermediate</option><option>advanced</option>'
              : '<option>ordered</option><option>unordered</option>'
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
          <label for="${field}">${field.charAt(0).toUpperCase() + field.slice(1)}</label>
          <select id="${field}">
            <option value="1">⭐ 1</option>
            <option value="2">⭐⭐ 2</option>
            <option value="3" selected>⭐⭐⭐ 3</option>
            <option value="4">⭐⭐⭐⭐ 4</option>
            <option value="5">⭐⭐⭐⭐⭐ 5</option>
          </select>
        `;
      } else {
        group.innerHTML = `
          <label for="${field}">${field.charAt(0).toUpperCase() + field.slice(1)}</label>
          <input type="text" id="${field}" placeholder="${field}...">
        `;
      }
      
      section.appendChild(group);
    });
    
    container.appendChild(section);
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
      ${tag}
      <button onclick="removeTag('${tag}')">×</button>
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
          metadata.items = input.value.split('\n').filter(i => i.trim());
        } else if (field === 'is_private') {
          metadata.is_private = input.checked;
        } else {
          metadata[field] = input.value;
        }
      }
    });
  }
  
  return metadata;
}

async function savePost() {
  const title = document.getElementById('title').value;
  const content = tinymce.get('content')?.getContent() || document.getElementById('content').value;
  
  if (!title || !content) {
    showAlert('Title and content are required', 'error');
    return;
  }
  
  const postData = {
    type_id: document.getElementById('postType').value,
    title,
    slug: document.getElementById('slug').value || generateSlug(title),
    excerpt: document.getElementById('excerpt').value,
    content,
    status: 'draft',
    tags,
    metadata: getMetadata()
  };
  
  try {
    const token = localStorage.getItem('auth_token');
    const url = currentPostId 
      ? `${API_URL}/posts/${currentPostId}`
      : `${API_URL}/posts`;
    
    const response = await fetch(url, {
      method: currentPostId ? 'PUT' : 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(postData)
    });
    
    if (!response.ok) throw new Error('Save failed');
    
    const data = await response.json();
    currentPostId = data.id;
    showAlert('Post saved as draft', 'success');
  } catch (error) {
    showAlert(error.message || 'Save failed', 'error');
  }
}

async function publishPost() {
  const title = document.getElementById('title').value;
  const content = tinymce.get('content')?.getContent() || document.getElementById('content').value;
  
  if (!title || !content) {
    showAlert('Title and content are required', 'error');
    return;
  }
  
  const postData = {
    type_id: document.getElementById('postType').value,
    title,
    slug: document.getElementById('slug').value || generateSlug(title),
    excerpt: document.getElementById('excerpt').value,
    content,
    status: 'published',
    tags,
    metadata: getMetadata(),
    is_featured: document.getElementById('isFeatured').checked
  };
  
  try {
    const token = localStorage.getItem('auth_token');
    const url = currentPostId 
      ? `${API_URL}/posts/${currentPostId}`
      : `${API_URL}/posts`;
    
    const response = await fetch(url, {
      method: currentPostId ? 'PUT' : 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(postData)
    });
    
    if (!response.ok) throw new Error('Publish failed');
    
    const data = await response.json();
    currentPostId = data.id;
    showAlert('Post published! 🚀', 'success');
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
    tinymce.get('content')?.setContent('');
    tags = [];
    renderTags();
    currentPostId = null;
  }
}

function generateSlug(title) {
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

function logout() {
  localStorage.removeItem('auth_token');
  window.location.href = 'login.html';
}

// Check auth on page load
window.addEventListener('load', () => {
  const token = localStorage.getItem('auth_token');
  if (!token) {
    window.location.href = 'login.html';
  }
  changePostType(); // Initialize type-specific fields
});

// Allow Enter to add tag
document.addEventListener('keypress', (e) => {
  if (e.key === 'Enter' && e.target.id === 'tagInput') {
    e.preventDefault();
    addTag();
  }
});
