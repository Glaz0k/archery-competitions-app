import type { PropsWithChildren } from "react";
import {
  IconArrowLeft,
  IconCircleDashedCheck,
  IconFileTypePdf,
  IconRefresh,
} from "@tabler/icons-react";
import { Link, type To } from "react-router";
import { ActionIcon, Group, LoadingOverlay, Stack, Title, Tooltip } from "@mantine/core";
import { ControlsCard } from "../cards/ControlsCard";

export interface TabsBarProps {
  title?: string;
  subtitle?: string;
  onRefresh?: () => unknown;
  onComplete?: () => unknown;
  onExport?: () => unknown;
  backTo?: To;
  loading?: boolean;
}

export function TabsBar({
  title,
  subtitle,
  backTo,
  onRefresh,
  onComplete,
  onExport,
  loading = false,
  children,
}: PropsWithChildren<TabsBarProps>) {
  return (
    <ControlsCard pos="relative">
      <LoadingOverlay visible={loading} />
      <Group>
        {backTo && (
          <Tooltip label="Назад">
            <ActionIcon component={Link} to={backTo}>
              <IconArrowLeft />
            </ActionIcon>
          </Tooltip>
        )}
        <Stack flex={1} gap={0}>
          {title ? <Title order={2}>{title}</Title> : "Загрузка..."}
          {subtitle && <Title order={3}>{subtitle}</Title>}
        </Stack>
        {onComplete && (
          <Tooltip label="Завершить">
            <ActionIcon onClick={onComplete}>
              <IconCircleDashedCheck />
            </ActionIcon>
          </Tooltip>
        )}
        {onExport && (
          <Tooltip label="Экспорт PDF">
            <ActionIcon onClick={onExport}>
              <IconFileTypePdf />
            </ActionIcon>
          </Tooltip>
        )}
        {onRefresh && (
          <Tooltip label="Обновить">
            <ActionIcon onClick={onRefresh}>
              <IconRefresh />
            </ActionIcon>
          </Tooltip>
        )}
        {children}
      </Group>
    </ControlsCard>
  );
}
