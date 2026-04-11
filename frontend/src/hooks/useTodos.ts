import { useCallback, useEffect, useMemo, useState } from "react";
import axios from "axios";
import {
  createTodoFromText,
  createTodo,
  fetchTodos,
  removeTodo,
  updateTodo,
} from "@/services/todoService";
import type { Todo } from "@/types/todo";

const ERROR_MESSAGES = {
  FETCH: "Failed to load todos",
  CREATE: "Failed to add todo",
  WORKFLOW: "Failed to create task from text",
  UPDATE: "Failed to update todo",
  DELETE: "Failed to delete todo",
} as const;

export function useTodos() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadTodos = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchTodos();
      setTodos(data);
    } catch (loadError) {
      console.error(loadError);
      setError(ERROR_MESSAGES.FETCH);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadTodos();
  }, [loadTodos]);

  const addTodo = useCallback(async (body: string) => {
    try {
      setError(null);
      const createdTodo = await createTodo({ body, completed: false });
      setTodos((current) => [createdTodo, ...current]);
    } catch (createError) {
      console.error(createError);
      setError(ERROR_MESSAGES.CREATE);
      throw createError;
    }
  }, []);

  const addTodoWithWorkflow = useCallback(async (message: string) => {
    try {
      setError(null);
      const result = await createTodoFromText({
        message,
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || "UTC",
      });

      setTodos((current) => [result.todo, ...current]);

      if (result.partial) {
        setError(result.message);
      }
    } catch (workflowError) {
      console.error(workflowError);

      if (axios.isAxiosError(workflowError)) {
        const apiError = workflowError.response?.data as
          | { error?: string }
          | undefined;
        setError(apiError?.error || ERROR_MESSAGES.WORKFLOW);
      } else {
        setError(ERROR_MESSAGES.WORKFLOW);
      }

      throw workflowError;
    }
  }, []);

  const toggleTodo = useCallback(
    async (id: string) => {
      const target = todos.find((todo) => todo._id === id);
      if (!target) {
        return;
      }

      try {
        setError(null);
        const updatedTodo = await updateTodo(id, {
          completed: !target.completed,
        });
        setTodos((current) =>
          current.map((todo) => (todo._id === id ? updatedTodo : todo)),
        );
      } catch (updateError) {
        console.error(updateError);
        setError(ERROR_MESSAGES.UPDATE);
        throw updateError;
      }
    },
    [todos],
  );

  const deleteTodo = useCallback(async (id: string) => {
    try {
      setError(null);
      await removeTodo(id);
      setTodos((current) => current.filter((todo) => todo._id !== id));
    } catch (deleteError) {
      console.error(deleteError);
      setError(ERROR_MESSAGES.DELETE);
      throw deleteError;
    }
  }, []);

  const completedCount = useMemo(
    () => todos.filter((todo) => todo.completed).length,
    [todos],
  );

  return {
    todos,
    loading,
    error,
    totalCount: todos.length,
    completedCount,
    addTodo,
    addTodoWithWorkflow,
    toggleTodo,
    deleteTodo,
    refresh: loadTodos,
  };
}
