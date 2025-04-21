import React from "react";
import { IconCheck, IconTrashX, IconX } from "@tabler/icons-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ActionIcon, LoadingOverlay, Table, useMantineTheme } from "@mantine/core";
import { putCompetitor } from "../../api/competitors/competition";
import { COMPETITOR_QUERY_KEYS } from "../../api/queryKeys";

interface CompetitorRowProps {
  competitionId: number;
  id: number;
  timeMoscow: string;
  fullName: string;
  birthDate: string;
  identity: string;
  bow: string | null;
  rank: string | null;
  region: string | null;
  federation: string | null;
  club: string | null;
  isActive: boolean;
  onDelete: () => void;
  isDeleting: boolean;
}

export default function CompetitorRow(rowElement: CompetitorRowProps) {
  const theme = useMantineTheme();
  const queryClient = useQueryClient();

  const { mutate: toggleCompetitor, isPending: isCompetitorToggling } = useMutation({
    mutationFn: () => putCompetitor(rowElement.competitionId, rowElement.id, !rowElement.isActive),
    onSuccess: (editedDetail) => {
      queryClient.setQueryData(
        COMPETITOR_QUERY_KEYS.allByCompetition(editedDetail.competitionId),
        (old: any[]) =>
          old.map((detail) =>
            detail?.competitor?.id === editedDetail?.competitor?.id ? editedDetail : detail
          )
      );
    },
  });

  const isRowLoading = isCompetitorToggling || rowElement.isDeleting;

  return (
    <Table.Tr
      bg={
        isCompetitorToggling
          ? `${theme.colors.yellow[8]}33`
          : !rowElement.isActive
            ? `${theme.colors.red[8]}33`
            : undefined
      }
    >
      <Table.Td>
        {!rowElement.isActive ? (
          <ActionIcon
            variant="transparent"
            onClick={(_) => toggleCompetitor()}
            loading={isRowLoading}
          >
            <IconCheck />
          </ActionIcon>
        ) : (
          <ActionIcon
            variant="transparent"
            onClick={(_) => toggleCompetitor()}
            loading={isRowLoading}
          >
            <IconX />
          </ActionIcon>
        )}
      </Table.Td>
      <Table.Td>{rowElement.timeMoscow}</Table.Td>
      <Table.Td>{rowElement.fullName}</Table.Td>
      <Table.Td>{rowElement.birthDate}</Table.Td>
      <Table.Td>{rowElement.identity}</Table.Td>
      <Table.Td>{rowElement.bow || "Не указан"}</Table.Td>
      <Table.Td>{rowElement.rank || "б/р"}</Table.Td>
      <Table.Td>{rowElement.region || "Не указан"}</Table.Td>
      <Table.Td>{rowElement.federation || "Отсутствует"}</Table.Td>
      <Table.Td>{rowElement.club || "Отсутствует"}</Table.Td>
      <Table.Td>
        <ActionIcon
          variant="transparent"
          loading={isRowLoading}
          onClick={(_) => rowElement.onDelete()}
        >
          <IconTrashX />
        </ActionIcon>
      </Table.Td>
    </Table.Tr>
  );
}
