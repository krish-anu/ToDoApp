import React, { useState } from "react";

interface TodoFormProps {
  onAddTodo: (body: string) => Promise<void>;
}

const TodoForm: React.FC<TodoFormProps> = ({ onAddTodo }) => {
  const [newTodo, setNewTodo] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    const trimmedTodo = newTodo.trim();
    if (!trimmedTodo) {
      return;
    }

    setIsSubmitting(true);
    try {
      await onAddTodo(trimmedTodo);
      setNewTodo("");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="mb-6">
      <div className="flex flex-col gap-3 sm:flex-row">
        <input
          type="text"
          value={newTodo}
          onChange={(event) => setNewTodo(event.target.value)}
          placeholder="What needs to be done?"
          className="w-full rounded-xl border border-slate-300 px-4 py-3 shadow-sm focus:border-cyan-600 focus:outline-none focus:ring-2 focus:ring-cyan-200 disabled:opacity-60"
          disabled={isSubmitting}
        />

        <button
          type="submit"
          disabled={!newTodo.trim() || isSubmitting}
          className="rounded-xl bg-slate-900 px-6 py-3 font-semibold text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {isSubmitting ? "Adding..." : "Add Task"}
        </button>
      </div>
    </form>
  );
};

export default TodoForm;
