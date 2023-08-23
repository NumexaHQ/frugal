import { Avatar, Box } from "@chakra-ui/react";

export const columnsDataCheck = [
  {
    Header: "NAME",
    accessor: "name",
  },
  {
    Header: "PROGRESS",
    accessor: "progress",
  },
  {
    Header: "QUANTITY",
    accessor: "quantity",
  },
  {
    Header: "DATE",
    accessor: "date",
  },
];
export const columnsModelData = [
  {
    Header: "NAME",
    accessor: "name",
  },
  {
    Header: "COUNT",
    accessor: "count",
  },
];

export const columnsUsersUsageStats = [
  {
    Header: "USER",
    accessor: (row) => {
      return (
        <Box w="100%" display="flex" alignItems="center">
          <Box>
            <Avatar
              _hover={{ cursor: "arrow" }}
              color="white"
              name={row.email}
              bg="brand.500"
              size="sm"
              w="25px"
              h="25px"
            />
          </Box>
          <Box>&nbsp;&nbsp; {row.email}</Box>
        </Box>
      );
    },
  },
  {
    Header: "REQUESTS",
    accessor: "total_request",
  },
  {
    Header: "COST",
    accessor: (row) => {
      if (row.cost !== 0) {
        // trim the cost to 6 decimal places
        return `$ ${row.cost.toFixed(6)}`;
      }
      return `$ ${row.cost}`;
    },
  },
];
