import { Button, Modal, NativeSelect } from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import { useForm } from "@mantine/form";
import CompetitionStage from "../../enums/CompetitionStage";

export default function AddCompetitionModal({ isOpened, onClose, onSubmit, isLoading }) {
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
    onClose?.();
  };

  const actionsOnSubmit = async (competitionFormValues) => {
    if (await onSubmit(competitionFormValues)) {
      competitionForm.reset();
    }
  };

  return (
    <Modal opened={isOpened} onClose={actionsOnClose} title="Новое соревнование" centered>
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
        <Button type="submit" loading={isLoading} loaderProps={{ type: "dots" }}>
          Добавить
        </Button>
      </form>
    </Modal>
  );
}
