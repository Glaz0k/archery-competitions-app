import { DeleteModal, type DeleteModalProps } from "../DeleteModal";

export type DeleteCupModalProps = Omit<DeleteModalProps, "confirmationText" | "title">;

export function DeleteCupModal(props: DeleteCupModalProps) {
  return (
    <DeleteModal
      title="Удаление кубка"
      confirmationText="Вы уверены, что хотите удалить кубок? Вместе с ним удалятся также все связанные соревнования и дивизионы."
      {...props}
    />
  );
}
