import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import type { Range, RangeGroup, ShootOut } from "../../../entities";
import { PLACES_QUERY_KEYS } from "../../query-keys/sparringPlaces";
import type { RangeEdit } from "../shared/types";
import { updateRangeGroup } from "../shared/utils";
import { sparringPlacesApi } from "./api";
import type { ShootOutEdit } from "./types";

export const useSparringPlace = (placeId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: PLACES_QUERY_KEYS.element(placeId),
    queryFn: async () => await sparringPlacesApi.getPlace(placeId),
    enabled,
  });
};

export const usePlaceRangeGroup = (placeId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: PLACES_QUERY_KEYS.rangeGroup(placeId),
    queryFn: async () => await sparringPlacesApi.getRangeGroup(placeId),
    enabled,
  });
};

export const useEditPlaceRange = (onSuccess?: () => unknown) => {
  const queryClient = useQueryClient();
  return useMutation<Range, Error, [number, RangeEdit]>({
    mutationFn: async ([placeId, data]) => await sparringPlacesApi.putRange(placeId, data),
    onSuccess: (edited, [placeId]) => {
      queryClient.invalidateQueries({
        queryKey: PLACES_QUERY_KEYS.element(placeId),
      });
      queryClient.setQueryData(
        PLACES_QUERY_KEYS.rangeGroup(placeId),
        (prev: RangeGroup | undefined) => updateRangeGroup(prev, edited)
      );
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось изменить серию",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useCompletePlaceRange = (onSuccess?: () => unknown) => {
  const queryClient = useQueryClient();
  return useMutation<Range, Error, [number, number]>({
    mutationFn: async ([placeId, rangeOrd]) => await sparringPlacesApi.endRange(placeId, rangeOrd),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: PLACES_QUERY_KEYS.all,
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось завершить серию",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useEditShootOut = (onSuccess?: () => unknown) => {
  const queryClient = useQueryClient();
  return useMutation<ShootOut, Error, [number, ShootOutEdit]>({
    mutationFn: async ([placeId, data]) => await sparringPlacesApi.putShootOut(placeId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: PLACES_QUERY_KEYS.all,
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось изменить перестрелку",
        message: error.message,
        color: "red",
      });
    },
  });
};
