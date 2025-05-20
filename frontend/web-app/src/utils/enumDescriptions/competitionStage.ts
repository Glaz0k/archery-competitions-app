import { CompetitionStage } from "../../entities";

export const competitionStageDescriptions: Record<CompetitionStage, string> = {
  [CompetitionStage.STAGE_1]: "I этап",
  [CompetitionStage.STAGE_2]: "II этап",
  [CompetitionStage.STAGE_3]: "III этап",
  [CompetitionStage.FINAL]: "Финал",
};

export const getCompetitionStageDescription = (stage: CompetitionStage): string => {
  return competitionStageDescriptions[stage] || "Неизвестно";
};
