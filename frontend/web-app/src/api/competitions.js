import BowClass from "../enums/BowClass";
import CompetitionStage from "../enums/CompetitionStage";
import GroupGender from "../enums/GroupGender";
import GroupState from "../enums/GroupState";
import apiMock from "./mocks";

export async function getCompetition(id) {
  await apiMock();
  const competition = {
    id: id,
    stage: CompetitionStage.STAGE_1,
    startDate: new Date(2023, 10, 21),
    endDate: new Date(2023, 10, 22),
    isEnded: false,
  };
  return competition;
}

export async function putCompetition({ id, startDate, endDate }) {
  await apiMock();
  const competition = {
    id: id,
    stage: CompetitionStage.STAGE_1,
    startDate: startDate,
    endDate: endDate,
    isEnded: false,
  };
  return competition;
}

export async function deleteCompetition(id) {
  console.log(id);
  await apiMock();
  return true;
}

export async function postEndCompetition(id) {
  await apiMock();
  const competition = {
    id: id,
    stage: CompetitionStage.STAGE_1,
    startDate: new Date(2025, 4, 14),
    endDate: new Date(2025, 4, 15),
    isEnded: true,
  };
  return competition;
}

export async function postIndividualGroup(id, { bow, identity }) {
  await apiMock();
  const individualGroup = {
    id: Math.floor(Math.random() * 10000),
    competitionId: id,
    bow: BowClass.valueOf(bow),
    identity: GroupGender.valueOf(identity),
    state: GroupState.CREATED,
  };
  return individualGroup;
}

export async function getIndividualGroups(id) {
  await apiMock();
  const data = [
    {
      bow: "classic",
      identity: "male",
      state: "completed",
    },
    {
      bow: "classic",
      identity: "female",
      state: "created",
    },
    {
      bow: "block",
      identity: "male",
      state: "qualification_start",
    },
    {
      bow: "block",
      identity: "female",
      state: "qualification_end",
    },
    {
      bow: "classic_newbie",
      identity: null,
      state: "final_start",
    },
    {
      bow: "3D_classic",
      identity: null,
      state: "completed",
    },
  ].map((entry, index) => {
    return {
      id: index,
      competition_id: id,
      ...entry,
    };
  });
  return data.map(mapToIndividualGroup);
}

function mapToIndividualGroup({ id, competition_id, bow, identity, state }) {
  return {
    id: Number(id),
    competitionId: Number(competition_id),
    bow: BowClass.valueOf(bow),
    identity: identity != null ? GroupGender.valueOf(identity) : GroupGender.UNITED,
    state: GroupState.valueOf(state),
  };
}
