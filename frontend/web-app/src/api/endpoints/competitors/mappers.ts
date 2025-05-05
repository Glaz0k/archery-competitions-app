import { formatISO, parseISO } from "date-fns";
import type { Competitor, CompetitorCompetitionDetail } from "../../../entities";
import type {
  CompetitorAdd,
  CompetitorAPI,
  CompetitorAPIAdd,
  CompetitorAPIEdit,
  CompetitorAPIToggle,
  CompetitorCompetitionDetailAPI,
  CompetitorEdit,
  CompetitorToggle,
} from "./types";

export const mapToCompetitorAPIAdd = (request: CompetitorAdd): CompetitorAPIAdd => {
  return {
    competitor_id: request.id,
  };
};

export const mapToCompetitorAPIEdit = (request: CompetitorEdit): CompetitorAPIEdit => {
  return {
    full_name: request.fullName,
    birth_date: formatISO(request.birthDate, { representation: "date" }),
    identity: request.identity,
    bow: request.bow,
    rank: request.rank,
    region: request.region,
    federation: request.federation,
    club: request.club,
  };
};

export const mapToCompetitor = (request: CompetitorAPI): Competitor => {
  return {
    id: request.id,
    fullName: request.full_name.trim(),
    birthDate: parseISO(request.birth_date),
    identity: request.identity,
    bow: request.bow,
    rank: request.rank,
    region: request.region?.trim() ?? null,
    federation: request.federation?.trim() ?? null,
    club: request.club?.trim() ?? null,
  };
};

export const mapToCompetitorCompetitionDetail = (
  response: CompetitorCompetitionDetailAPI
): CompetitorCompetitionDetail => {
  return {
    competitionId: response.competition_id,
    competitor: mapToCompetitor(response.competitor),
    createdAt: parseISO(response.created_at),
    isActive: response.is_active,
  };
};

export const mapToCompetitorAPIToggle = (request: CompetitorToggle): CompetitorAPIToggle => {
  return {
    is_active: request.isActive,
  };
};
