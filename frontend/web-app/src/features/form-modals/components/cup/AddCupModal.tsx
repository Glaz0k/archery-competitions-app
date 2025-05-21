import { TextInput } from "@mantine/core";
import type { TransformedValues } from "@mantine/form";
import { useSubmitCupForm } from "../..";
import { AddModal, type AddModalProps } from "../AddModal";

type AddCupForm = ReturnType<typeof useSubmitCupForm>;

export type AddCupModalProps = Omit<
  AddModalProps<AddCupForm["values"], TransformedValues<AddCupForm>>,
  "form" | "title" | "children"
>;

export function AddCupModal({ loading, ...others }: AddCupModalProps) {
  const cupForm = useSubmitCupForm();

  return (
    <AddModal form={cupForm} title="Новый кубок" loading={loading} {...others}>
      <TextInput
        withAsterisk
        label="Название"
        disabled={loading}
        key={cupForm.key("title")}
        {...cupForm.getInputProps("title")}
      />
      <TextInput
        label="Адрес проведения"
        disabled={loading}
        key={cupForm.key("address")}
        {...cupForm.getInputProps("address")}
      />
      <TextInput
        label="Сезон"
        disabled={loading}
        key={cupForm.key("season")}
        {...cupForm.getInputProps("season")}
      />
    </AddModal>
  );
}
