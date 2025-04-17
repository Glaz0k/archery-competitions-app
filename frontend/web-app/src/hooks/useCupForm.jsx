import { isNotEmpty, useForm } from "@mantine/form";
import { isEmptyString, isValidSeasonString } from "../validation/predicates";

export default function useCupForm() {
  return useForm({
    mode: "uncontrolled",
    initialValues: {
      title: "",
      address: "",
      season: "",
    },
    validate: {
      title: isNotEmpty("Название не должно быть пустым"),
      season: (value) =>
        isEmptyString(value) || isValidSeasonString(value) ? null : "Неверный сезон",
    },
  });
}
