import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import { type CompetitorGroupDetail, type Qualification } from "../../../entities";
import { COMPETITORS_QUERY_KEYS, INDIVIDUAL_GROUPS_QUERY_KEYS } from "../../query-keys";
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
    onError: (error) => {
      notifications.show({
        title: "Не удалось синхронизировать участников",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useQualification = (groupId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: INDIVIDUAL_GROUPS_QUERY_KEYS.qualification(groupId),
    queryFn: async () => await individualGroupsApi.getQualification(groupId),
    enabled,
  });
};

export const useStartQualification = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<Qualification, Error, number>({
    mutationFn: individualGroupsApi.startQualification,
    onSuccess: (started) => {
      queryClient.invalidateQueries({
        queryKey: INDIVIDUAL_GROUPS_QUERY_KEYS.element(started.groupId),
      });
      queryClient.setQueryData(
        INDIVIDUAL_GROUPS_QUERY_KEYS.qualification(started.groupId),
        started
      );
      notifications.show({
        title: "Квалификация успешно началась",
        message: "У всех участников стала активна первая серия первого раунда",
        color: "green",
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось начать квалификацию",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useEndQualification = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<Qualification, Error, number>({
    mutationFn: individualGroupsApi.endQualification,
    onSuccess: (ended) => {
      queryClient.invalidateQueries({
        queryKey: INDIVIDUAL_GROUPS_QUERY_KEYS.element(ended.groupId),
      });
      queryClient.setQueryData(INDIVIDUAL_GROUPS_QUERY_KEYS.qualification(ended.groupId), ended);
      notifications.show({
        title: "Квалификация успешно завершилась",
        message: "Для продолжения необходимо начать четвертьфинал",
        color: "green",
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось завершить квалификацию",
        message: error.message,
        color: "red",
      });
    },
  });
};
