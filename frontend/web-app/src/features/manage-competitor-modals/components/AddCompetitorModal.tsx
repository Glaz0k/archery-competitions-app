import { useState } from "react";
import { IconUser, IconUserPlus } from "@tabler/icons-react";
import { format } from "date-fns";
import {
  ActionIcon,
  Card,
  Group,
  Modal,
  Select,
  Stack,
  Text,
  ThemeIcon,
  Tooltip,
  type ComboboxData,
} from "@mantine/core";
import {
  useAddCompetitorToCompetition,
  useCompetitionCompetitors,
  useCompetitors,
} from "../../../api";
import type { Competitor } from "../../../entities";
import { getGenderDescription } from "../../../utils";
import { CancelButton, ConfirmButton } from "../../../widgets";
import { useAddCompetitorModalsStack } from "../hooks/useAddCompetitorModalsStack";
import { RegisterCompetitorModal } from "./RegisterCompetitorModal";

export interface AddCompetitorModalProps {
  stack: ReturnType<typeof useAddCompetitorModalsStack>;
  competitionId: number;
}

export function AddCompetitorModal({ competitionId, stack }: AddCompetitorModalProps) {
  const [pickedCompetitor, setPickedCompetitor] = useState<Competitor | null>(null);

  const { data: allCompetitors } = useCompetitors();

  const { data: addedDetails } = useCompetitionCompetitors(competitionId);

  const { mutate: addCompetitor, isPending: isCompetitorAdding } = useAddCompetitorToCompetition(
    stack.closeAll
  );

  const addedCompetitors = new Set<number>(addedDetails.map((val) => val.competitor.id));
  const pickableCompetitors = allCompetitors.filter(({ id }) => !addedCompetitors.has(id));
  const selectData: ComboboxData = pickableCompetitors.map(({ id, fullName }) => ({
    label: fullName,
    value: String(id),
  }));

  const handleChange = (value: string | null) => {
    setPickedCompetitor(value ? pickableCompetitors.find(({ id }) => id === Number(value))! : null);
  };

  const renderPicked = (
    <Card p={0}>
      <Group gap="sm" p="sm">
        <ThemeIcon size="sm">
          <IconUser />
        </ThemeIcon>
        <Stack gap={0} p={0}>
          {pickedCompetitor ? (
            <>
              <Text size="sm">{pickedCompetitor.fullName}</Text>
              <Text size="sm">
                {getGenderDescription(pickedCompetitor.identity) +
                  ", " +
                  format(pickedCompetitor.birthDate, "dd.MM.yyyy")}
              </Text>
            </>
          ) : (
            <Text>{"Участник не выбран"}</Text>
          )}
        </Stack>
      </Group>
    </Card>
  );

  const renderControls = (
    <Group flex={1}>
      <Tooltip label="Регистрация" position="left">
        <ActionIcon
          variant="filled"
          radius="md"
          size="lg"
          color="dark.3"
          loading={isCompetitorAdding}
          onClick={() => {
            stack.open("register-competitor");
          }}
        >
          <IconUserPlus />
        </ActionIcon>
      </Tooltip>
      <Group flex={1} justify="flex-end">
        <CancelButton label="Отменить" onClick={stack.closeAll} loading={isCompetitorAdding} />
        <ConfirmButton
          label="Добавить"
          disabled={!pickedCompetitor}
          loading={isCompetitorAdding}
          onClick={() => {
            addCompetitor([competitionId, { id: pickedCompetitor!.id }]);
          }}
        />
      </Group>
    </Group>
  );

  return (
    <Modal.Stack>
      <Modal title="Добавить участника в соревнование" {...stack.register("add-competitor")}>
        <Stack gap="md">
          <Select
            label="Выберите участника"
            placeholder="Фамилия Имя"
            limit={5}
            searchable
            clearable
            value={pickedCompetitor?.id.toString() ?? null}
            data={selectData}
            onChange={handleChange}
          />
          {renderPicked}
          {renderControls}
        </Stack>
      </Modal>
      <RegisterCompetitorModal {...stack.register("register-competitor")} />
    </Modal.Stack>
  );
}
