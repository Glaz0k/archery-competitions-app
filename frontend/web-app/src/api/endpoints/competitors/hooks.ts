import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import type { Competitor } from "../../../entities";
import { COMPETITORS_QUERY_KEYS } from "../../query-keys/competitors";
import { competitorsApi } from "./api";
import type { CompetitorEdit } from "./types";

export const useCompetitors = (enabled: boolean = true) => {
  return useQuery({
    queryKey: COMPETITORS_QUERY_KEYS.all,
    queryFn: competitorsApi.getCompetitors,
    initialData: [],
    enabled,
  });
};

export const useCompetitor = (competitorId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: COMPETITORS_QUERY_KEYS.element(competitorId),
    queryFn: async () => await competitorsApi.getCompetitor(competitorId),
    initialData: null,
    enabled,
  });
};

export const useEditCompetitor = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<Competitor, Error, [number, CompetitorEdit]>({
    mutationFn: async ([competitorId, data]) =>
      await competitorsApi.putCompetitor(competitorId, data),
    onSuccess: (edited) => {
      queryClient.invalidateQueries({
        queryKey: COMPETITORS_QUERY_KEYS.all,
      });
      queryClient.setQueryData(COMPETITORS_QUERY_KEYS.element(edited.id), edited);
      onSuccess?.();
    },
    onError: (error, [id]) => {
      notifications.show({
        title: `Не удалось изменить информацию о пользователе, id: ${id}`,
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useDeleteCompetitor = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<void, Error, number>({
    mutationFn: competitorsApi.deleteCompetitor,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: COMPETITORS_QUERY_KEYS.all,
      });
      onSuccess?.();
    },
  });
};
