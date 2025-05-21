import { Center, Flex, Group, LoadingOverlay, Text } from "@mantine/core";
import { usePlaceRangeGroup, useSparringPlace } from "../../../api";
import type { ShootOut } from "../../../entities";
import { NO_SCORE_VALUE } from "../../constants";
import { SPARRING_MAX_RANGES, SPARRING_PLACE_SIZE } from "../constants";

const CELL_PROPS = {
  px: SPARRING_PLACE_SIZE.px,
  py: SPARRING_PLACE_SIZE.py,
  style: {
    boxShadow: "1px 0px #263238, 0px 1px #263238",
  },
};

export interface SparringPlaceViewProps {
  placeId?: number;
  placeOrd?: number;
}

export function SparringPlaceView({ placeId, placeOrd }: SparringPlaceViewProps) {
  const { data: place, isFetching: isPlaceFetching } = useSparringPlace(placeId!, !!placeId);
  const { data: rangeGroup, isFetching: isRangeGroupFetching } = usePlaceRangeGroup(
    placeId!,
    !!placeId
  );

  const ranges = [...(rangeGroup?.ranges ?? [])];
  const renderRanges = ranges
    .sort(({ ordinal: a }, { ordinal: b }) => a - b)
    .map((range) => (
      <Center key={`${range.id}`} w={SPARRING_PLACE_SIZE.score} h="100%" {...CELL_PROPS}>
        <Text size="sm">{range.score ?? NO_SCORE_VALUE}</Text>
      </Center>
    ))
    .concat(
      Array(SPARRING_MAX_RANGES - (rangeGroup?.ranges.length ?? 0))
        .fill(0)
        .map((_, index) => (
          <Center key={index} w={SPARRING_PLACE_SIZE.score} h="100%" {...CELL_PROPS}>
            <Text size="sm" />
          </Center>
        ))
    );

  return (
    <Group gap={0} wrap="nowrap" flex={1} pos="relative">
      <LoadingOverlay
        visible={isPlaceFetching || isRangeGroupFetching}
        loaderProps={{
          type: "dots",
        }}
      />
      {place && (
        <>
          <Center w={SPARRING_PLACE_SIZE.shield} h="100%" {...CELL_PROPS}>
            <Text size="sm">??</Text>
          </Center>
          {placeOrd && (
            <Center w={SPARRING_PLACE_SIZE.place} h="100%" {...CELL_PROPS}>
              <Text size="sm">{placeOrd}</Text>
            </Center>
          )}
          <Flex
            h="100%"
            flex={1}
            justify="flex-start"
            align="center"
            direction="row"
            wrap="nowrap"
            {...CELL_PROPS}
          >
            <Text size="sm">{place.competitor.fullName}</Text>
          </Flex>
          {renderRanges}
          <Center w={SPARRING_PLACE_SIZE.score} h="100%" {...CELL_PROPS}>
            <Text size="sm">{formatShootOut(place.shootOut)}</Text>
          </Center>
          <Center w={SPARRING_PLACE_SIZE.total} h="100%" {...CELL_PROPS}>
            <Text size="sm">{place.score}</Text>
          </Center>
        </>
      )}
    </Group>
  );
}

function formatShootOut(shootOut: ShootOut | null): string {
  if (!shootOut) {
    return "";
  }
  if (shootOut.score === null) {
    return NO_SCORE_VALUE;
  }
  const value = shootOut.score;
  if (shootOut.priority == null) {
    return value;
  }
  return shootOut.priority ? `${value}+` : `${value}-`;
}
