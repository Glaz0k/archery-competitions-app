import { Modal, TextInput, Button, Stack, Text, Group } from "@mantine/core";
import { useForm } from "@mantine/form";

export function CupAddModal({
  opened = false,
  onClose = null,
  handleSubmit = null,
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

  return (
    <Modal
      opened={opened}
      onClose={() => {
        cupForm.reset();
        onClose();
      }}
      title="Новый кубок"
      centered
    >
      <form onSubmit={cupForm.onSubmit(handleSubmit)}>
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
        <TextInput
          label="Сезон"
          key={cupForm.key("season")}
          {...cupForm.getInputProps("season")}
        />
        <Button type="submit" loading={loading} loaderProps={{ type: "dots" }}>
          Добавить
        </Button>
      </form>
    </Modal>
  );
}

export function CupDeleteModal({
  opened = false,
  onClose = null,
  onDeny = null,
  onConfirm = null,
  loading = false,
}) {
  return (
    <Modal
      opened={opened}
      onClose={() => {
        onDeny();
        onClose();
      }}
      title="Удаление кубка"
      closeOnClickOutside={!loading}
      closeOnEscape={!loading}
      closeButtonProps={{ disabled: loading }}
      centered
    >
      <Stack align="end">
        <Text w="100%">
          Вы уверены, что хотите удалить кубок. Вместе с кубком удалятся также
          все связанные соревнования и группы.
        </Text>
        <Button loading={loading} onClick={onConfirm}>
          Удалить
        </Button>
      </Stack>
    </Modal>
  );
}
