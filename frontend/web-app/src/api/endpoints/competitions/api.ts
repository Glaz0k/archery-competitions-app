import {
  type Competition,
  type CompetitorCompetitionDetail,
  type IndividualGroup,
} from "../../../entities";
import apiClient from "../../axios/config";
import {
  mapToCompetitorAPIAdd,
  mapToCompetitorAPIToggle,
  mapToCompetitorCompetitionDetail,
} from "../competitors/mappers";
import { CompetitorCompetitionDetailAPISchema } from "../competitors/schemas";
import type { CompetitorAdd, CompetitorToggle } from "../competitors/types";
import { mapToIndividualGroup, mapToIndividualGroupAPICreate } from "../individual-groups/mappers";
import { IndividualGroupAPISchema } from "../individual-groups/schemas";
import type { IndividualGroupCreate } from "../individual-groups/types";
import { mapToCompetition, mapToCompetitionAPIEdit } from "./mappers";
import { CompetitionAPISchema } from "./schemas";
import type { CompetitionEdit } from "./types";

export const competitionsApi = {
  getCompetition: async (competitionId: number): Promise<Competition> => {
    const response = await apiClient.get(`/competitions/${competitionId}`);
    const validatedResponse = CompetitionAPISchema.parse(response.data);
    return mapToCompetition(validatedResponse);
  },
  putCompetition: async (competitionId: number, data: CompetitionEdit): Promise<Competition> => {
    const response = await apiClient.put(
      `/competitions/${competitionId}`,
      mapToCompetitionAPIEdit(data)
    );
    const validatedResponse = CompetitionAPISchema.parse(response.data);
    return mapToCompetition(validatedResponse);
  },
  deleteCompetition: async (competitionId: number): Promise<void> => {
    await apiClient.delete(`/competitions/${competitionId}`);
  },
  endCompetition: async (competitionId: number): Promise<Competition> => {
    const response = await apiClient.post(`/competitions/${competitionId}/end`);
    const validatedResponse = CompetitionAPISchema.parse(response.data);
    return mapToCompetition(validatedResponse);
  },
  postCompetitor: async (
    competitionId: number,
    data: CompetitorAdd
  ): Promise<CompetitorCompetitionDetail> => {
    const response = await apiClient.post(
      `/competitions/${competitionId}/competitors`,
      mapToCompetitorAPIAdd(data)
    );
    const validatedResponse = CompetitorCompetitionDetailAPISchema.parse(response.data);
    return mapToCompetitorCompetitionDetail(validatedResponse);
  },
  getCompetitors: async (competitionId: number): Promise<CompetitorCompetitionDetail[]> => {
    const response = await apiClient.get(`/competitions/${competitionId}/competitors`);
    const validatedResponse = CompetitorCompetitionDetailAPISchema.array().parse(response.data);
    return validatedResponse.map(mapToCompetitorCompetitionDetail);
  },
  putCompetitor: async (
    competitionId: number,
    competitorId: number,
    data: CompetitorToggle
  ): Promise<CompetitorCompetitionDetail> => {
    const response = await apiClient.put(
      `/competitions/${competitionId}/competitors/${competitorId}`,
      mapToCompetitorAPIToggle(data)
    );
    const validatedResponse = CompetitorCompetitionDetailAPISchema.parse(response.data);
    return mapToCompetitorCompetitionDetail(validatedResponse);
  },
  deleteCompetitor: async (competitionId: number, competitorId: number): Promise<void> => {
    await apiClient.delete(`/competitions/${competitionId}/competitors/${competitorId}`);
  },
  postIndividualGroup: async (
    competitionId: number,
    data: IndividualGroupCreate
  ): Promise<IndividualGroup> => {
    const response = await apiClient.post(
      `/competitions/${competitionId}/individual_groups`,
      mapToIndividualGroupAPICreate(data)
    );
    const validatedResponse = IndividualGroupAPISchema.parse(response.data);
    return mapToIndividualGroup(validatedResponse);
  },
  getIndidvidualGroups: async (competitionId: number): Promise<IndividualGroup[]> => {
    const response = await apiClient.get(`/competitions/${competitionId}/individual_groups`);
    const validatedResponse = IndividualGroupAPISchema.array().parse(response.data);
    return validatedResponse.map(mapToIndividualGroup);
  },
};
