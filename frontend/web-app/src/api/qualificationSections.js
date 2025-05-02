import apiMock from "./mocks";

export async function getSection(sectionId) {
  await apiMock();
  if (sectionId === 0) {
    return {
      id: sectionId,
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
    };
  } else {
    return {
      id: sectionId,
      competitor: {
        id: 0,
        fullName: "Тестов Тест",
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
    };
  }
}

export async function getRound(sectionId, roundOrdinal) {
  await apiMock();
  return {
    sectionId: sectionId,
    roundOrdinal: roundOrdinal,
    isActive: true,
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
  };
}

export async function getRanges(_sectionId, _roundOrdinal) {
  await apiMock();
  return {
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
  };
}

export async function putRange(_sectionId, _roundOrdinal, { rangeOrdinal, shots }) {
  await apiMock();
  return {
    id: 234098123,
    rangeOrdinal: rangeOrdinal,
    isActive: true,
    shots: shots,
    range_score: 29,
  };
}

export async function endRange(_sectionId, _roundOrdinal, rangeOrdinal) {
  await apiMock();
  return {
    id: 234098123,
    rangeOrdinal: rangeOrdinal,
    isActive: false,
    shots: [
      {
        shotOrdinal: 1,
        score: "8",
      },
      {
        shotOrdinal: 2,
        score: "10",
      },
      {
        shotOrdinal: 3,
        score: "M",
      },
    ],
    range_score: 18,
  };
}
