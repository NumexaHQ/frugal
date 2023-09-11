// authUtils.js

// Function to check if the JWT token exists and is valid
export const isAuthenticated = () => {
  const token = sessionStorage.getItem("jwtToken"); // Use sessionStorage instead of localStorage
  if (token) {
    // Decode the token and check if it's expired
    const decodedToken = decodeToken(token);
    return decodedToken.exp > Date.now() / 1000;
  }
  return false;
};

// Function to decode the JWT token
export const decodeToken = (token) => {
  return JSON.parse(atob(token.split(".")[1]));
};

export function getAuthHeader() {
  const authToken = sessionStorage.getItem("jwtToken");
  let auth;
  if (authToken) {
    auth = `Bearer ${authToken}`;
  } else {
    auth = "";
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

export const generateTimeParams = (timeFilter) => {
  if (timeFilter === "all") {
    return { from: "", to: "" }; // Return empty strings for "all" filter
  }
  const currentDate = new Date(); // Get the current date and time

  let fromDate, toDate;

  switch (timeFilter) {
    case "24h":
      // Calculate the date 24 hours ago from the current date
      fromDate = new Date(currentDate - 24 * 60 * 60 * 1000);
      toDate = currentDate;
      break;
    case "7d":
      // Calculate the date 7 days ago from the current date
      fromDate = new Date(currentDate - 7 * 24 * 60 * 60 * 1000);
      toDate = currentDate;
      break;
    case "1m":
      // Calculate the date 1 month ago from the current date
      fromDate = new Date(currentDate - 30 * 24 * 60 * 60 * 1000);
      toDate = currentDate;
      break;
    case "3m":
      // Calculate the date 3 months ago from the current date
      fromDate = new Date(currentDate - 90 * 24 * 60 * 60 * 1000);
      toDate = currentDate;
      break;
    default:
      // Default to a specific date range if an invalid filter is provided
      return { from: "", to: "" };
  }

  // Format the dates as ISO strings
  const fromParam = fromDate.toISOString();
  const toParam = toDate.toISOString();

  return { from: fromParam, to: toParam };
};
