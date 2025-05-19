import type { ReactNode } from "react";
import { Group, Modal, Stack } from "@mantine/core";
import type { UseFormReturnType } from "@mantine/form";
import { CancelButton, ConfirmButton } from "../../../widgets";

export interface AddModalProps<V, T> {
  form: UseFormReturnType<V, (values: V) => T>;
  title: string;
  opened: boolean;
  onClose: () => void;
  onSubmit: (values: T) => Promise<unknown>;
  loading: boolean;
  children: ReactNode;
}

export function AddModal<V, T>({
  form,
  title,
  opened,
  onClose,
  onSubmit,
  loading,
  children,
}: AddModalProps<V, T>) {
  const actionsOnClose = () => {
    form.reset();
    onClose();
  };

  const actionsOnSubmit = async (values: T) => {
    if (await onSubmit(values)) {
      form.reset();
    }
  };

  return (
    <Modal opened={opened} onClose={actionsOnClose} title={title}>
      <form onSubmit={form.onSubmit(actionsOnSubmit)}>
        <Stack gap="lg">
          <Stack gap="md">{children}</Stack>
          <Group gap="md" justify="flex-end">
            <CancelButton size="sm" label="Отменить" onClick={actionsOnClose} loading={loading} />
            <ConfirmButton size="sm" label="Добавить" type="submit" loading={loading} />
          </Group>
        </Stack>
      </form>
    </Modal>
  );
}
