import {
  Gender,
  Identity,
  type IndividualGroup,
  type Qualification,
  type QualificationRound,
  type QualificationRoundShrinked,
  type QualificationSection,
} from "../../../entities";
import { mapToCompetitorShrinked } from "../competitors/mappers";
import { mapToRangeGroup } from "../shared/mappers";
import type {
  IndividualGroupAPI,
  IndividualGroupAPICreate,
  IndividualGroupCreate,
  QualificationAPI,
  QualificationRoundAPI,
  QualificationRoundShrinkedAPI,
  QualificationSectionAPI,
} from "./types";

export const mapToIdentity = (identity: IndividualGroupAPI["identity"]): Identity => {
  switch (identity) {
    case Gender.MALE:
      return Identity.MALES;
    case Gender.FEMALE:
      return Identity.FEMALES;
    case null:
      return Identity.UNITED;
    default:
      throw new Error("Invalid identity");
  }
};

export const mapToIdentityAPI = (identity: Identity): IndividualGroupAPI["identity"] => {
  switch (identity) {
    case Identity.MALES:
      return Gender.MALE;
    case Identity.FEMALES:
      return Gender.FEMALE;
    case Identity.UNITED:
      return null;
    default:
      throw new Error("Invalid identity");
  }
};

export const mapToIndividualGroupAPICreate = (
  request: IndividualGroupCreate
): IndividualGroupAPICreate => {
  return {
    bow: request.bow,
    identity: mapToIdentityAPI(request.identity),
  };
};

export const mapToIndividualGroup = (response: IndividualGroupAPI): IndividualGroup => {
  return {
    id: response.id,
    competitionId: response.competition_id,
    bow: response.bow,
    identity: mapToIdentity(response.identity),
    state: response.state,
  };
};

export const mapToQualification = (response: QualificationAPI): Qualification => {
  return {
    groupId: response.group_id,
    distance: response.distance,
    roundCount: response.round_count,
    sections: response.sections.map(mapToSection),
  };
};

export const mapToSection = (response: QualificationSectionAPI): QualificationSection => {
  return {
    id: response.id,
    competitor: mapToCompetitorShrinked(response.competitor),
    place: response.place,
    rounds: response.rounds.map(mapToQualificationRoundShrinked),
    total: response.total,
    count10: response["10_s"],
    count9: response["9_s"],
    rankGained: response.rank_gained,
  };
};

export const mapToQualificationRoundShrinked = (
  response: QualificationRoundShrinkedAPI
): QualificationRoundShrinked => {
  return {
    ordinal: response.round_ordinal,
    isActive: response.is_active,
    totalScore: response.total_score,
  };
};

export const mapToQualificationRound = (response: QualificationRoundAPI): QualificationRound => {
  return {
    sectionId: response.section_id,
    ordinal: response.round_ordinal,
    isActive: response.is_active,
    rangeGroup: mapToRangeGroup(response.range_group),
  };
};
