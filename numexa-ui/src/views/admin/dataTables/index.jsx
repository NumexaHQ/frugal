// Chakra imports
import { Box } from "@chakra-ui/react";
import DateFilter from "components/DateFilter/date-filter";
import { useEffect, useState } from "react";
import { generateTimeParams } from "utils/utils";
import RequestTable from "views/admin/dataTables/components/new-request-table";
import { requestDataColumn } from "views/admin/dataTables/variables/columnsData";

import { connect } from "react-redux";

const RequestTables = ({ projectId, requests, getRequests }) => {
  // Chakra Color Mode

  const [queryparams, setQueryParams] = useState({});
  useEffect(() => {
    getRequests({ projectId });
  }, []);

  const handleTimeQueryParams = (timeRange) => {
    const queryparams = generateTimeParams(timeRange);
    setQueryParams(queryparams);
  };

  useEffect(() => {
    getRequests({ projectId, queryparams });
  }, [queryparams]);

  return (
    <Box pt={{ base: "130px", md: "80px", xl: "80px" }}>
      <Box
        mb="20px"
        fontSize="20px"
        fontWeight="bold"
        display="flex"
        justifyContent="flex-end"
      >
        <DateFilter queryParams={handleTimeQueryParams} />
      </Box>

      <RequestTable
        columnsData={requestDataColumn}
        tableData={requests}
        title={"All Requests"}
      />
      {/* <SimpleGrid
        mb='20px'
        columns={{ sm: 1, md: 2 }}
        spacing={{ base: "20px", xl: "20px" }}>
        <DevelopmentTable
          columnsData={columnsDataDevelopment}
          tableData={tableDataDevelopment}
        />
        <CheckTable columnsData={columnsDataCheck} tableData={tableDataCheck} />
        <ColumnsTable
          columnsData={requestDataColumn}
          tableData={tableDataColumns}
        />
        <ComplexTable
          columnsData={columnsDataComplex}
          tableData={tableDataComplex}
        />
      </SimpleGrid> */}
    </Box>
  );
};

const mapState = (state) => ({
  requests: state.ListRequests.requests || [],
  projectId: state.CommonState.projectID,
});

const mapDispatch = (dispatch) => ({
  getRequests: dispatch.ListRequests.getProviderRequests,
});

export default connect(mapState, mapDispatch)(RequestTables);
