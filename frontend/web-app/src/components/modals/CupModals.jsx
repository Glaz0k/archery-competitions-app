import { Button, Modal, Stack, Text, TextInput } from "@mantine/core";
import { useForm } from "@mantine/form";

export function CupAddModal({
  opened = false,
  onClose = () => {},
  handleSubmit = (_cupFormValues) => true,
  loading = false,
}) {
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
    onClose();
  };

  const actionsOnSubmit = async (cupFormValues) => {
    if (await handleSubmit(cupFormValues)) {
      cupForm.reset();
    }
  };

  return (
    <Modal opened={opened} onClose={actionsOnClose} title="Новый кубок" centered>
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
        <Button type="submit" loading={loading} loaderProps={{ type: "dots" }}>
          Добавить
        </Button>
      </form>
    </Modal>
  );
}

export function CupDeleteModal({
  opened = false,
  onClose = () => {},
  onConfirm = () => {},
  loading = false,
}) {
  return (
    <Modal opened={opened} onClose={onClose} title="Удаление кубка" centered>
      <Stack align="end">
        <Text w="100%">
          Вы уверены, что хотите удалить кубок. Вместе с ним удалятся также все связанные
          соревнования и группы.
        </Text>
        <Button loading={loading} onClick={onConfirm}>
          Удалить
        </Button>
      </Stack>
    </Modal>
  );
}
