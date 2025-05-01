import { Gender, Identity, type IndividualGroup } from "../../../entities";
import type { IndividualGroupAPI, IndividualGroupAPICreate, IndividualGroupCreate } from "./types";

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
