import { Inter } from "next/font/google";
import { metadata } from "@/components/common/metaData";
export { metadata };

import Header from "@/components/common/header";
import BackgroundImage from "@/components/common/backgroundImage";

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
        {children}
        <BackgroundImage></BackgroundImage>
      </body>
    </html>
  );
}
