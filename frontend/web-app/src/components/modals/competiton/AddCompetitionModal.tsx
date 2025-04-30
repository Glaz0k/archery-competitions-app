import { NativeSelect } from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import type { TransformedValues } from "@mantine/form";
import { CompetitionStage } from "../../../entities";
import { useCreateCompetitionForm } from "../../../hooks";
import { getCompetitionStageDescription } from "../../../utils";
import { AddFormModal, type AddFormModalProps } from "../AddFormModal";

type AddCompetitionForm = ReturnType<typeof useCreateCompetitionForm>;

export type AddCompetitionModalProps = Omit<
  AddFormModalProps<AddCompetitionForm["values"], TransformedValues<AddCompetitionForm>>,
  "form" | "title" | "children"
>;

export default function AddCompetitionModal({
  opened,
  onClose,
  onSubmit,
  loading,
}: AddCompetitionModalProps) {
  const competitionForm = useCreateCompetitionForm();

  return (
    <AddFormModal
      form={competitionForm}
      title="Новое соревнование"
      opened={opened}
      onClose={onClose}
      onSubmit={onSubmit}
      loading={loading}
    >
      <NativeSelect
        label="Этап"
        data={Object.values(CompetitionStage).map((stage) => {
          return {
            label: getCompetitionStageDescription(stage),
            value: stage,
          };
        })}
        disabled={loading}
        key={competitionForm.key("stage")}
        {...competitionForm.getInputProps("stage")}
      />
      <DatePickerInput
        label="Дата начала"
        clearable
        disabled={loading}
        key={competitionForm.key("startDate")}
        {...competitionForm.getInputProps("startDate")}
      />
      <DatePickerInput
        label="Дата окончания"
        clearable
        disabled={loading}
        key={competitionForm.key("endDate")}
        {...competitionForm.getInputProps("endDate")}
      />
    </AddFormModal>
  );
}
