// chakra imports
import { Box, Flex, Stack } from "@chakra-ui/react";
//   Custom components
import Links from "components/sidebar/components/Links";
import SidebarCard from "components/sidebar/components/SidebarCard";
import SidebarBrand from "./Brand";

// FUNCTIONS

function SidebarContent(props) {
  const { routes } = props;
  // SIDEBAR

  const firstStackRoutes = routes.slice(0, 4);
  const secondStackRoutes = routes.slice(4, 10);
  return (
    <Flex
      direction="column"
      height="100%"
      pt="25px"
      px="16px"
      borderRadius="30px"
    >
      <SidebarBrand />
      <Stack direction="column" mb="auto" mt="8px">
        <Box ps="20px" pe={{ md: "16px", "2xl": "1px" }}>
          <Links routes={firstStackRoutes} />
        </Box>
      </Stack>

      <Stack direction="column" mb="auto" mt="8px">
        <Box ps="20px" pe={{ md: "16px", "2xl": "1px" }}>
          Settings
          <Links routes={secondStackRoutes} />
        </Box>
      </Stack>

      <Box mr="15px" ml="5px" mt="60px" mb="40px" borderRadius="30px">
        <SidebarCard />
      </Box>
    </Flex>
  );
}

export default SidebarContent;
