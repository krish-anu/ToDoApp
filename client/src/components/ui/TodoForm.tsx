import { Box, Flex, Button, Text, Container, Input } from "@chakra-ui/react";
import { IoMoon } from "react-icons/io5";
import { LuSun } from "react-icons/lu";
import { useColorMode } from "./useColorMode";
import { useEffect, useRef, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const BASE_URL = "http://localhost:5000/api"; // Your backend URL

export default function Navbar() {
  const { colorMode, toggleColorMode } = useColorMode();
  const [newTodo, setNewTodo] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);
  const queryClient = useQueryClient();

  // Mutation for adding a new todo
  const { mutate: addTodo, isPending } = useMutation({
    mutationFn: async (body: string) => {
      const res = await fetch(`${BASE_URL}/todos`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ body }),
      });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || "Failed to add todo");
      }
      return res.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] }); // refresh todos list
      setNewTodo("");
      inputRef.current?.focus();
    },
  });

  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  return (
    <Container maxW="900px">
      <Box
        bg={colorMode === "light" ? "gray.400" : "gray.700"}
        px={4}
        my={4}
        borderRadius="md"
      >
        <Flex h={16} alignItems="center" justifyContent="space-between">
          <Flex align="center" gap={3}>
            <Text fontSize="xl" fontWeight="bold">
              Daily Tasks
            </Text>
            <Button onClick={toggleColorMode}>
              {colorMode === "light" ? <IoMoon /> : <LuSun size={20} />}
            </Button>
          </Flex>
        </Flex>

        <Flex mt={4} gap={3}>
          <Input
            ref={inputRef}
            type="text"
            value={newTodo}
            onChange={(e) => setNewTodo(e.target.value)}
            placeholder="Add a new task"
            disabled={isPending}
          />
          <Button
            colorScheme="teal"
            onClick={() => {
              if (newTodo.trim()) {
                addTodo(newTodo.trim());
              }
            }}
            loading={isPending}
          >
            Add
          </Button>
        </Flex>
      </Box>
    </Container>
  );
}
