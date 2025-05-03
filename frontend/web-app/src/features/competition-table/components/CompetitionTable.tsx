import { useState } from "react";
import { Card, LoadingOverlay, Table, useMantineTheme } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useCompetitionCompetitors, useRemoveCompetitorFromCompetition } from "../../../api";
import { CenterCard } from "../../../widgets";
import { RemoveCompetitorModal } from "../../form-modals";
import { CompetitionTableHead, CompetitionTableRow } from "./CompetitionTableRow";

export function CompetitionTable({ competitionId }: { competitionId: number }) {
  const theme = useMantineTheme();

  const [removingId, setRemovingId] = useState<number | null>(null);
  const [isRemoveOpened, controlRemove] = useDisclosure();

  const {
    data: details,
    isFetching: isDetailsFetching,
    isError: isDetailsError,
  } = useCompetitionCompetitors(competitionId, !Number.isNaN(competitionId));

  const cancelRemove = () => {
    setRemovingId(null);
    controlRemove.close();
  };

  const { mutate: removeCompetitor, isPending: isCompetitorRemoving } =
    useRemoveCompetitorFromCompetition(cancelRemove);

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
      <Card p={0} pos="relative">
        <LoadingOverlay visible={isDetailsFetching} />
        <Table.ScrollContainer minWidth={500}>
          <Table
            tabularNums
            withColumnBorders
            highlightOnHover
            highlightOnHoverColor={`${theme.colors.gray[0]}33`}
          >
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
          </Table>
        </Table.ScrollContainer>
      </Card>
    </>
  );
}
