// Chakra imports
import { Flex, useColorModeValue } from "@chakra-ui/react";
import { NumexaLogo } from "components/icons/Icons";
import { HSeparator } from "components/separator/Separator";
// Custom components

export function SidebarBrand() {
  //   Chakra color mode
  let logoColor = useColorModeValue("navy.700", "white");

  return (
    <Flex align="center" direction="column">
      <NumexaLogo h="50px" w="50px" my="0px" color={logoColor} />
      <HSeparator mb="20px" />
    </Flex>
  );
}

export default SidebarBrand;
