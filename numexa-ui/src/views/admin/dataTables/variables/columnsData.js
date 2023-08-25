import { Badge } from "@chakra-ui/react";

export const columnsDataDevelopment = [
  {
    Header: "NAME",
    accessor: "name",
  },
  {
    Header: "TECH",
    accessor: "tech",
  },
  {
    Header: "DATE",
    accessor: "date",
  },
  {
    Header: "PROGRESS",
    accessor: "progress",
  },
];

export const columnsDataCheck = [
  {
    Header: "NAME",
    accessor: "name",
  },
  {
    Header: "PROGRESS",
    accessor: "progress",
  },
  {
    Header: "QUANTITY",
    accessor: "quantity",
  },
  {
    Header: "DATE",
    accessor: "date",
  },
  {
    Header: "COST",
    accessor: "cost",
  },
];

export const columnsDataColumns = [
  {
    Header: "NAME",
    accessor: "name",
  },
  {
    Header: "PROGRESS",
    accessor: "progress",
  },
  {
    Header: "QUANTITY",
    accessor: "quantity",
  },
  {
    Header: "DATE",
    accessor: "date",
  },
  {
    Header: "COST",
    accessor: "cost",
  },
  {
    Header: "STATUS",
    accessor: "status",
  },
];

export const columnsDataComplex = [
  {
    Header: "NAME",
    accessor: "name",
  },
  {
    Header: "STATUS",
    accessor: "status",
  },
  {
    Header: "DATE",
    accessor: "date",
  },
  {
    Header: "PROGRESS",
    accessor: "progress",
  },
  {
    Header: "COST",
    accessor: "cost",
  },
];

export const requestDataColumn = [
  {
    Header: "ID",
    accessor: (row) => {
      // add middle ellipsis
      return `${row.id.slice(0, 4)}...${row.id.slice(-4)}`;
    },
  },
  {
    Header: "Initiated At",
    accessor: "initiated_at",
  },
  {
    Header: "Status code",
    accessor: (row) => {
      if (row.status_code === 0) {
        return <Badge colorScheme="yellow">{"Invalid"}</Badge>;
      }
      if (row.status_code === 200) {
        return (
          // need to make them pills
          <Badge colorScheme="green">
            {row.status_code} {row.status}
          </Badge>
        );
      } else {
        return (
          <Badge colorScheme="red">
            {row.status_code} {row.status}
          </Badge>
        );
      }
    },
  },
  {
    Header: "Model",
    accessor: "model",
  },
  {
    Header: "Prompt",
    accessor: "prompt",
  },
  {
    Header: "Cost",
    accessor: (row) => {
      if (row.cost !== 0) {
        // trim the cost to 6 decimal places
        return `$ ${row.cost.toFixed(6)}`;
      }
      return `$ ${row.cost}`;
    },
  },
  {
    Header: "Provider",
    accessor: "provider",
  },
];

export const apiKeycolumn = [
  {
    Header: "Name",
    accessor: "name",
  },
  {
    Header: "CREATED AT",
    accessor: "created_at",
  },
];