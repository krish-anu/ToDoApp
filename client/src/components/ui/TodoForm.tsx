import { Box, Flex, Button, Text, Container, Input } from "@chakra-ui/react";
import { IoMoon } from "react-icons/io5";
import { LuSun } from "react-icons/lu";
import { useColorMode } from "./useColorMode"; // Make sure Chakra is correctly installed
import { useEffect, useRef, useState } from "react";

export default function Navbar() {
  const { colorMode, toggleColorMode } = useColorMode();
  const [newTodo, setNewTodo] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

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
          />
          <Button
            colorScheme="teal"
            onClick={() => {
              if (newTodo.trim()) {
                console.log("New todo:", newTodo);
                setNewTodo("");
                inputRef.current?.focus(); // refocus after adding
              }
            }}
          >
            Add
          </Button>
        </Flex>
      </Box>
    </Container>
  );
}
