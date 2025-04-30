import { TextButton, type TextButtonProps } from "./TextButton";

export default function CancelButton(props: TextButtonProps) {
  return <TextButton {...props} variant="filled" color="red.8" />;
}
