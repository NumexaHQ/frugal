
import {
    Box,
    SimpleGrid,
    Stack
} from '@chakra-ui/react';
import MiniStatistics from "components/card/MiniStatistics";



import { connect } from 'react-redux';





const CachingPolicies = ({props}) => {

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
            value={90}
          />
          <MiniStatistics
            name='Total Cost Saved'
            value={90}
          />
            <MiniStatistics
            name='Total Time Saved'
            value={90}
            />


          </SimpleGrid>
            </Stack>
        </Box>

    )
}

const mapState = (state) => (
    {
        //Add States
    });

const mapDispatch = (dispatch) => ({
   // Add Methods
});


export default connect(mapState, mapDispatch)(CachingPolicies);