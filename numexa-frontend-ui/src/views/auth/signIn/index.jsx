
import React, { useEffect, useState } from "react";
import { useNavigate } from 'react-router-dom';

import { connect } from 'react-redux';

// Chakra imports
import {
  Box,
  Button,
  Flex,
  FormControl,
  FormLabel,
  Heading,
  Icon,
  Input,
  InputGroup,
  InputRightElement,
  Text,
  useColorModeValue
} from "@chakra-ui/react";
// Custom components
import { HSeparator } from "components/separator/Separator";
import DefaultAuth from "layouts/auth/Default";
// Assets
import illustration from "assets/img/auth/auth.png";
import { MdOutlineRemoveRedEye } from "react-icons/md";
import { RiEyeCloseLine } from "react-icons/ri";
import Register from 'views/auth/register/index';
const SignIn = ({ history, handleSignIn }) => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [password, setPassword] = useState("");
  const [username, setUsername] = useState("");
  const [isRegisterview, setIsRegisterview] = useState(false);
 

  const SignIn = async () => {
    setIsLoading(true);

    try {
      await handleSignIn({ username, password }); // Assuming handleSignIn handles successful login
      setIsLoading(false);
      navigate('/admin'); // Navigate to '/admin' route after successful sign-in
    } catch (error) {
      console.log("Sign In Error", error);
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (sessionStorage.getItem('jwtToken')) {
      navigate('/admin');
    }
  }, []);


  const textColor = useColorModeValue("navy.700", "white");
  const textColorSecondary = "gray.400";
  const textColorDetails = useColorModeValue("navy.700", "secondaryGray.600");
  const textColorBrand = useColorModeValue("brand.500", "white");
  const brandStars = useColorModeValue("brand.500", "brand.400");
  const googleBg = useColorModeValue("gray.200", "whiteAlpha.200");
  const googleText = useColorModeValue("navy.700", "white");
  const googleHover = useColorModeValue(
    { bg: "secondaryGray.300" },
    { bg: "whiteAlpha.300" }
  );
  const googleActive = useColorModeValue(
    { bg: "secondaryGray.300" },
    { bg: "whiteAlpha.200" }
  );
  const [show, setShow] = React.useState(false);
  const handleClick = () => setShow(!show);


  if (isRegisterview) {
    return <Register setIsRegisterview={setIsRegisterview} history={history} />
  }
  return (
    <DefaultAuth illustrationBackground={illustration} image={illustration}>
      <Flex
        maxW={{ base: "100%", md: "max-content" }}
        w='100%'
        mx={{ base: "auto", lg: "0px" }}
        me='auto'
        h='100%'
        alignItems='start'
        justifyContent='center'
        mb={{ base: "30px", md: "60px" }}
        px={{ base: "25px", md: "0px" }}
        mt={{ base: "40px", md: "14vh" }}
        flexDirection='column'>
        <Box me='auto'>
  <Heading color={textColor} fontSize='36px' mb='10px'>
    Sign In
  </Heading>
  <Flex
    flexDirection='column'
    justifyContent='center'
    alignItems='start'
    maxW='100%'
    mt='0px'
    mb='10px'>
    <Text color={textColorDetails} fontWeight='400' fontSize='14px'>
      Don't have an account, Sign Up? Click
      <Text
        color={textColorBrand}
        as='span'
        ms='5px'
        fontWeight='500'
        style={{ textDecoration: 'none', transition: 'text-decoration 0.3s ease' }}
        onMouseEnter={event => event.target.style.textDecoration = 'underline'}
        onMouseLeave={event => event.target.style.textDecoration = 'none'}
        onClick={() => setIsRegisterview(true)}>
        here
      </Text>
    </Text>
  </Flex>
</Box>
        <Flex
          zIndex='2'
          direction='column'
          w={{ base: "100%", md: "420px" }}
          maxW='100%'
          background='transparent'
          borderRadius='15px'
          mx={{ base: "auto", lg: "unset" }}
          me='auto'
          mb={{ base: "20px", md: "auto" }}>
          {/* <Button
            fontSize='sm'
            me='0px'
            mb='26px'
            py='15px'
            h='50px'
            borderRadius='16px'
            bg={googleBg}
            color={googleText}
            fontWeight='500'
            _hover={googleHover}
            _active={googleActive}
            _focus={googleActive}
            onClick={() => login()}
          >
            <Icon as={FcGoogle} w='20px' h='20px' me='10px' />
            Google

          </Button> */}
          <HSeparator mb={5} />
          {/* <Flex align='center' mb='25px'>
            <HSeparator />
            <Text color='gray.400' mx='14px'>
              or
            </Text>
            <HSeparator />
          </Flex> */}
          <FormControl>
            <FormLabel
              display='flex'
              ms='4px'
              fontSize='sm'
              fontWeight='500'
              color={textColor}
              mb='8px'>
              Email<Text color={brandStars}>*</Text>
            </FormLabel>
            <Input
              isRequired={true}
              variant='auth'
              fontSize='sm'
              ms={{ base: "0px", md: "0px" }}
              type='email'
              placeholder='jane@doe.io'
              mb='24px'
              fontWeight='500'
              size='lg'
              onChange={(e) => setUsername(e.target.value)}
            />
            <FormLabel
              ms='4px'
              fontSize='sm'
              fontWeight='500'
              color={textColor}
              display='flex'>
              Password<Text color={brandStars}>*</Text>
            </FormLabel>
            <InputGroup size='md'>
              <Input
                isRequired={true}
                fontSize='sm'
                placeholder='Min. 8 characters'
                mb='24px'
                size='lg'
                type={show ? "text" : "password"}
                variant='auth'
                onChange={(e) => setPassword(e.target.value)}
              />
              <InputRightElement display='flex' alignItems='center' mt='4px'>
                <Icon
                  color={textColorSecondary}
                  _hover={{ cursor: "pointer" }}
                  as={show ? RiEyeCloseLine : MdOutlineRemoveRedEye}
                  onClick={handleClick}
                />
              </InputRightElement>
            </InputGroup>
            {/* <Flex justifyContent='space-between' align='center' mb='24px'>
              <FormControl display='flex' alignItems='center'>
                <Checkbox
                  id='remember-login'
                  colorScheme='brandScheme'
                  me='10px'
                />
                <FormLabel
                  htmlFor='remember-login'
                  mb='0'
                  fontWeight='normal'
                  color={textColor}
                  fontSize='sm'>
                  Keep me logged in
                </FormLabel>
              </FormControl>
              <NavLink to='/auth/forgot-password'>
                <Text
                  color={textColorBrand}
                  fontSize='sm'
                  w='124px'
                  fontWeight='500'>
                  Forgot password?
                </Text>
              </NavLink>
            </Flex> */}
            <Button
              fontSize='sm'
              variant='brand'
              fontWeight='500'
              w='100%'
              h='50'
              mb='24px'
              onClick={SignIn}
              isLoading={isLoading}>
              Sign In
            </Button>
          </FormControl>
        </Flex>
      </Flex>
    </DefaultAuth>
  );
}

const mapState = (state) => ({

});

const mapDispatch = (dispatch) => ({
  handleSignIn: dispatch.Login.handleSignIn,
});

export default connect(mapState, mapDispatch)(SignIn);
