import { SportsRank } from "../../entities";

export const sportsRankDescriptions: Record<SportsRank, string> = {
  [SportsRank.CLASS_3]: "3",
  [SportsRank.CLASS_2]: "2",
  [SportsRank.CLASS_1]: "1",
  [SportsRank.MASTER_CANDIDATE]: "КМС",
  [SportsRank.MASTER]: "МС",
  [SportsRank.MASTER_INTERNATIONAL]: "МСМК",
  [SportsRank.MASTER_MERITED]: "ЗМС",
};

export const getSportsRankDescription = (rank: SportsRank): string => {
  return sportsRankDescriptions[rank] || "Неизвестно";
};
