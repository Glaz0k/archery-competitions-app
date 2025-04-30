import { TextInput } from "@mantine/core";
import type { TransformedValues } from "@mantine/form";
import { useSubmitCupForm } from "../../../hooks";
import { AddFormModal, type AddFormModalProps } from "../AddFormModal";

type AddCupForm = ReturnType<typeof useSubmitCupForm>;

export type AddCupModalProps = Omit<
  AddFormModalProps<AddCupForm["values"], TransformedValues<AddCupForm>>,
  "form" | "title" | "children"
>;

export default function AddCupModal({ opened, onClose, onSubmit, loading }: AddCupModalProps) {
  const cupForm = useSubmitCupForm();

  return (
    <AddFormModal
      form={cupForm}
      title="Новый кубок"
      opened={opened}
      onClose={onClose}
      onSubmit={onSubmit}
      loading={loading}
    >
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
    </AddFormModal>
  );
}
