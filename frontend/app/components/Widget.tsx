import type { Widget } from "~/lib/types";

export function Widget(props: { data: Widget }) {
  return (
    <div className="rounded border border-white m-2 p-2">
      <p>{props.data.name}</p>
    </div>
  );
}
