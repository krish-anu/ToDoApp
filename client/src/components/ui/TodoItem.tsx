import { Badge, Box, Flex, Spinner, Text } from "@chakra-ui/react";
import { FaCheckCircle } from "react-icons/fa";
import { MdDelete } from "react-icons/md";
import type { Todo } from "./TodoList";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const BASE_URL = "http://localhost:5000/api";

const TodoItem = ({ todo }: { todo: Todo }) => {
  const queryClient = useQueryClient();

  const { mutate: updateTodo, isPending: isUpdating } = useMutation({
    mutationKey: ["updateTodo"],
    mutationFn: async (id: string) => {
      if (todo.completed) throw new Error("Todo is already completed");

      const res = await fetch(`${BASE_URL}/todos/${id}`, {
        method: "PATCH",
      });

      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Update failed");
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
    },
    onError: (error: any) => {
      alert(`Update error: ${error.message || error}`);
    },
  });

  const { mutate: deleteTodo, isPending: isDeleting } = useMutation({
    mutationKey: ["deleteTodo"],
    mutationFn: async (id: string) => {
      const res = await fetch(`${BASE_URL}/todos/${id}`, {
        method: "DELETE",
      });

      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Delete failed");
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
    },
    onError: (error: any) => {
      alert(`Delete error: ${error.message || error}`);
    },
  });

  return (
    <Flex gap={2} alignItems="center">
      <Flex
        flex={1}
        alignItems="center"
        border="1px"
        borderColor="gray.600"
        p={2}
        borderRadius="lg"
        justifyContent="space-between"
      >
        <Text
          color={todo.completed ? "green.200" : "yellow.100"}
          textDecoration={todo.completed ? "line-through" : "none"}
        >
          {todo.body}
        </Text>
        {todo.completed ? (
          <Badge ml="1" colorScheme="green">
            Done
          </Badge>
        ) : (
          <Badge ml="1" colorScheme="yellow">
            In Progress
          </Badge>
        )}
      </Flex>

      <Flex gap={2} alignItems="center">
        <Box
          color="green.500"
          cursor="pointer"
          onClick={() => updateTodo(todo._id)}
        >
          {!isUpdating ? <FaCheckCircle size={20} /> : <Spinner size="sm" />}
        </Box>
        <Box
          color="red.500"
          cursor="pointer"
          onClick={() => deleteTodo(todo._id)}
        >
          {!isDeleting ? <MdDelete size={25} /> : <Spinner size="sm" />}
        </Box>
      </Flex>
    </Flex>
  );
};

export default TodoItem;
