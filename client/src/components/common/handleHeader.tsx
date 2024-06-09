"use client";

import { usePathname } from "next/navigation";
import AdminHeader from "../admin/common/header";
import UserHeader from "../user/common/header";

export default function Header() {
  const pathname = usePathname();
  const isAdminPage = pathname.startsWith("/admin");
  return (
    <header className="bg-green-600 text-white h-16 flex items-center shadow-lg fixed inset-0 z-40">
      {isAdminPage ? <AdminHeader /> : <UserHeader />}
    </header>
  );
}
