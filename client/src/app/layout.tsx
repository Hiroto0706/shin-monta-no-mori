import { Inter } from "next/font/google";
import { Analytics } from "@vercel/analytics/react";
import { metadata } from "@/components/common/metaData";
export { metadata };

import "@/styles/globals.css";

const inter = Inter({ subsets: ["latin"] });

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <Analytics />
      <html lang="jp" className="text-gray-600">
        <body className={inter.className}>
          <div>{children}</div>
          {/* <BackgroundImage /> */}
        </body>
      </html>
    </>
  );
}
