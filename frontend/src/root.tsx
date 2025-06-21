import "@/styles/index.css";
import { Outlet, Scripts, ScrollRestoration } from "react-router";
import { SWRConfig } from "swr";

import { fetcher } from "@/lib/api";

import { ThemeModeToggle } from "@/components/theme-mode-toggle";
import { ThemeProvider } from "@/components/theme-provider";

export default function Root() {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <title>ConformiTea</title>
      </head>
      <body>
        <SWRConfig
          value={{
            fetcher,
            revalidateOnFocus: true,
            revalidateOnReconnect: true,
          }}
        >
          <ThemeProvider defaultTheme="system" storageKey="conformitea-theme">
            <ThemeModeToggle className="fixed top-4 right-4 z-50" />
            <div id="root">
              <Outlet />
            </div>
          </ThemeProvider>
        </SWRConfig>
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}
