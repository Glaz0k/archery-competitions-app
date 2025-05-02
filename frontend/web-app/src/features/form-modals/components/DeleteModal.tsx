import { IconTrashX } from "@tabler/icons-react";
import { Modal, Stack, Text, type ModalProps } from "@mantine/core";
import { TextButton } from "../../../widgets";

export interface DeleteModalProps extends Omit<ModalProps, "children"> {
  confirmationText: string;
  onConfirm: () => void;
  loading: boolean;
}

export function DeleteModal({ confirmationText, onConfirm, loading, ...others }: DeleteModalProps) {
  return (
    <Modal {...others}>
      <Stack align="end" gap="lg">
        <Text w="100%">{confirmationText}</Text>
        <TextButton
          label="Удалить"
          variant="filled"
          color="red.6"
          loading={loading}
          onClick={onConfirm}
          rightSection={<IconTrashX />}
        />
      </Stack>
    </Modal>
  );
}
