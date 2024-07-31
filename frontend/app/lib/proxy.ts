import { z } from "zod";

export class ProxyFnError extends Error {
  constructor(public res: Response) {
    super(`failed to fetch ${res.url}: ${res.status} ${res.statusText}`);
  }
}

type FFRequestInit = RequestInit & { path: string | URL };

const extractValue = <T, U extends unknown[]>(
  fnOrValue: T | ((...args: U) => T),
  args: U
): T => (fnOrValue instanceof Function ? fnOrValue(...args) : fnOrValue);

type FFConfig<U extends unknown[] = [], Z extends z.ZodSchema = z.ZodSchema> = {
  [K in keyof FFRequestInit]:
    | FFRequestInit[K]
    | ((...args: U) => FFRequestInit[K]);
} & { schema: Z };

export default function ff<
  U extends unknown[] = [],
  Z extends z.ZodSchema = z.ZodSchema
>(cfg: FFConfig<U, Z>) {
  const { path, schema, ...init } = cfg;
  async function ffetch(...args: U): Promise<z.output<Z>> {
    const input = extractValue(path, args);
    const requestInit = Object.entries(init).reduce((acc, [key, value]) => {
      // @ts-expect-error we know this
      acc[key as keyof RequestInit] = extractValue(value, args);
      return acc;
    }, {} as RequestInit);
    const res = await fetch(input, requestInit);

    if (!res.ok) throw new ProxyFnError(res);
    const raw = await res.json();
    return await schema.parseAsync(raw);
  }

  return ffetch;
}
