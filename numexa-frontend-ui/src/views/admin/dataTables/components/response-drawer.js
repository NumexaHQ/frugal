import {
  Drawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerOverlay,
  useColorModeValue,
} from "@chakra-ui/react";
import React, { useMemo } from "react";
import ReactJson from "react-json-view";
import { connect } from "react-redux";
import DrawerTable from "./drawer-table";

// response Drawer
function ResponseDrawer(props) {
  const {
    columnsData,
    requestId,
    isOpen,
    onOpen,
    onClose,
    selectedRowData,
    data,
    drawerData,
  } = props;

  // useEffect(() => {
  //   props.getResponse({requestId: request_Id});
  // }, []);

  const columns = useMemo(() => columnsData, [columnsData]);
  // const data = useMemo(() => tableData, [tableData]) || [];

  const firstField = React.useRef();

  const textColor = useColorModeValue("secondaryGray.900", "white");
  const borderColor = useColorModeValue("gray.200", "whiteAlpha.100");
  return (
    <>
      <Drawer
        isOpen={isOpen}
        placement="right"
        initialFocusRef={firstField}
        onClose={onClose}
        size="lg"
      >
        <DrawerOverlay />
        <DrawerContent>
          <DrawerCloseButton />
          <DrawerHeader borderBottomWidth="1px">Information</DrawerHeader>

          <DrawerBody>
            <DrawerTable selectedRowData={selectedRowData} />
            <ReactJson src={data} theme="solarized" displayDataTypes={false} />
          </DrawerBody>

          <DrawerFooter borderTopWidth="1px"></DrawerFooter>
        </DrawerContent>
      </Drawer>
    </>
  );
}

const mapState = (state) => ({
  response: state.ListResponse.response || [],
});

const mapDispatch = (dispatch) => ({
  getResponse: dispatch.ListResponse.getResponse,
});

export default connect(mapState, mapDispatch)(ResponseDrawer);
