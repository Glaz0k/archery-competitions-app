import { DeleteModal, type DeleteModalProps } from "../DeleteModal";

export type DeleteIndividualGroupModalProps = Omit<DeleteModalProps, "confirmationText" | "title">;

export function DeleteIndividualGroupModal(props: DeleteIndividualGroupModalProps) {
  return (
    <DeleteModal
      title="Удаление дивизиона"
      confirmationText="Вы уверены, что хотите удалить дивизион? Вместе с ним удалится вся связанная информация, список участников, квалификация и финальная сетка."
      {...props}
    />
  );
}
