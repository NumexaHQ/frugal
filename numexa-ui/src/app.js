import { ChakraProvider } from "@chakra-ui/react";
import { ThemeEditorProvider } from "@hypertheme-editor/chakra-ui";
import "assets/css/App.css";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom"; // Import from react-router-dom v6
import theme from "theme/theme";
import SignIn from "views/auth/signIn";

import { GoogleOAuthProvider } from "@react-oauth/google";
import AdminLayouts from "layouts/admin";
import PrivateRoute from "privateroute";

const App = () => {
  return (
    <GoogleOAuthProvider clientId="656339371007-s0lcm9p57uurpil4prha4aqapjmafq5h.apps.googleusercontent.com">
      <ChakraProvider theme={theme}>
        <ThemeEditorProvider>
          <BrowserRouter>
            <Routes>
              <Route path="/auth" element={<SignIn />} />
              <Route
                path="/admin/*"
                element={
                  <PrivateRoute redirectTo="/auth">
                    <AdminLayouts />
                  </PrivateRoute>
                }
              />
              <Route path="/test/*" element={<Navigate to="/auth" />} />
              <Route path="/" element={<Navigate to="/auth" />} />
            </Routes>
          </BrowserRouter>
        </ThemeEditorProvider>
      </ChakraProvider>
    </GoogleOAuthProvider>
  );
};

export default App;
