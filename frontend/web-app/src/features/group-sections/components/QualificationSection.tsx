import { useEffect, useState } from "react";
import { Stack, Table, useMantineTheme } from "@mantine/core";
import { useEndQualification, useQualification, useStartQualification } from "../../../api";
import { GroupState, type Qualification, type QualificationSection } from "../../../entities";
import { getSportsRankDescription } from "../../../utils";
import { PageLoader, TableCard } from "../../../widgets";
import { NO_SCORE_VALUE } from "../../constants";
import { useGroupSections } from "../context/useGroupSections";
import { SectionTab } from "./SectionTab";
import { StartCard } from "./StartCard";

export function QualificationSection({ groupId }: { groupId: number }) {
  const {
    context: {
      info: { group },
    },
    setContext,
  } = useGroupSections();

  const hasGroup = group !== null;
  const hasStarted = hasGroup && group.state !== GroupState.CREATED;

  const [selectedSectionId, setSelectedSectionId] = useState<number | null>(null);

  const {
    data: qual,
    isFetching: isQualFetching,
    isLoading: isQualLoading,
    isError: isQualError,
    refetch: refreshQual,
  } = useQualification(groupId, hasStarted);

  const { mutate: startQual, isPending: isQualStarting } = useStartQualification();
  const { mutate: endQual, isPending: isQualEnding } = useEndQualification();

  const handleQualExport = () => {
    console.warn("handleQualExport temporary unavailable");
  };

  useEffect(() => {
    let refreshFn = undefined;
    let exportFn = undefined;
    let completeFn = undefined;
    if (group && !isQualError && !isQualLoading) {
      if (hasStarted) {
        refreshFn = refreshQual;
      }
      if (group.state === GroupState.QUAL_START) {
        completeFn = () => endQual(group.id);
      }
      const exportStates: GroupState[] = Object.values(GroupState).filter(
        (state) => state !== GroupState.CREATED && state !== GroupState.QUAL_START
      );
      if (exportStates.includes(group.state)) {
        exportFn = handleQualExport;
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
  }, [endQual, group, hasStarted, isQualError, isQualLoading, refreshQual, setContext]);

  if (!hasGroup) return;

  if (!hasStarted) {
    return (
      <StartCard
        title="Квалификация еще не началась"
        loading={isQualStarting}
        onStart={() => startQual(group.id)}
      />
    );
  }

  return (
    <PageLoader loading={isQualLoading || isQualEnding} error={isQualError}>
      <Stack>
        {qual && (
          <TableCard loading={isQualFetching}>
            <QualificationTableHead qual={qual} />
            <Table.Tbody>
              {qual.sections
                .sort(({ place: a }, { place: b }) => (a && b ? a - b : 0))
                .map((section) => (
                  <QualificationTableRow
                    key={section.id}
                    section={section}
                    selected={selectedSectionId === section.id}
                    setSelected={setSelectedSectionId}
                  />
                ))}
            </Table.Tbody>
          </TableCard>
        )}
        {selectedSectionId && <SectionTab sectionId={selectedSectionId} />}
      </Stack>
    </PageLoader>
  );
}

function QualificationTableHead({ qual }: { qual: Qualification }) {
  return (
    <Table.Thead>
      <Table.Tr>
        <Table.Th w={50}>Место</Table.Th>
        <Table.Th>Спортсмен</Table.Th>
        {Array(qual.roundCount)
          .fill(0)
          .map((_, index) => (
            <Table.Th w={100} key={`${qual.groupId}$round$${index + 1}`}>
              {`${qual.distance} ${NO_SCORE_VALUE} ${index + 1}`}
            </Table.Th>
          ))}
        <Table.Th w={100}>Итог</Table.Th>
        <Table.Th w={100}>10's</Table.Th>
        <Table.Th w={100}>9's</Table.Th>
        <Table.Th w={120}>Вып. разряд</Table.Th>
      </Table.Tr>
    </Table.Thead>
  );
}

interface QualificationTableRowProps {
  section: QualificationSection;
  selected: boolean;
  setSelected: (id: number) => unknown;
}

function QualificationTableRow({ section, selected, setSelected }: QualificationTableRowProps) {
  const theme = useMantineTheme();
  const rounds = [...section.rounds].sort(({ ordinal: a }, { ordinal: b }) => a - b);
  return (
    <Table.Tr
      onClick={() => setSelected(section.id)}
      bg={selected ? `${theme.colors.yellow[0]}33` : undefined}
    >
      <Table.Td>{section.place ?? NO_SCORE_VALUE}</Table.Td>
      <Table.Td>{section.competitor.fullName}</Table.Td>
      {rounds.map((round) => (
        <Table.Td
          key={`${section.id}$round${round.ordinal}`}
          bg={round.isActive ? `${theme.colors.yellow[8]}33` : undefined}
        >
          {round.totalScore ?? NO_SCORE_VALUE}
        </Table.Td>
      ))}
      <Table.Td>{section.total ?? NO_SCORE_VALUE}</Table.Td>
      <Table.Td>{section.count10 ?? NO_SCORE_VALUE}</Table.Td>
      <Table.Td>{section.count9 ?? NO_SCORE_VALUE}</Table.Td>
      <Table.Td>
        {section.rankGained ? getSportsRankDescription(section.rankGained) : NO_SCORE_VALUE}
      </Table.Td>
    </Table.Tr>
  );
}
