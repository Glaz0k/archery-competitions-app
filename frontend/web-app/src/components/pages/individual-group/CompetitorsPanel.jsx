import { useEffect } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns";
import { Table } from "@mantine/core";
import { getCompetitors, syncCompetitors } from "../../../api/competitors/individualGroup";
import { COMPETITOR_QUERY_KEYS } from "../../../api/queryKeys";
import TableCard from "../../cards/TableCard";

export default function CompetitorsPanel({ groupInfo, setGroupControls }) {
  const queryClient = useQueryClient();

  const group = groupInfo?.group;

  const { isFetching: isCompetitorsFetching, data: competitors } = useQuery({
    queryKey: COMPETITOR_QUERY_KEYS.allByGroup(group?.id),
    queryFn: () => getCompetitors(group?.id),
    initialData: [],
  });

  const { mutate: refreshCompetitors, isPending: isCompetitorsSyncing } = useMutation({
    mutationFn: () => syncCompetitors(group?.id),
    onSuccess: (synced) => {
      queryClient.setQueryData(COMPETITOR_QUERY_KEYS.allByGroup(group?.id), (_) => [...synced]);
    },
  });

  const isCompetitorsLoading = isCompetitorsFetching || isCompetitorsSyncing;

  useEffect(() => {
    setGroupControls({
      onExport: null,
      onRefresh: refreshCompetitors,
      onEnd: null,
    });
  }, [setGroupControls, refreshCompetitors]);

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
        <Table.Td>{competitor.identity.textValue}</Table.Td>
        <Table.Td>{competitor.bow?.textValue || "Не указан"}</Table.Td>
        <Table.Td>{competitor.rank?.textValue || "б/р"}</Table.Td>
        <Table.Td>{competitor.region || "Не указан"}</Table.Td>
        <Table.Td>{competitor.federation || "Отсутствует"}</Table.Td>
        <Table.Td>{competitor.club || "Отсутствует"}</Table.Td>
      </Table.Tr>
    );
  });

  return (
    <TableCard loading={isCompetitorsLoading}>
      <Table.Thead>{tableHead}</Table.Thead>
      <Table.Tbody>{tableRows}</Table.Tbody>
    </TableCard>
  );
}
