import axios from "axios";
import { getAuthHeader } from "utils/utils";

const BASE_URL = window.location.hostname;
const PROTOCOL = window.location.protocol;

const AUTH_BASE_URL = "http://localhost:8080";

const VIBE_BASE_URL = "http://localhost:8082";

// use base url for production else fallback to localhost
// const AUTH_BASE_URL =
//   process.env.REACT_APP_AUTH_BASE_URL ||
//   `${PROTOCOL}://${BASE_URL}:8080` ||
//   "http://localhost:8080";
// const VIBE_BASE_URL =
//   process.env.REACT_APP_VIBE_BASE_URL ||
//   `${PROTOCOL}://${BASE_URL}:8082` ||
//   "http://localhost:8082";

// const AUTH_BASE_URL = `${PROTOCOL}//${BASE_URL}/auth-service`;
// const VIBE_BASE_URL = `${PROTOCOL}//${BASE_URL}/api`;

export const CommonState = {
  state: {
    projectID: "",
  },
  reducers: {
    setProjectID(state, projectID) {
      return { ...state, projectID };
    },
  },
};

export const Login = {
  state: {},
  reducers: {
    ...CommonState.reducers,
    setSessionToken(state, sessionToken) {
      sessionStorage.setItem("jwtToken", sessionToken);
      return null;
    },
  },
  effects: (dispatch) => ({
    async handleSignIn(payload, state) {
      try {
        const credentials = {
          email: payload.username,
          password: payload.password,
        };
        console.log("base url", AUTH_BASE_URL);
        console.log("protocol", PROTOCOL);
        const response = await axios.post(
          `${AUTH_BASE_URL}/login`,
          credentials
        );
        const token = response.data.token;
        dispatch.Login.setSessionToken(token);
        dispatch.CommonState.setProjectID(response.data.project_id);
      } catch (error) {
        console.log("Sign In Error", error);
      } finally {
      }
    },
  }),
};

export const Register = {
  state: {},
  reducers: {
    ...CommonState.reducers,
    setSessionToken(state, sessionToken) {
      sessionStorage.setItem("jwtToken", sessionToken);
      return null;
    },
  },
  effects: (dispatch) => ({
    async handleRegister(payload, state) {
      console.log(payload);
      try {
        const credentials = {
          email: payload.username,
          password: payload.password,
        };
        const response = await axios.post(
          `${AUTH_BASE_URL}/register`,
          credentials
        );
        const token = response.data.token;
        dispatch.Register.setSessionToken(token);
        dispatch.CommonState.setProjectID(response.data.project_id);
      } catch (error) {
        console.log("Sign In Error", error);
      } finally {
      }
    },
  }),
};

export const ListApiKeys = {
  state: {
    apiKeys: [],
  },
  reducers: {
    setApiKeys(state, apiKeys) {
      return { ...state, apiKeys };
    },
  },
  effects: (dispatch) => ({
    async handleListApiKeys(payload, state) {
      var requestOptions = {
        method: "GET",
        headers: {
          Authorization: getAuthHeader(),
          "Content-Type": "application/json",
        },
        redirect: "follow",
      };
      try {
        const response = await fetch(
          `${AUTH_BASE_URL}/get_api_key`,
          requestOptions
        );
        const data = await response.json();

        dispatch.ListApiKeys.setApiKeys(data.keys);
      } catch (error) {
        console.log("List Api Keys Error", error);
      }
    },
  }),
};

export const ListRequests = {
  state: {
    requests: [],
  },
  reducers: {
    setRequests(state, requests) {
      return { ...state, requests };
    },
  },
  effects: (dispatch) => ({
    async getProviderRequests(payload, state) {
      var requestOptions = {
        method: "GET",
        headers: {
          Authorization: getAuthHeader(),
          "Content-Type": "application/json",
        },
        redirect: "follow",
      };
      try {
        const response = await fetch(
          `${VIBE_BASE_URL}/mng_request/${payload.projectId}`,
          requestOptions
        );
        const data = await response.json();

        dispatch.ListRequests.setRequests(data);
      } catch (error) {
        console.log("List Requests Error", error);
      }
    },
  }),
};

export const ListResponse = {
  state: {
    response: [],
  },
  reducers: {
    setResponse(state, response) {
      return { ...state, response };
    },
  },
  effects: (dispatch) => ({
    async getResponse(payload, state) {
      var requestOptions = {
        method: "GET",
        headers: {
          Authorization: getAuthHeader(),
          "Content-Type": "application/json",
        },
        redirect: "follow",
      };
      try {
        const response = await fetch(
          `${VIBE_BASE_URL}/mng_response/${payload.requestId}`,
          requestOptions
        );
        const data = await response.json();

        dispatch.ListResponse.setResponse(data);
      } catch (error) {
        console.log("List Requests Error", error);
      }
    },
  }),
};

