import {
  Table,
  TableCaption,
  TableContainer,
  Tbody,
  Td,
  Tr,
} from "@chakra-ui/react";

export default function DrawerTable(props) {
  const { selectedRowData, latency } = props;

  return (
    <TableContainer>
      <Table variant="striped" colorScheme="teal">
        <TableCaption></TableCaption>
        <Tbody>
          <Tr>
            <Td>RequestID</Td>
            <Td>{selectedRowData ? selectedRowData.id : ""}</Td>
          </Tr>
          <Tr>
            <Td>Prompt</Td>
            <Td>{selectedRowData ? selectedRowData.prompt : ""}</Td>
          </Tr>
          <Tr>
            <Td>CreatedAT</Td>
            <Td>{selectedRowData ? selectedRowData.initiated_at : ""}</Td>
          </Tr>
          <Tr>
            <Td>Latency</Td>
            <Td>{selectedRowData ? latency : ""}</Td>
          </Tr>
        </Tbody>
      </Table>
    </TableContainer>
  );
}
