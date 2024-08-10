"use client";

import Breadcrumb from "@/components/common/breadCrumb";
import ListIllustrations from "../listIllustrations";
import Link from "next/link";
import { Illustration } from "@/types/illustration";
import React from "react";
import useSidebarStore from "@/store/sidebar";
import { ChildCategory } from "@/types/category";

interface Props {
  illustrations: Illustration[];
  categoryID: number;
  childCategory: ChildCategory | null;
}

const IllustrationListByCategoryTemplate: React.FC<Props> = ({
  illustrations,
  categoryID,
  childCategory,
}) => {
  const { isShow } = useSidebarStore();

  return (
    <>
      <div
        className={`pl-0 duration-200 ${
          isShow ? "md:pl-[calc(4rem+14rem)]" : "md:pl-[calc(4rem)]"
        }`}
      >
        <Breadcrumb customString={childCategory?.name} />
        <h1 className="text-xl font-bold mb-6">
          {childCategory != null ? (
            <>{`『${childCategory?.name}』でカテゴリ検索`}</>
          ) : (
            <div>存在しないカテゴリを検索しています</div>
          )}
        </h1>

        {illustrations.length > 0 && childCategory != null ? (
          <ListIllustrations
            initialIllustrations={illustrations}
            fetchType={{ categoryID: categoryID }}
          />
        ) : (
          <div>
            イラストが見つかりませんでした{" "}
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

export default IllustrationListByCategoryTemplate;
