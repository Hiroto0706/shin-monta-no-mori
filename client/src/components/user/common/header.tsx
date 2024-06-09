"use client";

import Image from "next/image";
import React from "react";

interface Props {
  pathname: string;
}

const UserHeader: React.FC<Props> = ({ pathname }) => {
  return (
    <>
      {pathname !== "/" ? (
        <header className="bg-green-600 text-white h-16 flex items-center shadow-lg fixed inset-0 z-40">
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
        </header>
      ) : (
        <header className="bg-green-600 text-white h-80 z-40">
          <nav className="w-full flex justify-between items-center py-2 px-4">
            <a href="/" className="flex items-end">
              <Image
                src="/monta-no-mori-logo.svg"
                alt="もんたの森のロゴ"
                height={110} // 必須項目なのでとりあえず設定してるだけ
                width={110}
                style={{ height: "auto", objectFit: "contain" }}
              />
            </a>

            <div className="cursor-pointer w-16 h-16 rounded-full flex flex-col items-center justify-center hover:bg-white hover:bg-opacity-20 duration-200">
              <span className="w-8 h-1 bg-white block rounded-full mb-2"></span>
              <span className="w-8 h-1 bg-white block rounded-full mb-2"></span>
              <span className="w-8 h-1 bg-white block rounded-full"></span>
            </div>
          </nav>
        </header>
      )}
    </>
  );
};

export default UserHeader;
