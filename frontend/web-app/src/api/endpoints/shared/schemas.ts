import { isValid, parseISO } from "date-fns";
import { z } from "zod";

export const DateSchema = z
  .string()
  .trim()
  .regex(/^\d{4}-\d{2}-\d{2}$/)
  .date();

const isValidISO8601 = (value: string): boolean => {
  try {
    return isValid(parseISO(value));
  } catch {
    return false;
  }
};

export const DateTZSchema = z
  .string()
  .trim()
  .regex(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}([+-]\d{2}|Z)$/)
  .refine(isValidISO8601);
