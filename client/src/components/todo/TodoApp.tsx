import React from "react";
import Header from "@/components/layout/Header";
import TodoForm from "@/components/todo/TodoForm";
import TodoList from "@/components/todo/TodoList";
import { useTodos } from "@/hooks/useTodos";

const TodoApp: React.FC = () => {
  const {
    todos,
    loading,
    error,
    totalCount,
    completedCount,
    addTodo,
    toggleTodo,
    deleteTodo,
  } = useTodos();

  return (
    <main className="min-h-screen bg-gradient-to-br from-slate-100 via-cyan-50 to-blue-100 px-4 py-10">
      <section className="mx-auto w-full max-w-4xl space-y-6">
        <Header totalCount={totalCount} completedCount={completedCount} />

        <div className="rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-lg backdrop-blur-sm sm:p-8">
          <TodoForm onAddTodo={addTodo} />

          {error && (
            <div className="mb-4 rounded-lg border border-rose-300 bg-rose-50 px-4 py-3 text-rose-800">
              {error}
            </div>
          )}

          {loading ? (
            <div className="py-10 text-center text-slate-500">
              Loading tasks...
            </div>
          ) : (
            <TodoList
              todos={todos}
              onToggleTodo={toggleTodo}
              onDeleteTodo={deleteTodo}
            />
          )}
        </div>
      </section>
    </main>
  );
};

export default TodoApp;
