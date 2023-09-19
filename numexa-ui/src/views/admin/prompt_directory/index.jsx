// Chakra imports
import { Box } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { generateTimeParams } from "utils/utils";
import { PromptManagementColumn } from "views/admin/dataTables/variables/columnsData";

import PromptTable from "./components/prompt-table";

import { connect } from "react-redux";

const PromptManagement = ({
  projectId,
  getRequests,
  listPromptDirectory,
  promptDirectory,
}) => {
  const [queryparams, setQueryParams] = useState({});
  useEffect(() => {
    listPromptDirectory({ projectId });
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
      <PromptTable
        columnsData={PromptManagementColumn}
        tableData={promptDirectory}
        title={"Prompt Evaluation Table"}
      />
    </Box>
  );
};

const mapState = (state) => ({
  requests: state.ListRequests.requests || [],
  promptDirectory: state.ListPromptDirectory.promptDirectory || [],
  projectId: state.CommonState.projectID,
});

const mapDispatch = (dispatch) => ({
  getRequests: dispatch.ListRequests.getProviderRequests,
  listPromptDirectory: dispatch.ListPromptDirectory.listPromptDirectory,
});

export default connect(mapState, mapDispatch)(PromptManagement);
