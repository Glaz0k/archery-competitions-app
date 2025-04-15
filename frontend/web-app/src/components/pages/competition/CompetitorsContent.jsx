import { IconFileUpload } from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import { format } from "date-fns";
import { formatInTimeZone } from "date-fns-tz";
import { useParams } from "react-router";
import {
  ActionIcon,
  Box,
  Checkbox,
  LoadingOverlay,
  ScrollArea,
  Stack,
  Table,
  TableScrollContainer,
} from "@mantine/core";
import { useSet } from "@mantine/hooks";
import { getCompetitorsFromCompetition } from "../../../api/competitors/competition";
import { COMPETITOR_QUERY_KEYS } from "../../../api/queryKeys";
import MainBar from "../../bars/MainBar";
import EmptyCardSpace from "../../misc/EmptyCardSpace";

export default function CompetitorsContent() {
  const { competitionId } = useParams();

  const selectedCompetitorIds = useSet([]);

  const {
    data: competitorDetails,
    isFetching: isCompetitorDetailsLoading,
    refetch: refetchCompetitorDetails,
  } = useQuery({
    queryKey: COMPETITOR_QUERY_KEYS.allByCompetition(competitionId),
    queryFn: () => getCompetitorsFromCompetition(competitionId),
    initialData: [],
  });

  const tableData = competitorDetails.map((competitorDetail) => {
    const competitor = competitorDetail?.competitor;
    return {
      id: competitor?.id,
      timeMoscow: formatInTimeZone(
        competitorDetail?.createdAt,
        "Europe/Moscow",
        "dd.MM.yyyy HH:mm:ss"
      ),
      fullName: competitor?.fullName,
      birthDate: format(competitor?.birthDate, "dd.MM.yyyy"),
      identity: competitor?.identity?.textValue,
      bow: competitor?.bow?.textValue,
      rank: competitor?.rank?.textValue,
      region: competitor?.region,
      federation: competitor?.federation,
      club: competitor?.club,
    };
  });

  const tableHead = (
    <Table.Tr>
      <Table.Th w={40} />
      <Table.Th>{"Отметка времени (МСК)"}</Table.Th>
      <Table.Th>{"Фамилия, Имя"}</Table.Th>
      <Table.Th>{"Дата рождения"}</Table.Th>
      <Table.Th>{"Пол"}</Table.Th>
      <Table.Th>{"Класс лука"}</Table.Th>
      <Table.Th>{"Спортивный разряд/звание"}</Table.Th>
      <Table.Th>{"Регион"}</Table.Th>
      <Table.Th>{"Членство в спортивной федерации"}</Table.Th>
      <Table.Th>{"Клубная принадлежность"}</Table.Th>
    </Table.Tr>
  );

  const tableRows = tableData.map((rowElement) => (
    <Table.Tr key={rowElement.id}>
      <Table.Td>
        <Checkbox
          checked={selectedCompetitorIds.has(rowElement.id)}
          onChange={(event) => {
            event.currentTarget.checked
              ? selectedCompetitorIds.add(rowElement.id)
              : selectedCompetitorIds.delete(rowElement.id);
          }}
        />
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
    </Table.Tr>
  ));

  return (
    <Stack flex={1}>
      <MainBar
        title="Участники соревнования"
        onRefresh={() => {
          refetchCompetitorDetails();
          console.log("xui");
        }}
      >
        <ActionIcon>
          <IconFileUpload />
        </ActionIcon>
      </MainBar>
      <Stack pos="relative">
        <LoadingOverlay visible={isCompetitorDetailsLoading} />
        {competitorDetails.length !== 0 ? (
          <Table.ScrollContainer>
            <Table tabularNums withColumnBorders>
              <Table.Thead>{tableHead}</Table.Thead>
              <Table.Tbody>{tableRows}</Table.Tbody>
            </Table>
          </Table.ScrollContainer>
        ) : (
          <EmptyCardSpace label="Участники не найдены" />
        )}
      </Stack>
    </Stack>
  );
}
