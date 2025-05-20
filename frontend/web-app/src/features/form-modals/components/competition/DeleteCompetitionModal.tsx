import { DeleteModal, type DeleteModalProps } from "../DeleteModal";

export type DeleteCompetitionModalProps = Omit<DeleteModalProps, "confirmationText" | "title">;

export function DeleteCompetitionModal(props: DeleteCompetitionModalProps) {
  return (
    <DeleteModal
      title="Удаление соревнования"
      confirmationText="Вы уверены, что хотите удалить соревнование? Вместе с ним удалятся также все связанные дивизионы."
      {...props}
    />
  );
}
