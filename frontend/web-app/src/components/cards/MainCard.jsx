import { IconArrowLeft, IconCheck, IconEdit, IconTrashX, IconX } from "@tabler/icons-react";
import { ActionIcon, Button, Group, LoadingOverlay, Skeleton, Stack } from "@mantine/core";

export function MainCard({
  onEdit,
  isEditing,
  isLoading,
  onEditSubmit,
  onEditCancel,
  onDelete,
  children,
}) {
  return (
    <Stack w={300} align="start" pos="relative">
      <LoadingOverlay visible={isLoading} />
      <form onSubmit={onEditSubmit}>
        {children}
        <Group w="100%">
          <Group flex={1}>
            {!isEditing ? (
              <ActionIcon onClick={onEdit}>
                <IconEdit />
              </ActionIcon>
            ) : (
              <>
                <ActionIcon type="submit">
                  <IconCheck />
                </ActionIcon>
                <ActionIcon onClick={onEditCancel}>
                  <IconX />
                </ActionIcon>
              </>
            )}
          </Group>
          <ActionIcon onClick={onDelete}>
            <IconTrashX />
          </ActionIcon>
        </Group>
      </form>
    </Stack>
  );
}

export function MainCardSkeleton({ children }) {
  return (
    <Stack w={300} align="start">
      <Skeleton visible w={200}>
        <Button>{"Назад"}</Button>
      </Skeleton>
      {children}
      <Group w="100%">
        <Group flex={1}>
          <Skeleton visible>
            <ActionIcon>
              <IconTrashX />
            </ActionIcon>
          </Skeleton>
          <Skeleton visible>
            <ActionIcon>
              <IconTrashX />
            </ActionIcon>
          </Skeleton>
        </Group>
        <Skeleton visible>
          <ActionIcon>
            <IconTrashX />
          </ActionIcon>
        </Skeleton>
      </Group>
    </Stack>
  );
}
