import useSWR from "swr";
import { useRouter } from "next/navigation";
import { HTTP_METHOD } from "next/dist/server/web/http";
import { useEffect, useState } from "react";

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
    `${process.env.NEXT_PUBLIC_AUTH_URL}/${path
      .split("/")
      .filter((x) => !!x)
      .join("/")}`,
    opts
  );
  return res.json() as T;
};

interface User {
  id: string;
  name: string;
  email: string;
}

export function useAuth() {
  const router = useRouter();
  const [token, setStateToken] = useState<string>();

  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    if (!storedToken) router.push("/login");
    setToken(storedToken);
  }, []);

  const { data: user, error } = useSWR(
    ["/me", token],
    ([url, token]) => !!token && fetcher<User>(url, { token })
  );

  const setToken = (newToken?: string) => {
    if (newToken) {
      localStorage.setItem("token", newToken);
      setStateToken(newToken);
    } else {
      localStorage.removeItem("token");
      setStateToken(undefined);
      router.push("/login");
    }
  };

  return {
    token: token,
    user,
    isAuthenticated: !!token && !error && !!user,
    error,
    setToken,
  };
}