export const TotalRequests = {
  state: {
    totalRequest: 0,
  },
  reducers: {
    setTotalRequest(state, totalRequest) {
      return { ...state, totalRequest };
    },
  },
  effects: (dispatch) => ({
    async getTotalRequest(payload, state) {
      var requestOptions = {
        method: "GET",
        headers: {
          Authorization: getAuthHeader(),
          "Content-Type": "application/json",
        },
        redirect: "follow",
      };
      try {
        const response = await fetch(
          `${VIBE_BASE_URL}/total_requests/${payload.projectId}`,
          requestOptions
        );
        const data = await response.json();
        dispatch.TotalRequests.setTotalRequest(data.total_requests);
      } catch (error) {
        console.log("List Requests Error", error);
      }
    },
  }),
};

export const AvgLatency = {
  state: {
    avgLatency: 0,
  },
  reducers: {
    setAvgLatency(state, avgLatency) {
      return { ...state, avgLatency };
    },
  },
  effects: (dispatch) => ({
    async getAvgLatency(payload, state) {
      var requestOptions = {
        method: "GET",
        headers: {
          Authorization: getAuthHeader(),
          "Content-Type": "application/json",
        },
        redirect: "follow",
      };
      try {
        const response = await fetch(
          `${VIBE_BASE_URL}/avg_latency/${payload.projectId}`,
          requestOptions
        );
        const data = await response.json();
        dispatch.AvgLatency.setAvgLatency(data.avg_latency);
      } catch (error) {
        console.log("List Requests Error", error);
      }
    },
  }),
};

export const AvgTokens = {
  state: {
    avgTokens: [],
  },
  reducers: {
    setAvgTokens(state, avgTokens) {
      return { ...state, avgTokens };
    },
  },
  effects: (dispatch) => ({
    async getAvgTokens(payload, state) {
      var requestOptions = {
        method: "GET",
        headers: {
          Authorization: getAuthHeader(),
          "Content-Type": "application/json",
        },
        redirect: "follow",
      };
      try {
        const response = await fetch(
          `${VIBE_BASE_URL}/avg_prompt_tokens/${payload.projectId}`,
          requestOptions
        );
        const data = await response.json();
        dispatch.AvgTokens.setAvgTokens([data]);
      } catch (error) {
        console.log("List Requests Error", error);
      }
    },
  }),
};

export const ModelDistribution = {
  state: {
    modelDistribution: [],
  },
  reducers: {
    setModelDistribution(state, modelDistribution) {
      return { ...state, modelDistribution };
    },
  },
  effects: (dispatch) => ({
    async getModelDistribution(payload, state) {
      var requestOptions = {
        method: "GET",
        headers: {
          Authorization: getAuthHeader(),
          "Content-Type": "application/json",
        },
        redirect: "follow",
      };
      try {
        const response = await fetch(
          `${VIBE_BASE_URL}/unique_models/${payload.projectId}`,
          requestOptions
        );
        const data = await response.json();
        dispatch.ModelDistribution.setModelDistribution(data);
      } catch (error) {
        console.log("List Requests Error", error);
      }
    },
  }),
};

export const UsersUsageStat = {
  state: {
    usersUsageStat: [],
  },
  reducers: {
    setUsersUsageStat(state, usersUsageStat) {
      return { ...state, usersUsageStat };
    },
  },
  effects: (dispatch) => ({
    async getUsersUsageStat(payload, state) {
      var requestOptions = {
        method: "GET",
        headers: {
          Authorization: getAuthHeader(),
          "Content-Type": "application/json",
        },
        redirect: "follow",
      };
      try {
        const response = await fetch(
          `${VIBE_BASE_URL}/user_requests_stats/${payload.projectId}`,
          requestOptions
        );
        const data = await response.json();
        dispatch.UsersUsageStat.setUsersUsageStat(data);
      } catch (error) {
        console.log("List Requests Error", error);
      }
    },
  }),
};
export const GenerateKey = {
  state: {
    codeContent: "",
  },
  reducers: {
    setCodeContent(state, codeContent) {
      return { ...state, codeContent };
    },
  },
  effects: (dispatch) => ({
    async genKey(payload, state) {
      var requestOptions = {
        method: "POST",
        headers: {
          Authorization: getAuthHeader(),
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          project_id: payload.projectId,
          name: payload.keyName,
        }),
        redirect: "follow",
      };

      fetch(`${AUTH_BASE_URL}/generate_api_key`, requestOptions)
        .then((response) => response.json())
        .then((data) => dispatch.GenerateKey.setCodeContent(data.key));
    },
  }),
};