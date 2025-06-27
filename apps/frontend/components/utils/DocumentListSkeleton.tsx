import { Skeleton } from "@/components/ui/skeleton";

export default function DocumentListSkeleton() {
  return (
    <div className="space-y-2">
      {Array.from({ length: 5 }).map((_, index) => (
        <Skeleton className="w-full h-12" />
      ))}
    </div>
  );
}
