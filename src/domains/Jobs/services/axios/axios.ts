import type { AxiosInstance, AxiosResponse } from "axios";
import axios from "axios";

const axiosInstance: AxiosInstance = axios.create({
  baseURL: "http://206.189.48.157", // Replace with your backend URL
  timeout: 10000, // Timeout for requests (10 seconds)
});

// Add response interceptor for global error handling
axiosInstance.interceptors.response.use(
  (response: AxiosResponse) => {
    return response; // Return the response directly if successful
  },
  (error) => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          console.error("Unauthorized access");
          break;
        case 403:
          console.error("Access forbidden");
          break;
        default:
          console.error(error);
      }
    }
    return Promise.reject(error);
  },
);

export default axiosInstance;
