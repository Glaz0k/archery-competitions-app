import { useEffect, useState } from "react";
import {
  IconArrowLeft,
  IconCheck,
  IconEdit,
  IconPlus,
  IconRefresh,
  IconX,
} from "@tabler/icons-react";
import { useNavigate, useParams } from "react-router";
import {
  ActionIcon,
  Badge,
  Button,
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
import { useDisclosure } from "@mantine/hooks";
import { deleteCompetition } from "../../api/competitions";
import { getCompetitions, getCup, postCompetition, putCup } from "../../api/cups";
import { formatCompetitionDateRange, stageToTitle } from "../../helper/competitons";
import { LinkCard, LinkCardSkeleton } from "../cards/LinkCard";
import EmptyCardSpace from "../misc/EmptyCardSpace";
import { CompetitionAddModal, CompetitionDeleteModal } from "../modals/CompetitionModals";

export default function CupPage() {
  const { cupId } = useParams();
  const [cup, setCup] = useState(null);
  const [cupLoading, setCupLoading] = useState(true);

  const [isOpenedCompetitionDel, competitionDelControl] = useDisclosure(false);
  const [competitionDeletion, setCompetitionDeletion] = useState(false);
  const [competitionDeletionId, setCompetitionDeletionId] = useState(null);

  const [isOpenedCompetitionAdd, competitionAddControl] = useDisclosure(false);
  const [competitionSubmission, setCompetitionSubmission] = useState(false);

  const [isCupEditing, setCupEditing] = useState(false);
  const [editedCup, setEditedCup] = useState({
    id: null,
    title: null,
    address: null,
    season: null,
  });
  const [isEditedCupSubmitting, setEditedCupSubmitting] = useState(false);

  const [competitions, setCompetitions] = useState([]);
  const [competitionsLoading, setCompetitionsLoading] = useState(true);

  const navigate = useNavigate();

  const theme = useMantineTheme();

  const readCup = async (cupId) => {
    try {
      setCupLoading(true);
      const data = await getCup(cupId);
      setCup(data);
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setCupLoading(false);
    }
    return true;
  };

  const editCup = async () => {
    try {
      setEditedCupSubmitting(true);
      const data = await putCup(editedCup);
      setCup(data);
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setEditedCupSubmitting(false);
    }
    return true;
  };

  const readCompetitions = async (cupId) => {
    try {
      setCompetitionsLoading(true);
      const data = await getCompetitions(cupId);
      setCompetitions(data);
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setCompetitionsLoading(false);
    }
    return true;
  };

  const removeCompetition = async () => {
    try {
      setCompetitionDeletion(true);
      await deleteCompetition(competitionDeletionId);
      setCompetitions(
        competitions.filter((competition) => competition.id !== competitionDeletionId)
      );
    } catch (err) {
      console.log(err);
      return false;
    } finally {
      setCompetitionDeletion(false);
    }
    return true;
  };

  const createCompetition = async (newCompetition) => {
    try {
      setCompetitionSubmission(true);
      const createdCompetition = await postCompetition({
        cupId: cupId,
        ...newCompetition,
      });
      setCompetitions([createdCompetition, ...competitions]);
      competitionAddControl.close();
    } catch (err) {
      console.log(err);
      return false;
    } finally {
      setCompetitionSubmission(false);
    }
    return true;
  };

  const handleRefresh = () => {
    readCompetitions(cupId);
  };

  const handleExport = () => {
    console.warn("handleExport temporary unavailable");
  };

  const handleCompetitionDeletion = async () => {
    if (await removeCompetition()) {
      competitionDelControl.close();
    }
  };

  const confirmCompetitionDeletion = (competitionId) => {
    setCompetitionDeletionId(competitionId);
    competitionDelControl.open();
  };

  const denyCompetitionDeletion = () => {
    setCompetitionDeletionId(null);
  };

  const handleCompetitionSubmission = async ({ stage, startDate, endDate }) => {
    const newCompetition = {
      stage: stage,
      startDate: startDate,
      endDate: endDate,
    };
    await createCompetition(newCompetition);
  };

  const handleCupEditing = () => {
    setEditedCup({ ...cup });
    setCupEditing(true);
  };

  const handleCupEditingSubmit = async () => {
    if (await editCup()) {
      setCupEditing(false);
    }
  };

  useEffect(() => {
    readCup(cupId);
    readCompetitions(cupId);
  }, [cupId]);

  const skeletonLength = 4;

  return (
    <>
      <CompetitionDeleteModal
        opened={isOpenedCompetitionDel}
        onClose={competitionDelControl.close}
        onDeny={denyCompetitionDeletion}
        onConfirm={handleCompetitionDeletion}
        loading={competitionDeletion}
      />

      <CompetitionAddModal
        opened={isOpenedCompetitionAdd}
        onClose={competitionAddControl.close}
        handleSubmit={handleCompetitionSubmission}
        loading={competitionSubmission}
      />

      <LoadingOverlay visible={cupLoading} />
      {cup && (
        <Group align="start" flex={1}>
          <Stack w={300} align="start" pos="relative">
            <LoadingOverlay visible={isEditedCupSubmitting} />
            <Button
              onClick={() => {
                navigate("/cups");
              }}
              leftSection={<IconArrowLeft />}
            >
              Назад
            </Button>
            <Title order={2} contentEditable={isCupEditing}>
              {isCupEditing ? editedCup.title : cup.title}
            </Title>
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
            <Group>
              {!isCupEditing && (
                <ActionIcon onClick={handleCupEditing}>
                  <IconEdit />
                </ActionIcon>
              )}
              {isCupEditing && (
                <>
                  <ActionIcon onClick={handleCupEditingSubmit}>
                    <IconCheck />
                  </ActionIcon>
                  <ActionIcon onClick={() => setCupEditing(false)}>
                    <IconX />
                  </ActionIcon>
                </>
              )}
            </Group>
          </Stack>
          <Stack flex={1}>
            <Group>
              <Title order={2} flex={1}>
                Соревнования
              </Title>
              <ActionIcon onClick={handleRefresh}>
                <IconRefresh />
              </ActionIcon>
              <ActionIcon onClick={competitionAddControl.open}>
                <IconPlus />
              </ActionIcon>
            </Group>
            {competitionsLoading ? (
              Array(skeletonLength)
                .fill(0)
                .map((_, index) => (
                  <LinkCardSkeleton key={index} isTagged isExport isDelete>
                    <Skeleton height={rem(theme.fontSizes.md)} width={400} />
                  </LinkCardSkeleton>
                ))
            ) : competitions.length > 0 ? (
              competitions.map((competition, index) => (
                <LinkCard
                  key={index}
                  title={stageToTitle(competition.stage)}
                  to={"/competitions/" + competition.id}
                  tag={
                    competition.is_ended ? (
                      <Badge leftSection={<IconCheck />}>Завершено</Badge>
                    ) : null
                  }
                  onExport={competition.is_ended ? handleExport : null}
                  onDelete={() => {
                    confirmCompetitionDeletion(competition.id);
                  }}
                >
                  <Text>{formatCompetitionDateRange(competition)}</Text>
                </LinkCard>
              ))
            ) : (
              <EmptyCardSpace label="Соревнования не найдены" />
            )}
          </Stack>
        </Group>
      )}
    </>
  );
}
