import { GroupState } from "../../entities";

export const groupStateDescriptions: Record<GroupState, string> = {
  [GroupState.CREATED]: "Создан",
  [GroupState.QUAL_START]: "Квалификация началась",
  [GroupState.QUAL_END]: "Квалификация закончилась",
  [GroupState.QUARTERFINAL_START]: "Четвертьфинал начался",
  [GroupState.SEMIFINAL_START]: "Полуфинал начался",
  [GroupState.FINAL_START]: "Финал начался",
  [GroupState.COMPLETED]: "Закончен",
};

export const getGroupStateDescription = (state: GroupState): string => {
  return groupStateDescriptions[state] || "Неизвестно";
};
