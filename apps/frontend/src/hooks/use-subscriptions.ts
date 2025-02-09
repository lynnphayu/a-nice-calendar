import useSWR from "swr";
import useSWRMutation from "swr/mutation";
import { Subscription } from "../types/subscription";
import { HTTP_METHOD } from "next/dist/server/web/http";
import { useAuth } from "./use-auth";

export const fetcher = async <T = unknown>(
  path: string,
  {
    method,
    data,
    token,
  }: {
    method?: HTTP_METHOD;
    data?: unknown;
    token?: string;
  } = {}
) => {
  const opts = {};
  opts["method"] = method ? method : "GET";
  if (data) {
    opts["body"] = JSON.stringify(data);
  }
  if (token) {
    opts["headers"] = { Authorization: `Bearer ${token}` };
  }
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_API_URL}/${path
      .split("/")
      .filter((x) => !!x)
      .join("/")}`,
    opts
  );
  return res.json() as T;
};

// interface UseSubscriptionsReturn {
//   fetch: SWRResponse<Subscription[]>;
//   add: SWRMutationResponse<Subscription>;
//   remove: SWRMutationResponse<Subscription>;
// }

export function useSWRSubscriptions() {
  const { token } = useAuth();
  const fetch = useSWR<Subscription[], Error, string[]>(
    token && [`/subscriptions`, token],
    ([url, token]) => fetcher(url, { token })
  );

  const add = useSWRMutation(
    `/subscriptions`,
    async (url, { arg }: { arg: Omit<Subscription, "userId"> }) =>
      fetcher<Subscription>(url, {
        method: "POST",
        data: arg,
        token,
      })
  );

  const remove = useSWRMutation(
    `/subscriptions`,
    async (url, { arg }: { arg: { uuid: string } }) =>
      fetcher<Subscription>(`${url}/${arg.uuid}`, {
        method: "DELETE",
        token,
      })
  );

  const update = useSWRMutation(
    `/subscriptions`,
    async (url, { arg }: { arg: Omit<Subscription, "id" | "userId"> }) =>
      fetcher<Subscription>(`${url}/${arg.uuid}`, {
        method: "PUT",
        token,
        data: arg,
      })
  );

  return {
    fetch,
    add,
    remove,
    update,
  };
}
