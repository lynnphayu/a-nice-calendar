"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { SubscriptionTable } from "@/app/subscriptions/_components/subscription-table";
import { Button } from "@/components/ui/button";
import { ArrowLeftToLine } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { fetcher, useSWRSubscriptions } from "@/hooks/use-subscriptions";
import useSWR from "swr";
import { useAuth } from "@/hooks/use-auth";
import { Subscription } from "@/types/subscription";
import { format } from "date-fns";

export default function AllSubscriptions() {
  const { token } = useAuth();
  const params = useSearchParams();
  const { fetch } = useSWRSubscriptions();
  const date = params.get("date");

  const subs = useSWR<Subscription[], Error, string[]>(
    token &&
      date && [
        `/subscriptions?start_date_from=${format(
          new Date(date),
          "yyyy-MM-dd"
        )}&start_date_to=${format(new Date(date), "yyyy-MM-dd")}`,
        token,
      ],
    ([url, token]) => fetcher(url, { token })
  );

  const router = useRouter();

  return (
    <Card className="border rounded-xl bg-background backdrop-blur supports-[backdrop-filter]:bg-background/50">
      <CardHeader>
        <CardTitle className="text-foreground">All Subscriptions</CardTitle>
      </CardHeader>
      <CardContent>
        <SubscriptionTable
          subscriptions={(date ? subs.data : fetch.data) || []}
        />
        <div className="flex justify-end mt-4">
          <Button
            type="button"
            variant="outline"
            onClick={() => router.push("/")}
          >
            <ArrowLeftToLine />
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
