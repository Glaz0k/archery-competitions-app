import { useEffect, useState } from "react";
import { useMutation, useQueries, useQueryClient } from "@tanstack/react-query";
import { Link, Outlet, useLocation, useNavigate, useParams } from "react-router";
import { Group, LoadingOverlay, NativeSelect, Stack, Text, Title } from "@mantine/core";
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
import { DEFAULT_ENUM_VALUE, defaultAndEnumValues } from "../../helper/defaultAndEnumValues";
import useCompetitionForm from "../../hooks/useCompetitionForm";
import { TextButton } from "../buttons/TextButton";
import { MainCard } from "../cards/MainCard";
import PrimaryCard from "../cards/PrimaryCard";
import DeleteCompetitionModal from "../modals/competiton/DeleteCompetitionModal";

export default function CompetitionPage() {
  const { cupId, competitionId } = useParams();
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const [isCompetitionEditing, setCompetitionEditing] = useState(false);

  const [isOpenedCompetitionDel, competitionDelControl] = useDisclosure(false);

  const isCompetitorsPage = useLocation().pathname.endsWith("/competitors");
  const [individualGroupFilter, setIndividualGroupFilter] = useState({
    bow: DEFAULT_ENUM_VALUE,
    identity: DEFAULT_ENUM_VALUE,
    state: DEFAULT_ENUM_VALUE,
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
  const isMainInfoNotFound = mainQuery.some((query) => query.error?.response?.status === 404);

  const { mutate: editCompetition, isPending: isEditedCompetitionSubmitting } = useMutation({
    mutationFn: (editedCompetition) => putCompetition(editedCompetition),
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
    const { stage, ...others } = competition;
    editCompetitionForm.setValues({
      stage: stage.value,
      ...others,
    });
    setCompetitionEditing(true);
  };

  useEffect(() => {
    if (isMainInfoLoadError && isMainInfoNotFound) {
      navigate("/not-found");
    }
  }, [isMainInfoLoadError, isMainInfoNotFound, navigate]);

  useEffect(() => {
    let title = "ArcheryManager - ";
    if (cup && competition) {
      title += cup.title + " | " + competition.stage.textValue;
      if (competition.isEnded) {
        title += " - Завершён";
      }
      if (isCompetitorsPage) {
        title += " | Участники";
      }
    } else {
      title += "Кубок";
    }
    setWebTitle(title);
  }, [cup, competition, isCompetitorsPage]);

  const editCompetitionForm = useCompetitionForm();
  const editCompetitionFormStructure = (
    <>
      <DatePickerInput
        w="100%"
        label="Дата начала"
        key={editCompetitionForm.key("startDate")}
        {...editCompetitionForm.getInputProps("startDate")}
      />
      <DatePickerInput
        w="100%"
        label="Дата окончания"
        key={editCompetitionForm.key("endDate")}
        {...editCompetitionForm.getInputProps("endDate")}
      />
    </>
  );

  return (
    <>
      <DeleteCompetitionModal
        isOpened={isOpenedCompetitionDel}
        onClose={competitionDelControl.close}
        onConfirm={removeCompetition}
        isLoading={isCompetitionDeleting}
      />
      <Group align="start" display="flex" flex={1} style={{ overflow: "hidden" }} gap="lg">
        <Stack gap="md">
          <MainCard
            onEdit={handleCompetitionEditing}
            isEditing={isCompetitionEditing}
            isLoading={isMainInfoLoading || isEditedCompetitionSubmitting}
            onEditSubmit={editCompetitionForm.onSubmit((values) =>
              editCompetition({
                id: competitionId,
                ...values,
              })
            )}
            onEditCancel={() => setCompetitionEditing(false)}
            onDelete={competitionDelControl.open}
          >
            <Stack gap={0}>
              <Title order={2}>{competition?.stage.textValue || "Этап"}</Title>
              <Text>
                {cup?.title || "Название"}
                {cup?.season && ", сезон " + cup.season}
              </Text>
            </Stack>
            {isCompetitionEditing ? (
              editCompetitionFormStructure
            ) : (
              <>
                <DatePickerInput
                  w="100%"
                  disabled
                  label="Дата начала"
                  value={competition?.startDate}
                  onChange={() => {}}
                />
                <DatePickerInput
                  w="100%"
                  disabled
                  label="Дата окончания"
                  value={competition?.endDate}
                  onChange={() => {}}
                />
              </>
            )}
          </MainCard>
          {!isCompetitorsPage && (
            <PrimaryCard>
              <Stack w={300} align="start" pos="relative">
                <NativeSelect
                  w="100%"
                  label="Класс лука"
                  data={defaultAndEnumValues(BowClass).map((bowClass) => {
                    return {
                      label: bowClass.textValue,
                      value: bowClass.value,
                    };
                  })}
                  onChange={(e) =>
                    setIndividualGroupFilter({
                      ...individualGroupFilter,
                      bow: e.currentTarget.value,
                    })
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
            </PrimaryCard>
          )}
          {!isCompetitorsPage && (
            <PrimaryCard>
              <LoadingOverlay visible={isMainInfoLoading || isEditedCompetitionSubmitting} />
              <Stack w={300} align="start" pos="relative">
                <TextButton
                  label="Таблица участников"
                  w="100%"
                  component={Link}
                  to={"/cups/" + cupId + "/competitions/" + competitionId + "/competitors"}
                />
                <TextButton
                  label="Завершить"
                  w="100%"
                  disabled={competition?.isEnded}
                  onClick={endCompetition}
                  loading={isCompetitionEnding}
                />
              </Stack>
            </PrimaryCard>
          )}
        </Stack>
        <Outlet context={individualGroupFilter} />
      </Group>
    </>
  );
}
