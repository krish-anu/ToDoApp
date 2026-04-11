import api from "@/api";
import type {
  CreateTodoFromTextInput,
  CreateTodoFromTextResponse,
  CreateTodoInput,
  Todo,
  UpdateTodoInput,
} from "@/types/todo";

export async function fetchTodos(): Promise<Todo[]> {
  const response = await api.get<Todo[]>("/api/todos");
  return Array.isArray(response.data) ? response.data : [];
}

export async function createTodo(payload: CreateTodoInput): Promise<Todo> {
  const response = await api.post<Todo>("/api/todos", {
    body: payload.body,
    completed: payload.completed ?? false,
  });
  return response.data;
}

export async function updateTodo(
  id: string,
  payload: UpdateTodoInput,
): Promise<Todo> {
  const response = await api.patch<Todo>(`/api/todos/${id}`, payload);
  return response.data;
}

export async function removeTodo(id: string): Promise<void> {
  await api.delete(`/api/todos/${id}`);
}

export async function createTodoFromText(
  payload: CreateTodoFromTextInput,
): Promise<CreateTodoFromTextResponse> {
  const response = await api.post<CreateTodoFromTextResponse>(
    "/api/workflows/task-from-text",
    payload,
  );
  return response.data;
}
