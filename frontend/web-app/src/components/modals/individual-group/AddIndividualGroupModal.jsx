import { NativeSelect } from "@mantine/core";
import BowClass from "../../../enums/BowClass";
import GroupGender from "../../../enums/GroupGender";
import useIndividualGroupForm from "../../../hooks/useIndividualGroupForm";
import AddFormModal from "../AddFormModal";

export default function AddIndividualGroupModal({ isOpened, onClose, onSubmit, isLoading }) {
  const individualGroupForm = useIndividualGroupForm();

  return (
    <AddFormModal
      addForm={individualGroupForm}
      title="Новая индивидуальная группа"
      isOpened={isOpened}
      onClose={onClose}
      onSubmit={onSubmit}
      isLoading={isLoading}
    >
      <NativeSelect
        label="Тип лука"
        data={Object.values(BowClass).map((bowClass) => {
          return {
            label: bowClass.textValue,
            value: bowClass.value,
          };
        })}
        key={individualGroupForm.key("bow")}
        {...individualGroupForm.getInputProps("bow")}
      />
      <NativeSelect
        label="Пол"
        data={Object.values(GroupGender).map((gender) => {
          return {
            label: gender.textValue,
            value: gender.value,
          };
        })}
        key={individualGroupForm.key("identity")}
        {...individualGroupForm.getInputProps("identity")}
      />
    </AddFormModal>
  );
}
