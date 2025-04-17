import { Modal, Stack } from "@mantine/core";
import CancelButton from "../buttons/CancelButton";

export default function DeletionModal({
  title,
  isOpened,
  onClose,
  onConfirm,
  isLoading,
  children,
}) {
  return (
    <Modal opened={isOpened} onClose={onClose} title={title}>
      <Stack align="end" gap="lg">
        <Stack gap="md">{children}</Stack>
        <CancelButton label="Удалить" loading={isLoading} onClick={onConfirm} />
      </Stack>
    </Modal>
  );
}
