export default function AdminLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="pl-20 pt-16">
      <div className="p-12">{children}</div>
    </div>
  );
}
