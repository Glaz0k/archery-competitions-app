import { Modal, type ModalProps } from "@mantine/core";

export function RegisterCompetitorModal(props: Omit<ModalProps, "title">) {
  return <Modal title="Зарегистрировать нового участника" {...props}></Modal>;
}
