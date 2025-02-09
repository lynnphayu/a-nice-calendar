"use client";

import { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuth } from "@/providers/auth-provider";

export default function LoginCallback() {
  const router = useRouter();
  const { setToken } = useAuth();
  const searchParams = useSearchParams();

  useEffect(() => {
    const token = searchParams.get("token");
    if (!token) {
      router.push("/login");
    } else {
      setToken(token);
      router.push("/");
    }
  }, [searchParams]);

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
