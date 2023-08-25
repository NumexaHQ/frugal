import {
  Button,
  Flex,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tr,
  useColorModeValue,
  useDisclosure,
} from "@chakra-ui/react";
import React, { useMemo } from "react";
import {
  useExpanded,
  useGlobalFilter,
  usePagination,
  useSortBy,
  useTable,
} from "react-table";

import { useState } from "react";

import Card from "components/card/Card";
import { formatDateTime } from "utils/utils";

import { connect } from "react-redux";
import NoData from "./noData";
import ResponseDrawer from "./response-drawer";

// ... other imports

function RequestTable(props) {
  const { columnsData, tableData, getResponse } = props;

  const columns = useMemo(() => columnsData, [columnsData]);
  const data = useMemo(() => tableData, [tableData]) || [];

  const { isOpen, onOpen, onClose } = useDisclosure(); // Initialize useDisclosure

  const tableInstance = useTable(
    {
      columns,
      data,
      initialState: {
        pageIndex: 0,
        pageSize: 10,
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
    page, // Access the current page
    prepareRow,
    canPreviousPage,
    canNextPage,
    pageOptions,
    pageCount,
    gotoPage,
    nextPage,
    previousPage,
    state: { pageIndex },
    // ... other hooks
  } = tableInstance;

  const textColor = useColorModeValue("secondaryGray.900", "white");
  const borderColor = useColorModeValue("gray.200", "whiteAlpha.100");

  const [selectedRowData, setSelectedRowData] = useState(null);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false); // State to hold selected row data

  const handleRowClick = (event, row) => {
    event.stopPropagation();

    // Update the selected row data state with the clicked row's data
    setSelectedRowData(row.original);
    setIsDrawerOpen(true);
    getResponse({ requestId: row.original.id });
  };

  return (
    <Card
      direction="column"
      w="100%"
      px="0px"
      overflowX={{ sm: "scroll", lg: "hidden" }}
    >
      {data.length === 0 ? (
        <NoData />
      ) : (
        <>
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

          <Table
            {...getTableProps()}
            variant="simple"
            color="gray.500"
            mb="24px"
          >
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
                      onClick={(event) => handleRowClick(event, row)}
                      style={{ cursor: "pointer" }}
                    >
                      {row.cells.map((cell, index) => {
                        let data = "";
                        if (cell.column.Header === "CREATED AT") {
                          data = (
                            <Text
                              color={textColor}
                              fontSize="sm"
                              fontWeight="700"
                            >
                              {formatDateTime(cell.value)}
                            </Text>
                          );
                        } else {
                          data = (
                            <Text
                              color={textColor}
                              fontSize="sm"
                              fontWeight="700"
                            >
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
                  </React.Fragment>
                );
              })}
            </Tbody>
          </Table>

          {/* Pagination controls */}
          <Flex justify="space-between" align="center" mt="15px">
            <Text color={textColor} fontSize="14px" fontWeight="700" ml={8}>
              Page {pageIndex + 1} of {pageOptions.length}
            </Text>
            <Flex mr={8}>
              <Button
                onClick={() => gotoPage(0)}
                disabled={!canPreviousPage}
                mr="10px"
              >
                {"<<"}
              </Button>
              <Button
                onClick={previousPage}
                disabled={!canPreviousPage}
                mr="10px"
              >
                {"<"}
              </Button>
              <Button onClick={nextPage} disabled={!canNextPage} mr="10px">
                {">"}
              </Button>
              <Button
                onClick={() => gotoPage(pageCount - 1)}
                disabled={!canNextPage}
              >
                {">>"}
              </Button>
            </Flex>
          </Flex>
        </>
      )}

      {isDrawerOpen && (
        <ResponseDrawer
          isOpen={isDrawerOpen}
          onClose={() => setIsDrawerOpen(false)}
          columnsData={columnsData}
          onOpen={() => setIsDrawerOpen(true)}
          selectedRowData={selectedRowData}
          data={props.response}
        />
      )}
    </Card>
  );
}

const mapState = (state) => ({
  response: state.ListResponse.response,
});

const mapDispatch = (dispatch) => ({
  getResponse: dispatch.ListResponse.getResponse,
});

export default connect(mapState, mapDispatch)(RequestTable);
