import GroupState from "../enums/GroupState";

export default function groupStateComparator(stageA, stageB) {
  const valuesMap = new Map(Object.values(GroupState).map((state, index) => [state.value, index]));
  return valuesMap.get(stageA.value) - valuesMap.get(stageB.value);
}
