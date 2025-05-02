import SparringState from "../enums/SparringState";
import apiMock from "./mocks";

export async function getIndividualGroup(groupId) {
  return {
    id: groupId,
    competitionId: 0,
    bow: BowClass.CLASSIC,
    identity: GroupGender.MALE,
    state: GroupState.FINAL_START,
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

export async function getFinalGrid(groupId) {
  await apiMock();
  return {
    groupId: groupId,
    quarterfinal: {
      sparring1: {
        id: 22340981,
        topPlace: {
          id: 3458761,
          competitor: {
            id: 759265,
            fullName: "Иванов Иван",
          },
          rangeGroup: {
            id: 1297461324,
            rangesMaxCount: 10,
            rangeSize: 3,
            ranges: [
              {
                id: 234098123,
                rangeOrdinal: 1,
                isActive: false,
                shots: [
                  {
                    shotOrdinal: 1,
                    score: "10",
                  },
                  {
                    shotOrdinal: 2,
                    score: "9",
                  },
                  {
                    shotOrdinal: 3,
                    score: "X",
                  },
                ],
                rangeScore: 29,
              },
              {
                id: 23987614,
                rangeOrdinal: 2,
                isActive: true,
                shots: [
                  {
                    shotOrdinal: 1,
                    score: "8",
                  },
                  {
                    shotOrdinal: 2,
                    score: null,
                  },
                  {
                    shotOrdinal: 3,
                    score: null,
                  },
                ],
                rangeScore: 8,
              },
              {
                id: 43563568,
                rangeOrdinal: 3,
                isActive: false,
                shots: null,
                rangeScore: null,
              },
            ],
            totalScore: 37,
          },
          isActive: true,
        },
        botPlace: null,
        state: SparringState.TOP_WIN,
      },
      sparring2: null,
      sparring3: null,
      sparring4: null,
    },
    semifinal: null,
    final: null,
  };
}

export async function startQuarterfinal(groupId) {
  return await getFinalGrid(groupId);
}

export async function startSemifinal(groupId) {
  return await getFinalGrid(groupId);
}

export async function startFinal(groupId) {
  return await getFinalGrid(groupId);
}

export async function endFinal(groupId) {
  return await getFinalGrid(groupId);
}
