// Chakra imports
import { Box, Flex, Spinner } from "@chakra-ui/react";
import DateFilter from "components/DateFilter/date-filter";
import { useEffect, useState } from "react";
import { generateTimeParams } from "utils/utils";
import RequestTable from "views/admin/dataTables/components/new-request-table";
import { requestDataColumn } from "views/admin/dataTables/variables/columnsData";

import { connect } from "react-redux";
// import Pagination from "./components/pagination";

import {
  Pagination,
  PaginationContainer,
  PaginationNext,
  PaginationPage,
  PaginationPageGroup,
  PaginationPrevious,
  usePagination,
} from "@ajna/pagination";

const RequestTables = ({
  projectId,
  requests,
  getRequests,
  getTotalRequestsCount,
  totalRequestsCount,
}) => {
  const [queryparams, setQueryParams] = useState({});
  useEffect(() => {
    setLoading(true);
    getTotalRequestsCount({ projectId });
    getRequests({ projectId });
    setLoading(false);
  }, []);

  let totalPages = Math.ceil(totalRequestsCount / 10);

  const { currentPage, setCurrentPage, pagesCount, pages } = usePagination({
    pagesCount: totalPages,
    initialState: { currentPage: 1 },
  });

  const handleTimeQueryParams = (timeRange) => {
    const queryparams = generateTimeParams(timeRange);
    setCurrentPage(1);
    setQueryParams(queryparams);
  };

  // loading
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    getTotalRequestsCount({ projectId, queryparams });
    getRequests({ projectId, queryparams, currentPage });
    setLoading(false);
  }, [queryparams, currentPage, totalPages]);

  return (
    <Box pt={{ base: "130px", md: "80px", xl: "80px" }}>
      {loading && (
        <Flex justify="center" align="center">
          <Spinner
            thickness="4px"
            speed="0.65s"
            emptyColor="gray.200"
            color="brand.500"
            size="xl"
            mt="50px"
          />
        </Flex>
      )}

      <Flex justify="space-between" align="center" mt="15px">
        <Box
          mb="20px"
          fontSize="20px"
          fontWeight="bold"
          display="flex"
          justifyContent="flex-start"
        >
          <Pagination
            pagesCount={pagesCount}
            currentPage={currentPage}
            onPageChange={setCurrentPage}
          >
            <PaginationContainer>
              <PaginationPrevious
                _hover={{
                  bg: "brand.400",
                  textColor: "white",
                }}
                _current={{
                  bg: "brand.500",
                  textColor: "white",
                }}
              >
                Previous
              </PaginationPrevious>
              <PaginationPageGroup>
                {pages.map((page) => (
                  <PaginationPage
                    key={`pagination_page_${page}`}
                    page={page}
                    w={10}
                    fontSize="sm"
                    _hover={{
                      bg: "brand.400",
                      textColor: "white",
                    }}
                    _current={{
                      bg: "brand.500",
                      textColor: "white",
                    }}
                  />
                ))}
              </PaginationPageGroup>
              <PaginationNext
                _hover={{
                  bg: "brand.400",
                  textColor: "white",
                }}
                _current={{
                  bg: "brand.500",
                  textColor: "white",
                }}
              >
                Next
              </PaginationNext>
            </PaginationContainer>
          </Pagination>
        </Box>
        <Box
          mb="20px"
          fontSize="20px"
          fontWeight="bold"
          display="flex"
          justifyContent="flex-end"
        >
          <DateFilter
            justifyContent="flex-end"
            queryParams={handleTimeQueryParams}
          />
        </Box>
      </Flex>

      <RequestTable
        columnsData={requestDataColumn}
        tableData={requests}
        title={"All Requests"}
      />
    </Box>
  );
};

const mapState = (state) => ({
  requests: state.ListRequests.requests || [],
  totalRequestsCount: state.TotalRequests.totalRequest,
  projectId: state.CommonState.projectID,
});

const mapDispatch = (dispatch) => ({
  getRequests: dispatch.ListRequests.getProviderRequests,
  getTotalRequestsCount: dispatch.TotalRequests.getTotalRequest,
});

export default connect(mapState, mapDispatch)(RequestTables);
