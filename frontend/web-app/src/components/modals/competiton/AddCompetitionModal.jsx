import { NativeSelect } from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import CompetitionStage from "../../../enums/CompetitionStage";
import useCompetitionForm from "../../../hooks/useCompetitionForm";
import AddFormModal from "../AddFormModal";

export default function AddCompetitionModal({ isOpened, onClose, onSubmit, isLoading }) {
  const competitionForm = useCompetitionForm();

  return (
    <AddFormModal
      addForm={competitionForm}
      title="Новое соревнование"
      isOpened={isOpened}
      onClose={onClose}
      onSubmit={onSubmit}
      isLoading={isLoading}
    >
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
    </AddFormModal>
  );
}
