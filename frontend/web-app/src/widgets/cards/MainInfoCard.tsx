import type { PropsWithChildren } from "react";
import { IconCheck, IconEdit, IconTrashX, IconX } from "@tabler/icons-react";
import { ActionIcon, Group, LoadingOverlay, Stack, Tooltip } from "@mantine/core";
import { ControlsCard } from "./ControlsCard";

export interface MainInfoCardProps {
  onEdit: () => unknown;
  onFormSubmit: () => unknown;
  onCancel: () => unknown;
  onDelete: () => unknown;
  editing?: boolean;
  loading?: boolean;
}

export function MainInfoCard({
  onEdit,
  onFormSubmit,
  onCancel,
  onDelete,
  editing = false,
  loading = false,
  children,
}: PropsWithChildren<MainInfoCardProps>) {
  return (
    <ControlsCard>
      <LoadingOverlay visible={loading} />
      <form onSubmit={onFormSubmit}>
        <Stack align="start" pos="relative" gap="md">
          {children}
          <Group w="100%">
            <Group flex={1}>
              {!editing ? (
                <Tooltip label="Редактировать">
                  <ActionIcon onClick={onEdit}>
                    <IconEdit />
                  </ActionIcon>
                </Tooltip>
              ) : (
                <>
                  <Tooltip label="Подтвердить">
                    <ActionIcon type="submit">
                      <IconCheck />
                    </ActionIcon>
                  </Tooltip>
                  <Tooltip label="Отменить">
                    <ActionIcon onClick={onCancel}>
                      <IconX />
                    </ActionIcon>
                  </Tooltip>
                </>
              )}
            </Group>
            {!editing && (
              <Tooltip label="Удалить">
                <ActionIcon onClick={onDelete}>
                  <IconTrashX />
                </ActionIcon>
              </Tooltip>
            )}
          </Group>
        </Stack>
      </form>
    </ControlsCard>
  );
}
