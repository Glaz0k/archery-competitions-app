import { NativeSelect } from "@mantine/core";
import type { TransformedValues } from "@mantine/form";
import { BowClass, Identity } from "../../../../entities";
import { getBowClassDescription, getIdentityDescription } from "../../../../utils";
import { useCreateIndividualGroupForm } from "../../hooks/useCreateIndividualGroupForm";
import { AddModal, type AddModalProps } from "../AddModal";

type AddIndividualGroupForm = ReturnType<typeof useCreateIndividualGroupForm>;

export type AddIndividualGroupModalProps = Omit<
  AddModalProps<AddIndividualGroupForm["values"], TransformedValues<AddIndividualGroupForm>>,
  "form" | "title" | "children"
>;

export function AddIndividualGroupModal({ loading, ...others }: AddIndividualGroupModalProps) {
  const groupForm = useCreateIndividualGroupForm();

  return (
    <AddModal form={groupForm} title="Новый дивизион" loading={loading} {...others}>
      <NativeSelect
        label="Класс лука"
        data={Object.values(BowClass).map((bow) => {
          return {
            label: getBowClassDescription(bow),
            value: bow,
          };
        })}
        disabled={loading}
        key={groupForm.key("bow")}
        {...groupForm.getInputProps("bow")}
      />
      <NativeSelect
        label="Пол"
        data={Object.values(Identity).map((identity) => {
          return {
            label: getIdentityDescription(identity),
            value: identity,
          };
        })}
        disabled={loading}
        key={groupForm.key("identity")}
        {...groupForm.getInputProps("identity")}
      />
    </AddModal>
  );
}
