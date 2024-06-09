"use client";

import { usePathname } from "next/navigation";
import AdminSidebar from "@/components/admin/common/sidebar";
import UserSidebar from "../user/common/sidebar";

function Sidebar() {
  const pathname = usePathname();
  const isAdminPage = pathname.startsWith("/admin");

  return (
    <>
      {pathname !== "/" && (
        <div className="w-20 h-full fixed inset-0 z-30 bg-gray-100">
          <div className="pt-16">
            <ul className="flex flex-col items-center mt-2">
              {isAdminPage ? <AdminSidebar /> : <UserSidebar />}
            </ul>
          </div>
        </div>
      )}
    </>
  );
}

export default Sidebar;
