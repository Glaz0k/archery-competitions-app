import { TextInput } from "@mantine/core";
import useCupForm from "../../../hooks/useCupForm";
import AddFormModal from "../AddFormModal";

export default function AddCupModal({ isOpened, onClose, onSubmit, isLoading }) {
  const cupForm = useCupForm();

  return (
    <AddFormModal
      title="Новый кубок"
      addForm={cupForm}
      isOpened={isOpened}
      onClose={onClose}
      onSubmit={onSubmit}
      isLoading={isLoading}
    >
      <TextInput
        withAsterisk
        label="Название"
        key={cupForm.key("title")}
        {...cupForm.getInputProps("title")}
      />
      <TextInput
        label="Адрес проведения"
        key={cupForm.key("address")}
        {...cupForm.getInputProps("address")}
      />
      <TextInput label="Сезон" key={cupForm.key("season")} {...cupForm.getInputProps("season")} />
    </AddFormModal>
  );
}
