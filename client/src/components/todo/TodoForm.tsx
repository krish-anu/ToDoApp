import React, { useState } from "react";

interface TodoFormProps {
  onAddTodo: (body: string) => Promise<void>;
  onAddTodoWithWorkflow: (body: string) => Promise<void>;
}

const TodoForm: React.FC<TodoFormProps> = ({
  onAddTodo,
  onAddTodoWithWorkflow,
}) => {
  const [newTodo, setNewTodo] = useState("");
  const [submitMode, setSubmitMode] = useState<"manual" | "workflow" | null>(
    null,
  );

  const isSubmitting = submitMode !== null;

  const submitTodo = async (mode: "manual" | "workflow") => {
    const trimmedTodo = newTodo.trim();
    if (!trimmedTodo) {
      return;
    }

    setSubmitMode(mode);
    try {
      if (mode === "workflow") {
        await onAddTodoWithWorkflow(trimmedTodo);
      } else {
        await onAddTodo(trimmedTodo);
      }
      setNewTodo("");
    } finally {
      setSubmitMode(null);
    }
  };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    await submitTodo("manual");
  };

  return (
    <form onSubmit={handleSubmit} className="mb-6">
      <div className="flex flex-col gap-3">
        <input
          type="text"
          value={newTodo}
          onChange={(event) => setNewTodo(event.target.value)}
          placeholder="What needs to be done? Try: Pay rent tomorrow at 6 PM and remind me 1 hour before"
          className="w-full rounded-xl border border-slate-300 px-4 py-3 shadow-sm focus:border-cyan-600 focus:outline-none focus:ring-2 focus:ring-cyan-200 disabled:opacity-60"
          disabled={isSubmitting}
        />

        <div className="flex flex-col gap-3 sm:flex-row">
          <button
            type="submit"
            disabled={!newTodo.trim() || isSubmitting}
            className="rounded-xl bg-slate-900 px-6 py-3 font-semibold text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {submitMode === "manual" ? "Adding..." : "Add Task"}
          </button>

          <button
            type="button"
            onClick={() => {
              void submitTodo("workflow");
            }}
            disabled={!newTodo.trim() || isSubmitting}
            className="rounded-xl border border-cyan-600 bg-cyan-50 px-6 py-3 font-semibold text-cyan-900 transition hover:bg-cyan-100 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {submitMode === "workflow" ? "Understanding..." : "Smart Add"}
          </button>
        </div>

        <p className="text-sm text-slate-500">
          Smart Add uses the workflow endpoint to extract title, date, priority,
          tags, and reminder.
        </p>
      </div>
    </form>
  );
};

export default TodoForm;
