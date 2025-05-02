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
  const { cupId: paramCupId, competitionId: paramCompetitionId } = useParams();
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
    isFetching: isCupFetching,
    isLoading: isCupLoading,
    isError: isCupError,
    error: cupError,
  } = useCup(cupId);
  const {
    data: competition,
    isFetching: isCompetitionFetching,
    isLoading: isCompetitionLoading,
    isError: isCompetitionError,
    error: competitionError,
  } = useCompetition(competitionId);

  const isMainInfoFetching = isCupFetching || isCompetitionFetching;
  const isMainInfoLoading = isCupLoading || isCompetitionLoading;
  const isMainInfoError = isCupError || isCompetitionError;

  const { mutate: updateCompetition, isPending: isCompetitionUpdating } = useUpdateCompetition(() =>
    setCompetitionEditing(false)
  );

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
    if (cup && competition) {
      if (cup.id !== competition.cupId) {
        navigate("/404");
      }
    }
  });

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
      if (isMainInfoFetching) {
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
  }, [isMainInfoFetching, cup, competition, isCompetitorsPage]);

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
      <DatePickerInput w="100%" disabled label="Дата начала" value={competition?.startDate} />
      <DatePickerInput w="100%" disabled label="Дата окончания" value={competition?.endDate} />
    </>
  );

  const competitionTitle = competition
    ? getCompetitionStageDescription(competition.stage)
    : "Загрузка...";
  const competitionSubtitle = cup
    ? cup.title + (cup.season ? ", сезон " + cup.season : "")
    : "Загрузка...";

  const renderCompetitionPageControls = !isCompetitorsPage ? (
    <>
      <GroupFilterControls />
      <ControlsCard>
        <LoadingOverlay visible={isMainInfoFetching || isCompetitionUpdating} />
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

  if (isMainInfoLoading) {
    return <LoadingOverlay visible />;
  }

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
            loading={isMainInfoFetching || isCompetitionUpdating}
          >
            <Stack gap={0}>
              <Title order={2}>{competitionTitle}</Title>
              <Text>{competitionSubtitle}</Text>
              {cup?.address && <Text size="sm">{cup.address}</Text>}
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
