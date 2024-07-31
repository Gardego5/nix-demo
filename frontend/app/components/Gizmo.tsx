import type { Gizmo } from "~/lib/types";

export function Gizmo({
  data,
  className,
  ...rest
}: { data: Gizmo } & React.ComponentProps<"a">) {
  return (
    <a
      className={["block rounded border border-white m-2 p-2", className].join(
        " ",
      )}
      href={data.id.toString()}
      {...rest}
    >
      <p className="text-xl">{data.name}</p>
      <p>{data.description}</p>
    </a>
  );
}
