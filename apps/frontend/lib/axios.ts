import axios from "axios";
import { auth } from "./firebase";
import { toast } from "sonner";

const axiosInstance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
});

axiosInstance.interceptors.request.use(
  async (config) => {
    const user = auth.currentUser;
    if (user) {
      const token = await user.getIdToken();
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.error("Axios error:", error);

    const errorMessage =
      error.response?.data?.error || "An unexpected error occurred.";

    toast.error(errorMessage);

    return Promise.reject(error);
  }
);

export default axiosInstance;
