import { Text } from "@mantine/core";
import DeletionModal from "../DeletionModal";

export default function DeleteIndividualGroupModal({ isOpened, onClose, onConfirm, isLoading }) {
  return (
    <DeletionModal
      title="Удаление индивидуальной группы"
      isOpened={isOpened}
      onClose={onClose}
      onConfirm={onConfirm}
      isLoading={isLoading}
    >
      <Text>
        {
          "Вы уверены, что хотите удалить индивидуальную группу? Вместе с ней удалится вся связанная информация, такая как список участников, квалификация и финальная сетка."
        }
      </Text>
    </DeletionModal>
  );
}
