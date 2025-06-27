import { Skeleton } from "@/components/ui/skeleton";

const MainSkeleton = () => {
  return (
    <div className="flex h-screen w-full flex-col items-center mt-10">
      <header className="flex h-16 w-full flex-shrink-0 items-center justify-between  px-4 md:px-6">
        <Skeleton className="h-8 w-32" />
        <Skeleton className="h-10 w-10 rounded-full" />
      </header>
      <main className="flex-1 w-full  mx-auto py-8 px-4 md:px-6">
        <div className="space-y-4">
          <Skeleton className="h-10 w-full" />
          <Skeleton className="h-20 w-full" />
          <Skeleton className="h-20 w-full" />
          <Skeleton className="h-20 w-full" />
          <Skeleton className="h-20 w-full" />
          <Skeleton className="h-20 w-full" />
        </div>
      </main>
    </div>
  );
};

export default MainSkeleton;
