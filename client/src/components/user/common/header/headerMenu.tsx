"use client";

import { Category } from "@/types/category";
import { Character } from "@/types/character";
import Image from "next/image";
import { useState } from "react";

type Props = {
  characters: Character[];
  categories: Category[];
};

const HeaderMenu: React.FC<Props> = ({ characters, categories }) => {
  const [isCategoryOpen, setIsCategoryOpen] = useState(false);
  const [isCharacterOpen, setIsCharacterOpen] = useState(false);

  return (
    <>
      <div className="fixed inset-0 bg-green-600 z-30 text-white overflow-y-scroll">
        <nav className="mt-20 mb-8 mx-4">
          <div className="border-b-2 border-white mb-4">
            <a
              href="/"
              className="text-lg py-2 px-1 mb-2 block cursor-pointer hover:bg-white hover:bg-opacity-30 rounded-lg duration-200"
            >
              TOP
            </a>
            <a
              href="/illustrations"
              className="text-lg py-2 px-1 mb-2 block cursor-pointer hover:bg-white hover:bg-opacity-30 rounded-lg duration-200"
            >
              すべてのイラスト
            </a>
          </div>

          <div className="mb-2">
            <div
              className="flex justify-between mb-2 py-2 px-1 cursor-pointer hover:bg-white hover:bg-opacity-30 rounded-lg duration-200"
              onClick={() => setIsCategoryOpen(!isCategoryOpen)}
            >
              <div className="flex items-center">
                <Image
                  src="/icon/menu/category.png"
                  alt="categoryアイコン"
                  width={32}
                  height={32}
                />
                <span className="ml-2 text-lg">カテゴリ</span>
              </div>
              <Image
                className={`duration-200 ${
                  isCategoryOpen ? "-rotate-90" : "rotate-90"
                }`}
                src="/icon/menu/arrow.png"
                alt="arrowアイコン"
                width={28}
                height={28}
              />
            </div>
            {isCategoryOpen && (
              <ul className="ml-2 pl-4 border-l-2 border-white">
                {categories.map((category) => (
                  <div key={category.ParentCategory.id}>
                    <div className="text-lg font-bold flex items-center">
                      <Image
                        src={category.ParentCategory.src}
                        alt={category.ParentCategory.name}
                        width={24}
                        height={24}
                      />
                      <span className="ml-2">{category.ParentCategory.name}</span>
                    </div>
                    <div className="flex flex-wrap mb-4">
                      {category.ChildCategory.map((cc) => (
                        <a
                          key={cc.id}
                          href={`/illustrations/category/${cc.id}`}
                          className="text-sm ml-1 mr-2 py-1 px-2 cursor-pointer hover:bg-white hover:bg-opacity-30 rounded-full duration-200"
                        >
                          # {cc.name}
                        </a>
                      ))}
                    </div>
                  </div>
                ))}
              </ul>
            )}
          </div>

          <div className="mb-2">
            <div
              className="flex justify-between mb-2 py-2 px-1 cursor-pointer hover:bg-white hover:bg-opacity-30 rounded-lg duration-200"
              onClick={() => setIsCharacterOpen(!isCharacterOpen)}
            >
              <div className="flex items-center">
                <Image
                  src="/icon/menu/character.png"
                  alt="characterアイコン"
                  width={32}
                  height={32}
                />
                <span className="ml-2 text-lg">キャラクター</span>
              </div>
              <Image
                className={`duration-200 ${
                  isCharacterOpen ? "-rotate-90" : "rotate-90"
                }`}
                src="/icon/menu/arrow.png"
                alt="arrowアイコン"
                width={28}
                height={28}
              />
            </div>
            {isCharacterOpen && (
              <ul className="ml-2 pl-4 border-l-2 border-white">
                {characters.map((character) => (
                  <a
                    key={character.id}
                    href={`/illustrations/character/${character.id}`}
                    className="flex items-center py-1 pl-1 pr-2 mb-2 cursor-pointer hover:bg-white hover:bg-opacity-30 duration-200 rounded-full"
                  >
                    <Image
                      className="rounded-full shadow bg-white"
                      src={character.src}
                      alt={character.name}
                      width={32}
                      height={32}
                    />
                    <span className="ml-2">{character.name}</span>
                  </a>
                ))}
              </ul>
            )}
          </div>
        </nav>
      </div>
    </>
  );
};

export default HeaderMenu;
