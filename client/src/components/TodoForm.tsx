import React, { useState } from "react";

interface TodoFormProps {
  onAddTodo: (body: string) => void;
}

const TodoForm: React.FC<TodoFormProps> = ({ onAddTodo }) => {
  const [newTodo, setNewTodo] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newTodo.trim()) return;

    setIsSubmitting(true);
    try {
      await onAddTodo(newTodo.trim());
      setNewTodo("");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="mb-8">
      <div className="flex flex-col sm:flex-row gap-4 items-stretch sm:items-end">
        {/* Input field */}
        <div className="flex-1">
          <input
            type="text"
            value={newTodo}
            onChange={(e) => setNewTodo(e.target.value)}
            placeholder="📝 What needs to be done?"
            className="w-full px-5 py-3 rounded-xl border border-gray-300 bg-white/70 backdrop-blur-md shadow-inner focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition duration-200 text-gray-800 placeholder-gray-400 disabled:opacity-60"
            disabled={isSubmitting}
          />
        </div>

        {/* Submit button */}
        <button
          type="submit"
          disabled={!newTodo.trim() || isSubmitting}
          className="px-6 py-3 rounded-xl bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold shadow-lg hover:from-indigo-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-indigo-400 disabled:opacity-60 disabled:cursor-not-allowed transition duration-200"
        >
          {isSubmitting ? (
            <div className="flex items-center justify-center gap-2">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
              Adding...
            </div>
          ) : (
            "Add Task"
          )}
        </button>
      </div>
    </form>
  );
};

export default TodoForm;
