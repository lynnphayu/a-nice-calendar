import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import {
  FormField,
  FormItem,
  FormControl,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  BillingCycle,
  subscriptionFormSchema,
  SubscriptionFormValues,
} from "@/lib/schema";
import { cn } from "@/lib/utils";
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from "@radix-ui/react-popover";

import { useRouter } from "next/navigation";
import { format } from "date-fns";
import { Button } from "@/components/ui/button";
import { Form, FormProvider, useForm } from "react-hook-form";
import { Calendar } from "@/components/ui/calendar";
import { CalendarIcon } from "lucide-react";
import { useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";

export default function SubscriptionForm({
  defaultValues,
  onSubmit,
  onDelete,
}: {
  defaultValues: SubscriptionFormValues;
  onSubmit: (data: SubscriptionFormValues) => void;
  onDelete: () => void;
}) {
  const [open, setOpen] = useState(false);
  const router = useRouter();

  const logos = [
    {
      name: "Netflix",
      value: "https://cdn.simpleicons.org/netflix/E50914",
    },
    {
      name: "Spotify",
      value: "https://cdn.simpleicons.org/spotify/1DB954",
    },
    {
      name: "Discord",
      value: "https://cdn.simpleicons.org/discord/5865F2",
    },
    {
      name: "Amazon",
      value: "https://cdn.simpleicons.org/amazon/FF9900",
    },
    {
      name: "YouTube",
      value: "https://cdn.simpleicons.org/youtube/FF0000",
    },
    {
      name: "Slack",
      value: "https://cdn.simpleicons.org/slack/4A154B",
    },
    {
      name: "Apple",
      value: "https://cdn.simpleicons.org/apple/999999",
    },
  ];

  const form = useForm<SubscriptionFormValues>({
    resolver: zodResolver(subscriptionFormSchema),
    values: defaultValues,
  });

  return (
    <Card className="border rounded-xl bg-background backdrop-blur supports-[backdrop-filter]:bg-background/50">
      <CardHeader className="pb-2 px-6 pt-6">
        <CardTitle>Edit Subscription</CardTitle>
      </CardHeader>
      <CardContent className="p-6">
        <FormProvider {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit, (data) => console.log(data))}
            className="space-y-6"
          >
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="uuid"
                render={({ field }) => (
                  <FormItem hidden>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormControl>
                      <Input placeholder="Name" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="price"
                render={({ field }) => (
                  <FormItem>
                    <FormControl>
                      <Input type="number" placeholder="Price" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />
            </div>
            <FormField
              control={form.control}
              name="billingCycle"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <div className="flex flex-wrap gap-2">
                      {[
                        { label: "Monthly", value: "monthly" },
                        { label: "Yearly", value: "yearly" },
                        { label: "Quarterly", value: "quarterly" },
                      ].map((cycle) => (
                        <button
                          key={cycle.value}
                          type="button"
                          onClick={() =>
                            form.setValue(
                              "billingCycle",
                              cycle.value as BillingCycle
                            )
                          }
                          className={cn(
                            "px-3 py-2.5 border rounded-lg text-sm transition-all inline-flex items-center justify-center",
                            field.value === cycle.value
                              ? "bg-accent text-primary focus:ring-2 focus:ring-offset-2  focus:ring-primary"
                              : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                          )}
                        >
                          {cycle.label}
                        </button>
                      ))}
                    </div>
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="startDate"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <div className="relative w-[280px]">
                      <Popover open={open} onOpenChange={setOpen}>
                        <PopoverTrigger asChild>
                          <Button
                            variant={"outline"}
                            className={cn(
                              `focus:ring-2 focus:ring-offset-2  focus:ring-primary`,
                              !field.value && "text-muted-foreground"
                            )}
                          >
                            {field.value
                              ? format(field.value, "MMM d, yyyy")
                              : "Start date"}
                            <CalendarIcon className="ml-auto h-4 w-4 text-muted-foreground" />
                          </Button>
                        </PopoverTrigger>
                        <PopoverContent
                          className="w-auto p-0 bg-popover/80 backdrop-blur-sm border rounded-lg shadow-xl"
                          align="start"
                        >
                          <Calendar
                            mode="single"
                            selected={form.getValues("startDate")}
                            onSelect={(date) => {
                              form.setValue("startDate", date as Date);
                              setOpen(false);
                            }}
                            disabled={(date) =>
                              date > new Date() || date < new Date("1900-01-01")
                            }
                            className="p-3"
                          />
                        </PopoverContent>
                      </Popover>
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="logo"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <div className="flex flex-col gap-4">
                      <div className="flex flex-wrap gap-4 justify-start">
                        {logos.map((logo) => (
                          <button
                            key={logo.value}
                            type="button"
                            className={`w-10 h-10 rounded-md transition-all background border p-1 ${
                              field.value === logo.value
                                ? "focus:ring-2 focus:ring-offset-2  focus:ring-primary scale-110 bg-accent"
                                : "hover:scale-110"
                            }`}
                            onClick={() => form.setValue("logo", logo.value)}
                            title={logo.name}
                          >
                            <img
                              src={logo.value}
                              alt={logo.name}
                              className="w-full h-full object-contain"
                            />
                          </button>
                        ))}
                      </div>
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="flex justify-end gap-3 pt-2">
              {defaultValues.uuid && (
                <Button type="button" variant="destructive" onClick={onDelete}>
                  Delete
                </Button>
              )}
              <Button
                type="button"
                variant="outline"
                onClick={() => router.push("/subscriptions")}
              >
                Cancel
              </Button>
              <Button type="submit" variant="outline">
                Save
              </Button>
            </div>
          </form>
        </FormProvider>
      </CardContent>
    </Card>
  );
}
