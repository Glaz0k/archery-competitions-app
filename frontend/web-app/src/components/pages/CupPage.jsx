import { useEffect, useState } from "react";
import { IconCheck } from "@tabler/icons-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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
import { deleteCompetition } from "../../api/competitions";
import { deleteCup, getCompetitions, getCup, postCompetition, putCup } from "../../api/cups";
import { COMPETITION_QUERY_KEYS, CUP_QUERY_KEYS } from "../../api/queryKeys";
import CompetitionStage from "../../enums/CompetitionStage";
import { formatCompetitionDateRange } from "../../helper/competitons";
import useCupForm from "../../hooks/useCupForm";
import MainBar from "../bars/MainBar";
import { LinkCard, LinkCardSkeleton } from "../cards/LinkCard";
import { MainCard } from "../cards/MainCard";
import EmptyCardSpace from "../misc/EmptyCardSpace";
import AddCompetitionModal from "../modals/competiton/AddCompetitionModal";
import DeleteCompetitionModal from "../modals/competiton/DeleteCompetitionModal";
import DeleteCupModal from "../modals/cup/DeleteCupModal";

const SKELETON_LENGTH = 4;

export default function CupPage() {
  const { cupId } = useParams();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const theme = useMantineTheme();

  const [competitionDeletingId, setCompetitionDeletingId] = useState(null);

  const [isCupEditing, setCupEditing] = useState(false);

  const [isOpenedCompetitionDel, competitionDelControl] = useDisclosure(false);
  const [isOpenedCompetitionAdd, competitionAddControl] = useDisclosure(false);
  const [isOpenedCupDel, cupDelControl] = useDisclosure(false);

  const [webTitle, setWebTitle] = useState(null);
  useDocumentTitle(webTitle);

  const {
    data: cup,
    isFetching: isCupLoading,
    isError: isCupReadError,
    error: cupReadError,
  } = useQuery({
    queryKey: CUP_QUERY_KEYS.element(cupId),
    queryFn: () => getCup(cupId),
    initialData: null,
  });

  const { mutate: editCup, isPending: isEditedCupSubmitting } = useMutation({
    mutationFn: (editedCup) => putCup(editedCup),
    onSuccess: (editedCup) => {
      queryClient.setQueryData(CUP_QUERY_KEYS.element(cupId), editedCup);
      setCupEditing(false);
    },
  });

  const { mutate: removeCup, isPending: isCupDeleting } = useMutation({
    mutationFn: () => deleteCup(cupId),
    onSuccess: () => {
      navigate("/cups");
      cupDelControl.close();
    },
  });

  const { mutateAsync: createCompetition, isPending: isCompetitonSubmitting } = useMutation({
    mutationFn: (newCompetition) => postCompetition(cupId, newCompetition),
    onSuccess: (createdCompetition) => {
      queryClient.setQueryData(COMPETITION_QUERY_KEYS.allByCup(cupId), (old) => [
        createdCompetition,
        ...(old || []),
      ]);
      competitionAddControl.close();
    },
  });

  const {
    data: competitions,
    isFetching: isCompetitionsLoading,
    refetch: refetchCompetitions,
  } = useQuery({
    queryKey: COMPETITION_QUERY_KEYS.allByCup(cupId),
    queryFn: () => getCompetitions(cupId),
    initialData: [],
  });

  const { mutate: removeCompetition, isPending: isCompetitionDeleting } = useMutation({
    mutationFn: () => deleteCompetition(competitionDeletingId),
    onSuccess: () => {
      queryClient.setQueryData(COMPETITION_QUERY_KEYS.allByCup(cupId), (old) => {
        return old.filter((competition) => competition.id !== competitionDeletingId);
      });
      competitionDelControl.close();
      setCompetitionDeletingId(null);
    },
  });

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
    editCupForm.setValues({ ...cup });
    setCupEditing(true);
  };

  useEffect(() => {
    if (isCupReadError && cupReadError.response?.status === 404) {
      navigate("/not-found");
    }
  }, [isCupReadError, cupReadError, navigate]);

  useEffect(() => {
    if (cup) {
      setWebTitle("ArcheryManager - " + cup.title);
    } else {
      setWebTitle("ArcheryManager - Кубок");
    }
  }, [cup]);

  const editCupForm = useCupForm();

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

  return (
    <>
      <DeleteCompetitionModal
        isOpened={isOpenedCompetitionDel}
        onClose={denyCompetitionDeletion}
        onConfirm={removeCompetition}
        isLoading={isCompetitionDeleting}
      />
      <AddCompetitionModal
        isOpened={isOpenedCompetitionAdd}
        onClose={competitionAddControl.close}
        onSubmit={({ stage, ...others }) => {
          createCompetition({
            stage: CompetitionStage.valueOf(stage),
            ...others,
          });
        }}
        isLoading={isCompetitonSubmitting}
      />
      <DeleteCupModal
        isOpened={isOpenedCupDel}
        onClose={cupDelControl.close}
        onConfirm={removeCup}
        isLoading={isCupDeleting}
      />
      <Group align="start" flex={1} gap="lg">
        <MainCard
          onEdit={handleCupEditing}
          isEditing={isCupEditing}
          isLoading={isCupLoading || isEditedCupSubmitting}
          onEditSubmit={editCupForm.onSubmit((values) => editCup({ id: cupId, ...values }))}
          onEditCancel={() => setCupEditing(false)}
          onDelete={cupDelControl.open}
        >
          {isCupEditing ? (
            editCupFormStructure
          ) : (
            <>
              <Title order={2}>{cup?.title || "Название"}</Title>
              <TextInput w="100%" disabled label="Адрес" defaultValue={cup?.address || ""} />
              <TextInput w="100%" disabled label="Сезон" defaultValue={cup?.season || ""} />
            </>
          )}
        </MainCard>
        <Flex direction="column" flex={1} h="100%" gap="lg">
          <MainBar
            title={"Соревнования"}
            onRefresh={refetchCompetitions}
            onAdd={competitionAddControl.open}
            onBack={() => navigate("/cups")}
          />
          <Stack flex={1}>
            {isCompetitionsLoading ? (
              Array(SKELETON_LENGTH)
                .fill(0)
                .map((_, index) => (
                  <LinkCardSkeleton key={index} isTagged isExport isDelete>
                    <Skeleton height={theme.fontSizes.sm} width={250} />
                  </LinkCardSkeleton>
                ))
            ) : competitions.length > 0 ? (
              competitions.map(({ id, stage, startDate, endDate, isEnded }, index) => (
                <LinkCard
                  key={index}
                  title={stage.textValue}
                  to={"/cups/" + cupId + "/competitions/" + id}
                  tag={
                    isEnded ? (
                      <Badge leftSection={<IconCheck />}>
                        <Text tt="capitalize">{"Завершено"}</Text>
                      </Badge>
                    ) : null
                  }
                  onExport={isEnded ? () => handleExport(id) : null}
                  onDelete={() => confirmCompetitionDeletion(id)}
                >
                  <Text size="sm">{formatCompetitionDateRange({ startDate, endDate })}</Text>
                </LinkCard>
              ))
            ) : (
              <EmptyCardSpace label="Соревнования не найдены" />
            )}
          </Stack>
        </Flex>
      </Group>
    </>
  );
}
