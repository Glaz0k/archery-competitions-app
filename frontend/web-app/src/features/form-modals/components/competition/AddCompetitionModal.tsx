import { NativeSelect } from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import type { TransformedValues } from "@mantine/form";
import { useCreateCompetitionForm } from "../..";
import { CompetitionStage } from "../../../../entities";
import { getCompetitionStageDescription } from "../../../../utils";
import { AddModal, type AddModalProps } from "../AddModal";

type AddCompetitionForm = ReturnType<typeof useCreateCompetitionForm>;

export type AddCompetitionModalProps = Omit<
  AddModalProps<AddCompetitionForm["values"], TransformedValues<AddCompetitionForm>>,
  "form" | "title" | "children"
>;

export function AddCompetitionModal({ loading, ...others }: AddCompetitionModalProps) {
  const competitionForm = useCreateCompetitionForm();

  return (
    <AddModal form={competitionForm} title="Новое соревнование" loading={loading} {...others}>
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
    </AddModal>
  );
}
