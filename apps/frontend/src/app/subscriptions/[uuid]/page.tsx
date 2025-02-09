"use client";

import { useParams, useRouter } from "next/navigation";

import {
  BillingCycle,
  fromBillingCycle,
  SubscriptionFormValues,
  ToBillingCycle as toBillingCycle,
} from "@/lib/schema";
import { fetcher, useSWRSubscriptions } from "@/hooks/use-subscriptions";
import SubscriptionForm from "../_components/subscription-form";
import { Subscription } from "@/types/subscription";
import useSWR from "swr";
import { useAuth } from "@/hooks/use-auth";
import { startOfDay } from "date-fns";

export default function EditSubscription() {
  const params = useParams();
  const router = useRouter();
  const { token } = useAuth();
  const { update, remove } = useSWRSubscriptions();
  const subs = useSWR(
    token && [`/subscriptions/${params?.uuid}`, token],
    ([url, token]) =>
      fetcher<Subscription>(url, { token }).then((data) => ({
        ...data,
        billingCycle: toBillingCycle(data.billingCycle),
        startDate: new Date(data.startDate),
      }))
  );

  function onSubmit(data: SubscriptionFormValues) {
    update.trigger({
      ...data,
      billingCycle: fromBillingCycle(data.billingCycle),
    } as Omit<Subscription, "id" | "userId">);
    router.push("/subscriptions");
  }

  function onDelete() {
    if (!subs.data) return;
    remove.trigger({ uuid: subs.data.uuid });
    router.push("/subscriptions");
  }

  return (
    <SubscriptionForm
      defaultValues={
        subs.data || {
          uuid: "",
          name: "",
          price: 0,
          billingCycle: BillingCycle.Monthly,
          startDate: startOfDay(new Date()),
        }
      }
      onSubmit={onSubmit}
      onDelete={onDelete}
    />
  );
}
