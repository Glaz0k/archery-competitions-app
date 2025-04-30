import { formatISO, parseISO } from "date-fns";
import type { Competition } from "../../../entities";
import type {
  CompetitionAPI,
  CompetitionAPICreate,
  CompetitionAPIEdit,
  CompetitionCreate,
  CompetitionEdit,
} from "./types";

export const mapToCompetitionAPICreate = (request: CompetitionCreate): CompetitionAPICreate => {
  return {
    stage: request.stage,
    start_date: request.startDate ? formatISO(request.startDate, { representation: "date" }) : null,
    end_date: request.endDate ? formatISO(request.endDate, { representation: "date" }) : null,
  };
};

export const mapToCompetitionAPIEdit = (request: CompetitionEdit): CompetitionAPIEdit => {
  return {
    start_date: request.startDate ? formatISO(request.startDate, { representation: "date" }) : null,
    end_date: request.endDate ? formatISO(request.endDate, { representation: "date" }) : null,
  };
};

export const mapToCompetition = (response: CompetitionAPI): Competition => {
  return {
    id: response.id,
    stage: response.stage,
    startDate: response.start_date ? parseISO(response.start_date) : null,
    endDate: response.end_date ? parseISO(response.end_date) : null,
    isEnded: response.is_ended,
  };
};
