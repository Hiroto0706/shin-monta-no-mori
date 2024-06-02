"use client";

import Image from "next/image";
import { useEffect, useState } from "react";
import { FetchIllustrationsResponse } from "@/types/illustration";
import { Character } from "@/types/character";
import { truncateText } from "@/utils/text";
import { Category } from "@/types/category";

type Props = {
  characters: Character[];
  categories: Category[];
};

const SearchBox: React.FC<Props> = ({
  characters,
  categories,
}) => {
  const displayLimit = 4;
  const charactersIDs = characters.map((c) => c.id);
  const [checkedCharactersIDs, setCheckedCharactersIDs] = useState<number[]>(
    []
  );
  const categoriesIDs = categories.flatMap((c) =>
    c.ChildCategory.map((child) => child.id)
  );
  const [checkedCategoriesIDs, setCheckedCategoriesIDs] = useState<number[]>(
    []
  );
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
    setCheckedCharactersIDs((prev) => {
      if (prev.includes(id)) {
        return prev.filter((charId) => charId !== id);
      }
      return [...prev, id];
    });
  };

  const handleCategoriesSelect = (id: number) => {
    setCheckedCategoriesIDs((prev) => {
      if (prev.includes(id)) {
        return prev.filter((cateId) => cateId !== id);
      }
      return [...prev, id];
    });
  };

  return (
    <div className="flex flex-col flex-col-reverse lg:flex-row justify-between">
      <form className="flex flex-wrap">
        <div className="lg:mr-3 mb-6 lg:mb-0 w-full lg:w-72">
          <input
            type="text"
            placeholder="タイトル検索"
            className="border-2 border-gray-200 py-3 px-4 rounded-md w-full"
          />
        </div>

        <div className="pr-2 lg:pr-0 lg:mr-3 mb-6 lg:mb-0 relative w-1/2 lg:w-52">
          <div
            onClick={() => setShowCharacterModal(!showCharacterModal)}
            className="border-2 border-gray-200 py-3 px-4 rounded bg-white flex justify-between cursor-pointer character-modal"
          >
            <div>
              {checkedCharactersIDs.length > 0 ? (
                <div>
                  {checkedCharactersIDs.slice(0, displayLimit).map((id) => (
                    <span
                      key={id}
                      className="bg-gray-200 px-1 mr-1 rounded border-2 border-gray-400"
                    >
                      {id}
                    </span>
                  ))}
                  {checkedCharactersIDs.length > displayLimit && (
                    <span>...</span>
                  )}
                </div>
              ) : (
                <span>キャラクター</span>
              )}
            </div>
            <Image
              className={`duration-100 ${
                !showCharacterModal ? "rotate-90" : "-rotate-90"
              }`}
              src="/icon/arrow.png"
              alt="arrowアイコン"
              width={20}
              height={20}
            />
          </div>
          {showCharacterModal && (
            <div className="absolute top-16 left-0 bg-white border-2 border-gray-300 p-4 rounded w-60 z-50 shadow-md character-modal-content">
              {charactersIDs.map((charId) => (
                <div
                  key={charId}
                  className="flex items-center mb-2 rounded hover:bg-gray-200 cursor-pointer"
                >
                  <input
                    type="checkbox"
                    checked={checkedCharactersIDs.includes(charId)}
                    onChange={() => handleCharacterSelect(charId)}
                    className="mx-2 cursor-pointer"
                    id={`character-${charId}`}
                  />
                  <label
                    htmlFor={`character-${charId}`}
                    className="cursor-pointer"
                  >
                    {truncateText(characters.find((c) => c.id == charId)?.name)}
                  </label>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="pl-2 lg:pl-0 lg:mr-3 mb-6 lg:mb-0 relative w-1/2 lg:w-52">
          <div
            onClick={() => setShowCategoryModal(!showCategoryModal)}
            className="border-2 border-gray-200 py-3 px-4 rounded bg-white flex justify-between cursor-pointer category-modal"
          >
            <div>
              {checkedCategoriesIDs.length > 0 ? (
                <div>
                  {checkedCategoriesIDs.slice(0, displayLimit).map((id) => (
                    <span
                      key={id}
                      className="bg-gray-200 px-1 mr-1 rounded border-2 border-gray-400"
                    >
                      {id}
                    </span>
                  ))}
                  {checkedCategoriesIDs.length > displayLimit && (
                    <span>...</span>
                  )}
                </div>
              ) : (
                <span>子カテゴリ</span>
              )}
            </div>{" "}
            <Image
              className={`duration-100 ${
                !showCategoryModal ? "rotate-90" : "-rotate-90"
              }`}
              src="/icon/arrow.png"
              alt="arrowアイコン"
              width={20}
              height={20}
            />
          </div>
          {showCategoryModal && (
            <div className="absolute top-16 left-0 bg-white border-2 border-gray-300 p-4 rounded w-60 z-50 shadow-md  category-modal-content">
              {categoriesIDs.map((cateId) => (
                <div
                  key={cateId}
                  className="flex items-center mb-2 rounded hover:bg-gray-200 cursor-pointer"
                >
                  <input
                    type="checkbox"
                    checked={checkedCategoriesIDs.includes(cateId)}
                    onChange={() => handleCategoriesSelect(cateId)}
                    className="mx-2 cursor-pointer"
                    id={`category-${cateId}`}
                  />
                  <label
                    htmlFor={`category-${cateId}`}
                    className="cursor-pointer"
                  >
                    {truncateText(
                      categories
                        .flatMap(
                          (c) =>
                            c.ChildCategory.find((cc) => cc.id == cateId)
                              ?.name || ""
                        )
                        .join("")
                    )}
                  </label>
                </div>
              ))}
            </div>
          )}
        </div>

        <button className="flex justify-center items-center lg:justify-start bg-green-600 text-white rounded-md font-bold py-3 pl-4 pr-3 lg:mb-0 w-full lg:w-auto hover:opacity-70 duration-200">
          <span className="mr-1">検索</span>
          <Image
            src="/icon/search.png"
            alt="searchアイコン"
            width={24}
            height={24}
          />
        </button>
      </form>
    </div>
  );
};

export default SearchBox;
