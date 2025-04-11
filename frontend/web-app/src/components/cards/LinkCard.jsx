import {
  ActionIcon,
  Card,
  Group,
  Stack,
  Skeleton,
  ThemeIcon,
  Title,
  useMantineTheme,
  rem,
} from "@mantine/core";
import {
  IconInfoCircle,
  IconTrashX,
  IconFileTypePdf,
} from "@tabler/icons-react";
import { Link } from "react-router";

export function LinkCard({
  key,
  title,
  to,
  tag = null,
  onExport = null,
  onDelete = null,
  children,
}) {
  return (
    <Card key={key}>
      <Group>
        <ThemeIcon>
          <IconInfoCircle />
        </ThemeIcon>
        <Stack flex={1}>
          <Link to={to}>
            <Title order={3}>{title}</Title>
          </Link>
          <Stack>{children}</Stack>
        </Stack>
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
    </Card>
  );
}

export function LinkCardSkeleton({
  key,
  isExport = false,
  isDelete = false,
  children,
}) {
  const theme = useMantineTheme();

  const iconSize = 40;

  return (
    <Card key={"skeleton-" + key}>
      <Group>
        <Skeleton circle height={iconSize} />
        <Stack flex={1}>
          <Skeleton
            height={rem(theme.headings.sizes.h3.fontSize)}
            width={200}
          />
          <Stack>{children}</Stack>
        </Stack>
        {isExport && <Skeleton circle height={iconSize} />}
        {isDelete && <Skeleton circle height={iconSize} />}
      </Group>
    </Card>
  );
}
