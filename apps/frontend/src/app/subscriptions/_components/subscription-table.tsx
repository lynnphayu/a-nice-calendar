import { format } from "date-fns";
import { Subscription } from "@/types/subscription";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
  SortingState,
  getSortedRowModel,
} from "@tanstack/react-table";
import { useState } from "react";
import { ArrowUpDown, Edit, Edit2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";

interface SubscriptionTableProps {
  subscriptions: Subscription[];
}

export function SubscriptionTable({ subscriptions }: SubscriptionTableProps) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const router = useRouter();

  const columns: ColumnDef<Subscription>[] = [
    {
      accessorKey: "name",
      header: "Service",
      cell: ({ row }) => (
        <div className="flex items-center gap-3 max-w-48">
          <div className="w-8 h-8 rounded-lg p-1">
            <img
              src={row.original.logo}
              alt={row.original.name}
              className="w-full h-full object-contain"
            />
          </div>
          <span className="font-medium">{row.original.name}</span>
        </div>
      ),
    },
    {
      accessorKey: "price",
      header: ({ column }) => {
        return (
          <Button
            className="w-16"
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Price
            <ArrowUpDown />
          </Button>
        );
      },
      cell: ({ row }) => (
        <div className="font-medium">
          ${(row.original.price || 0).toFixed(2)}
        </div>
      ),
    },
    {
      accessorKey: "startDate",
      header: ({ column }) => {
        return (
          <Button
            className="w-24"
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Start Date
            <ArrowUpDown />
          </Button>
        );
      },
      cell: ({ row }) => (
        <div className="text-muted-foreground">
          {format(new Date(row.original.startDate), "d MMM yyyy")}
        </div>
      ),
    },
    {
      accessorKey: "billingCycle",
      header: ({ column }) => {
        return (
          <Button
            className="w-24"
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Billing Cycle
            <ArrowUpDown />
          </Button>
        );
      },
      cell: ({ row }) => (
        <span className="px-3 py-1 rounded-full text-sm bg-muted">
          {((cycle) => {
            switch (cycle) {
              case 1:
                return "Monthly";
              case 2:
                return "Bi-monthly";
              case 3:
                return "Quarterly";
              case 12:
                return "Yearly";
            }
          })(row.original.billingCycle)}
        </span>
      ),
    },
    {
      accessorKey: "actions",
      header: "",
      cell: ({ row }) => (
        <Button
          variant="ghost"
          onClick={() => router.push(`/subscriptions/${row.original.uuid}`)}
          className="p-0"
        >
          <span className="sr-only">Edit</span>
          <Edit2 />
        </Button>
      ),
    },
  ];

  const table = useReactTable({
    data: subscriptions,
    columns,
    state: {
      sorting,
    },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  return (
    <div className="relative">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
