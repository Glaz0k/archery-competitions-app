import type { CompetitorGroupDetail, IndividualGroup } from "../../../entities";
import apiClient from "../../axios/config";
import { mapToCompetitorGroupDetail } from "../competitors/mappers";
import { CompetitorGroupDetailAPISchema } from "../competitors/schemas";
import { mapToIndividualGroup } from "./mappers";
import { IndividualGroupAPISchema } from "./schemas";

export const individualGroupsApi = {
  getGroup: async (groupId: number): Promise<IndividualGroup> => {
    const response = await apiClient.get(`/individual_groups/${groupId}`);
    const validatedResponse = IndividualGroupAPISchema.parse(response.data);
    return mapToIndividualGroup(validatedResponse);
  },
  deleteGroup: async (groupId: number): Promise<void> => {
    await apiClient.delete(`/individual_groups/${groupId}`);
  },
  getCompetitors: async (groupId: number): Promise<CompetitorGroupDetail[]> => {
    const response = await apiClient.get(`/individual_groups/${groupId}/competitors`);
    const validatedResponse = CompetitorGroupDetailAPISchema.array().parse(response);
    return validatedResponse.map(mapToCompetitorGroupDetail);
  },
  syncCompetitors: async (groupId: number): Promise<CompetitorGroupDetail[]> => {
    const response = await apiClient.post(`/individual_groups/${groupId}/competitors/sync`);
    const validatedResponse = CompetitorGroupDetailAPISchema.array().parse(response);
    return validatedResponse.map(mapToCompetitorGroupDetail);
  },
};
