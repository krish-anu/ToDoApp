import {
  Box,
  Flex,
  Button,
  Text,
  Container,

} from "@chakra-ui/react";
import { IoMoon } from "react-icons/io5";
import { LuSun } from "react-icons/lu";
import {useColorMode} from "./useColorMode"
export default function Navbar() {
  const { colorMode, toggleColorMode } = useColorMode();

  return (
    <Container maxW="900px">
      <Box
        bg={colorMode === "light" ? "green.400" : "red.700"}
        px={4}
        my={4}
        borderRadius="5"
      >
        <Flex h={16} alignItems="center" justifyContent="space-between">
          <Flex align="center" gap={3}>
            <Text fontSize="xl">Daily Tasks</Text>
            <Button onClick={toggleColorMode}>
              {colorMode === "light" ? <IoMoon /> : <LuSun size={20} />}
            </Button>
          </Flex>
        </Flex>
      </Box>
    </Container>
  );
}
