import { Group, Modal, Stack } from "@mantine/core";
import CancelButton from "../buttons/CancelButton";
import ConfirmButton from "../buttons/ConfirmButton";

export default function AddFormModal({
  addForm,
  title,
  isOpened,
  onClose,
  onSubmit,
  isLoading,
  children,
}) {
  const actionsOnClose = () => {
    addForm.reset();
    onClose?.();
  };

  const actionsOnSubmit = async (addFormValues) => {
    if (await onSubmit(addFormValues)) {
      addForm.reset();
    }
  };

  return (
    <Modal opened={isOpened} onClose={actionsOnClose} title={title}>
      <form onSubmit={addForm.onSubmit(actionsOnSubmit)}>
        <Stack gap="lg">
          <Stack gap="md">{children}</Stack>
          <Group gap="md" justify="flex-end">
            <CancelButton size="sm" label="Отменить" onClick={actionsOnClose} loading={isLoading} />
            <ConfirmButton size="sm" label="Добавить" type="submit" loading={isLoading} />
          </Group>
        </Stack>
      </form>
    </Modal>
  );
}
