import {
  IconArrowLeft,
  IconCircleDashedCheck,
  IconFileTypePdf,
  IconRefresh,
} from "@tabler/icons-react";
import { ActionIcon, Group, LoadingOverlay, Stack, Title, Tooltip } from "@mantine/core";
import { TextButton } from "../buttons/TextButton";
import PrimaryCard from "../cards/PrimaryCard";

export default function NavigationBar({
  title,
  subTitle,
  onRefresh,
  onExport,
  onBack,
  onEnd,
  loading,
  children,
}) {
  return (
    <PrimaryCard pos="relative">
      <LoadingOverlay visible={loading} />
      <Group gap="md">
        {onBack && (
          <ActionIcon onClick={onBack}>
            <IconArrowLeft />
          </ActionIcon>
        )}
        <Stack flex={1} gap="xs">
          {title && <Title order={2}>{title}</Title>}
          {subTitle && <Title order={3}>{subTitle}</Title>}
        </Stack>
        {onEnd && (
          <Tooltip label="Завершить">
            <ActionIcon onClick={onEnd}>
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
          <ActionIcon onClick={onRefresh}>
            <IconRefresh />
          </ActionIcon>
        )}
        {children}
      </Group>
    </PrimaryCard>
  );
}
