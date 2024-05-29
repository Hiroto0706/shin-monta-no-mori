"use client";

import Image from "next/image";
import { Illustration } from "@/types/illustration";
import { useState, useEffect } from "react";

interface IllustrationsProps {
  illustrations: Illustration[];
}

export default function Illustrations({ illustrations }: IllustrationsProps) {
  const [characters, setCharacters] = useState<number[]>([]);
  const [categories, setCategories] = useState<number[]>([]);
  const [showCharacterModal, setShowCharacterModal] = useState(false);
  const [showCategoryModal, setShowCategoryModal] = useState(false);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        (event.target as HTMLElement).closest(".character-modal") === null &&
        (event.target as HTMLElement).closest(".character-modal-content") ===
          null
      ) {
        setShowCharacterModal(false);
      }
      if (
        (event.target as HTMLElement).closest(".category-modal") === null &&
        (event.target as HTMLElement).closest(".category-modal-content") ===
          null
      ) {
        setShowCategoryModal(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const handleCharacterSelect = (id: number) => {
    setCharacters((prev) => {
      if (prev.includes(id)) {
        return prev.filter((charId) => charId !== id);
      }
      return [...prev, id];
    });
  };

  const handleCategoriesSelect = (id: number) => {
    setCategories((prev) => {
      if (prev.includes(id)) {
        return prev.filter((cateId) => cateId !== id);
      }
      return [...prev, id];
    });
  };

  if (!illustrations || illustrations.length === 0) {
    return <div>No illustrations available</div>;
  }

  return (
    <div>
      <div className="flex flex-col flex-col-reverse lg:flex-row justify-between">
        <form className="lg:flex">
          <div className="lg:mr-3 mb-6 lg:mb-0 w-full lg:w-auto">
            <input
              type="text"
              placeholder="タイトル検索"
              className="border-2 border-gray-200 py-2 px-4 rounded-md w-full"
            />
          </div>

          <div className="lg:mr-3 mb-6 lg:mb-0 relative">
            <div
              onClick={() => setShowCharacterModal(!showCharacterModal)}
              className="border-2 border-gray-200 py-2 px-4 rounded bg-white flex justify-between cursor-pointer w-full lg:w-44 character-modal"
            >
              <span>キャラクター</span>
              <Image
                className={!showCharacterModal ? `-rotate-90` : `rotate-90`}
                src="/icon/arrow.png"
                alt="arrowアイコン"
                width={20}
                height={20}
              />
            </div>
            {showCharacterModal && (
              <div className="absolute top-12 bg-white border-2 border-gray-300 p-4 rounded w-48 z-50 shadow-md character-modal-content">
                {[1, 2, 3].map((charId) => (
                  <div
                    key={charId}
                    className="flex items-center mb-2 rounded hover:bg-gray-200 cursor-pointer"
                  >
                    <input
                      type="checkbox"
                      checked={characters.includes(charId)}
                      onChange={() => handleCharacterSelect(charId)}
                      className="mx-2 cursor-pointer"
                      id={`character-${charId}`}
                    />
                    <label
                      htmlFor={`character-${charId}`}
                      className="cursor-pointer"
                    >
                      Character {charId}
                    </label>
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="lg:mr-3 mb-6 lg:mb-0 relative">
            <div
              onClick={() => setShowCategoryModal(!showCategoryModal)}
              className="border-2 border-gray-200 py-2 px-4 rounded bg-white flex justify-between cursor-pointer w-full lg:w-44 category-modal"
            >
              <span>カテゴリ</span>
              <Image
                className={!showCategoryModal ? `-rotate-90` : `rotate-90`}
                src="/icon/arrow.png"
                alt="arrowアイコン"
                width={20}
                height={20}
              />
            </div>
            {showCategoryModal && (
              <div className="absolute top-12 bg-white border-2 border-gray-300 p-4 rounded w-48 z-50 shadow-md  category-modal-content">
                {[1, 2, 3].map((cateId) => (
                  <div
                    key={cateId}
                    className="flex items-center mb-2 rounded hover:bg-gray-200 cursor-pointer"
                  >
                    <input
                      type="checkbox"
                      checked={categories.includes(cateId)}
                      onChange={() => handleCategoriesSelect(cateId)}
                      className="mx-2 cursor-pointer"
                      id={`category-${cateId}`}
                    />
                    <label
                      htmlFor={`category-${cateId}`}
                      className="cursor-pointer"
                    >
                      Category {cateId}
                    </label>
                  </div>
                ))}
              </div>
            )}
          </div>

          <button className="flex justify-center items-center lg:justify-start bg-green-600 text-white rounded-md font-bold py-2 pl-3 pr-2 lg:mb-0 w-full lg:w-auto">
            <span className="mr-2">検索</span>
            <Image
              src="/icon/search.png"
              alt="searchアイコン"
              width={24}
              height={24}
            />
          </button>
        </form>

        <a
          href="illustrations/new"
          className="flex items-center bg-green-600 text-white py-2.5 font-bold rounded-md mb-6 ml-auto lg:mb-0 w-36 justify-center"
        >
          + イラスト追加
        </a>
      </div>

      <div className="my-12">
        <ul className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
          {illustrations.map((illustration, index) => (
            <li key={index} className="p-2 border-2 border-gray-200 rounded-xl">
              <div className="flex items-center">
                <span className="ml-4 font-bold text-2xl">
                  {illustration.Image.title}
                </span>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
