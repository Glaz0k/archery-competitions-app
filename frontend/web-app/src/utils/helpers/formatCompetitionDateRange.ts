import { format } from "date-fns";

export function formatCompetitionDateRange(startDate: Date | null, endDate: Date | null) {
  if (!startDate && !endDate) {
    return "Не указано";
  }

  if (!startDate) {
    return `Не указано — ${format(endDate!, "d MMMM yyyy")}`;
  }

  if (!endDate) {
    return `${format(startDate, "d MMMM yyyy")} — Не указано`;
  }

  if (format(startDate, "yyyy-MM") === format(endDate, "yyyy-MM")) {
    return `${format(startDate, "d")}-${format(endDate, "d MMMM yyyy")}`;
  } else if (format(startDate, "yyyy") === format(endDate, "yyyy")) {
    return `${format(startDate, "d MMMM")}-${format(endDate, "d MMMM yyyy")}`;
  }
  return `${format(startDate, "d MMMM yyyy")} — ${format(endDate, "d MMMM yyyy")}`;
}
