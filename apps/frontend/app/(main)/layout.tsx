"use client";

import MainSkeleton from "@/components/utils/MainSkeleton";
import { useAuth } from "@/hooks/useAuth";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { user } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (user === null) {
      router.push("/login");
    }
  }, [user, router]);

  if (user === undefined) {
    return <MainSkeleton />;
  }

  if (!user) {
    return null;
  }

  return <>{children}</>;
}
