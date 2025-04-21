import { useState } from "react";
import { IconFileUpload } from "@tabler/icons-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns";
import { formatInTimeZone } from "date-fns-tz";
import { useNavigate, useParams } from "react-router";
import { ActionIcon, Card, LoadingOverlay, Stack, Table } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { deleteCompetitor, getCompetitors } from "../../../api/competitors/competition";
import { COMPETITOR_QUERY_KEYS } from "../../../api/queryKeys";
import MainBar from "../../bars/MainBar";
import NotFoundCard from "../../cards/NotFoundCard";
import TableCard from "../../cards/TableCard";
import CompetitorRow from "../../misc/CompetitorRow";
import DeleteCompetitorModal from "../../modals/competitor/DeleteCompetitorModal";

export default function CompetitorsContent() {
  const { competitionId } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [competitorDeletingId, setCompetitorDeletingId] = useState(null);

  const [isOpenedCompetitorDel, competitorDelControl] = useDisclosure();

  const {
    data: competitorDetails,
    isFetching: isCompetitorDetailsLoading,
    refetch: refetchCompetitorDetails,
  } = useQuery({
    queryKey: COMPETITOR_QUERY_KEYS.allByCompetition(competitionId),
    queryFn: () => getCompetitors(competitionId),
    initialData: [],
  });

  const { mutate: removeCompetitor, isPending: isCompetitorDeleting } = useMutation({
    mutationFn: () => deleteCompetitor(competitionId, competitorDeletingId),
    onSuccess: () => {
      queryClient.setQueryData(COMPETITOR_QUERY_KEYS.allByCompetition(competitionId), (old) =>
        old.filter((detail) => detail?.competitor?.id !== competitorDeletingId)
      );
      setCompetitorDeletingId(null);
      competitorDelControl.close();
    },
  });

  const handleCompetitorDeleting = (id) => {
    setCompetitorDeletingId(id);
    competitorDelControl.open();
  };

  const handleCompetitorAdd = () => {
    // TODO
    console.warn("handleCompetitorAdd temporary unavailable");
  };

  const handleExcelTableUpload = () => {
    // TODO
    console.warn("handleExcelTableUpload temporary unavailable");
  };

  const tableData = competitorDetails.map((competitorDetail) => {
    const competitor = competitorDetail?.competitor;
    return {
      competitionId: Number(competitionId),
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
      isActive: competitorDetail.isActive,
      onDelete: () => handleCompetitorDeleting(competitorDetail?.competitor?.id),
      isDeleting: competitorDetail?.competitor?.id === competitorDeletingId && isCompetitorDeleting,
    };
  });

  const tableHead = (
    <Table.Tr>
      <Table.Th />
      <Table.Th>{"Отметка времени (МСК)"}</Table.Th>
      <Table.Th>{"Фамилия, Имя"}</Table.Th>
      <Table.Th>{"Дата рождения"}</Table.Th>
      <Table.Th>{"Пол"}</Table.Th>
      <Table.Th>{"Класс лука"}</Table.Th>
      <Table.Th>{"Разряд"}</Table.Th>
      <Table.Th>{"Регион"}</Table.Th>
      <Table.Th>{"Федерация"}</Table.Th>
      <Table.Th>{"Клуб"}</Table.Th>
      <Table.Th />
    </Table.Tr>
  );

  const tableRows = tableData.map((rowElement) => (
    <CompetitorRow key={rowElement.id} {...rowElement} />
  ));

  return (
    <>
      <DeleteCompetitorModal
        isOpened={isOpenedCompetitorDel}
        onClose={competitorDelControl.close}
        onConfirm={removeCompetitor}
        isLoading={isCompetitorDeleting}
      />
      <Stack flex={1} style={{ overflow: "hidden" }} gap="lg" miw={500}>
        <MainBar
          title="Участники соревнования"
          onRefresh={refetchCompetitorDetails}
          onAdd={handleCompetitorAdd}
          onBack={() => navigate("..")}
        >
          <ActionIcon onClick={handleExcelTableUpload}>
            <IconFileUpload />
          </ActionIcon>
        </MainBar>
        {competitorDetails.length !== 0 ? (
          <TableCard loading={isCompetitorDetailsLoading}>
            <Table.Thead>{tableHead}</Table.Thead>
            <Table.Tbody>{tableRows}</Table.Tbody>
          </TableCard>
        ) : isCompetitorDetailsLoading ? (
          <Card h={500} pos="relative">
            <LoadingOverlay visible />
          </Card>
        ) : (
          <NotFoundCard label="Участники не найдены" />
        )}
      </Stack>
    </>
  );
}
