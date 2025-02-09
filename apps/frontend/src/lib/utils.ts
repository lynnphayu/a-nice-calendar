import { clsx, type ClassValue } from "clsx";
import { HTTP_METHOD } from "next/dist/server/web/http";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const fetcher = async (
  url: string,
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
  const res = await fetch(url, opts);
  return res.json();
};
