"use client";

import { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuth } from "@/hooks/use-auth";

export default function LoginCallback() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setToken } = useAuth();

  useEffect(() => {
    const token = searchParams.get("token");
    if (!token) {
      router.push("/");
      return;
    }

    const verifyToken = async () => {
      try {
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_AUTH_URL}/passwordless/verify`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ token }),
          }
        );

        if (!response.ok) {
          throw new Error("Token verification failed");
        }

        const data = await response.json();
        setToken(data.token);
        router.push("/");
      } catch (error) {
        router.push("/");
      }
    };

    verifyToken();
  }, [router, searchParams]);

  return (
    <Card className="w-full max-w-[400px] mx-auto mt-8 border rounded-xl bg-background backdrop-blur supports-[backdrop-filter]:bg-background/50">
      <CardHeader>
        <CardTitle className="text-foreground">Verifying login...</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-muted-foreground">
          Please wait while we verify your login.
        </p>
      </CardContent>
    </Card>
  );
}
