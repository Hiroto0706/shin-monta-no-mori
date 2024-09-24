import AbsolutePromotion from "@/components/user/common/absolutePromotion";
import Header from "@/components/user/common/header/handleHeader";
import UserSidebar from "@/components/user/common/sidebar/sidebar";

export default function UserLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <Header />
      <UserSidebar />
      <div className="pt-16">
        <div className="p-4 md:p-12">{children}</div>
      </div>

      <AbsolutePromotion />
    </>
  );
}
