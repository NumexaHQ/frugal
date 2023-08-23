import {
  Flex,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tr,
  useColorModeValue,
} from "@chakra-ui/react";
import React, { useMemo, useState } from "react";
import {
  useExpanded,
  useGlobalFilter,
  usePagination,
  useSortBy,
  useTable,
} from "react-table";

// Custom components
import Card from "components/card/Card";
import { formatDateTime } from "utils/utils";
import ResponseTable from "views/admin/dataTables/components/response-table";

export const responsetDataColumn = [
  {
    Header: "ID",
    accessor: "id",
  },
  {
    Header: "CREATED AT",
    accessor: "created",
  },
  {
    Header: "Model",
    accessor: "model",
  },
  {
    Header: "Message",
    accessor: "choices[0].text",
  },
  // {
  //   Header: "Prompt",
  //   accessor: (row) => {
  //     const requestBody = JSON.parse(row.request_body);
  //     return requestBody.prompt;
  //   }
  // },
  {
    Header: "Completion tokens",
    accessor: "usage.completion_tokens",
  },
  {
    Header: "Prompt Tokens",
    accessor: "usage.prompt_tokens",
  },
  {
    Header: "Total tokens",
    accessor: "usage.total_tokens",
  },
];

export default function RequestTable(props) {
  const { columnsData, tableData } = props;

  const columns = useMemo(() => columnsData, [columnsData]);
  const data = useMemo(() => tableData, [tableData]) || [];

  const [pageSize, setPageSize] = useState(15);

  const tableInstance = useTable(
    {
      columns,
      data,
      initialState: {
        pageIndex: 0,
        pageSize: 15,
        expanded: {}, // Initial state of expanded rows
      },
    },
    useGlobalFilter,
    useSortBy,
    useExpanded,
    usePagination
  );

  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    page,
    prepareRow,

    state: { pageIndex, expanded },
    gotoPage,
    nextPage,
    previousPage,
    canNextPage,
    canPreviousPage,
    pageCount,
  } = tableInstance;

  const textColor = useColorModeValue("secondaryGray.900", "white");
  const borderColor = useColorModeValue("gray.200", "whiteAlpha.100");

  const handleRowClick = (row) => {
    tableInstance.toggleRowExpanded(row.id);
  };

  return (
    <Card
      direction="column"
      w="100%"
      px="0px"
      overflowX={{ sm: "scroll", lg: "hidden" }}
    >
      <Flex px="25px" justify="space-between" mb="20px" align="center">
        <Text
          color={textColor}
          fontSize="22px"
          fontWeight="700"
          lineHeight="100%"
        >
          {props.title}
        </Text>
      </Flex>
      <Table {...getTableProps()} variant="simple" color="gray.500" mb="24px">
        <Thead>
          {headerGroups.map((headerGroup, index) => (
            <Tr {...headerGroup.getHeaderGroupProps()} key={index}>
              {headerGroup.headers.map((column, index) => (
                <Th
                  {...column.getHeaderProps(column.getSortByToggleProps())}
                  pe="10px"
                  key={index}
                  borderColor={borderColor}
                >
                  <Flex
                    justify="space-between"
                    align="center"
                    fontSize={{ sm: "10px", lg: "12px" }}
                    color="gray.400"
                  >
                    {column.render("Header")}
                  </Flex>
                </Th>
              ))}
            </Tr>
          ))}
        </Thead>
        <Tbody {...getTableBodyProps()}>
          {page.map((row, index) => {
            prepareRow(row);
            return (
              <React.Fragment key={index}>
                <Tr
                  {...row.getRowProps()}
                  key={index}
                  onClick={() => handleRowClick(row)}
                  style={{ cursor: "pointer" }}
                >
                  {row.cells.map((cell, index) => {
                    let data = "";
                    if (cell.column.Header === "CREATED AT") {
                      data = (
                        <Text color={textColor} fontSize="sm" fontWeight="700">
                          {formatDateTime(cell.value)}
                        </Text>
                      );
                    } else {
                      data = (
                        <Text color={textColor} fontSize="sm" fontWeight="700">
                          {cell.value}
                        </Text>
                      );
                    }
                    return (
                      <Td
                        {...cell.getCellProps()}
                        key={index}
                        fontSize={{ sm: "14px" }}
                        minW={{ sm: "150px", md: "200px", lg: "auto" }}
                        borderColor="transparent"
                      >
                        {data}
                      </Td>
                    );
                  })}
                </Tr>
                {row.isExpanded ? (
                  <Tr>
                    <Td colSpan={columns.length}>
                      {/* <Text fontWeight="700" mb="8px">
                        Details for row with ID: {row.original.request_id}
                      </Text> */}
                      <ResponseTable
                        columnsData={responsetDataColumn}
                        tableData={[]}
                        requestId={row.original.request_id}
                      />
                    </Td>
                  </Tr>
                ) : null}
              </React.Fragment>
            );
          })}
        </Tbody>
      </Table>
    </Card>
  );
}
