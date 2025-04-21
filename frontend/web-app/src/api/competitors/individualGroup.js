import BowClass from "../../enums/BowClass";
import Gender from "../../enums/competitor/Gender";
import SportsRank from "../../enums/competitor/SportsRank";
import apiMock from "../mocks";

export async function getCompetitors(groupId) {
  await apiMock();
  const data = [
    {
      competitionId: groupId,
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
    },
    {
      competitionId: groupId,
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
    },
    {
      competitionId: groupId,
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
    },
    {
      competitionId: groupId,
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
    },
    {
      competitionId: groupId,
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
    },
  ];
  return data;
}

export async function syncCompetitors(groupId) {
  return await getCompetitors(groupId);
}
