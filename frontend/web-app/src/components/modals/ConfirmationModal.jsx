import { Modal, Stack, Text } from "@mantine/core";
import { ConfirmButton } from "../../widgets";

export default function ConfirmationModal({ title, text, opened, onConfirm, onClose, loading }) {
  return (
    <Modal opened={opened} onClose={onClose} title={title}>
      <Stack align="end" gap="lg">
        <Stack w="100%">
          <Text>{text}</Text>
        </Stack>
        <ConfirmButton size="sm" label="Подтвердить" loading={loading} onClick={onConfirm} />
      </Stack>
    </Modal>
  );
}
