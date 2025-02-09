import * as z from "zod";
import { v4 as uuidv4 } from "uuid";

export enum BillingCycle {
  Monthly = "monthly",
  Yearly = "yearly",
  Quarterly = "quarterly",
}

export const subscriptionFormSchema = z.object({
  uuid: z.string().default(uuidv4()),
  name: z.string().min(1, { message: "Name is required" }),
  price: z.coerce
    .number({ required_error: "Price is required" })
    .min(0.01, { message: "Price must be greater than 0" }),
  billingCycle: z.nativeEnum(BillingCycle, {
    required_error: "Billing cycle is required",
  }),
  startDate: z.date({
    required_error: "Start date is required",
  }),
  logo: z.string({
    required_error: "Please select a logo",
  }),
});

export type SubscriptionFormValues = z.infer<typeof subscriptionFormSchema>;

export const ToBillingCycle = (cycle) => {
  switch (cycle) {
    case 1:
      return BillingCycle.Monthly;
    case 3:
      return BillingCycle.Quarterly;
    case 12:
      return BillingCycle.Yearly;
  }
};

export const fromBillingCycle = (cycle) => {
  switch (cycle) {
    case BillingCycle.Monthly:
      return 1;
    case BillingCycle.Quarterly:
      return 3;
    case BillingCycle.Yearly:
      return 12;
  }
};
