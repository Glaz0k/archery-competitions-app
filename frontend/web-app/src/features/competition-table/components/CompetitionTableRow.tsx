import { IconCheck, IconMinus, IconX } from "@tabler/icons-react";
import { format } from "date-fns";
import { formatInTimeZone } from "date-fns-tz";
import { ActionIcon, Table, Tooltip, useMantineTheme } from "@mantine/core";
import { useToggleCompetitor } from "../../../api";
import type { CompetitorCompetitionDetail } from "../../../entities";
import {
  getBowClassDescription,
  getGenderDescription,
  getSportsRankDescription,
} from "../../../utils";

export function CompetitionTableHead() {
  return (
    <Table.Thead>
      <Table.Tr>
        <Table.Th />
        <Table.Th>Отметка времени (МСК)</Table.Th>
        <Table.Th>Фамилия, Имя</Table.Th>
        <Table.Th>Дата рождения</Table.Th>
        <Table.Th>Пол</Table.Th>
        <Table.Th>Класс лука</Table.Th>
        <Table.Th>Разряд</Table.Th>
        <Table.Th>Регион</Table.Th>
        <Table.Th>Федерация</Table.Th>
        <Table.Th>Клуб</Table.Th>
        <Table.Th />
      </Table.Tr>
    </Table.Thead>
  );
}

export interface CompetitionTableRowProps {
  detail: CompetitorCompetitionDetail;
  deleting: boolean;
  onDelete: (competitorId: number) => void;
}

export function CompetitionTableRow({ detail, deleting, onDelete }: CompetitionTableRowProps) {
  const theme = useMantineTheme();
  const { mutate: toggle, isPending: isToggling } = useToggleCompetitor();

  const handleToggle = () =>
    toggle([detail.competitionId, detail.competitor.id, { isActive: !detail.isActive }]);

  const competitor = detail.competitor;

  return (
    <Table.Tr
      bg={
        isToggling
          ? `${theme.colors.yellow[8]}33`
          : !detail.isActive
            ? `${theme.colors.red[8]}33`
            : undefined
      }
    >
      <ActivatingTd active={detail.isActive} loading={isToggling} toggleFn={handleToggle} />
      <Table.Td>
        {formatInTimeZone(detail.createdAt, "Europe/Moscow", "dd.MM.yyyy HH:mm:ss")}
      </Table.Td>
      <Table.Td>{competitor.fullName}</Table.Td>
      <Table.Td>{format(competitor.birthDate, "dd.MM.yyyy")}</Table.Td>
      <Table.Td>{getGenderDescription(competitor.identity)}</Table.Td>
      <Table.Td>{competitor.bow ? getBowClassDescription(competitor.bow) : "Не указан"}</Table.Td>
      <Table.Td>{competitor.rank ? getSportsRankDescription(competitor.rank) : "б/р"}</Table.Td>
      <Table.Td>{competitor.region || "Не указан"}</Table.Td>
      <Table.Td>{competitor.federation || "Не указана"}</Table.Td>
      <Table.Td>{competitor.club || "Не указан"}</Table.Td>
      <DeletingTd loading={deleting} deleteFn={() => onDelete(competitor.id)} />
    </Table.Tr>
  );
}

interface ActivatingTdProps {
  active: boolean;
  loading: boolean;
  toggleFn: () => void;
}

function ActivatingTd({ active, loading, toggleFn }: ActivatingTdProps) {
  return (
    <Table.Td>
      {!active ? (
        <Tooltip label="Активировать">
          <ActionIcon variant="transparent" onClick={toggleFn} loading={loading} size="md">
            <IconCheck />
          </ActionIcon>
        </Tooltip>
      ) : (
        <Tooltip label="Деактивировать">
          <ActionIcon variant="transparent" onClick={toggleFn} loading={loading} size="md">
            <IconX />
          </ActionIcon>
        </Tooltip>
      )}
    </Table.Td>
  );
}

interface DeletingTdProps {
  loading: boolean;
  deleteFn: () => void;
}

function DeletingTd({ loading, deleteFn }: DeletingTdProps) {
  return (
    <Table.Td>
      <Tooltip label="Исключить">
        <ActionIcon variant="transparent" loading={loading} onClick={deleteFn} size="md">
          <IconMinus />
        </ActionIcon>
      </Tooltip>
    </Table.Td>
  );
}
