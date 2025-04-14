import { format } from "date-fns";

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
