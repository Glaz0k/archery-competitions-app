import { Text } from "@mantine/core";
import CompetitionStage from "../../../enums/CompetitionStage";
import DeletionModal from "../DeletionModal";

export default function DeleteCompetitionModal({ isOpened, onClose, onConfirm, isLoading }) {
  return (
    <DeletionModal
      title="Удаление соревнования"
      isOpened={isOpened}
      onClose={onClose}
      onConfirm={onConfirm}
      isLoading={isLoading}
    >
      <Text>
        {
          "Вы уверены, что хотите удалить соревнование? Вместе с ним удалятся также все связанные группы."
        }
      </Text>
    </DeletionModal>
  );
}
