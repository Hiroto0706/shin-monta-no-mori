"use client";

import { usePathname } from "next/navigation";
import AdminSidebar from "@/components/admin/common/sidebar";

function Sidebar() {
  const pathname = usePathname();

  return (
    <>
      {pathname == "/admin" && (
        <div className="w-16 h-full fixed inset-0 z-30 bg-gray-100">
          <div className="pt-16">
            <ul className="flex flex-col items-center mt-2">
              <AdminSidebar />
              {/* {isAdminPage ? <AdminSidebar /> : <UserSidebar />} */}
            </ul>
          </div>
        </div>
      )}
    </>
  );
}

export default Sidebar;
