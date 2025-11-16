import axios, { AxiosError } from "axios";

export const getContentType = () =>
  ({
    "Content-Type": "application/json",
  } as const);

export const errorCatch = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<{ message?: string | string[] }>;
    const message = axiosError.response?.data?.message;

    if (message) {
      return Array.isArray(message) ? String(message[0]) : String(message);
    }
    return axiosError.message || "An unexpected error occurred";
  }

  if (error instanceof Error) return error.message;
  if (typeof error === "string") return error;

  return "An unexpected error occurred";
};

export const isNetworkError = (error: unknown): boolean => {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError;
    return !axiosError.response && !!axiosError.request;
  }
  return false;
};

export const isServerError = (error: unknown): boolean => {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError;
    const status = axiosError.response?.status;
    return typeof status === "number" ? status >= 500 && status < 600 : false;
  }
  return false;
};

export const isClientError = (error: unknown): boolean => {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError;
    const status = axiosError.response?.status;
    return typeof status === "number" ? status >= 400 && status < 500 : false;
  }
  return false;
};
