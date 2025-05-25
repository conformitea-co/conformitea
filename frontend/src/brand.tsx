import { toast } from "sonner";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Toaster } from "@/components/ui/sonner";
import { ThemeModeToggle } from "@/components/theme-mode-toggle";

export default function Brand() {
  const [inputValue, setInputValue] = useState("");
  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col items-center justify-center px-4 py-12">
      <div className="w-full max-w-2xl space-y-10">
        {/* Hero Section */}
        <section className="text-center space-y-2">
          <img
            src="/images/conformitea.svg"
            alt="Conformitea Logo"
            className="mx-auto h-16 w-16 mb-2"
          />
          <h1 className="text-4xl font-extrabold tracking-tight">
            Conformi<span className="text-tea ml-[1px]">Tea</span> UI Theme
          </h1>
          <p className="text-lg text-muted-foreground max-w-xl mx-auto">
            A modern, minimal, and accessible UI theme. Built to be easy on the
            eyes, with clear primary, accent, and destructive colors.
          </p>
          <ThemeModeToggle />
        </section>

        {/* Color Palette Showcase */}
        <section>
          <h2 className="text-xl font-semibold mb-4">Color Palette</h2>
          <div className="grid grid-cols-3 gap-4 mb-2">
            <div className="flex flex-col items-center">
              <span className="block w-12 h-12 rounded-full bg-primary border-2 border-border" />
              <span className="text-xs mt-2">Primary</span>
            </div>
            <div className="flex flex-col items-center">
              <span className="block w-12 h-12 rounded-full bg-accent border-2 border-border" />
              <span className="text-xs mt-2">Accent</span>
            </div>
            <div className="flex flex-col items-center">
              <span className="block w-12 h-12 rounded-full bg-destructive border-2 border-border" />
              <span className="text-xs mt-2">Destructive</span>
            </div>
          </div>
          <div className="grid grid-cols-3 gap-4">
            <div className="flex flex-col items-center">
              <span className="block w-12 h-12 rounded-full bg-background border-2 border-border" />
              <span className="text-xs mt-2">Background</span>
            </div>
            <div className="flex flex-col items-center">
              <span className="block w-12 h-12 rounded-full bg-card border-2 border-border" />
              <span className="text-xs mt-2">Card</span>
            </div>
            <div className="flex flex-col items-center">
              <span className="block w-12 h-12 rounded-full bg-muted border-2 border-border" />
              <span className="text-xs mt-2">Muted</span>
            </div>
          </div>
        </section>

        {/* Button Showcase */}
        <section>
          <h2 className="text-xl font-semibold mb-4 mt-8">Buttons</h2>
          <div className="flex flex-wrap gap-3 justify-center">
            <Button>Primary</Button>
            <Button variant="destructive">Destructive</Button>
            <Button variant="outline">Outline</Button>
            <Button variant="secondary">Secondary</Button>
            <Button variant="ghost">Ghost</Button>
            <Button variant="link">Link</Button>
          </div>
        </section>

        {/* Card Example */}
        <section>
          <h2 className="text-xl font-semibold mb-4 mt-8">Card Example</h2>
          <Card className="max-w-md mx-auto">
            <CardHeader>
              <CardTitle>Minimal Card</CardTitle>
              <CardDescription>
                Cards use the card and card-foreground colors.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p>
                This card demonstrates the theme's card styling and spacing.
              </p>
            </CardContent>
            <CardFooter>
              <Button size="sm">Action</Button>
            </CardFooter>
          </Card>
        </section>

        {/* Input Example */}
        <section>
          <h2 className="text-xl font-semibold mb-4 mt-8">Input Example</h2>
          <Input
            placeholder="Type something..."
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            className="mb-2"
          />
          <div className="text-sm text-muted-foreground mb-2">
            Value: {inputValue}
          </div>
        </section>

        {/* Toast Example */}
        <section>
          <h2 className="text-xl font-semibold mb-4 mt-8">
            Toast Notification
          </h2>
          <Button onClick={() => toast("This is a toast!")}>Show Toast</Button>
          <Toaster position="top-center" richColors />
        </section>
      </div>
    </div>
  );
}
