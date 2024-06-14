"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import React, { FormEvent, useState } from "react";
import SearchBox from "./searchBox";

const UserHeader: React.FC = () => {
  const router = useRouter();
  const [name, setName] = useState("");

  const searchIllustrations = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const queryParams: { [key: string]: string } = {};

    if (name) {
      queryParams.q = name;
    }

    const queryString = new URLSearchParams(queryParams).toString();
    router.push(`/illustrations?${queryString}`);
    router.refresh();
  };

  return (
    <>
      <div className="bg-green-600 text-white h-16 flex items-center shadow-lg fixed inset-0 z-40">
        <nav className="w-full h-16 flex justify-between items-center py-2 px-4">
          <a href="/" className="flex items-end mr-4">
            <Image
              src="/monta-no-mori-logo.svg"
              alt="もんたの森のロゴ"
              height={110} // 必須項目なのでとりあえず設定してるだけ
              width={110}
              style={{ height: "auto", objectFit: "contain" }}
            />
          </a>

          <SearchBox />

          <div className="cursor-pointer w-12 h-12 rounded-full flex flex-col items-center justify-center hover:bg-white hover:bg-opacity-20 duration-200 ml-4">
            <span className="w-7 h-1 bg-white block rounded-full mb-1.5"></span>
            <span className="w-7 h-1 bg-white block rounded-full mb-1.5"></span>
            <span className="w-7 h-1 bg-white block rounded-full"></span>
          </div>
        </nav>
      </div>
    </>
  );
};

export default UserHeader;
