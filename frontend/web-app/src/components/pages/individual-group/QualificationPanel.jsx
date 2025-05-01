import { useCallback, useEffect, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  Card,
  Center,
  LoadingOverlay,
  Stack,
  Table,
  TableTbody,
  TableThead,
  useMantineTheme,
} from "@mantine/core";
import {
  endQualification,
  getQualification,
  startQualification,
} from "../../../api/individualGroups";
import { INDIVIDUAL_GROUP_QUERY_KEYS } from "../../../api/queryKeys";
import GroupState from "../../../enums/GroupState";
import groupStateComparator from "../../../helper/groupStateComparator";
import { TextButton } from "../../../widgets/buttons/TextButton";
import NotFoundCard from "../../cards/NotFoundCard";
import TableCard from "../../cards/TableCard";
import QualificationSection from "../../containers/QualificationSection";

export default function QualificationPanel({ groupInfo, setGroupControls }) {
  const theme = useMantineTheme();
  const queryClient = useQueryClient();

  const group = groupInfo?.group;
  const [selectedSectionId, setSelectedSectionId] = useState(null);

  const {
    data: qualification,
    isFetching: isQualificationFetching,
    refetch: refreshQualification,
  } = useQuery({
    queryKey: INDIVIDUAL_GROUP_QUERY_KEYS.qualification(group?.id),
    queryFn: () => getQualification(group?.id),
    initialData: null,
  });

  const { mutate: initQualification, isPending: isQualificationStarting } = useMutation({
    mutationFn: () => startQualification(group?.id),
    onSuccess: (started) => {
      queryClient.invalidateQueries(INDIVIDUAL_GROUP_QUERY_KEYS.element(group?.id));
      queryClient.setQueryData(
        INDIVIDUAL_GROUP_QUERY_KEYS.qualification(group?.id),
        (_) => started
      );
    },
  });

  const { mutate: closeQualification, isPending: isQualificationEnding } = useMutation({
    mutationFn: () => endQualification(group?.id),
    onSuccess: (started) => {
      queryClient.invalidateQueries(INDIVIDUAL_GROUP_QUERY_KEYS.element(group?.id));
      queryClient.setQueryData(
        INDIVIDUAL_GROUP_QUERY_KEYS.qualification(group?.id),
        (_) => started
      );
    },
  });

  const isQualificationLoading =
    isQualificationFetching || isQualificationStarting || isQualificationEnding;
  const roundCount = qualification?.roundCount || 0;
  const sections = qualification?.sections || [];

  const handleGroupExport = useCallback(() => {
    console.warn("handleGroupExport temporary unavailable");
  }, []);

  useEffect(() => {
    if (group != null) {
      setGroupControls({
        onExport: handleGroupExport,
        onRefresh: () => {
          refreshQualification(), setSelectedSectionId(null);
        },
        onEnd:
          groupStateComparator(GroupState.QUAL_START, group.state) === 0
            ? closeQualification
            : null,
      });
    }
  }, [setGroupControls, handleGroupExport, refreshQualification, group, closeQualification]);

  const tableHead = (
    <Table.Tr>
      <Table.Th w={50}>{"Место"}</Table.Th>
      <Table.Th>{"Спортсмен"}</Table.Th>
      {Array(roundCount)
        .fill(0)
        .map((_, index) => (
          <Table.Th w={100} key={index}>
            {qualification?.distance + "-" + (index + 1)}
          </Table.Th>
        ))}
      <Table.Th w={100}>{"Итог"}</Table.Th>
      <Table.Th w={100}>{"10's"}</Table.Th>
      <Table.Th w={100}>{"9's"}</Table.Th>
      <Table.Th w={50}>{"Выполненный разряд"}</Table.Th>
    </Table.Tr>
  );

  const tableRows = [...(sections || [])]
    .sort((a, b) => a.place - b.place)
    .map((section) => {
      const rounds = [...section.rounds]
        .sort((a, b) => a.roundOrdinal - b.roundOrdinal)
        .concat(Array(Math.min(roundCount - section.rounds.length, 0)).fill(null));

      return (
        <Table.Tr
          onClick={() => setSelectedSectionId(section.id)}
          bg={section.id === selectedSectionId ? `${theme.colors.yellow[0]}33` : undefined}
          key={section.id}
        >
          <Table.Td>{section.place}</Table.Td>
          <Table.Td>{section.competitor.fullName}</Table.Td>
          {rounds.map((round, index) => (
            <Table.Td key={index} bg={round?.isOngoing ? `${theme.colors.yellow[8]}33` : undefined}>
              {round?.total || ""}
            </Table.Td>
          ))}
          <Table.Td>{section.total || 0}</Table.Td>
          <Table.Td>{section.count10 || 0}</Table.Td>
          <Table.Td>{section.count9 || 0}</Table.Td>
          <Table.Td>{section.rankGained?.textValue || ""}</Table.Td>
        </Table.Tr>
      );
    });

  return sections.length !== 0 ? (
    <Stack gap="md">
      <TableCard loading={isQualificationLoading}>
        <TableThead>{tableHead}</TableThead>
        <TableTbody>{tableRows}</TableTbody>
      </TableCard>
      {selectedSectionId != null && <QualificationSection sectionId={selectedSectionId} />}
    </Stack>
  ) : isQualificationLoading ? (
    <Card h={400} pos="relative">
      <LoadingOverlay visible />
    </Card>
  ) : groupStateComparator(GroupState.CREATED, group.state) === 0 ? (
    <Center flex={1}>
      <Card p="md">
        <TextButton label="Начать квалификацию" onClick={initQualification} />
      </Card>
    </Center>
  ) : (
    <NotFoundCard label={"Произошла ошибка"} />
  );
}
