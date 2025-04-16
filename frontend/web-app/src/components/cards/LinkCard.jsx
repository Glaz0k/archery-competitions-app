import { IconFileTypePdf, IconInfoCircle, IconTrashX } from "@tabler/icons-react";
import { Link } from "react-router";
import {
  ActionIcon,
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
      <Group>
        <ThemeIcon>
          <IconInfoCircle />
        </ThemeIcon>
        <Stack flex={1}>
          <Title order={3}>
            <Link to={to}>{title}</Link>
          </Title>
          {children && <Stack>{children}</Stack>}
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

const iconSize = 40;
const tagSize = {
  width: 140,
  height: 30,
};

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
        <Stack flex={1}>
          <Skeleton height={rem(theme.headings.sizes.h3.fontSize)} width={200} />
          {children && <Stack>{children}</Stack>}
        </Stack>
        {isTagged && <Skeleton width={tagSize.width} height={tagSize.height} />}
        {isExport && <Skeleton circle height={iconSize} />}
        {isDelete && <Skeleton circle height={iconSize} />}
      </Group>
    </Card>
  );
}
