"use client";

import { usePathname } from "next/navigation";
import AdminHeader from "../admin/common/header";
import UserHeader from "../user/common/header";

export default function Header() {
  const pathname = usePathname();
  const isAdminPage = pathname.startsWith("/admin");
  return (
    <header>
      {isAdminPage ? <AdminHeader /> : <UserHeader pathname={pathname} />}
    </header>
  );
}
