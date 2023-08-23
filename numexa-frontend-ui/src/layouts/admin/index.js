import { Box, Portal, useDisclosure } from "@chakra-ui/react";
import Footer from "components/footer/FooterAdmin.js";
import Navbar from "components/navbar/NavbarAdmin.js";
import Sidebar from "components/sidebar/Sidebar.js";
import { SidebarContext } from "contexts/SidebarContext";
import { useState } from "react";
import { Navigate, Route, Routes } from "react-router-dom"; // Import from react-router-dom v6
import { routes } from "routes.js";

export default function Dashboard(props) {
  const { ...rest } = props;
  const getRoute = () => {
    return window.location.pathname !== "/admin/full-screen-maps";
  };
  const [fixed] = useState(false);
  const [toggleSidebar, setToggleSidebar] = useState(false);

  const getActiveRoute = (routes) => {
    let activeRoute = "";
    for (let i = 0; i < routes.length; i++) {
      if (routes[i].items) {
        let nestedActiveRoute = getActiveRoute(routes[i].items);
        if (nestedActiveRoute !== activeRoute) {
          return nestedActiveRoute;
        }
      } else {
        if (window.location.pathname === routes[i].path) {
          return routes[i].name;
        }
      }
    }
    return activeRoute;
  };

  const getActiveNavbar = (routes) => {
    let activeNavbar = false;
    for (let i = 0; i < routes.length; i++) {
      if (routes[i].items) {
        let nestedActiveNavbar = getActiveNavbar(routes[i].items);
        if (nestedActiveNavbar) {
          return nestedActiveNavbar;
        }
      } else {
        if (window.location.pathname === routes[i].path) {
          return routes[i].secondary;
        }
      }
    }
    return activeNavbar;
  };

  const getActiveNavbarText = (routes) => {
    let activeNavbarText = "";
    for (let i = 0; i < routes.length; i++) {
      if (routes[i].items) {
        let nestedActiveNavbarText = getActiveNavbarText(routes[i].items);
        if (nestedActiveNavbarText) {
          return nestedActiveNavbarText;
        }
      } else {
        if (window.location.pathname === routes[i].path) {
          return routes[i].messageNavbar;
        }
      }
    }
    return activeNavbarText;
  };

  const getRoutes = (routes) => {
    return routes.map((prop, key) => {
      if (prop.layout === "/admin") {
        return (
          <Route path={prop.path} element={<prop.component />} key={key} />
        );
      }
      if (prop.items) {
        return getRoutes(prop.items);
      } else {
        return null;
      }
    });
  };

  const { onOpen } = useDisclosure();

  return (
    <Box>
      <Box>
        <SidebarContext.Provider
          value={{
            toggleSidebar,
            setToggleSidebar,
          }}
        >
          <Sidebar routes={routes} display="none" {...rest} />
          <Box
            float="right"
            minHeight="100vh"
            height="100%"
            overflow="auto"
            position="relative"
            maxHeight="100%"
            w={{ base: "100%", xl: "calc( 100% - 290px )" }}
            maxWidth={{ base: "100%", xl: "calc( 100% - 290px )" }}
            transition="all 0.33s cubic-bezier(0.685, 0.0473, 0.346, 1)"
            transitionDuration=".2s, .2s, .35s"
            transitionProperty="top, bottom, width"
            transitionTimingFunction="linear, linear, ease"
          >
            <Portal>
              <Box>
                <Navbar
                  onOpen={onOpen}
                  logoText={"Numexa UI Dashboard"}
                  brandText={getActiveRoute(routes)}
                  secondary={getActiveNavbar(routes)}
                  message={getActiveNavbarText(routes)}
                  fixed={fixed}
                  {...rest}
                />
              </Box>
            </Portal>
            {getRoute() ? (
              <Box
                mx="auto"
                p={{ base: "20px", md: "30px" }}
                pe="20px"
                minH="100vh"
                pt="50px"
              >
                <Routes>
                  {getRoutes(routes)}
                  <Route
                    path="*"
                    element={<Navigate to="/admin/dashboard" />}
                  />
                </Routes>
              </Box>
            ) : null}
            <Box>
              <Footer />
            </Box>
          </Box>
        </SidebarContext.Provider>
      </Box>
    </Box>
  );
}
