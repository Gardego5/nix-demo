import { unstable_defineLoader as defineLoader } from "@remix-run/node";
import { useLoaderData, useParams } from "@remix-run/react";
import { Gizmo } from "~/components/Gizmo";
import { Widget } from "~/components/Widget";
import { getGizmoWidgets, getGizmos } from "~/lib/api";

export const loader = defineLoader(({ params: { gizmoId } }) =>
  Promise.all([getGizmos(), getGizmoWidgets(+(gizmoId ?? 0))])
);

export default function Index() {
  const [gizmos, widgets] = useLoaderData<typeof loader>();
  const { gizmoId } = useParams();

  return (
    <main className="grid grid-cols-3 gap-2 h-screen items-start">
      <section className="grid">
        {gizmos.map((gizmo) => (
          <Gizmo
            key={gizmo.id}
            data={gizmo}
            className="data-[highlighted='true']:bg-red-200"
            data-highlighted={gizmo.id === +(gizmoId ?? 0)}
          />
        ))}
      </section>

      <section className="grid">
        {widgets.map((widget) => (
          <Widget key={widget.id} data={widget} />
        ))}
      </section>
    </main>
  );
}
