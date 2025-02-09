"use client";

import React from "react";
import { Card, CardContent, CardTitle } from "../../components/ui/card";
import { Button } from "../../components/ui/button";
import { useRouter } from "next/navigation";
import { format, getDaysInMonth, isEqual } from "date-fns";
import { ArrowLeft, ArrowRight, List, Plus } from "lucide-react";
import { fetcher } from "@/hooks/use-subscriptions";
import useSWR from "swr";
import { Subscription } from "@/types/subscription";
import { useAuth } from "@/hooks/use-auth";

const Calendar: React.FC = () => {
  const [currentDate, setCurrentDate] = React.useState(new Date());
  const router = useRouter();
  const { token } = useAuth();
  const fetch = useSWR<Subscription[], Error, string[]>(
    token && [`/subscriptions`, token],
    ([url, token]) => fetcher(url, { token })
  );

  const subs = fetch.data || [];

  const getFirstDayOfMonth = (date: Date) =>
    new Date(date.getFullYear(), date.getMonth(), 1).getDay();

  const getMonthlyTotal = () =>
    subs.reduce((total, sub) => total + (sub.price || 0) / sub.billingCycle, 0);

  const getSubscriptionsForDay = (day: number) => {
    const targetDate = new Date(
      currentDate.getFullYear(),
      currentDate.getMonth(),
      day
    );
    return subs.filter((sub) => isEqual(new Date(sub.startDate), targetDate));
  };

  // Update renderCalendarDays function
  const renderCalendarDays = () => {
    const days = [];
    const daysInMonth = getDaysInMonth(currentDate);
    const firstDay = getFirstDayOfMonth(currentDate);

    // Add empty cells for days before the first day of the month
    for (let i = 0; i < firstDay; i++) {
      days.push(
        <div
          key={`empty-${i}`}
          className="h-16 sm:h-20 rounded-lg bg-muted/50"
        ></div>
      );
    }

    // Add cells for each day of the month
    for (let day = 1; day <= daysInMonth; day++) {
      const daySubscriptions = getSubscriptionsForDay(day);

      days.push(
        <Card
          key={day}
          className="h-16 sm:h-20 p-1 sm:p-2 relative border bg-card hover:bg-accent transition-colors"
        >
          <CardTitle className="text-card-foreground">{day}</CardTitle>
          <div
            className="absolute inset-0 cursor-pointer rounded-lg"
            onClick={() => {
              const date = new Date();
              date.setUTCFullYear(
                currentDate.getFullYear(),
                currentDate.getMonth(),
                day
              );
              return (
                daySubscriptions.length > 0 &&
                router.push(
                  `/subscriptions/?date=${format(date, "yyyy MM dd")}`
                )
              );
            }}
          >
            <div className="absolute bottom-2 left-2 flex flex-wrap gap-1 max-w-[calc(100%-8px)]">
              {daySubscriptions.map((sub, index) => (
                <div
                  key={index}
                  className="w-4 h-4 rounded-lg p-0.5 transition-colors"
                  title={`${sub.name} ($${sub.price})`}
                >
                  <img
                    src={sub.logo}
                    alt={sub.name}
                    className="w-full h-full object-contain"
                  />
                </div>
              ))}
            </div>
          </div>
        </Card>
      );
    }

    return days;
  };

  return (
    <Card className="mx-auto border rounded-xl bg-background backdrop-blur supports-[backdrop-filter]:bg-background/50">
      <CardContent>
        <div className="text-right my-6 bg-accent p-6 rounded-lg border border-border">
          <p className="text-sm text-muted-foreground">Monthly total</p>
          <p className="text-2xl font-semibold text-primary">
            ${getMonthlyTotal().toFixed(2)}
          </p>
        </div>
        <div className="flex flex-col sm:flex-row items-center justify-between space-y-4 sm:space-y-0 mb-6">
          <CardTitle className="text-foreground">
            {format(currentDate, "LLLL yyyy")}
          </CardTitle>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              onClick={() =>
                setCurrentDate(
                  new Date(
                    currentDate.getFullYear(),
                    currentDate.getMonth() - 1
                  )
                )
              }
            >
              <ArrowLeft />
            </Button>
            <Button
              variant="outline"
              onClick={() =>
                setCurrentDate(
                  new Date(
                    currentDate.getFullYear(),
                    currentDate.getMonth() + 1
                  )
                )
              }
            >
              <ArrowRight />
            </Button>
            <Button
              variant="outline"
              onClick={() => router.push(`/subscriptions/new`)}
            >
              <Plus />
            </Button>
            <Button
              variant="outline"
              onClick={() => router.push(`/subscriptions`)}
            >
              <List />
            </Button>
          </div>
        </div>
        <div className="grid grid-cols-7 gap-3 text-sm sm:text-base">
          {["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"].map((day) => (
            <div key={day} className="text-center font-medium">
              {day}
            </div>
          ))}
          {renderCalendarDays()}
        </div>
      </CardContent>
    </Card>
  );
};

export default Calendar;
