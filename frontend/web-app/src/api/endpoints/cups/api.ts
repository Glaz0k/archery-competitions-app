import { z } from "zod";
import { CupSchema, type Competition, type Cup } from "../../../entities";
import apiClient from "../../axios/config";
import { mapToCompetition, mapToCompetitionAPICreate } from "../competitions/mappers";
import { CompetitionAPISchema } from "../competitions/schemas";
import type { CompetitionCreate } from "../competitions/types";
import type { CupEdit } from "./types";

export const cupsApi = {
  postCup: async (data: CupEdit): Promise<Cup> => {
    const response = await apiClient.post("/cups", data);
    return CupSchema.parse(response.data);
  },
  getCup: async (cupId: number): Promise<Cup> => {
    const response = await apiClient.get(`/cups/${cupId}`);
    return CupSchema.parse(response.data);
  },
  getCups: async (): Promise<Cup[]> => {
    const response = await apiClient.get("/cups");
    return CupSchema.array().parse(response.data);
  },
  putCup: async (cupId: number, data: CupEdit): Promise<Cup> => {
    const response = await apiClient.put(`/cups/${cupId}`, data);
    return CupSchema.parse(response.data);
  },
  deleteCup: async (cupId: number): Promise<void> => {
    await apiClient.delete(`/cups/${cupId}`);
  },
  postCompetiton: async (cupId: number, data: CompetitionCreate): Promise<Competition> => {
    const response = await apiClient.post(
      `/cups/${cupId}/competitions`,
      mapToCompetitionAPICreate(data)
    );
    const validatedResponse = CompetitionAPISchema.parse(response.data);
    return mapToCompetition(validatedResponse);
  },
  getCompetitons: async (cupId: number): Promise<Competition[]> => {
    const response = await apiClient.get(`/cups/${cupId}/competitions`);
    const validatedResponse = z.array(CompetitionAPISchema).parse(response.data);
    return validatedResponse.map(mapToCompetition);
  },
};
