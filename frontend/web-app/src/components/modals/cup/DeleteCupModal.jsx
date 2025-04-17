import { Button, Modal, Stack, Text } from "@mantine/core";
import DeletionModal from "../DeletionModal";

export default function DeleteCupModal({ isOpened, onClose, onConfirm, isLoading }) {
  return (
    <DeletionModal
      title="Удаление кубка"
      isOpened={isOpened}
      onClose={onClose}
      onConfirm={onConfirm}
      isLoading={isLoading}
    >
      <Text>
        {
          "Вы уверены, что хотите удалить кубок? Вместе с ним удалятся также все связанные соревнования и группы."
        }
      </Text>
    </DeletionModal>
  );
}
