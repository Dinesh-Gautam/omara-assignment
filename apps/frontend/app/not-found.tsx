import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { PATHS } from "@/constants/config";
import Link from "next/link";

export default function NotFound() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-4xl font-bold">404</CardTitle>
          <CardDescription>
            Oops! The page you are looking for does not exist.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">
            It seems you have ventured into uncharted territory. Let us guide
            you back.
          </p>
        </CardContent>
        <CardFooter>
          <Button asChild>
            <Link href={PATHS.Dashboard}>Go to Dashboard</Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
