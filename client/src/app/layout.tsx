import { Inter } from "next/font/google";
import { metadata } from "@/components/common/metaData";
export { metadata };

import Header from "@/components/common/header";
import BackgroundImage from "@/components/common/backgroundImage";
import Sidebar from "@/components/admin/common/sidebar";

import "@/styles/globals.css";

const inter = Inter({ subsets: ["latin"] });

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="jp">
      <body className={inter.className}>
        <Header></Header>
        <Sidebar></Sidebar>
        <div className="pl-20 pt-20">
          {children}
          <BackgroundImage></BackgroundImage>
        </div>
      </body>
    </html>
  );
}
