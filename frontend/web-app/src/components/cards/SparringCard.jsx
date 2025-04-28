import { useContext } from "react";
import {
  Card,
  Center,
  Flex,
  Group,
  LoadingOverlay,
  rem,
  Stack,
  Text,
  useMantineTheme,
} from "@mantine/core";
import useSparringPlace from "../../hooks/useSparringPlace";
import useSparringPlaceRangeGroup from "../../hooks/useSparringPlaceRangeGroup";
import { FinalContext } from "../pages/individual-group/FinalContext";
import styles from "./SparringCard.module.css";

const RANGES_MAX_LENGTH = 5;
const placeSize = {
  shield: rem(25),
  place: rem(20),
  rangeScore: rem(25),
  totalScore: rem(35),
  pX: rem(5),
  pY: rem(2),
};
const cellStyle = {
  boxShadow: "1px 0px #263238, 0px 1px #263238",
};
const cellPlaceholder = "-";

export default function SparringCard({ sparring, ordinal }) {
  const { sparringSize } = useContext(FinalContext);

  return (
    <Card
      h={rem(sparringSize.heigth)}
      w={rem(sparringSize.width)}
      p={0}
      style={{
        border: "2px solid #263238",
      }}
    >
      <Stack gap={0} flex={1}>
        <SparringPlace
          placeId={sparring?.topPlace?.id}
          qualificationPlace={numberToPair(ordinal)[0]}
        />
        <SparringPlace
          placeId={sparring?.botPlace?.id}
          qualificationPlace={numberToPair(ordinal)[1]}
        />
      </Stack>
    </Card>
  );
}

function SparringPlace({ placeId, qualificationPlace }) {
  const theme = useMantineTheme();

  const { selectedPlaceId, setSelectedPlaceId } = useContext(FinalContext);
  const { place, isPlaceFetching } = useSparringPlace(placeId);
  const { rangeGroup, isRangeGroupFetching } = useSparringPlaceRangeGroup(placeId);

  const ranges = rangeGroup?.ranges || [];

  const renderRanges = ranges
    .concat(Array(RANGES_MAX_LENGTH - ranges.length).fill(0))
    .map((range, index) => (
      <Center
        key={index}
        w={placeSize.rangeScore}
        h="100%"
        px={placeSize.pX}
        py={placeSize.pY}
        style={cellStyle}
      >
        <Text size="sm">{range?.rangeScore || cellPlaceholder}</Text>
      </Center>
    ));

  return (
    <Group
      gap={0}
      wrap="nowrap"
      flex={1}
      pos="relative"
      className={styles.hoverBg}
      onClick={placeId != null ? () => setSelectedPlaceId(placeId) : undefined}
      bg={
        selectedPlaceId != null && selectedPlaceId === placeId
          ? `${theme.colors.yellow[0]}33`
          : undefined
      }
    >
      <LoadingOverlay
        visible={isPlaceFetching || isRangeGroupFetching}
        loaderProps={{
          type: "dots",
        }}
      />
      <Center w={placeSize.shield} h="100%" px={placeSize.pX} py={placeSize.pY} style={cellStyle}>
        <Text size="sm">{"??"}</Text>
      </Center>
      {qualificationPlace != null && (
        <Center w={placeSize.place} h="100%" px={placeSize.pX} py={placeSize.pY} style={cellStyle}>
          <Text size="sm">{qualificationPlace}</Text>
        </Center>
      )}
      <Flex
        h="100%"
        flex={1}
        justify="flex-start"
        align="center"
        direction="row"
        wrap="nowrap"
        px={placeSize.pX}
        py={placeSize.pY}
        style={cellStyle}
      >
        <Text size="sm">{place?.competitor.fullName}</Text>
      </Flex>
      {renderRanges}
      <Center
        w={placeSize.rangeScore}
        h="100%"
        px={placeSize.pX}
        py={placeSize.pY}
        style={cellStyle}
      >
        <Text size="sm">{getShootOutText(place?.shootOut)}</Text>
      </Center>
      <Center
        w={placeSize.totalScore}
        h="100%"
        px={placeSize.pX}
        py={placeSize.pY}
        style={cellStyle}
      >
        <Text size="sm">{place?.sparringScore || cellPlaceholder}</Text>
      </Center>
    </Group>
  );
}

function getShootOutText(shootOut) {
  let text = shootOut?.score || cellPlaceholder;
  if (shootOut?.priority == null) {
    return text;
  }
  return shootOut.priority ? text + "+" : text + "-";
}

function numberToPair(input) {
  const mapping = {
    1: [1, 8],
    2: [5, 4],
    3: [3, 6],
    4: [7, 2],
  };

  return mapping[input] || [null, null];
}
