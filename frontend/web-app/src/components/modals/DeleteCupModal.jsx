import { Button, Modal, Stack, Text } from "@mantine/core";

export default function DeleteCupModal({ isOpened, onClose, onConfirm, isLoading }) {
  return (
    <Modal opened={isOpened} onClose={onClose} title="Удаление кубка" centered>
      <Stack align="end">
        <Text w="100%">
          Вы уверены, что хотите удалить кубок. Вместе с ним удалятся также все связанные
          соревнования и группы.
        </Text>
        <Button loading={isLoading} onClick={onConfirm}>
          Удалить
        </Button>
      </Stack>
    </Modal>
  );
}
