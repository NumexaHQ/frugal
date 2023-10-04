import {
  Button,
  Flex,
  FormControl,
  FormLabel,
  HStack,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Table,
  Tag,
  TagCloseButton,
  TagLabel,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tooltip,
  Tr,
  useColorModeValue,
  useToast,
} from "@chakra-ui/react";
import React, { useMemo } from "react";
import {
  useExpanded,
  useGlobalFilter,
  usePagination,
  useSortBy,
  useTable,
} from "react-table";

import { FaPlus } from "react-icons/fa";

import { useState } from "react";

import Card from "components/card/Card";

import { connect } from "react-redux";
import { formatDateTime, humanizeDateTime } from "utils/utils";
import NoData from "./noData";
import ResponseDrawer from "./response-drawer";

// ... other imports

function RequestTable(props) {
  const { columnsData, tableData, getResponse, addtoPromptDirectory } = props;

  const columns = useMemo(() => columnsData, [columnsData]);
  const data = useMemo(() => tableData, [tableData]) || [];

  // State to manage modal open/close and input values
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [comment, setComment] = useState("");
  const [score, setScore] = useState(0);

  const [tags, setTags] = useState([]);
  const [tagInput, setTagInput] = useState("");

  const handleTagAdd = () => {
    if (tagInput.trim() !== "") {
      setTags([...tags, tagInput]);
      setTagInput("");
    }
  };

  const handleTagRemove = (tagToRemove) => {
    const updatedTags = tags.filter((tag) => tag !== tagToRemove);
    setTags(updatedTags);
  };

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
  const toast = useToast();

  const handleRowClick = (event, row) => {
    event.stopPropagation();

    // Update the selected row data state with the clicked row's data
    setSelectedRowData(row.original);
    setIsDrawerOpen(true);
    getResponse({ requestId: row.original.id });
  };

  const handleAddButtonClick = (event, row) => {
    event.stopPropagation();
    setSelectedRowData(row.original);
    setIsModalOpen(true);
    // Update the selected row data state with the clicked row's data
  };

  const handleModalSubmit = () => {
    // Handle modal submit here
    const payloadObj = {
      ...selectedRowData,
      custom_metadata: JSON.stringify(tags),
      comment: comment,
      // score should be a integer
      score: parseInt(score),
    };
    addtoPromptDirectory({ payloadObj: payloadObj });
    setIsModalOpen(false);
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
                        if (cell.column.Header === "Initiated At") {
                          data = (
                            <Text
                              color={textColor}
                              fontSize="sm"
                              fontWeight="700"
                            >
                              {formatDateTime(cell.value)} <br />(
                              {humanizeDateTime(cell.value)})
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
                        if (index === 0) {
                          // Add a button to the last cell in the row
                          return (
                            <Td
                              {...cell.getCellProps()}
                              key={index}
                              fontSize={{ sm: "14px" }}
                              minW={{ sm: "150px", md: "200px", lg: "auto" }}
                              borderColor="transparent"
                            >
                              {data}{" "}
                              <Tooltip label="Add to Prompt Management">
                                <Button
                                  onClick={(event) =>
                                    handleAddButtonClick(event, row)
                                  }
                                  size="sm"
                                  variant="outline"
                                  leftIcon={<FaPlus />}
                                >
                                  Add
                                </Button>
                              </Tooltip>
                            </Td>
                          );
                        } else {
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
                        }
                      })}
                    </Tr>
                  </React.Fragment>
                );
              })}
            </Tbody>
          </Table>
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
        // eslint-disable-next-line react/jsx-no-comment-textnodes
      )}
      {/* Modal for prompot managment */}
      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>
            Evaluate the Prompt by providing a Score, Comments, and Tags.
          </ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <FormControl>
              <FormLabel>Comment</FormLabel>
              <Input
                value={comment}
                onChange={(e) => setComment(e.target.value)}
                placeholder="Enter your comment"
              />
            </FormControl>
            <FormControl mt={4}>
              <FormLabel>Score</FormLabel>
              <Input
                type="number"
                value={score}
                onChange={(e) => setScore(e.target.value)}
                placeholder="Enter the score (default is 0)"
                max={10}
              />
            </FormControl>
            <FormControl mt={4}>
              <FormLabel>
                Custom Tags: Incorporate your own Evaluation Criteria by adding
                a tag.{" "}
              </FormLabel>
              <Input
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                placeholder="Enter a tag"
                onKeyPress={(e) => {
                  if (e.key === "Enter") {
                    handleTagAdd();
                  }
                }}
              />
            </FormControl>
            <HStack mt={2} spacing={2}>
              {tags.map((tag, index) => (
                <Tag key={index} size="md">
                  <TagLabel>{tag}</TagLabel>
                  <TagCloseButton onClick={() => handleTagRemove(tag)} />
                </Tag>
              ))}
            </HStack>
          </ModalBody>
          <ModalFooter>
            <Button colorScheme="blue" onClick={handleModalSubmit}>
              Add
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Card>
  );
}

const mapState = (state) => ({
  response: state.ListResponse.response,
});

const mapDispatch = (dispatch) => ({
  getResponse: dispatch.ListResponse.getResponse,
  addtoPromptDirectory: dispatch.AddtoPromptDirectory.addtoPromptDirectory,
});

export default connect(mapState, mapDispatch)(RequestTable);
