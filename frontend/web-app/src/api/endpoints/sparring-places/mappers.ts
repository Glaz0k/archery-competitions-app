import type { ShootOutAPIEdit, ShootOutEdit } from "./types";

export const mapToShootOutAPIEdit = (request: ShootOutEdit): ShootOutAPIEdit => {
  return {
    score: request.score,
    priority: request.priority,
  };
};
