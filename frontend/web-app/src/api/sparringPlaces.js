import apiMock from "./mocks";

const testRangeGroup = {
  id: 1297461324,
  rangesMaxCount: 5,
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

export async function getPlace(placeId) {
  await apiMock();
  return {
    id: placeId,
    competitor: {
      id: 759265,
      fullName: "Загороднова Анастасия",
    },
    rangeGroup: testRangeGroup,
    isActive: true,
    shootOut: {
      score: "7",
      priority: true,
    },
    sparringScore: 37,
  };
}

export async function getRangeGroup(_placeId) {
  await apiMock();
  return testRangeGroup;
}

export async function putRange(_placeId, { rangeOrdinal, shots }) {
  await apiMock();
  return {
    id: 234098123,
    rangeOrdinal: rangeOrdinal,
    isActive: true,
    shots: shots,
    range_score: 29,
  };
}

export async function endRange(_placeId, rangeOrdinal) {
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

export async function putShootOut(_placeId, { score, priority }) {
  await apiMock();
  return {
    id: 123124312,
    score: score,
    priority: priority,
  };
}
