"use client";

import { useRouter } from "next/navigation";
import {
  BillingCycle,
  fromBillingCycle,
  SubscriptionFormValues,
} from "@/lib/schema";
import { useSWRSubscriptions } from "@/hooks/use-subscriptions";
import SubscriptionForm from "../_components/subscription-form";
import { Subscription } from "@/types/subscription";

export default function NewSubs() {
  const { add } = useSWRSubscriptions();
  const router = useRouter();

  const defaultValues: SubscriptionFormValues = {
    name: "",
    price: 0,
    billingCycle: BillingCycle.Monthly,
    startDate: new Date(),
    logo: "",
  };

  function onSubmit(data: SubscriptionFormValues) {
    add.trigger({
      ...data,
      billingCycle: fromBillingCycle(data.billingCycle),
    } as Subscription);
    router.push("/subscriptions");
  }

  function onDelete() {
    router.push("/subscriptions");
  }

  return (
    <SubscriptionForm
      defaultValues={defaultValues}
      onSubmit={onSubmit}
      onDelete={onDelete}
    />
  );
}
