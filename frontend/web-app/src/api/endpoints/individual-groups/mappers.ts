import {
  Gender,
  Identity,
  type Final,
  type FinalGrid,
  type IndividualGroup,
  type Qualification,
  type QualificationRound,
  type QualificationRoundShrinked,
  type QualificationSection,
  type Quarterfinal,
  type Semifinal,
  type ShootOut,
  type Sparring,
  type SparringPlace,
} from "../../../entities";
import { mapToCompetitorShrinked } from "../competitors/mappers";
import { mapToRangeGroup } from "../shared/mappers";
import type {
  FinalAPI,
  FinalGridAPI,
  IndividualGroupAPI,
  IndividualGroupAPICreate,
  IndividualGroupCreate,
  QualificationAPI,
  QualificationRoundAPI,
  QualificationRoundShrinkedAPI,
  QualificationSectionAPI,
  QuarterfinalAPI,
  SemifinalAPI,
  ShootOutAPI,
  SparringAPI,
  SparringPlaceAPI,
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

export const mapToFinalGrid = (response: FinalGridAPI): FinalGrid => {
  return {
    groupId: response.group_id,
    quarterfinal: mapToQuarterfinal(response.quarterfinal),
    semifinal: response.semifinal ? mapToSemifinal(response.semifinal) : null,
    final: response.final ? mapToFinal(response.final) : null,
  };
};

export const mapToQuarterfinal = (response: QuarterfinalAPI): Quarterfinal => {
  return {
    sparring1: response.sparring_1 ? mapToSparring(response.sparring_1) : null,
    sparring2: response.sparring_2 ? mapToSparring(response.sparring_2) : null,
    sparring3: response.sparring_3 ? mapToSparring(response.sparring_3) : null,
    sparring4: response.sparring_4 ? mapToSparring(response.sparring_4) : null,
  };
};

export const mapToSemifinal = (response: SemifinalAPI): Semifinal => {
  return {
    sparring5: response.sparring_5 ? mapToSparring(response.sparring_5) : null,
    sparring6: response.sparring_6 ? mapToSparring(response.sparring_6) : null,
  };
};

export const mapToFinal = (response: FinalAPI): Final => {
  return {
    sparringGold: response.sparring_gold ? mapToSparring(response.sparring_gold) : null,
    sparringBronze: response.sparring_bronze ? mapToSparring(response.sparring_bronze) : null,
  };
};

export const mapToSparring = (response: SparringAPI): Sparring => {
  return {
    id: response.id,
    top: response.top_place ? mapToSparringPlace(response.top_place) : null,
    bot: response.bot_place ? mapToSparringPlace(response.bot_place) : null,
    state: response.state,
  };
};

export const mapToSparringPlace = (response: SparringPlaceAPI): SparringPlace => {
  return {
    id: response.id,
    competitor: mapToCompetitorShrinked(response.competitor),
    rangeGroup: mapToRangeGroup(response.range_group),
    isActive: response.is_active,
    score: response.sparring_score,
    shootOut: response.shoot_out ? mapToShootOut(response.shoot_out) : null,
  };
};

export const mapToShootOut = (response: ShootOutAPI): ShootOut => {
  return {
    id: response.id,
    score: response.score,
    priority: response.priority,
  };
};
