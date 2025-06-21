import type { User } from "@/types/auth";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

export const fetcher = async (url: string) => {
  const response = await fetch(`${API_URL}${url}`, {
    credentials: "include", // Important: include cookies
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!response.ok) {
    throw new ApiError(response.status, response.statusText);
  }

  return response.json();
};

// Type-safe API methods
export const api = {
  auth: {
    me: () => fetcher("/auth/me") as Promise<User>,
    logout: async () => {
      const response = await fetch(`${API_URL}/auth/logout`, {
        method: "POST",
        credentials: "include",
      });

      if (!response.ok) {
        throw new ApiError(response.status, response.statusText);
      }
    },
  },
};
