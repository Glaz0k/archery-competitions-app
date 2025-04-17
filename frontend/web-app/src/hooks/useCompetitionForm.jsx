import { isDate } from "date-fns";
import { useForm } from "@mantine/form";
import CompetitionStage from "../enums/CompetitionStage";

export default function useCompetitionForm() {
  return useForm({
    mode: "uncontrolled",
    initialValues: {
      stage: CompetitionStage.STAGE_1.value,
      startDate: null,
      endDate: null,
    },
    validate: {
      stage: (value) => (CompetitionStage.valueOf(value) != null ? null : "Неверный этап"),
      startDate: (value) => (isDate(value) ? null : "Неверная дата начала"),
      endDate: (value) => (isDate(value) ? null : "Неверная дата начала"),
    },
  });
}
