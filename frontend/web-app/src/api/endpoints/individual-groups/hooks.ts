import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import type { CompetitorGroupDetail } from "../../../entities";
import { COMPETITORS_QUERY_KEYS } from "../../query-keys/competitors";
import { INDIVIDUAL_GROUPS_QUERY_KEYS } from "../../query-keys/individualGroups";
import { individualGroupsApi } from "./api";

export const useIndividualGroup = (groupId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: INDIVIDUAL_GROUPS_QUERY_KEYS.element(groupId),
    queryFn: async () => await individualGroupsApi.getGroup(groupId),
    enabled,
  });
};

export const useDeleteIndividualGroup = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<void, Error, number>({
    mutationFn: individualGroupsApi.deleteGroup,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: INDIVIDUAL_GROUPS_QUERY_KEYS.all,
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось удалить дивизион",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useGroupCompetitors = (groupId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: COMPETITORS_QUERY_KEYS.allByGroup(groupId),
    queryFn: async () => await individualGroupsApi.getCompetitors(groupId),
    initialData: [],
    enabled,
  });
};

export const useSyncGroupCompetitors = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<CompetitorGroupDetail[], Error, number>({
    mutationFn: individualGroupsApi.syncCompetitors,
    onSuccess: (synced, groupId) => {
      queryClient.setQueryData(COMPETITORS_QUERY_KEYS.allByGroup(groupId), synced);
      onSuccess?.();
    },
  });
};
