import { Button, Modal, TextInput } from "@mantine/core";
import { useForm } from "@mantine/form";

export default function AddCupModal({ isOpened, onClose, onSubmit, isLoading }) {
  const cupForm = useForm({
    mode: "uncontrolled",
    initialValues: {
      title: "",
      address: "",
      season: "",
    },
  });

  const actionsOnClose = () => {
    cupForm.reset();
    onClose?.();
  };

  const actionsOnSubmit = async (cupFormValues) => {
    if (await onSubmit(cupFormValues)) {
      cupForm.reset();
    }
  };

  return (
    <Modal opened={isOpened} onClose={actionsOnClose} title="Новый кубок" centered>
      <form onSubmit={cupForm.onSubmit(actionsOnSubmit)}>
        <TextInput
          withAsterisk
          label="Название"
          key={cupForm.key("title")}
          {...cupForm.getInputProps("title")}
        />
        <TextInput
          label="Адрес проведения"
          key={cupForm.key("address")}
          {...cupForm.getInputProps("address")}
        />
        <TextInput label="Сезон" key={cupForm.key("season")} {...cupForm.getInputProps("season")} />
        <Button type="submit" loading={isLoading} loaderProps={{ type: "dots" }}>
          Добавить
        </Button>
      </form>
    </Modal>
  );
}
