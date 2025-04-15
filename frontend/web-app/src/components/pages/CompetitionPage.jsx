import { useEffect, useState } from "react";
import { useMutation, useQueries, useQueryClient } from "@tanstack/react-query";
import { Link, Outlet, useNavigate, useParams } from "react-router";
import { Button, Group, LoadingOverlay, NativeSelect, Stack, Text, Title } from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import { useDisclosure, useDocumentTitle } from "@mantine/hooks";
import {
  deleteCompetition,
  getCompetition,
  postEndCompetition,
  putCompetition,
} from "../../api/competitions";
import { getCup } from "../../api/cups";
import { COMPETITION_QUERY_KEYS, CUP_QUERY_KEYS } from "../../api/queryKeys";
import BowClass from "../../enums/BowClass";
import GroupGender from "../../enums/GroupGender";
import GroupState from "../../enums/GroupState";
import MainCard from "../cards/MainCard";
import DeleteCompetitionModal from "../modals/DeleteCompetitionModal";

function defaultAndEnumValues(enumObj) {
  return [
    {
      value: "default",
      textValue: "Не указано",
    },
    ...Object.values(enumObj),
  ];
}

export default function CompetitionPage() {
  const { cupId, competitionId } = useParams();
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const [isCompetitionEditing, setCompetitionEditing] = useState(false);
  const [editedCompetition, setEditedCompetition] = useState({
    id: null,
    stage: null,
    startDate: null,
    endDate: null,
    isEnded: null,
  });

  const [isOpenedCompetitionDel, competitionDelControl] = useDisclosure(false);

  const [isMainPage, setMainPage] = useState(true);
  const [individualGroupFilter, setIndividualGroupFilter] = useState({
    bow: "default",
    identity: "default",
    state: "default",
  });

  const [webTitle, setWebTitle] = useState(null);
  useDocumentTitle(webTitle);

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

  const { mutate: removeCompetition, isPending: isCompetitionDeleting } = useMutation({
    mutationFn: () => deleteCompetition(competitionId),
    onSuccess: () => {
      navigate("/cups/" + cupId);
      competitionDelControl.close();
    },
  });

  const { mutate: endCompetition, isPending: isCompetitionEnding } = useMutation({
    mutationFn: () => postEndCompetition(competitionId),
    onSuccess: (endedCompetiton) => {
      queryClient.setQueryData(COMPETITION_QUERY_KEYS.element(competitionId), endedCompetiton);
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

  useEffect(() => {
    let title = "ArcheryManager - ";
    if (cup && competition) {
      title += cup.title + "/" + competition.stage.textValue;
      if (competition.isEnded) {
        title += " - Завершён";
      }
    } else {
      title += "Кубок";
    }
    setWebTitle(title);
  }, [cup, competition]);

  if (cup == null || competition == null) {
    return <LoadingOverlay visible={true} />;
  }

  return (
    <>
      <DeleteCompetitionModal
        isOpened={isOpenedCompetitionDel}
        onClose={competitionDelControl.close}
        onConfirm={removeCompetition}
        isLoading={isCompetitionDeleting}
      />

      <LoadingOverlay visible={isMainInfoLoading} />
      <Group align="start" flex={1}>
        <Stack>
          <MainCard
            onBack={() => navigate("../")}
            onEdit={handleCompetitionEditing}
            isEditing={isCompetitionEditing}
            isLoading={isEditedCompetitionSubmitting}
            onEditSubmit={editCompetition}
            onEditCancel={() => setCompetitionEditing(false)}
            onDelete={competitionDelControl.open}
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
          {isMainPage && (
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
                onChange={(e) =>
                  setIndividualGroupFilter({ ...individualGroupFilter, bow: e.currentTarget.value })
                }
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
                onChange={(e) =>
                  setIndividualGroupFilter({
                    ...individualGroupFilter,
                    identity: e.currentTarget.value,
                  })
                }
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
                onChange={(e) =>
                  setIndividualGroupFilter({
                    ...individualGroupFilter,
                    state: e.currentTarget.value,
                  })
                }
              />
            </Stack>
          )}
          <Stack w={300} align="start" pos="relative">
            {isMainPage ? (
              <>
                <Button
                  w="100%"
                  component={Link}
                  to={"/cups/" + cupId + "/competitions/" + competitionId + "/competitors"}
                  onClick={() => setMainPage(false)}
                >
                  {"Таблица участников"}
                </Button>
                <Button
                  w="100%"
                  disabled={competition.isEnded}
                  onClick={endCompetition}
                  loading={isCompetitionEnding}
                >
                  {"Завершить"}
                </Button>
              </>
            ) : (
              <Button
                w="100%"
                component={Link}
                to={"/cups/" + cupId + "/competitions/" + competitionId}
                onClick={() => setMainPage(true)}
              >
                {"Список индивидульных групп"}
              </Button>
            )}
          </Stack>
        </Stack>
        <Outlet context={individualGroupFilter} />
      </Group>
    </>
  );
}
