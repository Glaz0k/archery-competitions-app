import { Button, Modal, NativeSelect, Stack, Text } from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import { useForm } from "@mantine/form";
import Stage from "../../enums/stage";
import { stageToTitle } from "../../helper/competitons";

export function CompetitionAddModal({
  opened = false,
  onClose = () => {},
  handleSubmit = (_competitionFormValues) => true,
  loading = false,
}) {
  const competitionForm = useForm({
    mode: "uncontrolled",
    initialValues: {
      stage: Stage.STAGE_1,
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
          data={[
            { label: stageToTitle(Stage.STAGE_1), value: Stage.STAGE_1 },
            { label: stageToTitle(Stage.STAGE_2), value: Stage.STAGE_2 },
            { label: stageToTitle(Stage.STAGE_3), value: Stage.STAGE_3 },
            { label: stageToTitle(Stage.FINAL), value: Stage.FINAL },
          ]}
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
