import {
  Box,
  Card,
  Center,
  Flex,
  Group,
  LoadingOverlay,
  ScrollArea,
  Stack,
  Title,
  useMantineTheme,
} from "@mantine/core";
import type { FinalGrid, Sparring } from "../../../entities";
import {
  GAP_BETWEEN_COLUMNS,
  GAP_BETWEEN_SPARRINGS,
  GAP_HEADER,
  SPARRING_SIZE,
} from "../constants";
import { SparringView } from "./SparringView";

export interface FinalGridViewProps {
  grid: FinalGrid;
  loading: boolean;
  selectSparring: (sparring: Sparring) => unknown;
}

export function FinalGridView({ grid, loading, selectSparring }: FinalGridViewProps) {
  return (
    <ScrollArea flex={1} style={{ overflow: "hidden" }} miw={500} offsetScrollbars>
      <LoadingOverlay visible={loading} />
      <Stack>
        <Group gap={GAP_BETWEEN_COLUMNS} align="stretch" justify="center" wrap="nowrap">
          <FinalHeader title="1/4 Финала" />
          <FinalHeader title="1/2 Финала" />
          <FinalHeader title="Финал" />
        </Group>
        <Group gap={0} align="stretch" justify="center" wrap="nowrap">
          <Stack gap={GAP_BETWEEN_SPARRINGS}>
            <Group gap={0} wrap="inherit">
              <Stack gap={GAP_HEADER}>
                <Stack h="100%" gap={GAP_BETWEEN_SPARRINGS}>
                  <SparringView
                    sparring={grid.quarterfinal.sparring1}
                    ordinal={1}
                    selectSparring={selectSparring}
                  />
                  <SparringView
                    sparring={grid.quarterfinal.sparring2}
                    ordinal={2}
                    selectSparring={selectSparring}
                  />
                </Stack>
                <GridMergeDivider
                  w={GAP_BETWEEN_COLUMNS}
                  h={GAP_BETWEEN_SPARRINGS + SPARRING_SIZE.heigth + 2}
                />
                <Stack gap={GAP_HEADER}>
                  <Center h="100%">
                    <SparringView
                      sparring={grid.semifinal?.sparring5}
                      selectSparring={selectSparring}
                    />
                  </Center>
                </Stack>
              </Stack>
            </Group>
            <Group gap={0} wrap="inherit">
              <Stack h="100%" gap={GAP_BETWEEN_SPARRINGS}>
                <SparringView
                  sparring={grid.quarterfinal.sparring3}
                  ordinal={3}
                  selectSparring={selectSparring}
                />
                <SparringView
                  sparring={grid.quarterfinal.sparring4}
                  ordinal={4}
                  selectSparring={selectSparring}
                />
              </Stack>
              <GridMergeDivider
                h={GAP_BETWEEN_SPARRINGS + SPARRING_SIZE.heigth + 2}
                w={GAP_BETWEEN_COLUMNS}
              />
              <Center h="100%">
                <SparringView
                  sparring={grid.semifinal?.sparring6}
                  selectSparring={selectSparring}
                />
              </Center>
            </Group>
          </Stack>
          <Stack gap={0} justify="center">
            <GridMergeDivider
              h={GAP_BETWEEN_SPARRINGS * 3 + SPARRING_SIZE.heigth + 2}
              w={GAP_BETWEEN_COLUMNS}
            />
          </Stack>
          <Stack gap={0}>
            <Stack flex={1} />
            <Stack gap={GAP_HEADER}>
              <SparringView sparring={grid.final?.sparringGold} selectSparring={selectSparring} />
            </Stack>
            <Stack flex={1} gap={GAP_HEADER} justify="flex-end">
              <FinalHeader title="Финал за 3-е место" />
              <SparringView sparring={grid.final?.sparringBronze} selectSparring={selectSparring} />
            </Stack>
          </Stack>
        </Group>
      </Stack>
    </ScrollArea>
  );
}

function FinalHeader({ title }: { title: string }) {
  return (
    <Stack gap={0} w={SPARRING_SIZE.width}>
      <Center>
        <Title order={2}>{title}</Title>
      </Center>
      <Card h={5} p={0} />
    </Stack>
  );
}

function GridMergeDivider({ w, h }: { w: number; h: number }) {
  const theme = useMantineTheme();
  return (
    <Flex w={w} h={h} align="flex-end" direction="row" wrap="nowrap">
      <Box
        h="100%"
        w="50%"
        style={{
          borderTop: `solid 3px ${theme.colors.secondary![9]}`,
          borderRight: `solid 3px ${theme.colors.secondary![9]}`,
          borderBottom: `solid 3px ${theme.colors.secondary![9]}`,
          borderStartEndRadius: theme.radius.md,
          borderEndEndRadius: theme.radius.md,
        }}
      />
      <Box
        h="51%"
        w="50%"
        style={{
          borderTop: `solid 3px ${theme.colors.secondary![9]}`,
        }}
      />
    </Flex>
  );
}
