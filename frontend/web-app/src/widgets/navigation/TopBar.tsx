import type { PropsWithChildren } from "react";
import { IconArrowLeft, IconCircleDashedCheck, IconPlus, IconRefresh } from "@tabler/icons-react";
import { ActionIcon, Group, LoadingOverlay, Stack, Title, Tooltip } from "@mantine/core";
import { TextButton } from "../buttons/TextButton";
import { ControlsCard } from "../cards/ControlsCard";

export interface TopBarProps {
  title?: string;
  subtitle?: string;
  onRefresh?: () => unknown;
  onAdd?: () => unknown;
  onBack?: () => unknown;
  onComplete?: () => unknown;
  loading?: boolean;
}

export function TopBar({
  title,
  subtitle,
  onRefresh,
  onAdd,
  onBack,
  onComplete,
  loading = false,
  children,
}: PropsWithChildren<TopBarProps>) {
  return (
    <ControlsCard pos="relative">
      <LoadingOverlay visible={loading} />
      <Group>
        {onBack && (
          <Tooltip label="Назад">
            <ActionIcon onClick={onBack}>
              <IconArrowLeft />
            </ActionIcon>
          </Tooltip>
        )}
        <Stack flex={1} gap="xs">
          {title ? <Title order={2}>{title}</Title> : "Загрузка..."}
          {subtitle && <Title order={3}>{subtitle}</Title>}
        </Stack>
        {children}
        {onRefresh && (
          <Tooltip label="Обновить">
            <ActionIcon onClick={onRefresh}>
              <IconRefresh />
            </ActionIcon>
          </Tooltip>
        )}
        {onAdd && (
          <Tooltip label="Добавить">
            <ActionIcon onClick={onAdd}>
              <IconPlus />
            </ActionIcon>
          </Tooltip>
        )}
        {onComplete && (
          <TextButton
            label="Завершить"
            leftSection={<IconCircleDashedCheck />}
            onClick={onComplete}
          />
        )}
      </Group>
    </ControlsCard>
  );
}
