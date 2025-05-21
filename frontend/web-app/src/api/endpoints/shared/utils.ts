import type { Range, RangeGroup } from "../../../entities";

export const updateRangeGroup = (
  prev: RangeGroup | undefined,
  edited: Range
): RangeGroup | undefined => {
  if (!prev) {
    return undefined;
  }
  return {
    ...prev,
    ranges: prev.ranges.map((range) => (range.ordinal !== edited.ordinal ? range : edited)),
  };
};
