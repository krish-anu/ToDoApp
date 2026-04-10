export interface Todo {
  _id: string;
  body: string;
  completed: boolean;
}

export interface CreateTodoInput {
  body: string;
  completed?: boolean;
}

export interface UpdateTodoInput {
  completed: boolean;
}
