import { useContext, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { PLACE_QUERY_KEYS } from "../../api/queryKeys";
import { putShootOut } from "../../api/sparringPlaces";
import BowClass from "../../enums/BowClass";
import useCompleteRange from "../../hooks/useCompleteRange";
import useEditRange from "../../hooks/useEditRange";
import useSparringPlace from "../../hooks/useSparringPlace";
import useSparringPlaceRangeGroup from "../../hooks/useSparringPlaceRangeGroup";
import NavigationBar from "../bars/NavigationBar";
import RangeCard from "../cards/RangeCard";
import ShootOutCard from "../cards/ShootOutCard";
import { FinalContext } from "../pages/individual-group/FinalContext";
import { GroupContext } from "../pages/individual-group/GroupContext";

export default function FinalSection() {
  const queryClient = useQueryClient();
  const groupContext = useContext(GroupContext);

  const { selectedPlaceId } = useContext(FinalContext);
  const { place, isPlaceFetching } = useSparringPlace(selectedPlaceId);
  const { rangeGroup, isRangeGroupFetching } = useSparringPlaceRangeGroup(selectedPlaceId);

  const rangeRegex =
    groupContext?.groupBow.value === BowClass.CLASSIC.value ||
    groupContext?.groupBow.value === BowClass.BLOCK.value
      ? /^(M|[6-9]|10|X)$/
      : /^(M|[1-9]|10|X)$/;

  return (
    <>
      <NavigationBar
        title={place?.competitor.fullName}
        onRefresh={() => {
          queryClient.invalidateQueries({
            queryKey: PLACE_QUERY_KEYS.element(selectedPlaceId),
          });
        }}
        loading={isPlaceFetching}
      />
      {rangeGroup?.ranges
        ?.sort((a, b) => a.rangeOrdinal - b.rangeOrdinal)
        .map((range) => (
          <FinalRangeCard
            key={range.id}
            placeId={selectedPlaceId}
            initialRange={range}
            rangeSize={rangeGroup?.rangeSize || 0}
            rangeRegex={rangeRegex}
            loading={isRangeGroupFetching}
          />
        ))}
      {place?.shootOut && (
        <FinalShootOutCard
          placeId={selectedPlaceId}
          initialShootOut={place.shootOut}
          rangeRegex={rangeRegex}
          loading={isRangeGroupFetching}
        />
      )}
    </>
  );
}

function FinalRangeCard({ placeId, initialRange, rangeSize, rangeRegex, loading }) {
  const [range, setRange] = useState(initialRange);
  const { asyncEditRange, isRangePending: isEditing } = useEditRange(
    { placeId: placeId },
    range.rangeOrdinal,
    setRange
  );
  const { completeRange, isRangePending: isCompleting } = useCompleteRange(
    { placeId: placeId },
    range.rangeOrdinal,
    setRange
  );

  return (
    <RangeCard
      range={range}
      rangeSize={rangeSize}
      rangeRegex={rangeRegex}
      loading={loading || isEditing || isCompleting}
      onShotsSubmit={asyncEditRange}
      onComplete={completeRange}
    />
  );
}

function FinalShootOutCard({ placeId, initialShootOut, rangeRegex, loading }) {
  const [shootOut, setShootOut] = useState(initialShootOut);
  const queryClient = useQueryClient();
  const { mutateAsync: asyncEditShootOut, isPending: isShootOutPending } = useMutation({
    mutationFn: (changeShootOut) => putShootOut(placeId, changeShootOut),
    onSuccess: (editedShootOut) => {
      queryClient.invalidateQueries({
        queryKey: PLACE_QUERY_KEYS.all,
      });
      setShootOut(editedShootOut);
    },
  });

  return (
    <ShootOutCard
      shootOut={shootOut}
      shootOutRegex={toShootOutRegex(rangeRegex)}
      loading={loading || isShootOutPending}
      onShootOutSubmit={asyncEditShootOut}
    />
  );
}

function toShootOutRegex(regex) {
  const regexString = regex.toString();
  const pattern = regexString.slice(1, -1);
  const modifiedPattern = pattern.replace(/\$$/, "[+-]?$");
  return new RegExp(modifiedPattern);
}
