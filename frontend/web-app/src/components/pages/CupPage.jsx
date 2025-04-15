import { useEffect, useState } from "react";
import { IconCheck, IconRefresh } from "@tabler/icons-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";
import {
  ActionIcon,
  Badge,
  Flex,
  Group,
  LoadingOverlay,
  rem,
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
import { formatCompetitionDateRange } from "../../helper/competitons";
import MainBar from "../bars/MainBar";
import { LinkCard, LinkCardSkeleton } from "../cards/LinkCard";
import MainCard from "../cards/MainCard";
import EmptyCardSpace from "../misc/EmptyCardSpace";
import AddCompetitionModal from "../modals/AddCompetitionModal";
import DeleteCompetitionModal from "../modals/DeleteCompetitionModal";
import DeleteCupModal from "../modals/DeleteCupModal";

const SKELETON_LENGTH = 4;

export default function CupPage() {
  const { cupId } = useParams();
  const queryClient = useQueryClient();
  const theme = useMantineTheme();
  const navigate = useNavigate();

  const [competitionDeletingId, setCompetitionDeletingId] = useState(null);

  const [isCupEditing, setCupEditing] = useState(false);
  const [editedCup, setEditedCup] = useState({
    id: null,
    title: null,
    address: null,
    season: null,
  });

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
    mutationFn: () => putCup(editedCup),
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
    setEditedCup({ ...cup });
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

  if (cup == null) {
    return <LoadingOverlay visible={true} />;
  }

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
        onSubmit={createCompetition}
        isLoading={isCompetitonSubmitting}
      />
      <DeleteCupModal
        isOpened={isOpenedCupDel}
        onClose={cupDelControl.close}
        onConfirm={removeCup}
        isLoading={isCupDeleting}
      />
      <LoadingOverlay visible={isCupLoading} />
      <Group align="start" flex={1}>
        <MainCard
          onBack={() => navigate("/cups")}
          onEdit={handleCupEditing}
          isEditing={isCupEditing}
          isLoading={isEditedCupSubmitting}
          onEditSubmit={editCup}
          onEditCancel={() => setCupEditing(false)}
          onDelete={cupDelControl.open}
        >
          {isCupEditing ? (
            <TextInput
              w="100%"
              label="Название"
              value={editedCup.title}
              onChange={(e) => setEditedCup({ ...editedCup, title: e.currentTarget.value })}
            />
          ) : (
            <Title order={2}>{cup.title}</Title>
          )}
          <TextInput
            w="100%"
            disabled={!isCupEditing}
            label="Адрес"
            value={isCupEditing ? editedCup.address : cup.address}
            onChange={(e) => setEditedCup({ ...editedCup, address: e.currentTarget.value })}
          />
          <TextInput
            w="100%"
            disabled={!isCupEditing}
            label="Сезон"
            value={isCupEditing ? editedCup.season : cup.season}
            onChange={(e) => setEditedCup({ ...editedCup, season: e.currentTarget.value })}
          />
        </MainCard>
        <Flex direction="column" flex={1} h="100%">
          <MainBar
            title={"Соревнования"}
            onRefresh={refetchCompetitions}
            onAdd={competitionAddControl.open}
          />
          <Stack flex={1}>
            {isCompetitionsLoading ? (
              Array(SKELETON_LENGTH)
                .fill(0)
                .map((_, index) => (
                  <LinkCardSkeleton key={index} isTagged isExport isDelete>
                    <Skeleton height={rem(theme.fontSizes.md)} width={250} />
                  </LinkCardSkeleton>
                ))
            ) : competitions.length > 0 ? (
              competitions.map(({ id, stage, startDate, endDate, isEnded }, index) => (
                <LinkCard
                  key={index}
                  title={stage.textValue}
                  to={"/cups/" + cupId + "/competitions/" + id}
                  tag={isEnded ? <Badge leftSection={<IconCheck />}>{"Завершено"}</Badge> : null}
                  onExport={isEnded ? () => handleExport(id) : null}
                  onDelete={() => confirmCompetitionDeletion(id)}
                >
                  <Text>{formatCompetitionDateRange({ startDate, endDate })}</Text>
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
