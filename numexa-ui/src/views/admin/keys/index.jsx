import {
  Badge,
  Box,
  Button,
  Flex,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Stack,
  Text,
  useClipboard,
  useDisclosure,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";

import ColumnsTable from "views/admin/dataTables/components/general-table";

import { apiKeycolumn } from "views/admin/dataTables/variables/columnsData";

import { MdContentCopy } from "react-icons/md";
import { connect } from "react-redux";

const ApiKeys = ({ apiKeys, handleListApiKeys, projectId, code, genKey }) => {
  const { hasCopied, onCopy } = useClipboard(code);
  const [keyName, setKeyName] = useState("");

  const {
    isOpen: isOpenMain,
    onOpen: onOpenMain,
    onClose: onCloseMain,
  } = useDisclosure();

  useEffect(() => {
    handleListApiKeys();
  }, []);

  const genKeyClickEvent = () => {
    onOpenMain();
    genKey({ projectId, keyName });
  };

  return (
    <Box pt={{ base: "130px", md: "80px", xl: "80px" }}>
      <Stack spacing={3}>
        <Text fontSize="xl" mb={8}>
          Your secret API keys are listed below. Please note that we do not
          display your secret API keys again after you generate them. Do not
          share your API key with others or expose it in the browser or other
          client-side code.
        </Text>
      </Stack>
      <ColumnsTable
        columnsData={apiKeycolumn}
        tableData={apiKeys}
        title={"Keys"}
      />
      <Text mt={5}>
        <strong>Generate API Key</strong>
      </Text>
      <Flex alignItems="center" mt={3}>
        <Input
          placeholder="your friendly key name"
          width="10px"
          flex="1"
          onChange={(e) => setKeyName(e.target.value)}
        />
        <Button colorScheme="brand" onClick={genKeyClickEvent} ml={3}>
          🔑 Generate
        </Button>
      </Flex>
      <Modal isOpen={isOpenMain} onClose={onCloseMain} size="md">
        <ModalOverlay />
        <ModalContent align="center" justify="center">
          <ModalHeader>Generated Secret API Key 🤫</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <p>
              <strong>
                This is the only time you'll be able to view your API key.
              </strong>
            </p>
            <p>
              Please save it in a secure and easily accessible place. In case
              you lose your API key, you'll need to generate a new one.
            </p>
            <Badge
              colorScheme="brandScheme"
              variant="outline"
              textTransform="none"
              onClick={onCopy}
            >
              {code}
            </Badge>
            <Button colorScheme="brandScheme" size="xs" onClick={onCopy} ml={1}>
              {hasCopied ? "Copied!" : <MdContentCopy />}
            </Button>
          </ModalBody>
          <ModalFooter></ModalFooter>
        </ModalContent>
      </Modal>
    </Box>
  );
};

const mapState = (state) => ({
  apiKeys: state.ListApiKeys.apiKeys || [],
  projectId: state.CommonState.projectID,
  code: state.GenerateKey.codeContent,
});

const mapDispatch = (dispatch) => ({
  handleListApiKeys: dispatch.ListApiKeys.handleListApiKeys,
  genKey: dispatch.GenerateKey.genKey,
});

export default connect(mapState, mapDispatch)(ApiKeys);
