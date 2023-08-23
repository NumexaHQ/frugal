import { Button, Flex, Link, Text, useColorModeValue } from "@chakra-ui/react";
import { LuExternalLink } from "react-icons/lu";

export default function SidebarDocs() {
  const bgColor = "linear-gradient(135deg, #17A2B8 0%, #007B8A 100%)";
  const borderColor = useColorModeValue("white", "navy.800");

  return (
    <Flex
      justify="center"
      direction="column"
      align="center"
      bg={bgColor}
      borderRadius="30px"
      position="relative"
    >
      <Flex
        direction="column"
        mb="12px"
        align="center"
        justify="center"
        px="15px"
      >
        <Text
          fontSize={{ base: "lg", xl: "18px" }}
          color="white"
          fontWeight="bold"
          lineHeight="150%"
          textAlign="center"
          px="10px"
          mt="10px"
          mb="6px"
        >
          🔐 Private beta
        </Text>
        <Text
          fontSize="14px"
          color={"white"}
          fontWeight="500"
          px="10px"
          mb="6px"
          textAlign="center"
        >
          Start Saving Cost & understand your LLM Usage with Numexa!
        </Text>
      </Flex>
      <Link href="https://docs.numexa.io/" target="blank">
        <Button
          bg="whiteAlpha.300"
          _hover={{ bg: "whiteAlpha.200" }}
          _active={{ bg: "whiteAlpha.100" }}
          mb={{ sm: "16px", xl: "24px" }}
          color={"white"}
          fontWeight="regular"
          fontSize="sm"
          minW="185px"
          mx="auto"
        >
          {/* icon */}
          <LuExternalLink />
          &nbsp; Read the Docs
        </Button>
      </Link>
    </Flex>
  );
}
