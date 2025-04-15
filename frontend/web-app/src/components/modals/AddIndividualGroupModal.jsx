import { Button, Modal, NativeSelect } from "@mantine/core";
import { useForm } from "@mantine/form";
import BowClass from "../../enums/BowClass";
import GroupGender from "../../enums/GroupGender";

export default function AddIndividualGroupModal({ isOpened, onClose, onSubmit, isLoading }) {
  const individualGroupForm = useForm({
    mode: "uncontrolled",
    initialValues: {
      bow: BowClass.CLASSIC.value,
      identity: GroupGender.MALE.value,
    },
  });

  const actionsOnClose = () => {
    individualGroupForm.reset();
    onClose?.();
  };

  const actionsOnSubmit = async (individualGroupFormValues) => {
    if (await onSubmit(individualGroupFormValues)) {
      individualGroupForm.reset();
    }
  };

  return (
    <Modal opened={isOpened} onClose={actionsOnClose} title="Новая индивидуальная группа" centered>
      <form onSubmit={individualGroupForm.onSubmit(actionsOnSubmit)}>
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
        <Button type="submit" loading={isLoading} loaderProps={{ type: "dots" }}>
          {"Добавить"}
        </Button>
      </form>
    </Modal>
  );
}
