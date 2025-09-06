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
    logout: async () => {
      const response = await fetch(`${API_URL}/auth/logout`, {
        method: "POST",
        credentials: "include",
      });

      if (!response.ok) {
        throw new ApiError(response.status, response.statusText);
      }
    },
    createSession: async (code: string, state?: string) => {
      const response = await fetch(`${API_URL}/auth/session`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ code, state }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new ApiError(response.status, error.message || response.statusText);
      }

      return response.json();
    },
  },
};
