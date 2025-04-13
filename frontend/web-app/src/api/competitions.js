import apiMock from "./mocks";

export async function putCompetition({ competitionId, startDate, endDate }) {
  await apiMock();
  const competition = {
    id: competitionId,
    stage: "I",
    startDate: startDate,
    endDate: endDate,
    isEnded: false,
  };
  return competition;
}

export async function deleteCompetition(cupId) {
  console.log(cupId);
  await apiMock();
  return true;
}

export async function endCompetition({ competitionId }) {
  await apiMock();
  const competition = {
    id: competitionId,
    stage: "I",
    startDate: new Date(2025, 4, 14),
    endDate: new Date(2025, 4, 14),
    isEnded: true,
  };
  return competition;
}
