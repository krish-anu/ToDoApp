import React, { useState } from "react";
import type { Todo } from "@/types/todo";

interface TodoItemProps {
  todo: Todo;
  onToggle: (id: string) => Promise<void>;
  onDelete: (id: string) => Promise<void>;
}

const PRIORITY_STYLES: Record<"low" | "medium" | "high", string> = {
  low: "border-emerald-200 bg-emerald-50 text-emerald-700",
  medium: "border-amber-200 bg-amber-50 text-amber-700",
  high: "border-rose-200 bg-rose-50 text-rose-700",
};

function formatDateTimeLabel(value?: string): string | null {
  if (!value) {
    return null;
  }

  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) {
    return null;
  }

  return parsed.toLocaleString([], {
    dateStyle: "medium",
    timeStyle: "short",
  });
}

function toPriorityLabel(priority: "low" | "medium" | "high"): string {
  return priority.charAt(0).toUpperCase() + priority.slice(1);
}

const TodoItem: React.FC<TodoItemProps> = ({ todo, onToggle, onDelete }) => {
  const [isDeleting, setIsDeleting] = useState(false);
  const [isToggling, setIsToggling] = useState(false);
  const dueAtLabel = formatDateTimeLabel(todo.due_at);
  const remindAtLabel = formatDateTimeLabel(todo.remind_at);
  const hasTags = Array.isArray(todo.tags) && todo.tags.length > 0;
  const priority = todo.priority;

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

      <div className="flex-1 text-left">
        <p
          className={[
            "break-words",
            todo.completed
              ? "text-slate-500 line-through"
              : "text-slate-800",
          ].join(" ")}
        >
          {todo.body}
        </p>

        {(dueAtLabel || remindAtLabel || priority || hasTags) && (
          <div className="mt-2 flex flex-wrap items-center gap-2 text-xs">
            {dueAtLabel && (
              <span className="rounded-full border border-sky-200 bg-sky-50 px-2 py-1 text-sky-700">
                Due {dueAtLabel}
              </span>
            )}

            {remindAtLabel && (
              <span className="rounded-full border border-violet-200 bg-violet-50 px-2 py-1 text-violet-700">
                Reminder {remindAtLabel}
              </span>
            )}

            {priority && (
              <span
                className={[
                  "rounded-full border px-2 py-1",
                  PRIORITY_STYLES[priority],
                ].join(" ")}
              >
                Priority {toPriorityLabel(priority)}
              </span>
            )}

            {hasTags &&
              todo.tags!.map((tag) => (
                <span
                  key={`${todo._id}-${tag}`}
                  className="rounded-full border border-slate-200 bg-slate-100 px-2 py-1 text-slate-700"
                >
                  #{tag}
                </span>
              ))}
          </div>
        )}
      </div>

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
