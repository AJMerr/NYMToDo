import { useEffect, useMemo, useRef, useState } from 'react';
import { getTodos, createTodo, updateTodo, deleteTodo } from './api';
import './styles.css';

export default function App() {
  const [todos, setTodos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [adding, setAdding] = useState(false);
  const [title, setTitle] = useState('');
  const [desc, setDesc] = useState('');
  const [error, setError] = useState('');
  const [editingId, setEditingId] = useState(null);
  const [editingTitle, setEditingTitle] = useState('');
  const [editingDesc, setEditingDesc] = useState('');
  const [expandedId, setExpandedId] = useState(null);
  const inputRef = useRef(null);

  useEffect(() => {
    (async () => {
      try {
        setLoading(true);
        const data = await getTodos();
        setTodos(Array.isArray(data) ? data : []);
      } catch (e) {
        setError(e.message || 'Failed to load todos');
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  async function onAdd(e) {
    e.preventDefault();
    if (!title.trim()) return;
    try {
      setAdding(true);
      const t = await createTodo(title.trim(), desc.trim());
      setTodos((prev) => [t, ...prev]);
      setTitle('');
      setDesc('');
      inputRef.current?.focus();
    } catch (e) {
      setError(e.message || 'Failed to add todo');
    } finally {
      setAdding(false);
    }
  }

  function startEdit(t) {
    setEditingId(t.id);
    setEditingTitle(t.title ?? '');
    setEditingDesc(t.description ?? '');
  }

  async function saveEdit(t) {
    const newTitle = editingTitle.trim();
    const newDesc = editingDesc.trim();
    if (!newTitle) return;
    try {
      const updated = await updateTodo(t.id, {
        id: t.id,
        title: newTitle,
        completed: t.completed ?? false,
        description: newDesc,
      });
      setTodos((prev) => prev.map((x) => (x.id === t.id ? updated : x)));
      setEditingId(null);
      setEditingTitle('');
      setEditingDesc('');
    } catch (e) {
      setError(e.message || 'Failed to update todo');
    }
  }

  async function toggleComplete(t) {
    try {
      const updated = await updateTodo(t.id, {
        id: t.id,
        title: t.title,
        completed: !Boolean(t.completed),
        description: t.description ?? '',
      });
      setTodos((prev) => prev.map((x) => (x.id === t.id ? updated : x)));
    } catch (e) {
      setError(e.message || 'Failed to toggle todo');
    }
  }

  async function remove(t) {
    try {
      await deleteTodo(t.id);
      setTodos((prev) => prev.filter((x) => x.id !== t.id));
      if (expandedId === t.id) setExpandedId(null);
    } catch (e) {
      setError(e.message || 'Failed to delete todo');
    }
  }

  const remaining = useMemo(() => todos.filter((t) => !t.completed).length, [todos]);

  return (
    <div className="shell">
      <header className="hero">
        <div className="heroTitle">
          <svg width="26" height="26" viewBox="0 0 24 24" className="logo" aria-hidden>
            <path d="M9 11l3 3L22 4" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
            <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V7a2 2 0 0 1 2-2h11" fill="none" stroke="currentColor" strokeWidth="2" />
          </svg>
          <h1>NYMToDo</h1>
        </div>
        <span className="pill">{remaining} open</span>
      </header>

      <form className="composer" onSubmit={onAdd}>
        <div className="fields">
          <input
            ref={inputRef}
            className="field"
            placeholder="Task title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            disabled={adding}
          />
          <textarea
            className="field"
            placeholder="Description (optional)"
            value={desc}
            onChange={(e) => setDesc(e.target.value)}
            disabled={adding}
            rows={3}
          />
        </div>
        <button className="btn primary" disabled={adding || !title.trim()} type="submit">
          {adding ? 'Adding…' : 'Add task'}
        </button>
      </form>

      {error && (
        <div className="toast" role="alert" onClick={() => setError('')}>
          {error}
        </div>
      )}

      <section className="panel">
        {loading ? (
          <div className="empty">Loading…</div>
        ) : todos.length === 0 ? (
          <div className="empty">No todos yet. Add your first one!</div>
        ) : (
          <ul className="list">
            {todos.map((t) => {
              const isEditing = editingId === t.id;
              const isOpen = expandedId === t.id;
              return (
                <li key={t.id} className={`row ${isOpen ? 'open' : ''}`}>
                  <label className="check">
                    <input
                      type="checkbox"
                      checked={Boolean(t.completed)}
                      onChange={() => toggleComplete(t)}
                    />
                    <span />
                  </label>

                  <div className="rowMain">
                    {isEditing ? (
                      <>
                        <input
                          className="editTitle"
                          value={editingTitle}
                          autoFocus
                          onChange={(e) => setEditingTitle(e.target.value)}
                          onKeyDown={(e) => {
                            if (e.key === 'Enter' && !e.shiftKey) {
                              e.preventDefault();
                              saveEdit(t);
                            }
                            if (e.key === 'Escape') {
                              setEditingId(null);
                              setEditingTitle('');
                              setEditingDesc('');
                            }
                          }}
                        />
                        <textarea
                          className="editDesc"
                          value={editingDesc}
                          onChange={(e) => setEditingDesc(e.target.value)}
                          rows={3}
                          placeholder="Edit description…"
                        />
                      </>
                    ) : (
                      <>
                        <button
                          className="titleBtn"
                          onClick={() => setExpandedId(isOpen ? null : t.id)}
                          title="Expand"
                        >
                          <span className={`caret ${isOpen ? 'rot' : ''}`}>▸</span>
                          <span className={`title ${t.completed ? 'done' : ''}`}>{t.title}</span>
                        </button>

                        <div className={`details ${isOpen ? 'show' : ''}`}>
                          {t.description?.trim() ? (
                            <p className="desc">{t.description}</p>
                          ) : (
                            <p className="desc muted">No description</p>
                          )}
                        </div>
                      </>
                    )}
                  </div>

                  <div className="rowActions">
                    {isEditing ? (
                      <>
                        <button className="btn subtle" onClick={() => saveEdit(t)} title="Save">
                          Save
                        </button>
                        <button
                          className="btn subtle"
                          onClick={() => {
                            setEditingId(null);
                            setEditingTitle('');
                            setEditingDesc('');
                          }}
                          title="Cancel"
                        >
                          Cancel
                        </button>
                      </>
                    ) : (
                      <>
                        <button className="icon" title="Edit" onClick={() => startEdit(t)}>✏️</button>
                        <button className="icon danger" title="Delete" onClick={() => remove(t)}>🗑️</button>
                      </>
                    )}
                  </div>
                </li>
              );
            })}
          </ul>
        )}
      </section>

      <footer className="foot">
        <span>Powered by Parsec · /api → backend</span>
      </footer>
    </div>
  );
}
