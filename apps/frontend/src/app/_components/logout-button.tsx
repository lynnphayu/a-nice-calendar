"use client";

import { Button } from "@/components/ui/button";
import { useAuth } from "@/providers/auth-provider";
import { LogOut } from "lucide-react";
import { useRouter } from "next/navigation";

export function LogoutButton() {
  const { setToken, token } = useAuth();

  const router = useRouter();

  const handleLogout = () => {
    setToken(undefined);
    router.push("/login");
  };

  return (
    token && (
      <Button variant="outline" onClick={handleLogout} className="border">
        <LogOut />
      </Button>
    )
  );
}
