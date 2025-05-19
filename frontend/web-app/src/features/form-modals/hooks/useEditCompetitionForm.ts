import { z } from "zod";
import { useForm, zodResolver } from "@mantine/form";
import type { CompetitionEdit } from "../../../api";

const EditCompetitionFormSchema = z.object({
  startDate: z.date().nullable(),
  endDate: z.date().nullable(),
});

export type EditCompetitionFormValues = z.infer<typeof EditCompetitionFormSchema>;

export const useEditCompetitionForm = () => {
  return useForm<EditCompetitionFormValues, (values: EditCompetitionFormValues) => CompetitionEdit>(
    {
      mode: "uncontrolled",
      initialValues: {
        startDate: null,
        endDate: null,
      },
      validate: zodResolver(EditCompetitionFormSchema),
      transformValues: (values) => ({
        startDate: values.startDate,
        endDate: values.endDate,
      }),
    }
  );
};
