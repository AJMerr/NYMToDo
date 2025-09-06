// src/api.js
const API = `${window.location.origin.replace(/\/$/, '')}/api`;

async function request(path, opts = {}) {
  const res = await fetch(`${API}${path}`, {
    headers: { 'Content-Type': 'application/json', ...(opts.headers || {}) },
    ...opts,
  });
  const text = await res.text().catch(() => '');
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}${text ? ` - ${text}` : ''}`);
  if (!text) return null;
  const ct = res.headers.get('content-type') || '';
  return ct.includes('application/json') ? JSON.parse(text) : text;
}

// map server <-> client shapes
const fromServer = (t) => t && ({
  id: t.id,
  title: t.title,
  completed: Boolean(t.is_completed),
  description: t.description ?? '',
});
const toServer = (t) => ({
  title: t.title,
  is_completed: Boolean(t.completed),
  description: t.description ?? '',
});

export const getTodos = async () => {
  const data = await request('/todos');
  return Array.isArray(data) ? data.map(fromServer) : [];
};

export const createTodo = async (title, description = '') => {
  const created = await request('/todos', {
    method: 'POST',
    body: JSON.stringify({ title, description }),
  });
  return fromServer(created);
};

export const updateTodo = async (id, patch) => {
  const updated = await request(`/todos/${encodeURIComponent(id)}`, {
    method: 'PATCH',
    body: JSON.stringify(toServer(patch)),
  });
  return fromServer(updated);
};

export const deleteTodo = (id) =>
  request(`/todos/${encodeURIComponent(id)}`, { method: 'DELETE' });
