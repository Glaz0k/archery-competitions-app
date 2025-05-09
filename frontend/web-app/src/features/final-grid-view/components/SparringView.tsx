import { Card, Stack } from "@mantine/core";
import type { Sparring } from "../../../entities";
import { SPARRING_SIZE } from "../constants";
import { SparringPlaceView } from "./SparringPlaceView";
import classes from "./SparringView.module.css";

export interface SparringViewProps {
  sparring?: Sparring | null;
  ordinal?: number;
  selectSparring: (sparring: Sparring) => unknown;
}

export function SparringView({ sparring, ordinal, selectSparring }: SparringViewProps) {
  const places = placesByOrdinal(ordinal);
  const isSelectable = sparring && sparring.top && sparring.bot;

  return (
    <Card
      h={SPARRING_SIZE.heigth}
      w={SPARRING_SIZE.width}
      p={0}
      style={{
        border: "2px solid #263238",
      }}
      onClick={isSelectable ? () => selectSparring(sparring) : undefined}
    >
      <Stack gap={0} flex={1} className={isSelectable ? classes.sparringCard : undefined}>
        <SparringPlaceView placeId={sparring?.top?.id} placeOrd={places?.top} />
        <SparringPlaceView placeId={sparring?.bot?.id} placeOrd={places?.bot} />
      </Stack>
    </Card>
  );
}

interface PlaceMapping {
  top: number;
  bot: number;
}

const PLACES_MAPPING: Record<number, PlaceMapping> = {
  1: { top: 1, bot: 8 },
  2: { top: 5, bot: 4 },
  3: { top: 3, bot: 6 },
  4: { top: 7, bot: 2 },
};

function placesByOrdinal(ordinal?: number): PlaceMapping | undefined {
  return ordinal ? PLACES_MAPPING[ordinal] : undefined;
}
