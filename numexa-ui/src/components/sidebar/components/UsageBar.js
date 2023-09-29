const { Stack, Progress, Box, Text, Badge, Flex } = require("@chakra-ui/react");

export default function UsageBar({ usage }) {
  let totalUsage = usage?.usage;
  let limit = usage?.limit;

  let usagePercentage = (totalUsage / limit) * 100;
  return (
    <Stack spacing={8} direction="column" mb="10px">
      <Box fontSize="sm" fontWeight="bold" textAlign="center" mt="0">
        {usage && (
          <Progress
            colorScheme="green"
            size="md"
            value={usagePercentage}
            width="100%"
          />
        )}
        {!usage && (
          <Progress colorScheme="green" size="md" value={0} width="100%" />
        )}
        {totalUsage}/{limit}
        <Text as="span" fontWeight="normal">
          &nbsp;requests used
        </Text>
        <Box>
          <Badge ml="10px" colorScheme="green">
            {usage?.plan}
          </Badge>
        </Box>
      </Box>
    </Stack>
  );
}
