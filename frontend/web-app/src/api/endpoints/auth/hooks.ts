import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AUTH_QUERY_KEYS } from "../../query-keys/auth";
import { authApi } from "./api";
import type { AuthData, Credentials } from "./types";

const AUTH_DATA_KEY = "auth_data";

export const useUser = (enabled: boolean = true) => {
  const queryClient = useQueryClient();
  return useQuery({
    queryKey: AUTH_QUERY_KEYS.user,
    queryFn: async (): Promise<AuthData | null> => {
      const cachedData = queryClient.getQueryData(AUTH_QUERY_KEYS.user);
      if (cachedData) return cachedData as AuthData;

      const storedData = localStorage.getItem(AUTH_DATA_KEY);
      if (storedData) {
        const authData = JSON.parse(storedData) as AuthData;
        queryClient.setQueryData(AUTH_QUERY_KEYS.user, authData);
        return authData;
      }

      return null;
    },
    enabled,
  });
};

export const useAdminSignIn = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<AuthData, Error, Credentials>({
    mutationFn: authApi.adminSignIn,
    onSuccess: (authData) => {
      queryClient.setQueryData(AUTH_QUERY_KEYS.user, authData);
      localStorage.setItem(AUTH_DATA_KEY, JSON.stringify(authData));
      onSuccess?.();
    },
  });
};

export const useLogout = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<void, Error, void>({
    mutationFn: authApi.logout,
    onSuccess: () => {
      queryClient.setQueryData(AUTH_QUERY_KEYS.user, null);
      localStorage.removeItem(AUTH_DATA_KEY);
      onSuccess?.();
    },
  });
};
