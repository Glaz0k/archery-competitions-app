import { useEffect, useState } from "react";
import {
  IconArrowLeft,
  IconCheck,
  IconEdit,
  IconPlus,
  IconRefresh,
  IconTrashX,
  IconX,
} from "@tabler/icons-react";
import { useMutation, useQueries, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";
import {
  ActionIcon,
  Button,
  Group,
  LoadingOverlay,
  NativeSelect,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import { getCompetition, putCompetition } from "../../api/competitions";
import { getCup } from "../../api/cups";
import { COMPETITION_QUERY_KEYS, CUP_QUERY_KEYS } from "../../api/queryKeys";
import BowClass from "../../enums/BowClass";
import GroupGender from "../../enums/GroupGender";
import GroupState from "../../enums/GroupState";
import MainBar from "../bars/MainBar";
import MainCard from "../cards/MainCard";

const defaultAndEnumValues = (enumObj) => {
  return [
    {
      value: null,
      textValue: "Не указано",
    },
    ...Object.values(enumObj),
  ];
};

export default function CompetitionPage() {
  const { cupId, competitionId } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [isCompetitionEditing, setCompetitionEditing] = useState(false);
  const [editedCompetition, setEditedCompetition] = useState({
    id: null,
    stage: null,
    startDate: null,
    endDate: null,
    isEnded: null,
  });

  const mainQuery = useQueries({
    queries: [
      {
        queryKey: CUP_QUERY_KEYS.element(cupId),
        queryFn: () => getCup(cupId),
        initialData: null,
      },
      {
        queryKey: COMPETITION_QUERY_KEYS.element(competitionId),
        queryFn: () => getCompetition(competitionId),
        initialData: null,
      },
    ],
  });
  const [{ data: cup }, { data: competition }] = mainQuery;
  const isMainInfoLoading = mainQuery.some((query) => query.isFetching);
  const isMainInfoLoadError = mainQuery.some((query) => query.isError);

  const { mutate: editCompetition, isPending: isEditedCompetitionSubmitting } = useMutation({
    mutationFn: () => putCompetition(editedCompetition),
    onSuccess: (editedCompetition) => {
      queryClient.setQueryData(COMPETITION_QUERY_KEYS.element(competitionId), editedCompetition);
      setCompetitionEditing(false);
    },
  });

  const handleCompetitionEditing = () => {
    setEditedCompetition({ ...competition });
    setCompetitionEditing(true);
  };

  useEffect(() => {
    if (isMainInfoLoadError && mainQuery.some((query) => query.error.response?.status === 404)) {
      navigate("/not-found");
    }
  }, [isMainInfoLoadError, mainQuery, navigate]);

  if (cup == null || competition == null) {
    return <LoadingOverlay visible={true} />;
  }

  return (
    <>
      <LoadingOverlay visible={isMainInfoLoading} />
      <Group align="start" flex={1}>
        <Stack>
          <MainCard
            onBack={() => navigate("/cups/" + cupId)}
            onEdit={handleCompetitionEditing}
            isEditing={isCompetitionEditing}
            isLoading={isEditedCompetitionSubmitting}
            onEditSubmit={editCompetition}
            onEditCancel={() => setCompetitionEditing(false)}
            onDelete={() => console.log("TODO")}
          >
            <Title order={2}>{competition.stage.textValue}</Title>
            <Text>
              {cup.title}, сезон {cup.season ? cup.season : "не указан"}
            </Text>
            <DatePickerInput
              w="100%"
              label="Дата начала"
              disabled={!isCompetitionEditing}
              value={isCompetitionEditing ? editedCompetition.startDate : competition.startDate}
              onChange={(date) => setEditedCompetition({ ...editedCompetition, startDate: date })}
            />
            <DatePickerInput
              w="100%"
              label="Дата окончания"
              disabled={!isCompetitionEditing}
              value={isCompetitionEditing ? editedCompetition.endDate : competition.endDate}
              onChange={(date) => setEditedCompetition({ ...editedCompetition, endDate: date })}
            />
          </MainCard>
          <Stack w={300} align="start" pos="relative">
            <NativeSelect
              w="100%"
              label="Тип лука"
              data={defaultAndEnumValues(BowClass).map((bowClass) => {
                return {
                  label: bowClass.textValue,
                  value: bowClass.value,
                };
              })}
            />
            <NativeSelect
              w="100%"
              label="Пол"
              data={defaultAndEnumValues(GroupGender).map((gender) => {
                return {
                  label: gender.textValue,
                  value: gender.value,
                };
              })}
            />
            <NativeSelect
              w="100%"
              label="Состояние"
              data={defaultAndEnumValues(GroupState).map((state) => {
                return {
                  label: state.textValue,
                  value: state.value,
                };
              })}
            />
          </Stack>
          <Stack w={300} align="start" pos="relative">
            <Button w="100%">Таблица участников</Button>
            <Button w="100%">Завершить</Button>
          </Stack>
        </Stack>
        <Stack flex={1}>
          <MainBar title={"Индивидуальные группы"}>
            <ActionIcon>
              <IconRefresh />
            </ActionIcon>
          </MainBar>
        </Stack>
      </Group>
    </>
  );
}
