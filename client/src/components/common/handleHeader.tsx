"use client";

import { usePathname, useSearchParams } from "next/navigation";
import AdminHeader from "../admin/common/header";
import UserHeader from "../user/common/header";

export default function Header() {
  const pathname = usePathname();
  const isAdminPage = pathname.startsWith("/admin");
  const searchParams = useSearchParams();
  const query = searchParams.get("q") || "";
  return (
    <header>
      {isAdminPage ? (
        <AdminHeader />
      ) : (
        <>{pathname != "/" ? <UserHeader query={query} /> : <></>}</>
      )}
    </header>
  );
}
