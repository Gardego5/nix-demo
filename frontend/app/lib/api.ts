import { z } from "zod";
import ff from "./proxy";
import { Gizmo, Widget, zGizmo, zWidget } from "./types";

const { API_BASE } = process.env;

export const getHealth = ff({
  path: new URL("/health", API_BASE),
  schema: z.object({ status: z.literal("ok") }),
});

export const getGizmos = ff({
  path: new URL("/gizmos", API_BASE),
  schema: zGizmo.array(),
});

export const postGizmo = ff({
  body: (data: Pick<Gizmo, "name" | "description">) => JSON.stringify(data),
  method: "POST",
  path: new URL("/gizmos", API_BASE),
  schema: zGizmo,
});

export const deleteGizmo = ff({
  method: "DELETE",
  path: (id: number) => new URL(`/gizmos/${id}`, API_BASE),
  schema: z.void(),
});

export const getGizmo = ff({
  path: (id: number) => new URL(`/gizmos/${id}`, API_BASE),
  schema: zGizmo,
});

export const getGizmoWidgets = ff({
  path: (id: number) => new URL(`/gizmos/${id}/widgets`, API_BASE),
  schema: zWidget.array(),
});

export const getWidget = ff({
  path: (gizmoId: number, widgetId: number) =>
    new URL(`/gizmos/${gizmoId}/widgets/${widgetId}`, API_BASE),
  schema: zWidget,
});

export const postWidget = ff({
  body: (_gizmoId: number, data: Pick<Widget, "name">) => JSON.stringify(data),
  method: "POST",
  path: (gizmoId) => new URL(`/gizmos/${gizmoId}/widgets`, API_BASE),
  schema: zWidget,
});

export const deleteWidget = ff({
  method: "DELETE",
  path: (gizmoId: number, widgetId: number) =>
    new URL(`/gizmos/${gizmoId}/widgets/${widgetId}`, API_BASE),
  schema: z.void(),
});
