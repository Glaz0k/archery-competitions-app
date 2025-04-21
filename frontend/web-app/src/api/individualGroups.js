import BowClass from "../enums/BowClass";
import GroupGender from "../enums/GroupGender";
import GroupState from "../enums/GroupState";
import apiMock from "./mocks";

export async function getIndividualGroup(id) {
  return {
    id: id,
    competitionId: 0,
    bow: BowClass.CLASSIC,
    identity: GroupGender.MALE,
    state: GroupState.SEMIFINAL_START,
  };
}

export async function deleteIndividualGroup(id) {
  console.log(id);
  await apiMock();
  return true;
}
