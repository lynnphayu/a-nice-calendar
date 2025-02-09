"use client";
import React, { createContext, useContext, useEffect, useState } from "react";
import { useAuth as useSWRAuth } from "../hooks/use-auth";

interface User {
  id: string;
  name: string;
  email: string;
}

interface AuthContextType {
  token?: string;
  setToken: (token: string | undefined) => void;
  isAuthenticated: boolean;
  user?: User;
  error?: Error;
}

const AuthContext = createContext<AuthContextType>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const auth = useSWRAuth();
  const [mounted, setMounted] = useState(false);
  useEffect(() => {
    setMounted(true);
  }, []);

  return mounted ? (
    <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>
  ) : (
    false
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
