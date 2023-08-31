
import {
  Box,
  Icon,
  SimpleGrid,
  Stack,
  useColorModeValue
} from '@chakra-ui/react';
import MiniStatistics from "components/card/MiniStatistics";


import { MdDataUsage, MdScore } from "react-icons/md";


import IconBox from "components/icons/IconBox";
import { useEffect } from "react";
import { connect } from 'react-redux';
import RequestTable from "views/admin/dataTables/components/new-request-table";
import {
  requestDataColumn
} from "views/admin/dataTables/variables/columnsData";


const CachingPolicies = ({ projectId, requests, getRequests}) => {

  const boxBg = useColorModeValue("secondaryGray.300", "whiteAlpha.100");
  const brandColor = useColorModeValue("brand.500", "white");
  useEffect(() => {
    getRequests({ projectId });
  }, []);
  const cachedRequests = requests.filter((request) => request.is_cache_hit === true);

  const TotalCostSaved = cachedRequests.reduce((acc, request) => {
    // return up to 6 decimal places
    return acc + parseFloat(request.cost.toFixed(6));
  }, 0);

    return (
        <Box pt={{ base: "130px", md: "80px", xl: "80px" }}>
            <Stack spacing={3}>
                {/* <Banner /> */}
                <SimpleGrid
          columns={{ base: 3, md: 3, lg: 3, "2xl": 3 }}
          gap='20px'
          mb='20px'>    
          <MiniStatistics
            name='Cache Hits'
            value={cachedRequests.length}
            startContent={
              <IconBox
                w='56px'
                h='56px'
                bg={boxBg}
                icon={
                  <Icon w='32px' h='32px' as={MdDataUsage} color={brandColor} />
                }
              />
            }
          />
          <MiniStatistics
            name='Total Cost Saved'
            value={TotalCostSaved.toFixed(6)}
            startContent={
              <IconBox
                w='56px'
                h='56px'
                bg={boxBg}
                icon={
                  <Icon w='32px' h='32px' as={MdScore} color={brandColor}/>
                }
              />
            }
            
          />
          </SimpleGrid>
            </Stack>
            <RequestTable
        columnsData={requestDataColumn}
        tableData={cachedRequests}
        title={"Cache Hit Entries"}
      />
        </Box>
    )
}

const mapState = (state) => (
    {
      requests: state.ListRequests.requests || [],
      projectId: state.CommonState.projectID,
    });

const mapDispatch = (dispatch) => ({
 getRequests: dispatch.ListRequests.getProviderRequests,
});


export default connect(mapState, mapDispatch)(CachingPolicies);