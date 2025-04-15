import { Button, Modal, Stack, Text } from "@mantine/core";

export default function DeleteIndividualGroupModal({ isOpened, onClose, onConfirm, isLoading }) {
  return (
    <Modal opened={isOpened} onClose={onClose} title="Удаление индивидуальной группы" centered>
      <Stack align="end">
        <Text w="100%">
          {
            "Вы уверены, что хотите удалить индивидуальную группу. Вместе с ней удалится также вся информация, такая как список участников, квалификация и финальная сетка."
          }
        </Text>
        <Button loading={isLoading} onClick={onConfirm}>
          {"Удалить"}
        </Button>
      </Stack>
    </Modal>
  );
}
