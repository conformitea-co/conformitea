import { Outlet, Scripts, ScrollRestoration } from "react-router";

import { ThemeProvider } from "@/components/theme-provider";
import "@/styles/index.css";
import { ThemeModeToggle } from "./components/theme-mode-toggle";

export default function Root() {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <title>ConformiTea</title>
      </head>
      <body>
        <ThemeProvider defaultTheme="system" storageKey="conformitea-theme">
          <ThemeModeToggle className="fixed top-4 right-4 z-50" />
          <div id="root">
            <Outlet />
          </div>
        </ThemeProvider>
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}
