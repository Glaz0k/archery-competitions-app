import { DeleteModal, type DeleteModalProps } from "../DeleteModal";

export type RemoveCompetitorModalProps = Omit<DeleteModalProps, "confirmationText" | "title">;

export function RemoveCompetitorModal(props: RemoveCompetitorModalProps) {
  return (
    <DeleteModal
      title="Исключение участника"
      confirmationText="Вы уверены, что хотите исключить участника? Вы можете добавить его вручную позже."
      {...props}
    />
  );
}
