import { useEffect, useState } from "react";
import { Stack } from "@mantine/core";
import {
  useEndFinal,
  useFinalGrid,
  useStartFinal,
  useStartQuarterfinal,
  useStartSemifinal,
} from "../../../api";
import { GroupState, type Sparring } from "../../../entities";
import { CenterCard, PageLoader } from "../../../widgets";
import { FinalGridView } from "../../final-grid-view";
import { useGroupSections } from "../context/useGroupSections";
import { SparringTab } from "./SparringTab";
import { StartCard } from "./StartCard";

export function FinalSection({ groupId }: { groupId: number }) {
  const {
    context: {
      info: { group },
    },
    setContext,
  } = useGroupSections();

  const hasGroup = group !== null;
  const hasStarted =
    hasGroup &&
    group.state !== GroupState.CREATED &&
    group.state !== GroupState.QUAL_START &&
    group.state !== GroupState.QUAL_END;

  const [selectedSparring, setSelectedSparring] = useState<Sparring | null>();

  const {
    data: grid,
    isFetching: isGrifFetching,
    isLoading: isGridLoading,
    isError: isGridError,
    refetch: refreshGrid,
  } = useFinalGrid(groupId, hasStarted);

  const { mutate: startQFinal, isPending: isQFinalStarting } = useStartQuarterfinal();
  const { mutate: startSFinal, isPending: isSFinalStarting } = useStartSemifinal();
  const { mutate: startFinal, isPending: isFinalStarting } = useStartFinal();
  const { mutate: endFinal, isPending: isFinalEnding } = useEndFinal();
  const isGridMutating = isQFinalStarting || isSFinalStarting || isFinalStarting || isFinalEnding;

  const handleGridExport = () => {
    console.warn("handleGridExport temporary unavailable");
  };

  useEffect(() => {
    let refreshFn = undefined;
    let exportFn = undefined;
    let completeFn = undefined;
    if (group && !isGridError && !isGridLoading) {
      if (hasStarted) {
        refreshFn = () => {
          setSelectedSparring(null);
          refreshGrid();
        };
      }
      switch (group.state) {
        case GroupState.QUARTERFINAL_START:
          completeFn = () => startSFinal(group.id);
          break;
        case GroupState.SEMIFINAL_START:
          completeFn = () => startFinal(group.id);
          break;
        case GroupState.FINAL_START:
          completeFn = () => endFinal(group.id);
          break;
        case GroupState.COMPLETED:
          exportFn = handleGridExport;
          break;
      }
    }
    setContext((prev) => ({
      ...prev,
      controls: {
        onRefresh: refreshFn,
        onComplete: completeFn,
        onExport: exportFn,
      },
    }));
  }, [
    endFinal,
    group,
    hasStarted,
    isGridError,
    isGridLoading,
    refreshGrid,
    setContext,
    startFinal,
    startSFinal,
  ]);

  if (!hasGroup) return;

  if (group.state === GroupState.CREATED || group.state === GroupState.QUAL_START) {
    return <CenterCard label="Пока нельзя начать финал" />;
  }

  if (!hasStarted) {
    return (
      <StartCard
        title="Финал еще не начался"
        loading={isQFinalStarting}
        onStart={() => startQFinal(group.id)}
      />
    );
  }

  return (
    <PageLoader loading={isGridLoading || isGridMutating} error={isGridError}>
      <Stack display="flex" flex={1} style={{ overflow: "hidden" }}>
        {grid && (
          <>
            <FinalGridView
              grid={grid}
              loading={isGrifFetching}
              selectSparring={setSelectedSparring}
            />
            {selectedSparring && <SparringTab sparring={selectedSparring} />}
          </>
        )}
      </Stack>
    </PageLoader>
  );
}
