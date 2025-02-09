import { BillingCycle } from "@/lib/schema";

export type Subscription = {
  id: number;
  uuid: string;
  name: string;
  price: number;
  billingCycle: number;
  startDate: Date;
  logo: string;
  userId: string;
};
