import { useEffect, useRef, useState } from "react";
import { IconPlayerPlay } from "@tabler/icons-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  ActionIcon,
  Box,
  Card,
  Center,
  Flex,
  Group,
  LoadingOverlay,
  rem,
  ScrollArea,
  Stack,
  Title,
  Tooltip,
  useMantineTheme,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {
  endFinal,
  getFinalGrid,
  startFinal,
  startQuarterfinal,
  startSemifinal,
} from "../../../api/individualGroups";
import { INDIVIDUAL_GROUP_QUERY_KEYS } from "../../../api/queryKeys";
import GroupState from "../../../enums/GroupState";
import { TextButton } from "../../buttons/TextButton";
import SparringCard from "../../cards/SparringCard";
import FinalSection from "../../containers/FinalSection";
import ConfirmationModal from "../../modals/ConfirmationModal";
import { FinalContext } from "./FinalContext";

const GAP_BETWEEN_COLUMNS = 80;
const GAP_BETWEEN_SPARRINGS = 60;
const GAP_HEADER = 20;

const GAP_BETWEEN_COLUMNS_REM = rem(GAP_BETWEEN_COLUMNS);
const GAP_BETWEEN_SPARRINGS_REM = rem(GAP_BETWEEN_SPARRINGS);
const GAP_HEADER_REM = rem(GAP_HEADER);

const sparringSize = {
  heigth: 60,
  width: 420,
};

