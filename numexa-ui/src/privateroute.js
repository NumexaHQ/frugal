import { Navigate } from 'react-router-dom'; // Import from react-router-dom v6
import { isAuthenticated } from 'utils/utils';

const PrivateRoute = ({ children, redirectTo }) => {
  var isAuth = isAuthenticated();
  return isAuth ? children : <Navigate to={redirectTo} />;
}

export default PrivateRoute;
