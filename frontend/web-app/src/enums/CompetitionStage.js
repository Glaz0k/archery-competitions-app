import createEnum from "../helper/createEnum";

const CompetitionStage = createEnum({
  STAGE_1: { value: "I", textValue: "I этап" },
  STAGE_2: { value: "II", textValue: "II этап" },
  STAGE_3: { value: "III", textValue: "III этап" },
  FINAL: { value: "F", textValue: "Финал" },
});

export default CompetitionStage;
