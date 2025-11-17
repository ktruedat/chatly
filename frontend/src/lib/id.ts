import { v4 as uuidv4 } from "uuid";

export function makeId(): string {
  try {
    if (
      typeof crypto !== "undefined" &&
      typeof (crypto as any).randomUUID === "function"
    ) {
      return (crypto as any).randomUUID();
    }
  } catch (e) {}

  return uuidv4();
}

export default makeId;
