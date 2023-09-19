// Chakra imports
import { Box, Flex, Icon, Text, useColorModeValue } from "@chakra-ui/react";
import BarChart from "components/charts/BarChart";
import { useEffect, useState } from "react";

// Custom components
import Card from "components/card/Card.js";

// Assets
import { RiArrowUpSFill } from "react-icons/ri";
export default function DailyTraffic(props) {
  const { ...rest } = props;

  console.log(props);

  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check if xaxisCategories is available and dataRequestsPerDay has data
    if (props.xaxisCategories?.length > 0) {
      setIsLoading(false); // Data is available, set isLoading to false
    } else {
      setIsLoading(true); // Data is not available, set isLoading to true
    }
  }, [props.xaxisCategories]);

  const barChartDataDailyTraffic = [
    {
      name: "Daily requests",
      data: [20, 30, 40, 20, 45, 50, 30],
    },
  ];

  const barChartOptionsDailyTraffic = {
    chart: {
      toolbar: {
        show: false,
      },
    },
    tooltip: {
      style: {
        fontSize: "12px",
        fontFamily: undefined,
      },
      onDatasetHover: {
        style: {
          fontSize: "12px",
          fontFamily: undefined,
        },
      },
      theme: "dark",
    },
    xaxis: {
      categories: props.xaxisCategories,
      show: false,
      labels: {
        show: true,
        style: {
          colors: "#A3AED0",
          fontSize: "14px",
          fontWeight: "500",
        },
      },
      axisBorder: {
        show: false,
      },
      axisTicks: {
        show: false,
      },
    },
    yaxis: {
      show: false,
      color: "black",
      labels: {
        show: true,
        style: {
          colors: "#CBD5E0",
          fontSize: "14px",
        },
      },
    },
    grid: {
      show: false,
      strokeDashArray: 5,
      yaxis: {
        lines: {
          show: true,
        },
      },
      xaxis: {
        lines: {
          show: false,
        },
      },
    },
    fill: {
      type: "gradient",
      gradient: {
        type: "vertical",
        shadeIntensity: 1,
        opacityFrom: 0.7,
        opacityTo: 0.9,
        colorStops: [
          [
            {
              offset: 0,
              color: "teal", // Change color to teal
              opacity: 1,
            },
            {
              offset: 100,
              color: "rgba(0, 128, 128, 0.28)", // Change color to teal with opacity
            },
          ],
        ],
      },
    },
    dataLabels: {
      enabled: false,
    },
    plotOptions: {
      bar: {
        borderRadius: 5,
        columnWidth: "20px",
      },
    },
  };

  // Chakra Color Mode
  const textColor = useColorModeValue("secondaryGray.900", "white");
  return (
    <Card align="center" direction="column" w="100%" {...rest}>
      <Flex justify="space-between" align="start" px="10px" pt="5px">
        <Flex flexDirection="column" align="start" me="20px">
          <Flex w="100%">
            <Text
              me="auto"
              color="secondaryGray.600"
              fontSize="sm"
              fontWeight="500"
            >
              Daily Requests
            </Text>
          </Flex>
        </Flex>
        <Flex align="center">
          <Icon as={RiArrowUpSFill} color="green.500" />
          {/* <Text color="green.500" fontSize="sm" fontWeight="700">
            +2.45%
          </Text> */}
        </Flex>
      </Flex>
      <Box h="240px" mt="auto">
        {isLoading ? (
          <Text>Loading...</Text>
        ) : props.xaxisCategories?.length > 0 ? (
          <BarChart
            chartData={props.dataRequestsPerDay}
            chartOptions={barChartOptionsDailyTraffic}
          />
        ) : (
          <Text>No data available</Text>
        )}
      </Box>
    </Card>
  );
}
