"use client";

import { usePathname } from "next/navigation";
import AdminSidebar from "@/components/admin/common/sidebar";
import UserSidebar from "@/components/user/common/sidebar";

function Sidebar() {
  const pathname = usePathname();
  const isAdminPage = pathname.startsWith("/admin");

  return (
    <div className="w-20 h-full fixed inset-0 z-30 border-r-2 border-gray-200 bg-gray-50">
      <div className="pt-16">
        <ul className="flex flex-col items-center">
          {isAdminPage ? <AdminSidebar /> : <UserSidebar />}
        </ul>
      </div>
    </div>
  );
}

export default Sidebar;
