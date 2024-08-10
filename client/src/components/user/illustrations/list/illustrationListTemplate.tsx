"use client";

import Breadcrumb from "@/components/common/breadCrumb";
import ListIllustrations from "../listIllustrations";
import Link from "next/link";
import { Illustration } from "@/types/illustration";
import React from "react";
import useSidebarStore from "@/store/sidebar";

interface Props {
  illustrations: Illustration[];
}

const IllustrationListTemplate: React.FC<Props> = ({ illustrations }) => {
  const { isShow } = useSidebarStore();

  return (
    <>
      <div
        className={`pl-0 duration-200 ${
          isShow ? "md:pl-[calc(4rem+14rem)]" : "md:pl-[calc(4rem)]"
        }`}
      >
        <Breadcrumb />
        <h1 className="text-xl font-bold mb-6">すべてのイラスト</h1>

        {illustrations.length > 0 ? (
          <ListIllustrations
            initialIllustrations={illustrations}
            fetchType={{}}
          />
        ) : (
          <div>
            イラストが見つかりませんでした
            <Link
              href="/"
              className="text-sm ml-4 underline border-blue-600 text-blue-600 cursor-pointer hover:text-blue-700 duration-200"
            >
              ホームに戻る
            </Link>
          </div>
        )}
      </div>
    </>
  );
};

export default IllustrationListTemplate;
