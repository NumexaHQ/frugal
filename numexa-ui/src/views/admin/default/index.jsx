// Chakra imports
import {
  Box,
  Icon,
  SimpleGrid,
  Switch,
  Text,
  useColorModeValue,
} from "@chakra-ui/react";
// Assets
// Custom components
import DateFilter from "components/DateFilter/date-filter";
import MiniStatistics from "components/card/MiniStatistics";
import IconBox from "components/icons/IconBox";
import { useEffect, useState } from "react";
import { MdCheckCircleOutline, MdHighlightOff } from "react-icons/md";
import { generateTimeParams } from "utils/utils";
import ComplexTable from "views/admin/default/components/ComplexTable";

import {
  columnsModelData,
  columnsUsersUsageStats,
} from "views/admin/default/variables/columnsData";

import { connect } from "react-redux";

function IndicatorDot({ isActive }) {
  return (
    <Box
      w="10px"
      h="10px"
      borderRadius="50%"
      bg={isActive ? "green.500" : "gray.300"}
      boxShadow={isActive ? "0 0 10px rgba(0, 255, 0, 0.8)" : "none"}
      transition="box-shadow 0.3s"
      mr={2}
    />
  );
}

function UserReports(props) {
  // Chakra Color Mode
  const { totalRequests, projectId, avgLatency, avgTokens } = props;

  const [realTimeMetrics, setRealTimeMetrics] = useState(false);

  const [queryparams, setQueryParams] = useState({});

  const handleTimeQueryParams = (timeRange) => {
    const queryparams = generateTimeParams(timeRange);
    console.log(timeRange);
    setQueryParams(queryparams);
  };

  console.log(queryparams);

  useEffect(() => {
    props.getTotalRequests({ projectId });
    props.getAvgLatency({ projectId });
    props.getAvgTokens({ projectId });
    props.geModels({ projectId });
    props.getUsersUsageStat({ projectId });
  }, []);

  useEffect(() => {
    props.getTotalRequests({ projectId, queryparams });
    props.getAvgLatency({ projectId, queryparams });
    props.getAvgTokens({ projectId, queryparams });
    props.geModels({ projectId, queryparams });
    props.getUsersUsageStat({ projectId, queryparams });
  }, [queryparams]);

  useEffect(() => {
    if (realTimeMetrics) {
      // Polling logic
      const interval = setInterval(() => {
        // Fetch metrics here
        props.getTotalRequests({ projectId });
        props.getAvgLatency({ projectId });
        props.getAvgTokens({ projectId });
        props.geModels({ projectId });
        props.getUsersUsageStat({ projectId });
      }, 2000); // Polling interval: 5 seconds

      return () => {
        clearInterval(interval); // Clear interval on component unmount or when real-time metrics are turned off
      };
    }
  }, [realTimeMetrics]);

  const brandColor = useColorModeValue("brand.500", "white");
  const boxBg = useColorModeValue("secondaryGray.300", "whiteAlpha.100");
  const red = useColorModeValue("red.500", "red.200");
  const green = useColorModeValue("green.500", "green.200");
  return (
    <>
      <Box pt={{ base: "130px", md: "80px", xl: "80px" }}>
        {/* <Box display="flex" alignItems="center" mb={5}>
          <IndicatorDot isActive={realTimeMetrics} />
          <Switch
            isChecked={realTimeMetrics}
            onChange={() => setRealTimeMetrics(!realTimeMetrics)}
          />
          <Text ml={2}>Live Metrics</Text>
        </Box> */}

        <Box
          display="flex"
          alignItems="center"
          justifyContent="space-between"
          mb={5}
        >
          <div>
            <Box display="flex" alignItems="center" mb={5}>
              <IndicatorDot isActive={realTimeMetrics} />
              <Switch
                isChecked={realTimeMetrics}
                onChange={() => setRealTimeMetrics(!realTimeMetrics)}
              />
              <Text ml={2}>Live Metrics</Text>
            </Box>
          </div>
          <div>
            <DateFilter queryParams={handleTimeQueryParams} />
          </div>
        </Box>

        <SimpleGrid
          columns={{ base: 4, md: 4, lg: 4, "2xl": 4 }}
          gap="20px"
          mb="20px"
        >
          <MiniStatistics name="Total Requests" value={totalRequests} />
          <MiniStatistics
            name="Total Success Responses"
            startContent={
              <IconBox
                w="56px"
                h="56px"
                bg={boxBg}
                icon={
                  <Icon
                    w="32px"
                    h="32px"
                    as={MdCheckCircleOutline}
                    color={green}
                  />
                }
              />
            }
            value={avgTokens[0] ? avgTokens[0].total_success : 0}
          />
          <MiniStatistics
            name="Total Failed Responses"
            startContent={
              <IconBox
                w="56px"
                h="56px"
                bg={boxBg}
                icon={
                  <Icon w="32px" h="32px" as={MdHighlightOff} color={red} />
                }
              />
            }
            value={avgTokens[0] ? avgTokens[0].total_failure : 0}
          />
          <MiniStatistics
            name="Average Latency/ request"
            value={`${avgLatency.toFixed(2)} seconds`}
          />
          <MiniStatistics
            name="Average Prompt Tokens/ request"
            value={avgTokens[0] ? avgTokens[0].avg_prompt_tokens : 0}
          />
          <MiniStatistics
            name="Average Total Tokens/ request"
            value={avgTokens[0] ? avgTokens[0].avg_total_tokens : 0}
          />
          <MiniStatistics
            name="Average completion Tokens/ request"
            value={avgTokens[0] ? avgTokens[0].avg_completion_tokens : 0}
          />
          <MiniStatistics
            name="Total cost"
            value={avgTokens[0] ? `$ ${avgTokens[0]?.total_cost}` : `$ 0`}
          />
        </SimpleGrid>
        <SimpleGrid columns={{ base: 1, md: 1, xl: 2 }} gap="20px" mb="20px">
          <ComplexTable
            columnsData={columnsModelData}
            tableData={props.modelDistribution}
            title={"🎛️ Requests Per Model"}
          />
          <ComplexTable
            columnsData={columnsUsersUsageStats}
            tableData={props.usersUsageStat}
            title={"👤 User stats"}
          />
        </SimpleGrid>
      </Box>
    </>
  );
}

const mapState = (state) => ({
  totalRequests: state.TotalRequests.totalRequest,
  projectId: state.CommonState.projectID,
  avgLatency: state.AvgLatency.avgLatency,
  avgTokens: state.AvgTokens.avgTokens || [],
  modelDistribution: state.ModelDistribution.modelDistribution || [],
  usersUsageStat: state.UsersUsageStat.usersUsageStat || [],
});

const mapDispatch = (dispatch) => ({
  getTotalRequests: dispatch.TotalRequests.getTotalRequest,
  getAvgLatency: dispatch.AvgLatency.getAvgLatency,
  getAvgTokens: dispatch.AvgTokens.getAvgTokens,
  geModels: dispatch.ModelDistribution.getModelDistribution,
  getUsersUsageStat: dispatch.UsersUsageStat.getUsersUsageStat,
});

export default connect(mapState, mapDispatch)(UserReports);
