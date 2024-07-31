import { z } from "zod";

export type Gizmo = z.infer<typeof zGizmo>;
export const zGizmo = z.object({
  id: z.number(),
  name: z.string(),
  description: z.string(),
});

export type Widget = z.infer<typeof zWidget>;
export const zWidget = z.object({
  id: z.number(),
  gizmoId: z.number(),
  name: z.string(),
});

export type ApiErrorResponse = z.infer<typeof zApiErrorResponse>;
export const zApiErrorResponse = z.object({
  error: z.string(),
  description: z.string(),
});
