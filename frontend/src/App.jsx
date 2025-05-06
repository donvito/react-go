import { useState, useEffect } from 'react'
import './App.css'

const API_URL = import.meta.env.DEV ? 'http://localhost:8080/api/todos' : '/api/todos'

function App() {
  const [todos, setTodos] = useState([])
  const [newTodoText, setNewTodoText] = useState('')
  const [editingTodo, setEditingTodo] = useState(null) // { id, text, completed }
  const [editText, setEditText] = useState('')

  useEffect(() => {
    fetchTodos()
  }, [])

  const fetchTodos = async () => {
    try {
      const response = await fetch(API_URL)
      if (!response.ok) throw new Error('Failed to fetch todos')
      const data = await response.json()
      setTodos(data || []) // Ensure todos is an array
    } catch (error) {
      console.error(error)
      setTodos([]) // Set to empty array on error
    }
  }

  const handleAddTodo = async (e) => {
    e.preventDefault()
    if (!newTodoText.trim()) return
    try {
      const response = await fetch(API_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ text: newTodoText, completed: false }),
      })
      if (!response.ok) throw new Error('Failed to add todo')
      // const addedTodo = await response.json()
      // setTodos([...todos, addedTodo]) // Backend returns the created todo, can add directly
      setNewTodoText('')
      fetchTodos() // Refetch all todos to get the new ID
    } catch (error) {
      console.error(error)
    }
  }

  const toggleComplete = async (todo) => {
    try {
      const response = await fetch(`${API_URL}/${todo.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ...todo, completed: !todo.completed }),
      })
      if (!response.ok) throw new Error('Failed to update todo')
      // const updatedTodo = await response.json()
      // setTodos(todos.map(t => t.id === updatedTodo.id ? updatedTodo : t))
      fetchTodos() // Refetch for simplicity
    } catch (error) {
      console.error(error)
    }
  }

  const handleDeleteTodo = async (id) => {
    try {
      const response = await fetch(`${API_URL}/${id}`, {
        method: 'DELETE',
      })
      if (!response.ok) throw new Error('Failed to delete todo')
      // setTodos(todos.filter(todo => todo.id !== id))
      fetchTodos() // Refetch for simplicity
    } catch (error) {
      console.error(error)
    }
  }

  const handleStartEdit = (todo) => {
    setEditingTodo(todo)
    setEditText(todo.text)
  }

  const handleCancelEdit = () => {
    setEditingTodo(null)
    setEditText('')
  }

  const handleSaveEdit = async () => {
    if (!editingTodo || !editText.trim()) return
    try {
      const response = await fetch(`${API_URL}/${editingTodo.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ...editingTodo, text: editText }),
      })
      if (!response.ok) throw new Error('Failed to update todo')
      // const updatedTodo = await response.json()
      // setTodos(todos.map(t => t.id === updatedTodo.id ? updatedTodo : t))
      setEditingTodo(null)
      setEditText('')
      fetchTodos() // Refetch for simplicity
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <div className="app">
      <h1>Todo List</h1>
      <form onSubmit={handleAddTodo} className="todo-form">
        <input
          type="text"
          value={newTodoText}
          onChange={(e) => setNewTodoText(e.target.value)}
          placeholder="Add a new todo"
        />
        <button type="submit">Add</button>
      </form>

      {editingTodo && (
        <div className="edit-form">
          <h3>Edit Todo</h3>
          <input 
            type="text"
            value={editText}
            onChange={(e) => setEditText(e.target.value)}
          />
          <button onClick={handleSaveEdit}>Save</button>
          <button onClick={handleCancelEdit}>Cancel</button>
        </div>
      )}

      <ul className="todo-list">
        {Array.isArray(todos) && todos.map(todo => (
          <li key={todo.id} className={`todo-item ${todo.completed ? 'completed' : ''}`}>
            <span onClick={() => toggleComplete(todo)} style={{ cursor: 'pointer' }}>
              {todo.text}
            </span>
            <div>
              <button onClick={() => handleStartEdit(todo)} className="edit-btn">Edit</button>
              <button onClick={() => handleDeleteTodo(todo.id)} className="delete-btn">Delete</button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  )
}

export default App
