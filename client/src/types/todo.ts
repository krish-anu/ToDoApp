export interface Todo {
  _id: string;
  body: string;
  completed: boolean;
  due_at?: string;
  remind_at?: string;
  priority?: "low" | "medium" | "high";
  tags?: string[];
}

export interface CreateTodoInput {
  body: string;
  completed?: boolean;
}

export interface UpdateTodoInput {
  completed: boolean;
}

export interface CreateTodoFromTextInput {
  message: string;
  user_id?: string;
  timezone?: string;
}

export interface Reminder {
  _id: string;
  todo_id: string;
  user_id: string;
  remind_at: string;
  status: string;
  created_at: string;
}

export interface ParsedTodoFromText {
  title: string;
  due_at: string | null;
  priority: "low" | "medium" | "high";
  tags: string[];
  remind_at: string | null;
}

export interface CreateTodoFromTextResponse {
  success: boolean;
  message: string;
  partial: boolean;
  todo: Todo;
  reminder?: Reminder;
  parsed: ParsedTodoFromText;
}
