import React from "react";
import TodoItem from "./TodoItem.tsx";
import type { Todo } from "./TodoApp.tsx";

interface TodoListProps {
  todos: Todo[];
  onToggleTodo: (id: string) => void;
  onDeleteTodo: (id: string) => void;
}

const TodoList: React.FC<TodoListProps> = ({
  todos,
  onToggleTodo,
  onDeleteTodo,
}) => {
  if (todos.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="text-6xl mb-4">📝</div>
        <h3 className="text-xl font-medium text-gray-700 mb-2">No tasks yet</h3>
        <p className="text-gray-500">Add your first task to get started!</p>
      </div>
    );
  }

  const incompleteTodos = todos.filter((todo) => !todo.completed);
  const completedTodos = todos.filter((todo) => todo.completed);

  return (
    <div className="space-y-6">
      {/* Active Tasks */}
      {incompleteTodos.length > 0 && (
        <div>
          <h3 className="text-lg font-semibold text-gray-700 mb-3">
            Active Tasks ({incompleteTodos.length})
          </h3>
          <div className="space-y-2">
            {incompleteTodos.map((todo) => (
              <TodoItem
                key={todo._id}
                todo={todo}
                onToggle={onToggleTodo}
                onDelete={onDeleteTodo}
              />
            ))}
          </div>
        </div>
      )}

      {/* Completed Tasks */}
      {completedTodos.length > 0 && (
        <div>
          <h3 className="text-lg font-semibold text-gray-700 mb-3">
            Completed ({completedTodos.length})
          </h3>
          <div className="space-y-2">
            {completedTodos.map((todo) => (
              <TodoItem
                key={todo._id}
                todo={todo}
                onToggle={onToggleTodo}
                onDelete={onDeleteTodo}
              />
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default TodoList;
