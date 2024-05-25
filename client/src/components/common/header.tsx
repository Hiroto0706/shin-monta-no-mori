"use client";

import { usePathname } from "next/navigation";
import UserHeader from "@/components/user/common/Header";
import AdminHeader from "@/components/admin/common/Header";

export default function Header() {
  const pathname = usePathname();
  const isAdminPage = pathname.startsWith("/admin");
  return <>{isAdminPage ? <AdminHeader /> : <UserHeader />}</>;
}
