import { useForm } from "@mantine/form";
import BowClass from "../enums/BowClass";
import GroupGender from "../enums/GroupGender";

export default function useIndividualGroupForm() {
  return useForm({
    mode: "uncontrolled",
    initialValues: {
      bow: BowClass.CLASSIC.value,
      identity: GroupGender.MALE.value,
    },
    validate: {
      bow: (value) => (BowClass.valueOf(value) != null ? null : "Неверный класс лука"),
      identity: (value) => (GroupGender.valueOf(value) != null ? null : "Неверный пол"),
    },
  });
}
