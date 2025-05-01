import type { PropsWithChildren, ReactNode } from "react";
import { IconFileTypePdf, IconInfoCircle, IconTrashX } from "@tabler/icons-react";
import { Link, type To } from "react-router";
import {
  ActionIcon,
  Anchor,
  Card,
  Group,
  rem,
  Skeleton,
  Stack,
  Title,
  useMantineTheme,
} from "@mantine/core";

const ICON_SIZE = "2.5rem";
const TAG_SIZE = {
  width: 140,
  height: 30,
};
const BUTTON_SIZE = "3rem";

export interface EntityCardProps {
  title: string;
  to: To;
  tag?: ReactNode;
  onExport?: () => void;
  onDelete?: () => void;
}

export interface EntityCardSkeletonProps {
  tagged?: boolean;
  exported?: boolean;
  deleted?: boolean;
}

export function EntityCard({
  title,
  to,
  tag,
  onExport,
  onDelete,
  children,
}: PropsWithChildren<EntityCardProps>) {
  return (
    <Card>
      <Group gap="md">
        <IconInfoCircle size={ICON_SIZE} />
        <Stack flex={1} gap="xs" align="flex-start">
          <Anchor component={Link} to={to}>
            <Title order={3}>{title}</Title>
          </Anchor>
          {children && <Stack gap={0}>{children}</Stack>}
        </Stack>
        <Group gap="md" wrap="nowrap">
          {tag}
          {onExport && (
            <ActionIcon onClick={onExport}>
              <IconFileTypePdf />
            </ActionIcon>
          )}
          {onDelete && (
            <ActionIcon onClick={onDelete}>
              <IconTrashX />
            </ActionIcon>
          )}
        </Group>
      </Group>
    </Card>
  );
}

export function EntityCardSkeleton({
  tagged = false,
  exported = false,
  deleted = false,
  children,
}: PropsWithChildren<EntityCardSkeletonProps>) {
  const theme = useMantineTheme();
  return (
    <Card>
      <Group>
        <Skeleton circle height={ICON_SIZE} />
        <Stack flex={1} gap="md">
          <Skeleton height={rem(theme.headings.sizes.h3.fontSize)} width={200}>
            <Title>{"placeholder"}</Title>
          </Skeleton>
          {children && <Stack gap="sm">{children}</Stack>}
        </Stack>
        {tagged && <Skeleton width={TAG_SIZE.width} height={TAG_SIZE.height} />}
        {exported && <Skeleton circle height={BUTTON_SIZE} />}
        {deleted && <Skeleton circle height={BUTTON_SIZE} />}
      </Group>
    </Card>
  );
}