export default function FinalPanel({ groupInfo, setGroupControls }) {
  const theme = useMantineTheme();
  const queryClient = useQueryClient();
  const group = groupInfo?.group;
  const [selectedPlaceId, setSelectedPlaceId] = useState(null);

  const [isConfirmationOpened, confirmationControl] = useDisclosure(false);
  const openConfirmationRef = useRef(confirmationControl.open);
  const [confirmationTitle, setConfirmationTitle] = useState("");
  const [confirmationText, setConfirmationText] = useState("");
  const [onConfirmFn, setOnConfirmFn] = useState(undefined);

  const {
    data: grid,
    isFetching: isGridFetching,
    refetch: refreshGrid,
  } = useQuery({
    queryKey: INDIVIDUAL_GROUP_QUERY_KEYS.finalGrid(group?.id),
    queryFn: () => getFinalGrid(group?.id),
    enabled: group?.id != null,
  });

  const { isPending: isQuarterfinalStarting, mutate: beginQuarterfinal } = useMutation({
    mutationFn: () => startQuarterfinal(group?.id),
    onSuccess: (newGrid) => {
      queryClient.setQueryData(INDIVIDUAL_GROUP_QUERY_KEYS.finalGrid(group?.id), newGrid);
      confirmationControl.close();
    },
  });

  const { isPending: isSemifinalStarting, mutate: beginSemifinal } = useMutation({
    mutationFn: () => startSemifinal(group?.id),
    onSuccess: (newGrid) => {
      queryClient.setQueryData(INDIVIDUAL_GROUP_QUERY_KEYS.finalGrid(group?.id), newGrid);
      confirmationControl.close();
    },
  });

  const { isPending: isFinalStarting, mutate: beginFinal } = useMutation({
    mutationFn: () => startFinal(group?.id),
    onSuccess: (newGrid) => {
      queryClient.setQueryData(INDIVIDUAL_GROUP_QUERY_KEYS.finalGrid(group?.id), newGrid);
      confirmationControl.close();
    },
  });

  const { isPending: isFinalEnding, mutate: completeFinal } = useMutation({
    mutationFn: () => endFinal(group?.id),
    onSuccess: (newGrid) => {
      queryClient.setQueryData(INDIVIDUAL_GROUP_QUERY_KEYS.finalGrid(group?.id), newGrid);
      confirmationControl.close();
    },
  });

  const isGridLoading =
    isGridFetching ||
    isQuarterfinalStarting ||
    isSemifinalStarting ||
    isFinalStarting ||
    isFinalEnding;

  const onQuarterFinalStart = () => {
    setConfirmationTitle("Начать четвертьфинал");
    setConfirmationText("Вы уверены что хотите начать четвертьфинал?");
    setOnConfirmFn(() => () => beginQuarterfinal());
    confirmationControl.open();
  };

  const onSemiFinalStart = () => {
    setConfirmationTitle("Начать полуфинал");
    setConfirmationText("Вы уверены что хотите закончить четвертьфинал и начать полуфинал?");
    setOnConfirmFn(() => () => beginSemifinal());
    confirmationControl.open();
  };

  const onFinalStart = () => {
    setConfirmationTitle("Начать финал");
    setConfirmationText("Вы уверены что хотите закончить полуфинал и начать финал?");
    setOnConfirmFn(() => () => beginFinal());
    confirmationControl.open();
  };

  useEffect(() => {
    const onFinalEnd = () => {
      setConfirmationTitle("Закончить финал");
      setConfirmationText("Вы уверены что хотите закончить финал?");
      setOnConfirmFn(() => () => completeFinal());
      openConfirmationRef.current();
    };

    setGroupControls({
      onExport: () => {
        console.warn("final export temporary unavailable");
      },
      onRefresh: () => {
        setSelectedPlaceId(null);
        refreshGrid();
      },
      onEnd: group?.state.value === GroupState.FINAL_START.value ? onFinalEnd : undefined,
    });
  }, [setGroupControls, refreshGrid, group, completeFinal]);

  if (grid == null) {
    if (isGridFetching) {
      return (
        <Center flex={1}>
          <Card>
            <Title order={1}>{"Загрузка..."}</Title>
          </Card>
        </Center>
      );
    }
    if (group?.state.value === GroupState.QUAL_END.value) {
      return (
        <Center flex={1}>
          <Card>
            <TextButton label="Начать четвертьфинал" onClick={onQuarterFinalStart} />
          </Card>
        </Center>
      );
    }
    return (
      <Center flex={1}>
        <Card>
          <Title order={1}>{"Финальная сетка не найдена"}</Title>
        </Card>
      </Center>
    );
  }

  return (
    <>
      <ConfirmationModal
        title={confirmationTitle}
        text={confirmationText}
        opened={isConfirmationOpened}
        onConfirm={onConfirmFn}
        onClose={confirmationControl.close}
        loading={isGridLoading}
      />
      <Stack gap="md" display="flex" flex={1} style={{ overflow: "hidden" }}>
        <FinalContext.Provider
          value={{
            selectedPlaceId: selectedPlaceId,
            setSelectedPlaceId: setSelectedPlaceId,
            sparringSize: sparringSize,
          }}
        >
          <ScrollArea flex={1} style={{ overflow: "hidden" }} miw={rem(500)} offsetScrollbars>
            <LoadingOverlay visible={isGridLoading} />
            <Stack>
              <Group gap={GAP_BETWEEN_COLUMNS_REM} align="stretch" justify="center" wrap="nowrap">
                <Stack gap={0} w={sparringSize.width}>
                  <Center>
                    <Title order={2}>{"1/4 Финала"}</Title>
                  </Center>
                  <HeaderDivider />
                </Stack>
                <Group gap="sm" w={sparringSize.width}>
                  {group?.state.value === GroupState.QUARTERFINAL_START.value && (
                    <Tooltip label={"Начать полуфинал"}>
                      <ActionIcon onClick={onSemiFinalStart} color={theme.colors.secondary[9]}>
                        <IconPlayerPlay />
                      </ActionIcon>
                    </Tooltip>
                  )}
                  <Stack gap={0} flex={1}>
                    <Center>
                      <Title order={2}>{"1/2 Финала"}</Title>
                    </Center>
                    <HeaderDivider />
                  </Stack>
                </Group>
                <Group gap="sm" w={sparringSize.width}>
                  {group?.state.value === GroupState.SEMIFINAL_START.value && (
                    <Tooltip label={"Начать финал"}>
                      <ActionIcon onClick={onFinalStart} color={theme.colors.secondary[9]}>
                        <IconPlayerPlay />
                      </ActionIcon>
                    </Tooltip>
                  )}
                  <Stack gap={0} flex={1}>
                    <Center>
                      <Title order={2}>{"Финал"}</Title>
                    </Center>
                    <HeaderDivider />
                  </Stack>
                </Group>
              </Group>
              <Group gap={0} align="stretch" justify="center" wrap="nowrap">
                <Stack gap={GAP_BETWEEN_SPARRINGS_REM}>
                  <Group gap={0} wrap="inherit">
                    <Stack gap={GAP_HEADER_REM}>
                      <Stack h="100%" gap={GAP_BETWEEN_SPARRINGS_REM}>
                        <SparringCard sparring={grid?.quarterfinal?.sparring1} ordinal={1} />
                        <SparringCard sparring={grid?.quarterfinal?.sparring2} ordinal={2} />
                      </Stack>
                    </Stack>
                    <MergeDivider
                      heigth={rem(GAP_BETWEEN_SPARRINGS + sparringSize.heigth + 2)}
                      width={GAP_BETWEEN_COLUMNS_REM}
                    />
                    <Stack gap={GAP_HEADER_REM}>
                      <Center h="100%">
                        <SparringCard sparring={grid?.semifinal?.sparring5} />
                      </Center>
                    </Stack>
                  </Group>
                  <Group gap={0} wrap="inherit">
                    <Stack h="100%" gap={GAP_BETWEEN_SPARRINGS_REM}>
                      <SparringCard sparring={grid?.quarterfinal?.sparring3} ordinal={3} />
                      <SparringCard sparring={grid?.quarterfinal?.sparring4} ordinal={4} />
                    </Stack>
                    <MergeDivider
                      heigth={rem(GAP_BETWEEN_SPARRINGS + sparringSize.heigth + 2)}
                      width={GAP_BETWEEN_COLUMNS_REM}
                    />
                    <Center h="100%">
                      <SparringCard sparring={grid?.semifinal?.sparring6} />
                    </Center>
                  </Group>
                </Stack>
                <Stack gap={0} justify="center">
                  <MergeDivider
                    heigth={rem(GAP_BETWEEN_SPARRINGS * 3 + sparringSize.heigth + 2)}
                    width={GAP_BETWEEN_COLUMNS_REM}
                  />
                </Stack>
                <Stack gap={0}>
                  <Stack flex={1} />
                  <Stack gap={GAP_HEADER_REM}>
                    <SparringCard sparring={grid?.final?.sparringGold} />
                  </Stack>
                  <Stack flex={1} gap={GAP_HEADER_REM} justify="flex-end">
                    <Stack gap={0}>
                      <Center>
                        <Title order={2}>{"Финал за 3-е место"}</Title>
                      </Center>
                      <HeaderDivider />
                    </Stack>
                    <SparringCard sparring={grid?.final?.sparringBronze} />
                  </Stack>
                </Stack>
              </Group>
            </Stack>
          </ScrollArea>
          {selectedPlaceId != null && <FinalSection />}
        </FinalContext.Provider>
      </Stack>
    </>
  );
}

function HeaderDivider() {
  return <Card h={rem(5)} p={0} />;
}

function MergeDivider({ heigth, width }) {
  const theme = useMantineTheme();
  return (
    <Flex h={heigth} w={width} align="flex-end" direction="row" wrap="nowrap">
      <Box
        h="100%"
        w="50%"
        style={{
          borderTop: `solid 3px ${theme.colors.secondary[9]}`,
          borderRight: `solid 3px ${theme.colors.secondary[9]}`,
          borderBottom: `solid 3px ${theme.colors.secondary[9]}`,
          borderStartEndRadius: theme.radius.md,
          borderEndEndRadius: theme.radius.md,
        }}
      />
      <Box
        h="51%"
        w="50%"
        style={{
          borderTop: `solid 3px ${theme.colors.secondary[9]}`,
        }}
      />
    </Flex>
  );
}
