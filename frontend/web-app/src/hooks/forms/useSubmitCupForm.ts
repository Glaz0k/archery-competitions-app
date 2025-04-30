import { z } from "zod";
import { useForm, zodResolver } from "@mantine/form";
import type { CupEdit } from "../../api";
import { SEASON_REGEX } from "../../entities";

const SubmitCupFormSchema = z.object({
  title: z.string().trim().nonempty({ message: "Название не должно быть пустым" }),
  address: z.string(),
  season: z.union(
    [z.string().trim().regex(SEASON_REGEX, "Неверный сезон"), z.string().trim().length(0)],
    {
      message: "Неверный сезон",
    }
  ),
});

export type SubmitCupFormValues = z.infer<typeof SubmitCupFormSchema>;

export const useSubmitCupForm = () => {
  return useForm<SubmitCupFormValues, (values: SubmitCupFormValues) => CupEdit>({
    mode: "uncontrolled",
    initialValues: {
      title: "",
      address: "",
      season: "",
    },
    validate: zodResolver(SubmitCupFormSchema),
    transformValues: (values) => ({
      title: values.title.trim(),
      address: values.address.trim() === "" ? null : values.address.trim(),
      season: values.season.trim() === "" ? null : values.season.trim(),
    }),
  });
};
