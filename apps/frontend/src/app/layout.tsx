import { ThemeProvider } from "../providers/theme-provider";
import { AuthProvider } from "../providers/auth-provider";
import "./globals.css";
import { ThemeToggle } from "./_components/theme-toggle";
import { LogoutButton } from "./_components/logout-button";
export const metadata = {
  title: "Subscription Tracker",
  description: "Track your subscriptions and manage recurring payments",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <AuthProvider>
          <ThemeProvider>
            <div className="gap-2 p-2">
              <div className="flex justify-end gap-2 p-4">
                <LogoutButton />
                <ThemeToggle />
              </div>
              <main>
                <div className="max-w-[720px] w-full mx-auto">{children}</div>
              </main>
            </div>
          </ThemeProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
