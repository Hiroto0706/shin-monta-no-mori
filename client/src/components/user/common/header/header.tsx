"use client";

import Image from "next/image";
import React, { useEffect, useState } from "react";
import SearchBox from "../searchBox";
import HeaderMenu from "./headerMenu";
import { usePathname } from "next/navigation";
import { Category } from "@/types/category";
import { Character } from "@/types/character";

type Props = {
  characters: Character[];
  categories: Category[];
};

const UserHeader: React.FC<Props> = ({ characters, categories }) => {
  const pathname = usePathname();
  const [query, setQuery] = useState("");
  const [handleHeader, setHandlerHeader] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const pathSegments = pathname.split("/");
    const lastSegment = pathSegments[pathSegments.length - 1];
    setQuery(lastSegment);
  }, [pathname]);

  useEffect(() => {
    const handleResize = () => {
      const pageWidth = window.innerWidth;
      const pageSizeMiddle = 768;
      if (pageWidth <= pageSizeMiddle) {
        setIsMobile(true);
      } else {
        setIsMobile(false);
        setHandlerHeader(false); // ウィンドウサイズが大きい場合にメニューを閉じる
      }
    };
    window.addEventListener("resize", handleResize);
    // 初回チェック
    handleResize();
    // クリーンアップ
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  return (
    <>
      <div className="bg-green-600 text-white h-16 flex items-center fixed inset-0 z-40">
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

          <SearchBox
            query={query}
            maxWidth={550}
            addClass="md:w-1/2 md:min-w-[400px]"
          />

          {/* smではハンバーガーメニューを表示 */}
          <div className="block md:hidden">
            <div
              onClick={() => setHandlerHeader(!handleHeader)}
              className="cursor-pointer w-12 h-12 rounded-full flex flex-col items-center justify-center hover:bg-white hover:bg-opacity-20 duration-200 ml-4"
            >
              <span className="w-7 h-0.5 bg-white block rounded-full mb-2"></span>
              <span className="w-7 h-0.5 bg-white block rounded-full mb-2"></span>
              <span className="w-7 h-0.5 bg-white block rounded-full"></span>
            </div>
          </div>

          {/* sm以外ではすべてのイラストを表示 */}
          <div className="hidden md:block">
            <a
              href="/illustrations"
              className="text-sm py-2 px-4 rounded-lg hover:bg-white hover:bg-opacity-30 duration-200 cursor-pointer"
            >
              すべてのイラスト
            </a>
          </div>
        </nav>
      </div>

      {isMobile && handleHeader && (
        <HeaderMenu characters={characters} categories={categories} />
      )}
    </>
  );
};

export default UserHeader;
