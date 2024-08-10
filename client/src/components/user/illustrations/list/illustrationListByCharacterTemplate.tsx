"use client";

import Breadcrumb from "@/components/common/breadCrumb";
import ListIllustrations from "../listIllustrations";
import Link from "next/link";
import { Illustration } from "@/types/illustration";
import React from "react";
import useSidebarStore from "@/store/sidebar";
import { Character } from "@/types/character";

interface Props {
  illustrations: Illustration[];
  characterID: number;
  character: Character | null;
}

const IllustrationListByCharacterTemplate: React.FC<Props> = ({
  illustrations,
  characterID,
  character,
}) => {
  const { isShow } = useSidebarStore();

  return (
    <>
      <div
        className={`pl-0 duration-200 ${
          isShow ? "md:pl-[calc(4rem+14rem)]" : "md:pl-[calc(4rem)]"
        }`}
      >
        <Breadcrumb customString={character?.name} />
        <h1 className="text-xl font-bold mb-6">
          {character != null ? (
            <>{`『${character.name}』でキャラクター検索`}</>
          ) : (
            <div>存在しないキャラクターを検索しています</div>
          )}
        </h1>

        {illustrations.length > 0 && character != null ? (
          <ListIllustrations
            initialIllustrations={illustrations}
            fetchType={{ characterID: characterID }}
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

export default IllustrationListByCharacterTemplate;
