import { useEffect, useState } from "react";
import { Card, LoadingOverlay, Table } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useCompetitionCompetitors, useRemoveCompetitorFromCompetition } from "../../../api";
import { CenterCard, TableCard } from "../../../widgets";
import { RemoveCompetitorModal } from "../../form-modals";
import { useTableControls } from "../context/useTableControls";
import { CompetitionTableHead, CompetitionTableRow } from "./CompetitionTableRow";

export function CompetitionTable({ competitionId }: { competitionId: number }) {
  const { setControls } = useTableControls();

  const [removingId, setRemovingId] = useState<number | null>(null);
  const [isRemoveOpened, controlRemove] = useDisclosure();

  const {
    data: details,
    isFetching: isDetailsFetching,
    isError: isDetailsError,
    refetch: refetchDetails,
  } = useCompetitionCompetitors(competitionId, !Number.isNaN(competitionId));

  const cancelRemove = () => {
    setRemovingId(null);
    controlRemove.close();
  };

  const { mutate: removeCompetitor, isPending: isCompetitorRemoving } =
    useRemoveCompetitorFromCompetition(cancelRemove);

  useEffect(() => {
    setControls((prev) => ({
      ...prev,
      refresh: () => {
        refetchDetails();
      },
    }));
  }, [refetchDetails, setControls]);

  if (details.length === 0 && isDetailsFetching) {
    return (
      <Card p={0} pos="relative" h={400}>
        <LoadingOverlay visible />
      </Card>
    );
  }

  if (isDetailsError) {
    return <CenterCard label="Произошла ошибка" />;
  }

  return (
    <>
      <RemoveCompetitorModal
        opened={isRemoveOpened}
        onConfirm={() => removeCompetitor([competitionId, removingId!])}
        onClose={cancelRemove}
        loading={isCompetitorRemoving}
      />
      <TableCard loading={isDetailsFetching}>
        <CompetitionTableHead />
        <Table.Tbody>
          {details.map((value) => (
            <CompetitionTableRow
              detail={value}
              deleting={isDetailsFetching && value.competitor.id === removingId}
              onDelete={(competitorId: number) => {
                setRemovingId(competitorId);
                controlRemove.open();
              }}
            />
          ))}
        </Table.Tbody>
      </TableCard>
    </>
  );
}
