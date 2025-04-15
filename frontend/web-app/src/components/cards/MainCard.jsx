import { IconArrowLeft, IconCheck, IconEdit, IconTrashX, IconX } from "@tabler/icons-react";
import { ActionIcon, Button, Group, LoadingOverlay, Stack } from "@mantine/core";

export default function MainCard({
  onBack,
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
      <Button onClick={onBack} leftSection={<IconArrowLeft />}>
        {"Назад"}
      </Button>
      {children}
      <Group w="100%">
        <Group flex={1}>
          {!isEditing ? (
            <ActionIcon onClick={onEdit}>
              <IconEdit />
            </ActionIcon>
          ) : (
            <>
              <ActionIcon onClick={onEditSubmit}>
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
    </Stack>
  );
}
