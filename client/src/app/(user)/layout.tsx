import UserSidebar from "@/components/user/common/sidebar/sidebar";

export default function UserLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <UserSidebar />
      <div className="pl-0 md:pl-[calc(4rem+14rem)] pt-16 duration-200">
        <div className="p-12">{children}</div>
      </div>
    </>
  );
}
