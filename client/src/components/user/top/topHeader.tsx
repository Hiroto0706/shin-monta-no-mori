"use client";

import Image from "next/image";
import SearchFormTop from "./searchForm";
import { Category, ChildCategory } from "@/types/category";
import HeaderMenu from "../common/header/headerMenu";
import { Character } from "@/types/character";
import { useEffect, useState } from "react";

type Props = {
  child_categories: ChildCategory[];
  characters: Character[];
  categories: Category[];
};

const TopHeader: React.FC<Props> = ({
  child_categories,
  characters,
  categories,
}) => {
  const [handleHeader, setHandlerHeader] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

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
      <div className="bg-green-600 text-white h-96">
        <nav
          className={`w-full h-16 flex justify-between items-center py-2 px-4 z-40 bg-green-600 ${
            !isMobile ? "absolute" : "fixed"
          }`}
        >
          <a href="/" className="flex items-end">
            <Image
              src="/monta-no-mori-logo.svg"
              alt="もんたの森のロゴ"
              height={110} // 必須項目なのでとりあえず設定してるだけ
              width={110}
              style={{ height: "auto", objectFit: "contain" }}
            />
          </a>

          {/* smではハンバーガーメニューを表示 */}
          <div className="block md:hidden">
            <div
              onClick={() => setHandlerHeader(!handleHeader)}
              className="cursor-pointer w-12 h-12 rounded-full flex flex-col items-center justify-center hover:bg-white hover:bg-opacity-20 duration-200"
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

        <div className="w-full h-96">
          <div className="h-full flex items-center justify-center flex-col">
            <SearchFormTop child_categories={child_categories} />
          </div>
        </div>
      </div>

      {isMobile && handleHeader && (
        <HeaderMenu characters={characters} categories={categories} />
      )}
    </>
  );
};

export default TopHeader;
