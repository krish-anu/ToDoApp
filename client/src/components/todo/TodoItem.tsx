import React, { useState } from "react";
import type { Todo } from "@/types/todo";

interface TodoItemProps {
  todo: Todo;
  onToggle: (id: string) => Promise<void>;
  onDelete: (id: string) => Promise<void>;
}

const TodoItem: React.FC<TodoItemProps> = ({ todo, onToggle, onDelete }) => {
  const [isDeleting, setIsDeleting] = useState(false);
  const [isToggling, setIsToggling] = useState(false);

  const handleToggle = async () => {
    setIsToggling(true);
    try {
      await onToggle(todo._id);
    } finally {
      setIsToggling(false);
    }
  };

  const handleDelete = async () => {
    setIsDeleting(true);
    try {
      await onDelete(todo._id);
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <article
      className={[
        "flex items-center gap-3 rounded-xl border p-4 shadow-sm transition",
        todo.completed
          ? "border-slate-200 bg-slate-50"
          : "border-slate-300 bg-white",
        isDeleting ? "opacity-50" : "",
      ].join(" ")}
    >
      <button
        type="button"
        onClick={handleToggle}
        disabled={isToggling || isDeleting}
        aria-pressed={todo.completed}
        aria-label={
          todo.completed
            ? `Mark ${todo.body} incomplete`
            : `Mark ${todo.body} complete`
        }
        className={[
          "h-5 w-5 rounded border-2 transition-colors",
          todo.completed
            ? "border-emerald-500 bg-emerald-500"
            : "border-slate-300 hover:border-cyan-600",
        ].join(" ")}
      >
        <span className="sr-only">Toggle status</span>
      </button>

      <p
        className={[
          "flex-1 text-left break-words",
          todo.completed ? "text-slate-500 line-through" : "text-slate-800",
        ].join(" ")}
      >
        {todo.body}
      </p>

      <button
        type="button"
        onClick={handleDelete}
        disabled={isDeleting || isToggling}
        className="rounded-lg border border-rose-300 px-3 py-1.5 text-rose-700 transition hover:bg-rose-50 disabled:opacity-50"
      >
        {isDeleting ? "Deleting..." : "Delete"}
      </button>
    </article>
  );
};

export default TodoItem;
