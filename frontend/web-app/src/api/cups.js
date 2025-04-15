import { parseISO } from "date-fns";
import CompetitionStage from "../enums/CompetitionStage";
import apiMock from "./mocks";

export async function postCup({ title, address, season }) {
  await apiMock();
  const cup = {
    id: Math.floor(Math.random() * 10000),
    title: title,
    address: address,
    season: season,
  };
  return cup;
}

export async function getCup(id) {
  await apiMock();
  const cup = {
    id: id,
    title: "Стрелы парадокса",
    address: 'г. Санкт-Петербург, ЛК "Парадокс Лучника"',
    season: "2023/2024",
  };
  return cup;
}

export async function getCups() {
  await apiMock();
  const data = [
    {
      id: 0,
      title: "Gink",
      address: "c. Linwood, st. Douglass Street 27",
      season: null,
    },
    {
      id: 1,
      title: "Blanet",
      address: null,
      season: "2023/2024",
    },
    {
      id: 2,
      title: "Ecraze",
      address: "c. Century, st. Roosevelt Place 28",
      season: null,
    },
    {
      id: 3,
      title: "Gracker",
      address: null,
      season: "2023/2024",
    },
    {
      id: 4,
      title: "Zytrex",
      address: null,
      season: "2023/2024",
    },
    {
      id: 5,
      title: "Nebulean",
      address: "c. Winfred, st. Glen Street 39",
      season: "2023/2024",
    },
    {
      id: 6,
      title: "Zialactic",
      address: "c. Thornport, st. Rockaway Avenue 62",
      season: null,
    },
    {
      id: 7,
      title: "Ozean",
      address: null,
      season: "2023/2024",
    },
    {
      id: 8,
      title: "Ecolight",
      address: null,
      season: null,
    },
    {
      id: 9,
      title: "Quantalia",
      address: "c. Abrams, st. Downing Street 66",
      season: null,
    },
    {
      id: 10,
      title: "Zaggle",
      address: null,
      season: null,
    },
    {
      id: 11,
      title: "Earbang",
      address: null,
      season: "2023/2024",
    },
    {
      id: 12,
      title: "Inventure",
      address: null,
      season: "2023/2024",
    },
    {
      id: 13,
      title: "Opticom",
      address: "c. Mathews, st. Allen Avenue 67",
      season: null,
    },
    {
      id: 14,
      title: "Zytrax",
      address: null,
      season: null,
    },
    {
      id: 15,
      title: "Isoswitch",
      address: "c. Alleghenyville, st. Hendrix Street 28",
      season: null,
    },
    {
      id: 16,
      title: "Quilch",
      address: null,
      season: null,
    },
    {
      id: 17,
      title: "Chillium",
      address: null,
      season: "2023/2024",
    },
    {
      id: 18,
      title: "Mediot",
      address: "c. Accoville, st. Reeve Place 70",
      season: null,
    },
    {
      id: 19,
      title: "Equitox",
      address: null,
      season: null,
    },
  ];
  return data.map(mapToCup);
}

function mapToCup({ id, title, address, season }) {
  return {
    id: Number(id),
    title: title,
    address: address,
    season: season,
  };
}

export async function putCup({ id, title, address, season }) {
  await apiMock();
  const cup = {
    id: id,
    title: title,
    address: address,
    season: season,
  };
  return cup;
}

export async function deleteCup(id) {
  console.log(id);
  await apiMock();
  return true;
}

export async function postCompetition(id, { stage, startDate, endDate }) {
  console.log(id);
  await apiMock();
  const competition = {
    id: Math.floor(Math.random() * 10000),
    stage: stage,
    startDate: startDate,
    endDate: endDate,
    isEnded: false,
  };
  return competition;
}

export async function getCompetitions(id) {
  console.log(id);
  await apiMock(1.1);
  const data = [
    {
      id: 0,
      stage: "I",
      start_date: "2025-05-15",
      end_date: "2025-05-16",
      is_ended: true,
    },
    {
      id: 1,
      stage: "II",
      start_date: "2025-06-08",
      end_date: null,
      is_ended: true,
    },
    {
      id: 2,
      stage: "III",
      start_date: null,
      end_date: null,
      is_ended: false,
    },
    {
      id: 3,
      stage: "F",
      start_date: "2026-01-14",
      end_date: "2025-10-01",
      is_ended: false,
    },
  ];
  return data.map(mapToCompetition);
}

function mapToCompetition({ id, stage, start_date, end_date, is_ended }) {
  return {
    id: Number(id),
    stage: CompetitionStage.valueOf(stage),
    startDate: start_date != null ? parseISO(start_date) : null,
    endDate: end_date != null ? parseISO(end_date) : null,
    isEnded: Boolean(is_ended),
  };
}
