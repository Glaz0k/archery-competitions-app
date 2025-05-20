import axios from "axios";

const apiClient = axios.create({
  baseURL: `${import.meta.env.VITE_BASE_API_URL}/api`,
  timeout: 10 * 1000,
  withCredentials: true,
});

export default apiClient;
