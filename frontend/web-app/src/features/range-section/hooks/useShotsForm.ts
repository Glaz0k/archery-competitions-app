import { z } from "zod";
import { FORM_INDEX, useForm, zodResolver } from "@mantine/form";
import { ShotSchema, type Shot } from "../../../entities";

export const useShotsForm = (initial: Shot[], shotRegex: RegExp) => {
  const ShotsFormSchema = z.object({
    shots: ShotSchema.extend({
      score: z.union([z.string().trim().max(0).nullable(), z.string().trim().regex(shotRegex)], {
        message: "Указан неверный счёт",
      }),
    }).array(),
  });
  type ShotsFormValues = z.infer<typeof ShotsFormSchema>;

  return useForm<ShotsFormValues, (values: ShotsFormValues) => Shot[]>({
    mode: "uncontrolled",
    initialValues: {
      shots: initial,
    },
    validateInputOnChange: [`shots.${FORM_INDEX}.score`],
    validate: zodResolver(ShotsFormSchema),
    transformValues: (values) =>
      values.shots.map((value) => ({
        ordinal: value.ordinal,
        score: value.score === "" ? null : value.score,
      })),
  });
};
