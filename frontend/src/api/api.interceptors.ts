import type { CreateAxiosDefaults } from "axios";
import axios from "axios";
import { errorCatch, getContentType } from "./api.helper";
import { SERVER_URL } from "@/config/api.config";

const options: CreateAxiosDefaults = {
  baseURL: SERVER_URL,
  headers: getContentType(),
  withCredentials: true,
};

const axiosClassic = axios.create(options);

axiosClassic.interceptors.response.use(
  (response) => response,
  (error) => {
    const message = errorCatch(error);
    return Promise.reject(new Error(message));
  }
);

export { axiosClassic };
