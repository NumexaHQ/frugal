
// Chakra imports
import { Box } from "@chakra-ui/react";
import { useEffect } from "react";
import RequestTable from "views/admin/dataTables/components/new-request-table";
import {
  requestDataColumn
} from "views/admin/dataTables/variables/columnsData";

import { connect } from 'react-redux';

const RequestTables = ({ projectId, requests, getRequests }) => {
  // Chakra Color Mode
  useEffect(() => {
    getRequests({ projectId });
  }, []);

  return (
    <Box pt={{ base: "130px", md: "80px", xl: "80px" }}>
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
}


const mapState = (state) => (
  {
    requests: state.ListRequests.requests || [],
    projectId: state.CommonState.projectID,
  });

const mapDispatch = (dispatch) => ({
  getRequests: dispatch.ListRequests.getProviderRequests,
});


export default connect(mapState, mapDispatch)(RequestTables);
