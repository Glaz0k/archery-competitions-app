import axios from "axios";

const apiClient = axios.create({
  baseURL: "/api/",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
    Authorization:
      "Bearer eyJhbGciOiJIUzI1NiJ9.eyJyb2xlIjoiYWRtaW4ifQ.hkRxKoem3fjgNx0h4WFtfc6IhSP09pFPNqi50C1zwak", // TEST ONLY
  },
});

export default apiClient;
