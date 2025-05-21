import { z } from "zod";
import { useForm, zodResolver } from "@mantine/form";
import type { CompetitionCreate } from "../../../api";
import { CompetitionStage } from "../../../entities";

const CreateCompetitionFormSchema = z.object({
  stage: z.nativeEnum(CompetitionStage),
  startDate: z.date().nullable(),
  endDate: z.date().nullable(),
});

export type CreateCompetitionFormValues = z.infer<typeof CreateCompetitionFormSchema>;

export const useCreateCompetitionForm = () => {
  return useForm<
    CreateCompetitionFormValues,
    (values: CreateCompetitionFormValues) => CompetitionCreate
  >({
    mode: "uncontrolled",
    initialValues: {
      stage: CompetitionStage.STAGE_1,
      startDate: null,
      endDate: null,
    },
    validate: zodResolver(CreateCompetitionFormSchema),
    transformValues: (values) => ({
      stage: values.stage,
      startDate: values.startDate,
      endDate: values.endDate,
    }),
  });
};
