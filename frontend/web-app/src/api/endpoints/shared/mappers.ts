import type { Range, RangeGroup, Shot } from "../../../entities";
import type { RangeAPI, RangeAPIEdit, RangeEdit, RangeGroupAPI, ShotAPI } from "./types";

export const mapToRangeGroup = (response: RangeGroupAPI): RangeGroup => {
  return {
    id: response.id,
    rangesMaxCount: response.ranges_max_count,
    rangeSize: response.range_size,
    type: response.type,
    ranges: response.ranges.map(mapToRange),
    totalScore: response.total_score,
  };
};

export const mapToRange = (response: RangeAPI): Range => {
  return {
    id: response.id,
    ordinal: response.range_ordinal,
    isActive: response.is_active,
    shots: response.shots ? response.shots.map(mapToShot) : null,
    score: response.range_score,
  };
};

export const mapToShot = (response: ShotAPI): Shot => {
  return {
    ordinal: response.shot_ordinal,
    score: response.score,
  };
};

export const mapToRangeAPIEdit = (request: RangeEdit): RangeAPIEdit => {
  return {
    range_ordinal: request.ordinal,
    shots: request.shots ? request.shots.map(mapToShotAPI) : null,
  };
};

export const mapToShotAPI = (request: Shot): ShotAPI => {
  return {
    shot_ordinal: request.ordinal,
    score: request.score,
  };
};
