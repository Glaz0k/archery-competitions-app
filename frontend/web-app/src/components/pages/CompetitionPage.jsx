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
import BowClass from "../../enums/BowClass";
import GroupGender from "../../enums/GroupGender";
import GroupState from "../../enums/GroupState";

export default function CompetitionPage() {
  const { cupId, competitionId } = useParams();
  const [cup, setCup] = useState(null);
  const [competition, setCompetition] = useState(null);
  const [isLoadingMainInfo, setLoadingMainInfo] = useState(true);

  const [isCompetitionEditing, setCompetitionEditing] = useState(false);
  const [editedCompetition, setEditedCompetition] = useState({
    id: null,
    stage: null,
    startDate: null,
    endDate: null,
    isEnded: null,
  });
  const [isEditedCompetitionSubmitting, setEditedCompetitionSubmitting] = useState(false);

  const navigate = useNavigate();

  const readMainInfo = async (cupId, competitionId) => {
    try {
      setLoadingMainInfo(true);
      const [cupData, competitionData] = await Promise.all([
        getCup(cupId),
        getCompetition(competitionId),
      ]);
      setCup(cupData);
      setCompetition(competitionData);
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setLoadingMainInfo(false);
    }
    return true;
  };

  const editCompetition = async () => {
    try {
      setEditedCompetitionSubmitting(true);
      const data = await putCompetition(editedCompetition);
      setCompetition(data);
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setEditedCompetitionSubmitting(false);
    }
    return true;
  };

  const handleCompetitionEditing = () => {
    setEditedCompetition({ ...competition });
    setCompetitionEditing(true);
  };

  const handleCompetitionEditingSubmit = async () => {
    if (await editCompetition()) {
      setCompetitionEditing(false);
    }
  };

  const defaultAndEnumValues = (enumObj) => {
    return [
      {
        value: null,
        textValue: "Не указано",
      },
      ...Object.values(enumObj),
    ];
  };

  useEffect(() => {
    readMainInfo(cupId, competitionId);
  }, [cupId, competitionId]);

  return (
    <>
      <LoadingOverlay visible={isLoadingMainInfo} />
      {cup && competition && (
        <Group align="start" flex={1}>
          <Stack>
            <Stack w={300} align="start" pos="relative">
              <LoadingOverlay visible={isEditedCompetitionSubmitting} />
              <Button
                onClick={() => {
                  navigate("/cups/" + cupId);
                }}
                leftSection={<IconArrowLeft />}
              >
                Назад
              </Button>
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
              <Group w="100%">
                <Group flex={1}>
                  {!isCompetitionEditing ? (
                    <ActionIcon onClick={handleCompetitionEditing}>
                      <IconEdit />
                    </ActionIcon>
                  ) : (
                    <>
                      <ActionIcon onClick={handleCompetitionEditingSubmit}>
                        <IconCheck />
                      </ActionIcon>
                      <ActionIcon onClick={() => setCompetitionEditing(false)}>
                        <IconX />
                      </ActionIcon>
                    </>
                  )}
                </Group>
                <ActionIcon>
                  <IconTrashX />
                </ActionIcon>
              </Group>
            </Stack>
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
            <Group>
              <Title order={2} flex={1}>
                Индивидуальные группы
              </Title>
              <ActionIcon>
                <IconRefresh />
              </ActionIcon>
              <ActionIcon>
                <IconPlus />
              </ActionIcon>
            </Group>
          </Stack>
        </Group>
      )}
    </>
  );
}
