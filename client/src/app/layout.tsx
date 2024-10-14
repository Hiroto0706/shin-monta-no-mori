import { Inter } from "next/font/google";
import { SpeedInsights } from "@vercel/speed-insights/next";
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
      <SpeedInsights />
      <html lang="jp" className="text-gray-600">
        <body className={inter.className}>
          <div>{children}</div>
          {/* <BackgroundImage /> */}
        </body>
      </html>
    </>
  );
}
