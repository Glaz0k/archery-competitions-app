import { z } from "zod";
import { useForm, zodResolver } from "@mantine/form";
import type { IndividualGroupCreate } from "../../../api";
import { BowClass, Identity } from "../../../entities";

const CreateIndividualGroupFormSchema = z.object({
  bow: z.nativeEnum(BowClass),
  identity: z.nativeEnum(Identity),
});

export type CreateIndividualGroupFormValues = z.infer<typeof CreateIndividualGroupFormSchema>;

export const useCreateIndividualGroupForm = () => {
  return useForm<
    CreateIndividualGroupFormValues,
    (values: CreateIndividualGroupFormValues) => IndividualGroupCreate
  >({
    mode: "uncontrolled",
    initialValues: {
      bow: BowClass.CLASSIC,
      identity: Identity.MALES,
    },
    validate: zodResolver(CreateIndividualGroupFormSchema),
    transformValues: (values) => ({
      bow: values.bow,
      identity: values.identity,
    }),
  });
};
