import { useEffect, useState } from "react";
import { IconCheck } from "@tabler/icons-react";
import axios from "axios";
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
import { formatCompetitionDateRange } from "../../helper/formatCompetitionDateRange";
import { useSubmitCupForm } from "../../hooks";
import { getCompetitionStageDescription } from "../../utils";
import MainBar from "../bars/MainBar";
import { LinkCard, LinkCardSkeleton } from "../cards/LinkCard";
import { MainCard } from "../cards/MainCard";
import NotFoundCard from "../cards/NotFoundCard";
import AddCompetitionModal from "../modals/competiton/AddCompetitionModal";
import DeleteCompetitionModal from "../modals/competiton/DeleteCompetitionModal";
import DeleteCupModal from "../modals/cup/DeleteCupModal";

const SKELETON_LENGTH = 4;

export default function CupPage() {
  const { cupId } = useParams();
  const navigate = useNavigate();
  const theme = useMantineTheme();

  const [competitionDeletingId, setCompetitionDeletingId] = useState(null);

  const [isCupEditing, setCupEditing] = useState(false);

  const [isOpenedCompetitionDel, competitionDelControl] = useDisclosure(false);
  const [isOpenedCompetitionAdd, competitionAddControl] = useDisclosure(false);
  const [isOpenedCupDel, cupDelControl] = useDisclosure(false);

  const [webTitle, setWebTitle] = useState("");
  useDocumentTitle(webTitle);

  const {
    data: cup,
    isFetching: isCupLoading,
    isError: isCupReadError,
    error: cupReadError,
  } = useCup(Number(cupId));

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

  const {
    data: competitions,
    isFetching: isCompetitionsLoading,
    refetch: refetchCompetitions,
    isError: isCompetitionsError,
    error: competitionsError,
  } = useCompetitions(Number(cupId), !!cup);

  const { mutate: removeCompetition, isPending: isCompetitionDeleting } = useDeleteCompetition(
    () => {
      competitionDelControl.close();
      setCompetitionDeletingId(null);
    }
  );

  const handleExport = (_id) => {
    console.warn("handleExport temporary unavailable");
  };

  const confirmCompetitionDeletion = (competitionId) => {
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
      if (axios.isAxiosError(cupReadError) && cupReadError.status === 404) {
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
        <LinkCardSkeleton key={index} isTagged isExport isDelete>
          <Skeleton height={theme.fontSizes.sm} width={250} />
        </LinkCardSkeleton>
      ));
  } else if (isCompetitionsError) {
    console.error(competitionsError.name + "\n" + competitionsError.message);
    renderContent = <NotFoundCard label="Произошла ошибка" />;
  } else if (competitions.length === 0) {
    renderContent = <NotFoundCard label="Соревнования не найдены" />;
  } else {
    renderContent = competitions.map(({ id, stage, startDate, endDate, isEnded }, index) => (
      <LinkCard
        key={cupId + "$" + index}
        title={getCompetitionStageDescription(stage)}
        to={id}
        tag={
          isEnded ? (
            <Badge leftSection={<IconCheck />} color={"green.8"}>
              <Text tt="capitalize">{"Завершено"}</Text>
            </Badge>
          ) : undefined
        }
        onExport={isEnded ? () => handleExport(id) : null}
        onDelete={() => confirmCompetitionDeletion(id)}
      >
        <Text size="sm">{formatCompetitionDateRange({ startDate, endDate })}</Text>
      </LinkCard>
    ));
  }

  return (
    <>
      <DeleteCompetitionModal
        isOpened={isOpenedCompetitionDel}
        onClose={denyCompetitionDeletion}
        onConfirm={() => removeCompetition(Number(competitionDeletingId))}
        isLoading={isCompetitionDeleting}
      />
      <AddCompetitionModal
        opened={isOpenedCompetitionAdd}
        onClose={competitionAddControl.close}
        onSubmit={(values) => createCompetition([Number(cupId), values])}
        loading={isCompetitonSubmitting}
      />
      <DeleteCupModal
        isOpened={isOpenedCupDel}
        onClose={cupDelControl.close}
        onConfirm={() => deleteCup(Number(cupId))}
        isLoading={isCupDeleting}
      />
      <Group align="start" flex={1} gap="lg">
        <MainCard
          onEdit={handleCupEditing}
          isEditing={isCupEditing}
          isLoading={isCupLoading || isEditedCupSubmitting}
          onEditSubmit={editCupForm.onSubmit((values) => updateCup([Number(cupId), values]))}
          onEditCancel={() => setCupEditing(false)}
          onDelete={cupDelControl.open}
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
                value={cup?.address || ""}
                onChange={() => {}}
              />
              <TextInput
                w="100%"
                disabled
                label="Сезон"
                value={cup?.season || ""}
                onChange={() => {}}
              />
            </>
          )}
        </MainCard>
        <Flex direction="column" flex={1} h="100%" gap="lg">
          <MainBar
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
