export const SERVER_URL: string = process.env.NEXT_PUBLIC_SERVER_URL || "";

export const API_URL = {
  root: (url = "") => `/api${url}`,
};
