// authUtils.js

// Function to check if the JWT token exists and is valid
export const isAuthenticated = () => {
  const token = sessionStorage.getItem('jwtToken'); // Use sessionStorage instead of localStorage
  if (token) {
    // Decode the token and check if it's expired
    const decodedToken = decodeToken(token);
    return decodedToken.exp > Date.now() / 1000;
  }
  return false;
};

// Function to decode the JWT token
export const decodeToken = (token) => {
  return JSON.parse(atob(token.split('.')[1]));
};

export function getAuthHeader() {
  const authToken = sessionStorage.getItem('jwtToken');
  let auth;
  if (authToken) {
    auth = `Bearer ${authToken}`;
  } else {
    auth = '';
  }
  return auth;
}

export function formatDateTime(timestamp) {
  const dateTime = new Date(timestamp);
  const options = {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "numeric",
    minute: "numeric",
    second: "numeric",
    timeZoneName: "short",
  };
  
  return new Intl.DateTimeFormat("en-US", options).format(dateTime);
}
