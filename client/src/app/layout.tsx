import { Inter } from "next/font/google";
import { metadata } from "@/components/common/metaData";
export { metadata };

import Header from "@/components/common/handleHeader";
import Sidebar from "@/components/common/handleSidebar";
import BackgroundImage from "@/components/common/backgroundImage";

import "@/styles/globals.css";

const inter = Inter({ subsets: ["latin"] });

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="jp" className="text-gray-600">
      <body className={inter.className}>
        <Header />
        <Sidebar />
        <div>{children}</div>
        <BackgroundImage />
      </body>
    </html>
  );
}
