import { Button, Modal, NativeSelect, Stack, Text } from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import { useForm } from "@mantine/form";
import CompetitionStage from "../../enums/CompetitionStage";

export function CompetitionAddModal({
  opened = false,
  onClose = () => {},
  handleSubmit = (_competitionFormValues) => true,
  loading = false,
}) {
  console.log(Object.values(CompetitionStage).length);

  const competitionForm = useForm({
    mode: "uncontrolled",
    initialValues: {
      stage: CompetitionStage.STAGE_1,
      startDate: null,
      endDate: null,
    },
  });

  const actionsOnClose = () => {
    competitionForm.reset();
    onClose();
  };

  const actionsOnSubmit = async (competitionFormValues) => {
    if (await handleSubmit(competitionFormValues)) {
      competitionForm.reset();
    }
  };

  return (
    <Modal opened={opened} onClose={actionsOnClose} title="Новое соревнование" centered>
      <form onSubmit={competitionForm.onSubmit(actionsOnSubmit)}>
        <NativeSelect
          label="Этап"
          data={Object.values(CompetitionStage).map((stage) => {
            return {
              label: stage.textValue,
              value: stage.value,
            };
          })}
          key={competitionForm.key("stage")}
          {...competitionForm.getInputProps("stage")}
        />
        <DatePickerInput
          label="Дата начала"
          clearable
          key={competitionForm.key("startDate")}
          {...competitionForm.getInputProps("startDate")}
        />
        <DatePickerInput
          label="Дата окончания"
          clearable
          key={competitionForm.key("endDate")}
          {...competitionForm.getInputProps("endDate")}
        />
        <Button type="submit" loading={loading} loaderProps={{ type: "dots" }}>
          Добавить
        </Button>
      </form>
    </Modal>
  );
}

export function CompetitionDeleteModal({
  opened = false,
  onClose = () => {},
  onConfirm = () => {},
  loading = false,
}) {
  return (
    <Modal opened={opened} onClose={onClose} title="Удаление соревнования" centered>
      <Stack align="end">
        <Text w="100%">
          Вы уверены, что хотите удалить соревнование. Вместе с ним удалятся также все связанные
          группы.
        </Text>
        <Button loading={loading} onClick={onConfirm}>
          Удалить
        </Button>
      </Stack>
    </Modal>
  );
}
