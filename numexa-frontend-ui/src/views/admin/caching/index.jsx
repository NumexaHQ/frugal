
import {
    Box,
    Stack
} from '@chakra-ui/react';




import { connect } from 'react-redux';

import Banner from './components/banner';




const CachingPolicies = ({props}) => {

    return (
        <Box pt={{ base: "130px", md: "80px", xl: "80px" }}>
            <Stack spacing={3}>
                <Banner />
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