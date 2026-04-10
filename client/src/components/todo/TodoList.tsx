import React from "react";
import type { Todo } from "@/types/todo";
import TodoItem from "@/components/todo/TodoItem";

interface TodoListProps {
  todos: Todo[];
  onToggleTodo: (id: string) => Promise<void>;
  onDeleteTodo: (id: string) => Promise<void>;
}

const TodoList: React.FC<TodoListProps> = ({
  todos,
  onToggleTodo,
  onDeleteTodo,
}) => {
  if (todos.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-slate-300 bg-slate-50 p-8 text-center">
        <h3 className="text-lg font-semibold text-slate-700">No tasks yet</h3>
        <p className="mt-1 text-slate-500">Add your first task to get started.</p>
      </div>
    );
  }

  const activeTodos = todos.filter((todo) => !todo.completed);
  const completedTodos = todos.filter((todo) => todo.completed);

  return (
    <div className="space-y-6">
      {activeTodos.length > 0 && (
        <section>
          <h3 className="mb-3 text-lg font-semibold text-slate-700">
            Active Tasks ({activeTodos.length})
          </h3>
          <div className="space-y-2">
            {activeTodos.map((todo) => (
              <TodoItem
                key={todo._id}
                todo={todo}
                onToggle={onToggleTodo}
                onDelete={onDeleteTodo}
              />
            ))}
          </div>
        </section>
      )}

      {completedTodos.length > 0 && (
        <section>
          <h3 className="mb-3 text-lg font-semibold text-slate-700">
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
        </section>
      )}
    </div>
  );
};

export default TodoList;
