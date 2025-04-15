import { Button, Modal, Stack, Text } from "@mantine/core";

export default function DeleteCompetitionModal({ isOpened, onClose, onConfirm, isLoading }) {
  return (
    <Modal opened={isOpened} onClose={onClose} title="Удаление соревнования" centered>
      <Stack align="end">
        <Text w="100%">
          Вы уверены, что хотите удалить соревнование. Вместе с ним удалятся также все связанные
          группы.
        </Text>
        <Button loading={isLoading} onClick={onConfirm}>
          Удалить
        </Button>
      </Stack>
    </Modal>
  );
}
