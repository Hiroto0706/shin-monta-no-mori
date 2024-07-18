"use client";

import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import React from "react";

interface Props {
  customString?: string;
}

const Breadcrumb: React.FC<Props> = ({ customString }) => {
  const pathname = usePathname();

  // フィルタリング用のリスト
  const excludedPaths = ["category", "character", "search"];

  // フィルタリングされたパス配列
  const pathArray = pathname
    .split("/")
    .filter((path) => path && !excludedPaths.includes(path));

  // フルパスを構築するための関数
  const generateHref = (index: number) => {
    const allPaths = pathname.split("/").filter((path) => path);
    const validPaths = allPaths.filter((path) => !excludedPaths.includes(path));
    const validPathIndex = allPaths.indexOf(validPaths[index]);
    return "/" + allPaths.slice(0, validPathIndex + 1).join("/");
  };

  return (
    <nav aria-label="breadcrumb" className="mb-4">
      <ol className="flex items-center">
        <li className="flex items-center">
          <Link
            href="/"
            className="flex items-center text-sm hover:opacity-70 duration-200"
          >
            <Image
              className="mr-1"
              src="/icon/breadCrumb/home.png"
              alt="homeアイコン"
              width={16}
              height={16}
            />
            <span className="underline font-bold">Home</span>
          </Link>
          {pathArray.length > 0 && <span className="mx-2 text-sm">＞</span>}
        </li>
        {pathArray.map((path, index) => {
          const href = generateHref(index);
          return (
            <li key={index} className="flex items-center">
              <Link
                href={href}
                className="underline font-bold text-sm hover:opacity-70 duration-200"
              >
                {index < pathArray.length - 1
                  ? decodeURIComponent(path)
                  : customString ?? decodeURIComponent(path)}
              </Link>
              {index < pathArray.length - 1 && (
                <span className="mx-2 text-sm">＞</span>
              )}
            </li>
          );
        })}
      </ol>
    </nav>
  );
};

export default Breadcrumb;
