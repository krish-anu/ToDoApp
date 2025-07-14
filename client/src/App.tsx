import React from "react";
import { Container, Stack } from "@chakra-ui/react";
import Navbar from "./components/ui/Navbar";
import TodoForm from "./components/ui/TodoForm";
import TodoList from "./components/ui/TodoList";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

// Create a QueryClient instance
const queryClient = new QueryClient();

export const BASE_URL =
  import.meta.env.MODE === "development" ? "http://localhost:5000/api" : "/api";

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Stack h="100vh">
        <Navbar />
        <Container>
          <TodoForm />
          <TodoList />
        </Container>
      </Stack>
    </QueryClientProvider>
  );
}

export default App;
