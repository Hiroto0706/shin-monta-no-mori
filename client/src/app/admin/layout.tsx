import AdminHeader from "@/components/admin/common/header";
import AdminSidebar from "@/components/admin/common/sidebar";

export default function AdminLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <AdminHeader />
      <AdminSidebar />
      <div className="pl-16 pt-16">
        <div className="p-2 md:p-12 text-sm md:text-md">{children}</div>
      </div>
    </>
  );
}
