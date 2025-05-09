import { IconRefresh } from "@tabler/icons-react";
import {
  ActionIcon,
  Card,
  Group,
  LoadingOverlay,
  Stack,
  Text,
  Tooltip,
  useMantineTheme,
} from "@mantine/core";
import { useSparringPlace } from "../../../api";
import { type Sparring, type SparringPlace } from "../../../entities";
import { ControlsCard, PageLoader } from "../../../widgets";
import { PlaceTab } from "./PlaceTab";

export function SparringTab({ sparring }: { sparring: Sparring }) {
  if (!sparring.top || !sparring.bot) {
    throw new Error("Selected sparring must not be empty");
  }

  const {
    data: topPlace,
    isLoading: isTopLoading,
    isFetching: isTopFetching,
    isError: isTopError,
    refetch: refreshTop,
  } = useSparringPlace(sparring.top.id);

  const {
    data: botPlace,
    isLoading: isBotLoading,
    isFetching: isBotFetching,
    isError: isBotError,
    refetch: refreshBot,
  } = useSparringPlace(sparring.bot.id);

  const isAnyLoading = isTopLoading || isBotLoading;
  const isAnyFetching = isTopFetching || isBotFetching;
  const isAnyError = isTopError || isBotError;

  return (
    <PageLoader loading={isAnyLoading} error={isAnyError}>
      {topPlace && botPlace && (
        <Stack>
          <SparringBar
            left={topPlace}
            right={botPlace}
            loading={isAnyFetching}
            refreshFn={() => {
              refreshTop();
              refreshBot();
            }}
          />
          <Group flex={1} grow>
            <PlaceTab place={topPlace} />
            <PlaceTab place={botPlace} />
          </Group>
        </Stack>
      )}
    </PageLoader>
  );
}

interface SparringBarProps {
  left: SparringPlace;
  right: SparringPlace;
  loading: boolean;
  refreshFn: () => unknown;
}

function SparringBar({ left, right, loading, refreshFn }: SparringBarProps) {
  const theme = useMantineTheme();
  return (
    <ControlsCard pos="relative">
      <LoadingOverlay visible={loading} />
      <Group flex={1}>
        <Group flex={1} justify="flex-end">
          <Text>{`${left.competitor.fullName} - ${left.score}`}</Text>
        </Group>
        <Card bg={theme.white} w={5} p={0} />
        <Group flex={1} justify="flex-start">
          <Group flex={1} justify="flex-start">
            <Text>{`${right.score} - ${right.competitor.fullName}`}</Text>
          </Group>
          <Tooltip label="Обновить">
            <ActionIcon onClick={refreshFn}>
              <IconRefresh />
            </ActionIcon>
          </Tooltip>
        </Group>
      </Group>
    </ControlsCard>
  );
}
