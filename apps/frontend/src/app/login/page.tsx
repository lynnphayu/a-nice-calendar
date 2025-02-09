"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { z } from "zod";

const loginFormSchema = z.object({
  email: z.string().email("Please enter a valid email address"),
});

type LoginFormValues = z.infer<typeof loginFormSchema>;

export default function Login() {
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState("");

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginFormSchema),
    defaultValues: {
      email: "",
    },
  });

  async function onSubmit(data: LoginFormValues) {
    setIsLoading(true);
    setMessage("");

    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_AUTH_URL}/passwordless/initiate`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ email: data.email }),
        }
      );

      if (!response.ok) {
        throw new Error("Login request failed");
      }

      setMessage("Check your email for the login link!");
    } catch (error) {
      setMessage("Failed to send login link. Please try again.");
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <Card className="w-full max-w-[400px] mx-auto mt-8 border rounded-xl bg-background backdrop-blur supports-[backdrop-filter]:bg-background/50">
      <CardHeader>
        <CardTitle className="text-foreground">Login</CardTitle>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input
                      placeholder="Enter your email"
                      type="email"
                      disabled={isLoading}
                      className="w-full border rounded-md focus:ring-2 focus:ring-offset-2"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Sending..." : "Login with Email"}
            </Button>
            {message && (
              <p
                className={`text-sm ${
                  message.includes("Failed")
                    ? "text-destructive"
                    : "text-primary"
                }`}
              >
                {message}
              </p>
            )}
            <div className="relative my-4">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-background px-2 text-muted-foreground">
                  Or continue with
                </span>
              </div>
            </div>
            <Button
              type="button"
              variant="outline"
              className="w-full bg-[#1DB954] hover:bg-[#1ed760] text-white hover:text-white border-0"
              onClick={() =>
                (window.location.href = `${process.env.NEXT_PUBLIC_AUTH_URL}/auth/spotify`)
              }
            >
              <svg
                className="mr-2 h-4 w-4"
                fill="currentColor"
                viewBox="0 0 24 24"
              >
                <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.371-.721.49-1.101.24-3.021-1.85-6.82-2.27-11.3-1.24-.418.12-.839-.15-.959-.57-.121-.421.149-.84.57-.961 4.91-1.121 9.021-.721 12.44 1.32.39.24.48.721.24 1.101zm1.47-3.27c-.301.42-.841.6-1.262.3-3.45-2.12-8.7-2.73-12.81-1.49-.481.15-1.02-.12-1.171-.6-.15-.48.12-1.021.6-1.171 4.671-1.41 10.47-.72 14.43 1.77.42.3.571.84.27 1.261zm.129-3.42c-4.14-2.46-11.07-2.69-15.06-1.49-.601.181-1.23-.181-1.411-.781-.18-.601.18-1.23.78-1.411 4.59-1.389 12.21-1.121 17.01 1.73.54.3.719 1.02.42 1.561-.301.519-1.02.699-1.561.42z" />
              </svg>
              Spotify
            </Button>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}
