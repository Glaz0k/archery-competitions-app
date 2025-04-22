import BowClass from "../enums/BowClass";
import SportsRank from "../enums/competitor/SportsRank";
import GroupGender from "../enums/GroupGender";
import GroupState from "../enums/GroupState";
import apiMock from "./mocks";

export async function getIndividualGroup(groupId) {
  return {
    id: groupId,
    competitionId: 0,
    bow: BowClass.CLASSIC,
    identity: GroupGender.MALE,
    state: GroupState.QUAL_START,
  };
}

export async function deleteIndividualGroup(groupId) {
  console.log(groupId);
  await apiMock();
  return true;
}

const testSections = [
  {
    id: 0,
    competitor: {
      id: 0,
      fullName: "Павлов Виталий",
    },
    place: 1,
    rounds: [
      {
        roundOrdinal: 1,
        isOngoing: false,
        total: 286,
      },
      {
        roundOrdinal: 2,
        isOngoing: false,
        total: 280,
      },
    ],
    total: 566,
    count10: 32,
    count9: 22,
    rankGained: SportsRank.MASTER_CANDIDATE,
  },
  {
    id: 1,
    competitor: {
      id: 1,
      fullName: "Левин Сергей",
    },
    place: 2,
    rounds: [
      {
        roundOrdinal: 1,
        isOngoing: false,
        total: 280,
      },
      {
        roundOrdinal: 2,
        isOngoing: false,
        total: 266,
      },
    ],
    total: 546,
    count10: 20,
    count9: 28,
    rankGained: SportsRank.MASTER_CANDIDATE,
  },
];

export async function getQualification(groupId) {
  await apiMock();
  if (groupId == 1) {
    return null;
  }
  return {
    groupId: groupId,
    distance: "18m",
    roundCount: 2,
    sections: testSections,
  };
}

export async function startQualification(groupId) {
  return await getQualification(groupId);
}

export async function endQualification(groupId) {
  return await getQualification(groupId);
}
