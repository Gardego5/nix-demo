import { unstable_defineLoader as defineLoader } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { Gizmo } from "~/components/Gizmo";
import { getGizmos } from "~/lib/api";

export const loader = defineLoader(getGizmos);

export default function Index() {
  const data = useLoaderData<typeof loader>();

  return (
    <section>
      {data.map((gizmo) => (
        <Gizmo key={gizmo.id} data={gizmo} />
      ))}
    </section>
  );
}
