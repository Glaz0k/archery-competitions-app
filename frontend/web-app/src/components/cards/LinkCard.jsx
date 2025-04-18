import { IconFileTypePdf, IconInfoCircle, IconTrashX } from "@tabler/icons-react";
import { Link } from "react-router";
import {
  ActionIcon,
  Anchor,
  Card,
  Group,
  rem,
  Skeleton,
  Stack,
  ThemeIcon,
  Title,
  useMantineTheme,
} from "@mantine/core";

export function LinkCard({ title, to, tag = null, onExport = null, onDelete = null, children }) {
  return (
    <Card>
      <Group gap="md">
        <LinkCardIcon />
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

const iconSize = "2.5rem";
const tagSize = {
  width: 140,
  height: 30,
};
const buttonSize = "3rem";

export function LinkCardSkeleton({
  isTagged = false,
  isExport = false,
  isDelete = false,
  children,
}) {
  const theme = useMantineTheme();
  return (
    <Card>
      <Group>
        <Skeleton circle height={iconSize} />
        <Stack flex={1} gap="md">
          <Skeleton height={rem(theme.headings.sizes.h3.fontSize)} width={200}>
            <Title>{"placeholder"}</Title>
          </Skeleton>
          {children && <Stack gap="sm">{children}</Stack>}
        </Stack>
        {isTagged && <Skeleton width={tagSize.width} height={tagSize.height} />}
        {isExport && <Skeleton circle height={buttonSize} />}
        {isDelete && <Skeleton circle height={buttonSize} />}
      </Group>
    </Card>
  );
}

function LinkCardIcon() {
  return <IconInfoCircle size={iconSize} />;
}
