import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import { type CompetitorGroupDetail, type FinalGrid, type Qualification } from "../../../entities";
import { COMPETITORS_QUERY_KEYS, INDIVIDUAL_GROUPS_QUERY_KEYS } from "../../query-keys";
import { PLACES_QUERY_KEYS } from "../../query-keys/sparringPlaces";
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
        message: "Доступен экспорт в PDF",
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

export const useFinalGrid = (groupId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: INDIVIDUAL_GROUPS_QUERY_KEYS.finalGrid(groupId),
    queryFn: async () => await individualGroupsApi.getFinalGrid(groupId),
    enabled,
  });
};

export const useStartQuarterfinal = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<FinalGrid, Error, number>({
    mutationFn: individualGroupsApi.startQuarterfinal,
    onSuccess: (grid, groupId) => {
      queryClient.setQueryData(INDIVIDUAL_GROUPS_QUERY_KEYS.finalGrid(groupId), grid);
      notifications.show({
        title: "Четвертьфинал успешно начался",
        message: "Спарринги доступны прошедшим участникам",
        color: "green",
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось начать четвертьфинал",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useStartSemifinal = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<FinalGrid, Error, number>({
    mutationFn: individualGroupsApi.startSemifinal,
    onSuccess: (grid, groupId) => {
      queryClient.invalidateQueries({
        queryKey: PLACES_QUERY_KEYS.all,
      });
      queryClient.setQueryData(INDIVIDUAL_GROUPS_QUERY_KEYS.finalGrid(groupId), grid);
      notifications.show({
        title: "Полуфинал успешно начался",
        message: "Спарринги доступны прошедшим участникам",
        color: "green",
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось начать полуфинал",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useStartFinal = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<FinalGrid, Error, number>({
    mutationFn: individualGroupsApi.startFinal,
    onSuccess: (grid, groupId) => {
      queryClient.invalidateQueries({
        queryKey: PLACES_QUERY_KEYS.all,
      });
      queryClient.setQueryData(INDIVIDUAL_GROUPS_QUERY_KEYS.finalGrid(groupId), grid);
      notifications.show({
        title: "Финал успешно начался",
        message: "Спарринги доступны прошедшим участникам",
        color: "green",
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось начать финал",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useEndFinal = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<FinalGrid, Error, number>({
    mutationFn: individualGroupsApi.endFinal,
    onSuccess: (grid, groupId) => {
      queryClient.invalidateQueries({
        queryKey: PLACES_QUERY_KEYS.all,
      });
      queryClient.setQueryData(INDIVIDUAL_GROUPS_QUERY_KEYS.finalGrid(groupId), grid);
      notifications.show({
        title: "Финал успешно завершился",
        message: "Доступен экспорт в PDF",
        color: "green",
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось завершить финал",
        message: error.message,
        color: "red",
      });
    },
  });
};
