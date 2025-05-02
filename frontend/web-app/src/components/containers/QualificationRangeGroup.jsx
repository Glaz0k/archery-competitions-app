import { useContext, useEffect, useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { endRange, getRanges, putRange } from "../../api/qualificationSections";
import { SECTION_QUERY_KEYS } from "../../api/queryKeys";
import RangeCard from "../cards/RangeCard";
import { GroupContext } from "../pages/individual-group/GroupContext";

export default function QualificationRangeGroup({ sectionId, roundOrdinal, setRefresh }) {
  const {
    data: rangeGroup,
    isFetching: isRangeGroupLoading,
    refetch: refreshRangeGroup,
  } = useQuery({
    queryKey: SECTION_QUERY_KEYS.ranges(sectionId, roundOrdinal),
    queryFn: () => getRanges(sectionId, roundOrdinal),
    initialData: null,
  });

  const rangeSize = rangeGroup?.rangeSize || 0;
  const ranges = [...(rangeGroup?.ranges || [])].sort((a, b) => a.rangeOrdinal - b.rangeOrdinal);

  useEffect(() => {
    setRefresh(() => () => refreshRangeGroup());
  }, [setRefresh, refreshRangeGroup]);

  return ranges?.map((range) => (
    <QualificationRangeCard
      key={range.id}
      sectionId={sectionId}
      roundOrdinal={roundOrdinal}
      initialRange={range}
      rangeSize={rangeSize}
      loading={isRangeGroupLoading}
    />
  ));
}

function QualificationRangeCard({ sectionId, roundOrdinal, initialRange, rangeSize, loading }) {
  const groupContext = useContext(GroupContext);

  const [range, setRange] = useState(initialRange);

  const { mutateAsync: editRange, isPending: isRangeEditing } = useMutation({
    mutationFn: (editedShots) =>
      putRange(sectionId, roundOrdinal, { rangeOrdinal: range.rangeOrdinal, shots: editedShots }),
    onSuccess: (fetchedRange) => {
      setRange(fetchedRange);
    },
  });

  const { mutate: completeRange, isPending: isRangeCompleting } = useMutation({
    mutationFn: () => endRange(sectionId, roundOrdinal, range.rangeOrdinal),
    onSuccess: (fetchedRange) => {
      setRange(fetchedRange);
    },
  });

  const rangeRegex =
    groupContext?.groupBow.value === BowClass.CLASSIC.value ||
    groupContext?.groupBow.value === BowClass.BLOCK.value
      ? /^(M|[6-9]|10|X)$/
      : /^(M|[1-9]|10|X)$/;

  return (
    <RangeCard
      range={range}
      rangeSize={rangeSize}
      rangeRegex={rangeRegex}
      loading={isRangeEditing || isRangeCompleting || loading}
      onShotsSubmit={editRange}
      onComplete={completeRange}
    />
  );
}
