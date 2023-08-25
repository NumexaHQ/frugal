// Chakra imports
import { SimpleGrid, Text, useColorModeValue } from "@chakra-ui/react";
// Custom components
import Card from "components/card/Card.js";
import Information from "views/admin/profile/components/Information";

// Assets
export default function GeneralInformation(props) {
  const { ...rest } = props;
  // Chakra Color Mode
  const textColorPrimary = useColorModeValue("secondaryGray.900", "white");
  const textColorSecondary = "gray.400";
  const cardShadow = useColorModeValue(
    "0px 18px 40px rgba(112, 144, 176, 0.12)",
    "unset"
  );

  const jwtToken = sessionStorage.getItem("jwtToken");
  const claims = jwtToken ? JSON.parse(atob(jwtToken.split(".")[1])) : null;

  return (
    <Card mb={{ base: "0px", "2xl": "20px" }} {...rest}>
      <Text
        color={textColorPrimary}
        fontWeight='bold'
        fontSize='2xl'
        mt='10px'
        mb='4px'>
        User Information
      </Text>
      <SimpleGrid columns='2' gap='20px'>
        <Information
          boxShadow={cardShadow}
          title='Username'
          value={claims ? claims["email"] : ""}
        />
        <Information
          boxShadow={cardShadow}
          title='Organization ID' 
          value={claims ? claims["organization_id"] : ""}
        />

      </SimpleGrid>
    </Card>
  );
}
