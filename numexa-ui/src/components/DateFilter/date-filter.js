// import { Box, Button, Center, HStack } from "@chakra-ui/react";
// import { useState } from "react";
// import "react-datepicker/dist/react-datepicker.css";

// const DateFilter = ({ onDateFilter, onTimeRangeFilter }) => {
//   const [startDate, setStartDate] = useState(null);
//   const [endDate, setEndDate] = useState(null);

//   const handleFilter = () => {
//     if (startDate && endDate) {
//       onDateFilter(startDate, endDate);
//     }
//   };

//   const handleTimeRangeFilter = (timeRange) => {
//     // Implement your time range filtering logic here
//     onTimeRangeFilter(timeRange);
//   };

//   return (
//     <Box>
//       <Center>
//         {/* <FormControl>
//           <FormLabel>Start Date</FormLabel>
//           <DatePicker
//             selected={startDate}
//             onChange={(date) => setStartDate(date)}
//             dateFormat="yyyy-MM-dd"
//           />
//         </FormControl>
//         <FormControl>
//           <FormLabel>End Date</FormLabel>
//           <DatePicker
//             selected={endDate}
//             onChange={(date) => setEndDate(date)}
//             dateFormat="yyyy-MM-dd"
//           />
//         </FormControl> */}
//         {/* <Button onClick={handleFilter} colorScheme="teal">
//           Filter
//         </Button> */}
//         <HStack spacing={4}>
//           <Button onClick={() => handleTimeRangeFilter("24h")}>24H</Button>
//           <Button onClick={() => handleTimeRangeFilter("7d")}>7D</Button>
//           <Button onClick={() => handleTimeRangeFilter("1m")}>1M</Button>
//           <Button onClick={() => handleTimeRangeFilter("3m")}>3M</Button>
//         </HStack>
//       </Center>
//     </Box>
//   );
// };

// export default DateFilter;

import { Box, Button, Center, HStack } from "@chakra-ui/react";
import { useState } from "react";
import "react-datepicker/dist/react-datepicker.css";

const DateFilter = ({ queryParams }) => {
  const [selectedTimeRange, setSelectedTimeRange] = useState("all");

  const handleTimeRangeFilter = (timeRange) => {
    setSelectedTimeRange(timeRange);
    queryParams(timeRange);
    // Update the selected time range
    // onTimeRangeFilter(timeRange);
  };

  return (
    <Box>
      <Center>
        <HStack spacing={4}>
          <Button
            onClick={() => handleTimeRangeFilter("all")}
            colorScheme={selectedTimeRange === "all" ? "teal" : undefined}
          >
            All
          </Button>
          <Button
            onClick={() => handleTimeRangeFilter("24h")}
            colorScheme={selectedTimeRange === "24h" ? "teal" : undefined}
          >
            24H
          </Button>
          <Button
            onClick={() => handleTimeRangeFilter("7d")}
            colorScheme={selectedTimeRange === "7d" ? "teal" : undefined}
          >
            7D
          </Button>
          <Button
            onClick={() => handleTimeRangeFilter("1m")}
            colorScheme={selectedTimeRange === "1m" ? "teal" : undefined}
          >
            1M
          </Button>
          <Button
            onClick={() => handleTimeRangeFilter("3m")}
            colorScheme={selectedTimeRange === "3m" ? "teal" : undefined}
          >
            3M
          </Button>
        </HStack>
      </Center>
    </Box>
  );
};

export default DateFilter;
