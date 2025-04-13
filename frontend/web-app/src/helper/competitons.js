import { format } from "date-fns";
import Stage from "../enums/stage";

export function stageToTitle(stage) {
  let title = "Error";
  switch (stage) {
    case Stage.STAGE_1:
    case Stage.STAGE_2:
    case Stage.STAGE_3:
      title = stage + " этап";
      break;
    case Stage.FINAL:
      title = "Финал";
      break;
    default: {
      console.error("Cannot convert stage :" + stage + " to title");
    }
  }
  return title;
}

export function formatCompetitionDateRange({ startDate, endDate }) {
  if (startDate === null && endDate === null) {
    return "Не указано";
  }

  if (startDate === null) {
    return `Не указано — ${format(endDate, "d MMMM yyyy")}`;
  }

  if (endDate === null) {
    return `${format(startDate, "d MMMM yyyy")} — Не указано`;
  }

  if (format(startDate, "yyyy-MM") === format(endDate, "yyyy-MM")) {
    return `${format(startDate, "d")}-${format(endDate, "d MMMM yyyy")}`;
  } else if (format(startDate, "yyyy") === format(endDate, "yyyy")) {
    return `${format(startDate, "d MMMM")}-${format(endDate, "d MMMM yyyy")}`;
  } else {
    return `${format(startDate, "d MMMM yyyy")} — ${format(endDate, "d MMMM yyyy")}`;
  }
}
