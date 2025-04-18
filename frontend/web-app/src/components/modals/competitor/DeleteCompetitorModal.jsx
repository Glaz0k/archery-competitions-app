import { Text } from "@mantine/core";
import DeletionModal from "../DeletionModal";

export default function DeleteCompetitorModal({ isOpened, onClose, onConfirm, isLoading }) {
  return (
    <DeletionModal
      title="Исключение участника"
      isOpened={isOpened}
      onClose={onClose}
      onConfirm={onConfirm}
      isLoading={isLoading}
    >
      <Text>
        {
          "Вы уверены, что хотите исключить участника? Он пропадёт из списка участников, но вы можете добавить его позже."
        }
      </Text>
    </DeletionModal>
  );
}
