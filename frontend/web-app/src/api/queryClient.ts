import { MutationCache, QueryCache, QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error) => {
      console.error(error);
    },
  }),
  mutationCache: new MutationCache({
    onError: (error) => {
      console.error(error);
    },
  }),
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 2,
      networkMode: "online",
    },
    mutations: {
      networkMode: "online",
    },
  },
});
