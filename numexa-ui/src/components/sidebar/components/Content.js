// chakra imports
import { Box, Flex, Stack } from "@chakra-ui/react";
//   Custom components
import Links from "components/sidebar/components/Links";
import SidebarCard from "components/sidebar/components/SidebarCard";
import { useEffect } from "react";
import { connect } from "react-redux";
import SidebarBrand from "./Brand";
import UsageBar from "./UsageBar";

// FUNCTIONS

function SidebarContent(props) {
  const { routes, projectId, usage, getUsage } = props;
  // SIDEBAR

  const firstStackRoutes = routes.slice(0, 5);
  const secondStackRoutes = routes.slice(5, 10);

  useEffect(() => {
    getUsage({ projectId });
  }, []);

  console.log("usage, projectId", usage, projectId);

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
      <Stack direction="column" mb="auto" mt="8px">
        <Box borderRadius="lg" pt="20px">
          <UsageBar usage={usage} />
        </Box>
      </Stack>

      <Box mr="15px" ml="5px" mt="30px" mb="40px" borderRadius="30px">
        <SidebarCard />
      </Box>
    </Flex>
  );
}

const mapState = (state) => ({
  projectId: state.CommonState.projectID,
  usage: state.Usage.usage,
});

const mapDispatch = (dispatch) => ({
  getUsage: dispatch.Usage.getUsage,
});

export default connect(mapState, mapDispatch)(SidebarContent);
