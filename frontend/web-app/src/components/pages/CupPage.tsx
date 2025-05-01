import { useEffect, useState } from "react";
import { IconCheck } from "@tabler/icons-react";
import { isAxiosError } from "axios";
import { useNavigate, useParams } from "react-router";
import {
  Badge,
  Flex,
  Group,
  Skeleton,
  Stack,
  Text,
  TextInput,
  Title,
  useMantineTheme,
} from "@mantine/core";
import { useDisclosure, useDocumentTitle } from "@mantine/hooks";
import {
  useCompetitions,
  useCreateCompetition,
  useCup,
  useDeleteCompetition,
  useDeleteCup,
  useUpdateCup,
} from "../../api";
import { APP_NAME } from "../../constants";
import {
  AddCompetitionModal,
  DeleteCompetitionModal,
  DeleteCupModal,
  useSubmitCupForm,
} from "../../features";
import { formatCompetitionDateRange, getCompetitionStageDescription } from "../../utils";
import {
  CenterCard,
  EntityCard,
  EntityCardSkeleton,
  MainInfoCard,
  SideBar,
  TopBar,
} from "../../widgets";

const SKELETON_LENGTH = 4;

export default function CupPage() {
  const { paramsCupId } = useParams();
  const cupId = Number(paramsCupId);

  const navigate = useNavigate();
  const theme = useMantineTheme();

  const [competitionDeletingId, setCompetitionDeletingId] = useState<number | null>(null);

  const [isCupEditing, setCupEditing] = useState(false);

  const [isOpenedCompetitionDel, competitionDelControl] = useDisclosure(false);
  const [isOpenedCompetitionAdd, competitionAddControl] = useDisclosure(false);
  const [isOpenedCupDel, cupDelControl] = useDisclosure(false);

  const [webTitle, setWebTitle] = useState<string>("");
  useDocumentTitle(webTitle);

  const {
    data: cup,
    isFetching: isCupLoading,
    isError: isCupReadError,
    error: cupReadError,
  } = useCup(cupId);

  const {
    data: competitions,
    isFetching: isCompetitionsLoading,
    refetch: refetchCompetitions,
    isError: isCompetitionsError,
    error: competitionsError,
  } = useCompetitions(cupId);

  const { mutate: updateCup, isPending: isEditedCupSubmitting } = useUpdateCup(() => {
    setCupEditing(false);
  });

  const { mutate: deleteCup, isPending: isCupDeleting } = useDeleteCup(() => {
    cupDelControl.close();
    navigate("..");
  });

  const { mutateAsync: createCompetition, isPending: isCompetitonSubmitting } =
    useCreateCompetition(() => {
      competitionAddControl.close();
    });

  const { mutate: removeCompetition, isPending: isCompetitionDeleting } = useDeleteCompetition(
    () => {
      competitionDelControl.close();
      setCompetitionDeletingId(null);
    }
  );

  const handleExport = (id: number) => {
    console.warn(`handleExport temporary unavailable ${id}`);
  };

  const confirmCompetitionDeletion = (competitionId: number) => {
    setCompetitionDeletingId(competitionId);
    competitionDelControl.open();
  };

  const denyCompetitionDeletion = () => {
    competitionDelControl.close();
    setCompetitionDeletingId(null);
  };

  const handleCupEditing = () => {
    editCupForm.setValues({
      title: cup?.title || "",
      address: cup?.address || "",
      season: cup?.season || "",
    });
    setCupEditing(true);
  };

  useEffect(() => {
    if (isCupReadError) {
      if (isAxiosError(cupReadError) && cupReadError.status === 404) {
        navigate("/404");
      }
    }
  }, [isCupReadError, cupReadError, navigate]);

  useEffect(() => {
    if (cup) {
      setWebTitle(APP_NAME + " - " + cup.title);
    } else if (isCupLoading) {
      setWebTitle(APP_NAME + " - Загрузка...");
    } else {
      setWebTitle(APP_NAME);
    }
  }, [cup, isCupLoading]);

  const editCupForm = useSubmitCupForm();

  const editCupFormStructure = (
    <>
      <TextInput
        w="100%"
        label="Название"
        key={editCupForm.key("title")}
        {...editCupForm.getInputProps("title")}
      />
      <TextInput
        w="100%"
        label="Адрес"
        key={editCupForm.key("address")}
        {...editCupForm.getInputProps("address")}
      />
      <TextInput
        w="100%"
        label="Сезон"
        key={editCupForm.key("season")}
        {...editCupForm.getInputProps("season")}
      />
    </>
  );

  let renderContent;
  if (isCompetitionsLoading) {
    renderContent = Array(SKELETON_LENGTH)
      .fill(0)
      .map((_, index) => (
        <EntityCardSkeleton key={index} tagged exported deleted>
          <Skeleton height={theme.fontSizes.sm} width={250} />
        </EntityCardSkeleton>
      ));
  } else if (isCompetitionsError) {
    console.error(competitionsError.name + "\n" + competitionsError.message);
    renderContent = <CenterCard label="Произошла ошибка" />;
  } else if (competitions.length === 0) {
    renderContent = <CenterCard label="Соревнования не найдены" />;
  } else {
    renderContent = competitions.map(({ id, stage, startDate, endDate, isEnded }) => (
      <EntityCard
        key={id}
        title={getCompetitionStageDescription(stage)}
        to={String(id)}
        tag={
          isEnded ? (
            <Badge leftSection={<IconCheck />} color={"green.8"}>
              <Text tt="capitalize">{"Завершено"}</Text>
            </Badge>
          ) : undefined
        }
        onExport={isEnded ? () => handleExport(id) : undefined}
        onDelete={() => confirmCompetitionDeletion(id)}
      >
        <Text size="sm">{formatCompetitionDateRange(startDate, endDate)}</Text>
      </EntityCard>
    ));
  }

  return (
    <>
      <DeleteCompetitionModal
        opened={isOpenedCompetitionDel}
        onClose={denyCompetitionDeletion}
        onConfirm={() => removeCompetition(competitionDeletingId!)}
        loading={isCompetitionDeleting}
      />
      <AddCompetitionModal
        opened={isOpenedCompetitionAdd}
        onClose={competitionAddControl.close}
        onSubmit={(values) => createCompetition([cupId, values])}
        loading={isCompetitonSubmitting}
      />
      <DeleteCupModal
        opened={isOpenedCupDel}
        onClose={cupDelControl.close}
        onConfirm={() => deleteCup(cupId)}
        loading={isCupDeleting}
      />
      <Group align="start" flex={1} gap="lg">
        <SideBar>
          <MainInfoCard
            onEdit={handleCupEditing}
            onFormSubmit={editCupForm.onSubmit((values) => updateCup([cupId, values]))}
            onCancel={() => setCupEditing(false)}
            onDelete={cupDelControl.open}
            editing={isCupEditing}
            loading={isCupLoading || isEditedCupSubmitting}
          >
            {isCupEditing ? (
              editCupFormStructure
            ) : (
              <>
                <Title order={2}>{cup?.title || "Загрузка..."}</Title>
                <TextInput
                  w="100%"
                  disabled
                  label="Адрес"
                  defaultValue={cup?.address || "Не указано"}
                />
                <TextInput
                  w="100%"
                  disabled
                  label="Сезон"
                  defaultValue={cup?.season || "Не указано"}
                />
              </>
            )}
          </MainInfoCard>
        </SideBar>
        <Flex direction="column" flex={1} h="100%" gap="lg">
          <TopBar
            title={"Соревнования"}
            onRefresh={refetchCompetitions}
            onAdd={competitionAddControl.open}
            onBack={() => navigate("..")}
          />
          <Stack flex={1}>{renderContent}</Stack>
        </Flex>
      </Group>
    </>
  );
}
