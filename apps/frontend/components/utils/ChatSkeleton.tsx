import { Skeleton } from "@/components/ui/skeleton";

const ChatSkeleton = () => {
  return (
    <div className="flex h-screen flex-col p-4">
      <div className="flex-grow rounded-lg border border-gray-200 bg-white p-4">
        <div className="space-y-6">
          <div className="flex items-start space-x-4">
            <Skeleton className="h-10 w-10 rounded-full" />
            <div className="space-y-2">
              <Skeleton className="h-4 w-48" />
              <Skeleton className="h-4 w-32" />
            </div>
          </div>
          <div className="flex items-start justify-end space-x-4">
            <div className="space-y-2 text-right">
              <Skeleton className="h-4 w-48" />
              <Skeleton className="h-4 w-32" />
            </div>
            <Skeleton className="h-10 w-10 rounded-full" />
          </div>
          <div className="flex items-start space-x-4">
            <Skeleton className="h-10 w-10 rounded-full" />
            <div className="space-y-2">
              <Skeleton className="h-4 w-64" />
              <Skeleton className="h-4 w-48" />
            </div>
          </div>
        </div>
      </div>
      <div className="mt-4 flex gap-2">
        <Skeleton className="h-12 flex-grow" />
        <Skeleton className="h-12 w-24" />
      </div>
    </div>
  );
};

export default ChatSkeleton;
