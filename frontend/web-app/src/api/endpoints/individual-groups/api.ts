import type { CompetitorGroupDetail, IndividualGroup, Qualification } from "../../../entities";
import type { FinalGrid } from "../../../entities/final-grid/types";
import apiClient from "../../axios/config";
import { mapToCompetitorGroupDetail } from "../competitors/mappers";
import { CompetitorGroupDetailAPISchema } from "../competitors/schemas";
import { mapToFinalGrid, mapToIndividualGroup, mapToQualification } from "./mappers";
import { FinalGridAPISchema, IndividualGroupAPISchema, QualificationAPISchema } from "./schemas";

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
    const validatedResponse = CompetitorGroupDetailAPISchema.array().parse(response.data);
    return validatedResponse.map(mapToCompetitorGroupDetail);
  },
  syncCompetitors: async (groupId: number): Promise<CompetitorGroupDetail[]> => {
    const response = await apiClient.post(`/individual_groups/${groupId}/competitors/sync`);
    const validatedResponse = CompetitorGroupDetailAPISchema.array().parse(response.data);
    return validatedResponse.map(mapToCompetitorGroupDetail);
  },
  getQualification: async (groupId: number): Promise<Qualification> => {
    const response = await apiClient.get(`/individual_groups/${groupId}/qualification`);
    const validatedResponse = QualificationAPISchema.parse(response.data);
    return mapToQualification(validatedResponse);
  },
  // WARNING, HARDCODED DATA
  startQualification: async (groupId: number): Promise<Qualification> => {
    const response = await apiClient.post(`/individual_groups/${groupId}/qualification/start`, {
      distance: "18m",
      round_count: 2,
      ranges_count: 10,
      range_size: 3,
    });
    const validatedResponse = QualificationAPISchema.parse(response.data);
    return mapToQualification(validatedResponse);
  },
  endQualification: async (groupId: number): Promise<Qualification> => {
    const response = await apiClient.post(`/individual_groups/${groupId}/qualification/end`);
    const validatedResponse = QualificationAPISchema.parse(response.data);
    return mapToQualification(validatedResponse);
  },
  getFinalGrid: async (groupId: number): Promise<FinalGrid> => {
    const response = await apiClient.get(`/individual_groups/${groupId}/final_grid`);
    const validatedResponse = FinalGridAPISchema.parse(response.data);
    return mapToFinalGrid(validatedResponse);
  },
  startQuarterfinal: async (groupId: number): Promise<FinalGrid> => {
    const response = await apiClient.post(`/individual_groups/${groupId}/quarterfinal/start`);
    const validatedResponse = FinalGridAPISchema.parse(response.data);
    return mapToFinalGrid(validatedResponse);
  },
  startSemifinal: async (groupId: number): Promise<FinalGrid> => {
    const response = await apiClient.post(`/individual_groups/${groupId}/semifinal/start`);
    const validatedResponse = FinalGridAPISchema.parse(response.data);
    return mapToFinalGrid(validatedResponse);
  },
  startFinal: async (groupId: number): Promise<FinalGrid> => {
    const response = await apiClient.post(`/individual_groups/${groupId}/final/start`);
    const validatedResponse = FinalGridAPISchema.parse(response.data);
    return mapToFinalGrid(validatedResponse);
  },
  endFinal: async (groupId: number): Promise<FinalGrid> => {
    const response = await apiClient.post(`/individual_groups/${groupId}/final/end`);
    const validatedResponse = FinalGridAPISchema.parse(response.data);
    return mapToFinalGrid(validatedResponse);
  },
};
