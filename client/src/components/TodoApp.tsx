import React, { useEffect, useState } from "react";
import axios from "axios";
import TodoForm from "./TodoForm.tsx";
import TodoList from "./TodoList.tsx";
import Header from "./Header.tsx";

export type Todo = {
  _id: string;
  body: string;
  completed: boolean;
};

const TodoApp: React.FC = () => {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchTodos();
  }, []);

  const fetchTodos = async () => {
    try {
      setLoading(true);
      setError(null);
      const res = await axios.get<Todo[]>("https://todoapp-3.onrender.com/api/todos");
      // Guard against null or unexpected responses from the API
      setTodos(Array.isArray(res.data) ? res.data : []);
    } catch (err) {
      setError(" Failed to fetch todos");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const addTodo = async (body: string) => {
    try {
      setError(null);
      const res = await axios.post<Todo>("http://localhost:5000/api/todos", {
        body,
        completed: false,
      });
      // Use functional update and guard previous value
      setTodos((prev) => [...(prev ?? []), res.data]);
    } catch (err) {
      setError(" Failed to add todo");
      console.error(err);
    }
  };

  const toggleTodo = async (id: string) => {
    try {
      setError(null);
      await axios.patch(`http://localhost:5000/api/todos/${id}`);
      // Use functional update and guard against null
      setTodos((prev) =>
        (prev ?? []).map((todo) =>
          todo._id === id ? { ...todo, completed: !todo.completed } : todo
        )
      );
    } catch (err) {
      setError(" Failed to update todo");
      console.error(err);
    }
  };

  const deleteTodo = async (id: string) => {
    try {
      setError(null);
      await axios.delete(`http://localhost:5000/api/todos/${id}`);
      // Use functional update and guard against null
      setTodos((prev) => (prev ?? []).filter((todo) => todo._id !== id));
    } catch (err) {
      setError(" Failed to delete todo");
      console.error(err);
    }
  };

  // Defensive usage in case `todos` becomes null due to unexpected API responses
  const safeTodos = todos ?? [];
  const completedCount = safeTodos.filter((todo) => todo.completed).length;
  const totalCount = safeTodos.length;

  return (
    <div className="min-h-screen w-full flex items-center justify-center bg-gradient-to-br from-slate-900 via-indigo-900 to-purple-900 px-4 py-12">
      <div className="w-full max-w-4xl sm:max-w-5xl rounded-3xl bg-white/10 backdrop-blur-xl p-6 sm:p-10 border border-white/20 shadow-2xl space-y-10">
        <div className="bg-white/50 backdrop-blur-xl border border-white/30 shadow-2xl rounded-3xl overflow-hidden transition-all duration-300 ">
          {/* Header */}
          <Header totalCount={totalCount} completedCount={completedCount} />

          {/* Main Content */}
          <div className="p-6 sm:p-8">
            <TodoForm onAddTodo={addTodo} />

            {/* Error Message */}
            {error && (
              <div className="mt-6 p-4 bg-red-100/70 border border-red-400 text-red-800 rounded-lg shadow-sm transition-all duration-300">
                {error}
              </div>
            )}

            {/* Loading Spinner */}
            {loading ? (
              <div className="flex justify-center items-center py-12">
                <div className="animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-indigo-600"></div>
              </div>
            ) : (
              <TodoList
                todos={todos}
                onToggleTodo={toggleTodo}
                onDeleteTodo={deleteTodo}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );

};

export default TodoApp;
