import { IconArrowLeft, IconCircleDashedCheck, IconPlus, IconRefresh } from "@tabler/icons-react";
import { ActionIcon, Group, LoadingOverlay, Stack, Title } from "@mantine/core";
import { TextButton } from "../buttons/TextButton";
import PrimaryCard from "../cards/PrimaryCard";

export default function MainBar({
  title,
  subTitle,
  onRefresh,
  onAdd,
  onBack,
  onEnd,
  loading,
  children,
}) {
  return (
    <PrimaryCard pos="relative">
      <LoadingOverlay visible={loading} />
      <Group>
        {onBack && (
          <ActionIcon onClick={onBack}>
            <IconArrowLeft />
          </ActionIcon>
        )}
        <Stack flex={1} gap="xs">
          {title && <Title order={2}>{title}</Title>}
          {subTitle && <Title order={3}>{subTitle}</Title>}
        </Stack>
        {children}
        {onRefresh && (
          <ActionIcon onClick={onRefresh}>
            <IconRefresh />
          </ActionIcon>
        )}
        {onAdd && (
          <ActionIcon onClick={onAdd}>
            <IconPlus />
          </ActionIcon>
        )}
        {onEnd && (
          <TextButton label="Завершить" leftSection={<IconCircleDashedCheck />} onClick={onEnd} />
        )}
      </Group>
    </PrimaryCard>
  );
}
