"use client";

import Image from "next/image";
import React from "react";

const UserHeader: React.FC = () => {
  return (
    <>
      <div className="bg-green-600 text-white h-16 flex items-center shadow-lg fixed inset-0 z-40">
        <nav className="w-full ml-4 mr-8 flex justify-between">
          <a href="/" className="flex items-end">
            <Image
              src="/monta-no-mori-logo.svg"
              alt="もんたの森のロゴ"
              height={110} // 必須項目なのでとりあえず設定してるだけ
              width={110}
              style={{ height: "auto", objectFit: "contain" }}
            />
          </a>
        </nav>
      </div>
    </>
  );
};

export default UserHeader;
