import type { Range, RangeGroup, ShootOut, SparringPlace } from "../../../entities";
import apiClient from "../../axios/config";
import { mapToShootOut, mapToSparringPlace } from "../individual-groups/mappers";
import { ShootOutAPISchema, SparringPlaceAPISchema } from "../individual-groups/schemas";
import { mapToRange, mapToRangeAPIEdit, mapToRangeGroup } from "../shared/mappers";
import { RangeAPISchema, RangeGroupAPISchema } from "../shared/schemas";
import type { RangeEdit } from "../shared/types";
import { mapToShootOutAPIEdit } from "./mappers";
import type { ShootOutEdit } from "./types";

export const sparringPlacesApi = {
  getPlace: async (placeId: number): Promise<SparringPlace> => {
    const response = await apiClient.get(`/sparring_places/${placeId}`);
    const validatedResponse = SparringPlaceAPISchema.parse(response.data);
    return mapToSparringPlace(validatedResponse);
  },
  getRangeGroup: async (placeId: number): Promise<RangeGroup> => {
    const response = await apiClient.get(`/sparring_places/${placeId}/ranges`);
    const validatedResponse = RangeGroupAPISchema.parse(response.data);
    return mapToRangeGroup(validatedResponse);
  },
  putRange: async (placeId: number, data: RangeEdit): Promise<Range> => {
    const response = await apiClient.put(
      `/sparring_places/${placeId}/ranges`,
      mapToRangeAPIEdit(data)
    );
    const validatedResponse = RangeAPISchema.parse(response.data);
    return mapToRange(validatedResponse);
  },
  endRange: async (placeId: number, rangeOrdinal: number): Promise<Range> => {
    const response = await apiClient.post(`/sparring_places/${placeId}/ranges/${rangeOrdinal}/end`);
    const validatedResponse = RangeAPISchema.parse(response.data);
    return mapToRange(validatedResponse);
  },
  putShootOut: async (placeId: number, data: ShootOutEdit): Promise<ShootOut> => {
    const response = await apiClient.put(
      `/sparring_places/${placeId}/shoot_out`,
      mapToShootOutAPIEdit(data)
    );
    const validatedResponse = ShootOutAPISchema.parse(response.data);
    return mapToShootOut(validatedResponse);
  },
};
