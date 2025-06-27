"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useEffect } from "react";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <Card className="w-full max-w-md border-red-400/50">
        <CardHeader>
          <CardTitle className="text-4xl font-bold">Error</CardTitle>
          <CardDescription>
            Oops! Something went wrong on our end.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">
            We are already working to fix the problem. Please try again later.
          </p>
        </CardContent>
        <CardFooter>
          <Button onClick={() => reset()}>Try Again</Button>
        </CardFooter>
      </Card>
    </div>
  );
}
