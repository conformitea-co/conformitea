import { useEffect } from "react";
import useSWR from "swr";

import type { User } from "@/types/auth";

import { useAuthStore } from "@/stores/auth-store";

import { ApiError, fetcher } from "@/lib/api";

export function useUser() {
  const { setUser, clearUser } = useAuthStore();

  const { data, error, isLoading, mutate } = useSWR<User>("/auth/me", fetcher, {
    revalidateOnFocus: true,
    revalidateOnReconnect: true,
    shouldRetryOnError: (error) => {
      // Don't retry on 401 errors
      if (error instanceof ApiError && error.status === 401) {
        return false;
      }

      return true;
    },
  });

  // Update Zustand store when user data changes
  useEffect(() => {
    if (data) {
      setUser(data);
    } else if (error instanceof ApiError && error.status === 401) {
      clearUser();
    }
  }, [data, error, setUser, clearUser]);

  return {
    user: data,
    isLoading,
    isError: error,
    mutate,
  };
}
