import React from "react";
import { TextButton, TextButtonProps } from "./TextButton";

export default function ConfirmButton({ ...props }: TextButtonProps) {
  return <TextButton {...props} variant="filled" color="green.8" />;
}
