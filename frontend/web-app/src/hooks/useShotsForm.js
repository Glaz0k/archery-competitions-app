import { FORM_INDEX, matches, useForm } from "@mantine/form";

export default function useShotsForm({ initialShots, shotRegex }) {
  return useForm({
    mode: "uncontrolled",
    initialValues: {
      shots: initialShots,
    },
    validateInputOnChange: [`shots.${FORM_INDEX}.score`],
    validate: {
      shots: {
        score: matches(shotRegex, "Неверный счёт"),
      },
    },
  });
}
