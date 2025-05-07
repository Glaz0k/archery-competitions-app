import { useEffect } from "react";
import { format } from "date-fns";
import { Table } from "@mantine/core";
import { useGroupCompetitors, useSyncGroupCompetitors } from "../../../api";
import { GroupState } from "../../../entities";
import {
  getBowClassDescription,
  getGenderDescription,
  getSportsRankDescription,
} from "../../../utils";
import { CenterCard, TableCard } from "../../../widgets";
import { useGroupSections } from "../context/useGroupSections";

export function CompetitorsSection({ groupId }: { groupId: number }) {
  const {
    context: {
      info: { group },
    },
    setContext,
  } = useGroupSections();

  const {
    isFetching: isCompetitorsFetching,
    data: competitors,
    isError: isCompetitorsError,
  } = useGroupCompetitors(groupId);

  const { mutate: syncCompetitors, isPending: isCompetitorsSyncing } = useSyncGroupCompetitors();

  const isCompetitorsUpdating = isCompetitorsFetching || isCompetitorsSyncing;

  useEffect(() => {
    let refteshFn = undefined;
    if (group && group.state === GroupState.CREATED) {
      refteshFn = () => syncCompetitors(group.id);
    }
    setContext((prev) => ({
      ...prev,
      controls: {
        onRefresh: refteshFn,
        onComplete: undefined,
        onExport: undefined,
      },
    }));
  }, [group, setContext, syncCompetitors]);

  const tableHead = (
    <Table.Tr>
      <Table.Th>{"Фамилия, Имя"}</Table.Th>
      <Table.Th>{"Дата рождения"}</Table.Th>
      <Table.Th>{"Пол"}</Table.Th>
      <Table.Th>{"Класс лука"}</Table.Th>
      <Table.Th>{"Разряд"}</Table.Th>
      <Table.Th>{"Регион"}</Table.Th>
      <Table.Th>{"Федерация"}</Table.Th>
      <Table.Th>{"Клуб"}</Table.Th>
    </Table.Tr>
  );

  const tableRows = competitors.map((detail) => {
    const competitor = detail.competitor;
    return (
      <Table.Tr key={competitor.id}>
        <Table.Td>{competitor.fullName}</Table.Td>
        <Table.Td>{format(competitor.birthDate, "dd.MM.yyyy")}</Table.Td>
        <Table.Td>{getGenderDescription(competitor.identity)}</Table.Td>
        <Table.Td>{competitor.bow ? getBowClassDescription(competitor.bow) : "Не указан"}</Table.Td>
        <Table.Td>{competitor.rank ? getSportsRankDescription(competitor.rank) : "б/р"}</Table.Td>
        <Table.Td>{competitor.region || "Не указан"}</Table.Td>
        <Table.Td>{competitor.federation || "Не указана"}</Table.Td>
        <Table.Td>{competitor.club || "Не указан"}</Table.Td>
      </Table.Tr>
    );
  });

  if (!isCompetitorsUpdating && isCompetitorsError) {
    return <CenterCard label="Произошла ошибка" />;
  }

  return (
    <TableCard loading={isCompetitorsUpdating}>
      <Table.Thead>{tableHead}</Table.Thead>
      <Table.Tbody>{tableRows}</Table.Tbody>
    </TableCard>
  );
}
