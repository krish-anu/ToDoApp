import React, { useState } from "react";
import type {Todo}  from "./TodoApp";

interface TodoItemProps {
  todo: Todo;
  onToggle: (id: string) => void;
  onDelete: (id: string) => void;
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
    <div
      className={`
        flex items-center gap-3 p-4 bg-white border rounded-lg shadow-sm hover:shadow-md transition-all duration-200
        ${todo.completed ? "bg-gray-50 border-gray-200" : "border-gray-300"}
        ${isDeleting ? "opacity-50 scale-95" : ""}
      `}
    >
      {/* Checkbox */}
      <button
        type="button"
        onClick={handleToggle}
        disabled={isToggling || isDeleting}
        aria-pressed={todo.completed}
        aria-label={todo.completed ? `Mark "${todo.body}" incomplete` : `Mark "${todo.body}" complete`}
        className={`
          flex items-center justify-center w-5 h-5 rounded border-2 transition-colors
          ${
            todo.completed
              ? "bg-green-500 border-green-500 text-white"
              : "border-gray-300 hover:border-indigo-500"
          }
          ${isToggling ? "opacity-50" : ""}
        `}
      >
        {todo.completed && (
          <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
            <path
              fillRule="evenodd"
              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
              clipRule="evenodd"
            />
          </svg>
        )}
      </button>

      {/* Todo Text */}
      <span
        className={`
          flex-1 text-left transition-all duration-200
          ${todo.completed ? "text-gray-500 line-through" : "text-gray-800"}
        `}
      >
        {todo.body}
      </span>

      {/* Actions */}
      <div className="flex gap-2">
        {!todo.completed && (
          <button
            onClick={handleToggle}
            disabled={isToggling || isDeleting}
            className="p-2 text-green-600 hover:bg-green-50 rounded-md transition-colors disabled:opacity-50"
            title="Mark as complete"
          >
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M5 13l4 4L19 7"
              />
            </svg>
          </button>
        )}

        {todo.completed && (
          <button
            onClick={handleToggle}
            disabled={isToggling || isDeleting}
            className="p-2 text-indigo-600 hover:bg-indigo-50 rounded-md transition-colors disabled:opacity-50"
            title="Mark as incomplete"
          >
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"
              />
            </svg>
          </button>
        )}

        <button
          onClick={handleDelete}
          disabled={isDeleting || isToggling}
          className="p-2 text-red-600 hover:bg-red-50 rounded-md transition-colors disabled:opacity-50"
          title="Delete task"
        >
          {isDeleting ? (
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-red-600"></div>
          ) : (
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
              />
            </svg>
          )}
        </button>
      </div>
    </div>
  );
};

export default TodoItem;
