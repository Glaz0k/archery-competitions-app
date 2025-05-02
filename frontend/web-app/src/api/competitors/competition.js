import { parseISO } from "date-fns";
import apiMock from "../mocks";

export async function addCompetitor(competitionId, competitorId) {
  await apiMock();
  const data = {
    competitionId: competitionId,
    competitor: {
      id: competitorId,
      fullName: "Иванов Иван",
      birthDate: new Date(1999, 8, 24),
      identity: Gender.MALE,
      bow: BowClass.CLASSIC,
      rank: SportsRank.CLASS_1,
      region: "Санкт-Петербург",
      federation: null,
      club: null,
    },
    isActive: true,
    createdAt: parseISO("2024-08-09T18:31:42+03"),
  };
  return data;
}

export async function getCompetitors(competitionId) {
  await apiMock();
  const data = [
    {
      competitionId: competitionId,
      competitor: {
        id: 0,
        fullName: "Уиллем Кумеш",
        birthDate: new Date(1999, 8, 24),
        identity: Gender.MALE,
        bow: BowClass.CLASSIC_NEWBIE,
        rank: SportsRank.CLASS_1,
        region: "Санкт-Петербург",
        federation: null,
        club: null,
      },
      isActive: true,
      createdAt: parseISO("2024-08-09T18:31:42+03"),
    },
    {
      competitionId: competitionId,
      competitor: {
        id: 1,
        fullName: "Синъя Когами",
        birthDate: new Date(1999, 8, 24),
        identity: Gender.MALE,
        bow: BowClass.BLOCK,
        rank: SportsRank.MASTER,
        region: "Москва",
        federation: "РОФСО СПФСЛ",
        club: null,
      },
      isActive: true,
      createdAt: parseISO("2024-07-09T18:31:42+03"),
    },
    {
      competitionId: competitionId,
      competitor: {
        id: 2,
        fullName: "Тисато Нисикиги",
        birthDate: new Date(1999, 8, 24),
        identity: Gender.FEMALE,
        bow: BowClass.CLASSIC,
        rank: SportsRank.CLASS_1,
        region: "Санкт-Петербург",
        federation: null,
        club: null,
      },
      isActive: true,
      createdAt: parseISO("2024-06-09T18:31:42+03"),
    },
    {
      competitionId: competitionId,
      competitor: {
        id: 3,
        fullName: "Рюко Матой",
        birthDate: new Date(1999, 8, 24),
        identity: Gender.FEMALE,
        bow: BowClass.CLASSIC,
        rank: SportsRank.MASTER,
        region: "Ленинградская область",
        federation: "РОО СФСЛЛО",
        club: "ССК Авангард",
      },
      isActive: true,
      createdAt: parseISO("2024-05-09T18:31:42+03"),
    },
    {
      competitionId: competitionId,
      competitor: {
        id: 4,
        fullName: "Нулл Нуллов",
        birthDate: new Date(1999, 8, 24),
        identity: Gender.MALE,
        bow: null,
        rank: null,
        region: null,
        federation: null,
        club: null,
      },
      isActive: false,
      createdAt: parseISO("2024-04-09T18:31:42+03"),
    },
  ];
  return data;
}

export async function putCompetitor(competitionId, competitorId, isActive) {
  console.log(competitorId);
  await apiMock();
  return {
    competitionId: competitionId,
    competitor: {
      id: competitorId,
      fullName: "Изменен Изменнов",
      birthDate: new Date(1999, 8, 24),
      identity: Gender.MALE,
      bow: null,
      rank: null,
      region: null,
      federation: null,
      club: null,
    },
    isActive: isActive,
    createdAt: parseISO("2024-04-09T18:31:42+03"),
  };
}

export async function deleteCompetitor(_competitionId, _competitorId) {
  await apiMock();
  return true;
}
