import { useEffect, useState } from "react";
import { isAxiosError } from "axios";
import { Link, Outlet, useLocation, useNavigate, useParams } from "react-router";
import { Group, LoadingOverlay, Stack, Text, Title } from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import { useDisclosure, useDocumentTitle } from "@mantine/hooks";
import {
  useCompetition,
  useCup,
  useDeleteCompetition,
  useEndCompetition,
  useUpdateCompetition,
} from "../../api";
import { APP_NAME } from "../../constants";
import {
  DeleteCompetitionModal,
  GroupFilterControls,
  GroupFilterProvider,
  useEditCompetitionForm,
} from "../../features";
import { getCompetitionStageDescription } from "../../utils";
import { ControlsCard, MainInfoCard, SideBar, TextButton } from "../../widgets";

export default function CompetitionPage() {
  const { paramCupId, paramCompetitionId } = useParams();
  const cupId = Number(paramCupId);
  const competitionId = Number(paramCompetitionId);

  const navigate = useNavigate();

  const [isCompetitionEditing, setCompetitionEditing] = useState(false);
  const [isOpenedCompetitionDel, competitionDelControl] = useDisclosure(false);
  const isCompetitorsPage = useLocation().pathname.endsWith("/competitors");

  const [webTitle, setWebTitle] = useState("");
  useDocumentTitle(webTitle);

  const {
    data: cup,
    isFetching: isCupLoading,
    isError: isCupError,
    error: cupError,
  } = useCup(cupId);
  const {
    data: competition,
    isFetching: isCompetitionLoading,
    isError: isCompetitionError,
    error: competitionError,
  } = useCompetition(cupId);

  const isMainInfoLoading = isCupLoading || isCompetitionLoading;
  const isMainInfoError = isCupError || isCompetitionError;

  const { mutate: updateCompetition, isPending: isCompetitionUpdating } = useUpdateCompetition();

  const { mutate: deleteCompetition, isPending: isCompetitionDeleting } = useDeleteCompetition(
    () => {
      competitionDelControl.close();
      navigate("..");
    }
  );

  const { mutate: endCompetition, isPending: isCompetitionEnding } = useEndCompetition();

  const handleCompetitionEdit = () => {
    editCompetitionForm.setValues({
      startDate: competition?.startDate,
      endDate: competition?.endDate,
    });
    setCompetitionEditing(true);
  };

  useEffect(() => {
    if (
      isMainInfoError &&
      ((isAxiosError(cupError) && cupError.status === 404) ||
        (isAxiosError(competitionError) && competitionError.status === 404))
    ) {
      navigate("/404");
    }
  }, [isMainInfoError, cupError, competitionError, navigate]);

  useEffect(() => {
    const titleFn = () => {
      let base = `${APP_NAME} - `;
      if (isMainInfoLoading) {
        return base + "Загрузка...";
      }
      if (cup && competition) {
        base += `${cup.title} | ${getCompetitionStageDescription(competition.stage)}`;
        if (isCompetitorsPage) {
          base += " | Участники";
        }
        return base;
      }
      return base + "Ошибка";
    };
    setWebTitle(titleFn());
  }, [isMainInfoLoading, cup, competition, isCompetitorsPage]);

  const editCompetitionForm = useEditCompetitionForm();
  const renderEditCompetitionForm = isCompetitionEditing ? (
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
  ) : (
    <>
      <DatePickerInput
        w="100%"
        disabled
        label="Дата начала"
        defaultValue={competition?.startDate}
      />
      <DatePickerInput
        w="100%"
        disabled
        label="Дата окончания"
        defaultValue={competition?.endDate}
      />
    </>
  );

  const competitionTitle = competition
    ? getCompetitionStageDescription(competition.stage)
    : "Загрузка...";
  const competitionSubtitle = cup
    ? cup.title + cup.season
      ? ", сезон " + cup.season
      : ""
    : "Загрузка...";

  const renderCompetitionPageControls = !isCompetitorsPage ? (
    <>
      <GroupFilterControls />
      <ControlsCard>
        <LoadingOverlay visible={isMainInfoLoading || isCompetitionUpdating} />
        <Stack pos="relative" justify="stretch">
          <TextButton label="Таблица участников" component={Link} to="competitors" />
          <TextButton
            label="Завершить"
            disabled={competition?.isEnded}
            onClick={() => endCompetition(competitionId)}
            loading={isCompetitionEnding}
          />
        </Stack>
      </ControlsCard>
    </>
  ) : undefined;

  return (
    <GroupFilterProvider>
      <DeleteCompetitionModal
        opened={isOpenedCompetitionDel}
        onClose={competitionDelControl.close}
        onConfirm={() => deleteCompetition(competitionId)}
        loading={isCompetitionDeleting}
      />
      <Group align="start" display="flex" flex={1} style={{ overflow: "hidden" }} gap="lg">
        <SideBar>
          <MainInfoCard
            onEdit={handleCompetitionEdit}
            onFormSubmit={editCompetitionForm.onSubmit((values) =>
              updateCompetition([competitionId, values])
            )}
            onCancel={() => setCompetitionEditing(false)}
            onDelete={competitionDelControl.open}
            editing={isCompetitionEditing}
            loading={isMainInfoLoading || isCompetitionUpdating}
          >
            <Stack gap={0}>
              <Title order={2}>{competitionTitle}</Title>
              <Text>{competitionSubtitle}</Text>
            </Stack>
            {renderEditCompetitionForm}
          </MainInfoCard>
          {renderCompetitionPageControls}
        </SideBar>
        <Outlet />
      </Group>
    </GroupFilterProvider>
  );
}
