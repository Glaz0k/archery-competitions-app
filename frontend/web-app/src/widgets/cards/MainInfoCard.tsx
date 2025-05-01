import type { PropsWithChildren } from "react";
import { IconCheck, IconEdit, IconTrashX, IconX } from "@tabler/icons-react";
import { ActionIcon, Group, LoadingOverlay, Stack } from "@mantine/core";
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
        <Stack w={300} align="start" pos="relative" gap="md">
          {children}
          <Group w="100%">
            <Group flex={1}>
              {!editing ? (
                <ActionIcon onClick={onEdit}>
                  <IconEdit />
                </ActionIcon>
              ) : (
                <>
                  <ActionIcon type="submit">
                    <IconCheck />
                  </ActionIcon>
                  <ActionIcon onClick={onCancel}>
                    <IconX />
                  </ActionIcon>
                </>
              )}
            </Group>
            <ActionIcon onClick={onDelete}>
              <IconTrashX />
            </ActionIcon>
          </Group>
        </Stack>
      </form>
    </ControlsCard>
  );
}
