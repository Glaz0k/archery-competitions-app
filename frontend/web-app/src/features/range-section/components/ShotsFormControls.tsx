import { IconCheck, IconCircleDashedCheck, IconEdit, IconX } from "@tabler/icons-react";
import { ActionIcon, Card, Group, Tooltip } from "@mantine/core";

export interface ShotsFormControlsProps {
  editing: boolean;
  active?: boolean;
  onEdit: () => unknown;
  onCancelEdit: () => unknown;
  onComplete?: () => unknown;
}

export function ShotsFormControls({
  editing,
  active,
  onEdit,
  onComplete,
  onCancelEdit,
}: ShotsFormControlsProps) {
  return (
    <Card p="xs">
      <Group gap="sm">
        {!editing ? (
          <>
            <Tooltip label="Редактировать">
              <ActionIcon onClick={onEdit}>
                <IconEdit />
              </ActionIcon>
            </Tooltip>
            {onComplete && (
              <Tooltip label="Завершить">
                <ActionIcon disabled={!active} onClick={onComplete}>
                  <IconCircleDashedCheck />
                </ActionIcon>
              </Tooltip>
            )}
          </>
        ) : (
          <>
            <Tooltip label="Отменить">
              <ActionIcon onClick={onCancelEdit}>
                <IconX />
              </ActionIcon>
            </Tooltip>
            <Tooltip label="Завершить">
              <ActionIcon type="submit">
                <IconCheck />
              </ActionIcon>
            </Tooltip>
          </>
        )}
      </Group>
    </Card>
  );
}
