import { FORM_INDEX, matches, useForm } from "@mantine/form";
import type { Shot } from "../../../entities";

export const useShotsForm = (initial: Shot[], shotRegex: RegExp) => {
  return useForm({
    mode: "uncontrolled",
    initialValues: {
      shots: initial,
    },
    validateInputOnChange: [`shots.${FORM_INDEX}.score`],
    validate: {
      shots: {
        score: matches(shotRegex, "Указан неверный счёт"),
      },
    },
  });
};
