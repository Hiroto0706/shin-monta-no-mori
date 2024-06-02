"use client";

import Image from "next/image";
import { useEffect, useState } from "react";
import { Character } from "@/types/character";
import { truncateText } from "@/utils/text";
import { Category } from "@/types/category";
import useSelectCharacters from "@/hooks/selectCharacters";
import useSelectCategories from "@/hooks/selectCategories";

type Props = {
  characters: Character[];
  categories: Category[];
};

const SearchBox: React.FC<Props> = ({ characters, categories }) => {
  const displayLimit = 3;
  const {
    checkedCharacters,
    showCharacterModal,
    handleCharacterSelect,
    toggleCharactersModal,
  } = useSelectCharacters();
  const {
    childCategories,
    checkedChildCategories,
    showCategoryModal,
    handleCategoriesSelect,
    toggleCategoriesModal,
  } = useSelectCategories(categories);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        (event.target as HTMLElement).closest(".character-modal") === null &&
        (event.target as HTMLElement).closest(".character-modal-content") ===
          null
      ) {
        toggleCharactersModal(false);
      }
      if (
        (event.target as HTMLElement).closest(".category-modal") === null &&
        (event.target as HTMLElement).closest(".category-modal-content") ===
          null
      ) {
        toggleCategoriesModal(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return (
    <div className="flex flex-col flex-col-reverse lg:flex-row justify-between">
      <form className="flex flex-wrap">
        <div className="lg:mr-3 mb-6 lg:mb-3 w-full lg:w-80">
          <input
            type="text"
            placeholder="タイトル検索"
            className="border-2 border-gray-200 py-3 px-4 rounded-md w-full"
          />
        </div>

        <div className="pr-2 lg:pr-0 lg:mr-3 mb-6 lg:mb-3 relative w-1/2 lg:w-80">
          <div
            onClick={() => toggleCharactersModal(!showCharacterModal)}
            className="border-2 border-gray-200 py-3 px-4 rounded bg-white flex justify-between cursor-pointer character-modal"
          >
            <div>
              {checkedCharacters.length > 0 ? (
                <div>
                  {checkedCharacters.slice(0, displayLimit).map((char) => (
                    <span
                      key={char.id}
                      className="bg-gray-200 py-2 px-2 mr-1 rounded-full border-gray-400 text-ellipsis overflow-hidden whitespace-nowrap"
                    >
                      # {truncateText(char.name, 4)}
                    </span>
                  ))}
                  {checkedCharacters.length > displayLimit && <span>...</span>}
                </div>
              ) : (
                <span className="text-gray opacity-50">
                  キャラクターを選択してください
                </span>
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
            <div className="relative w-full">
              <div className="absolute bg-white border-2 border-gray-300 p-4 rounded w-full z-50 shadow-md character-modal-content">
                {characters.map((char) => (
                  <div
                    key={char.id}
                    className="flex items-center rounded hover:bg-gray-200 cursor-pointer"
                  >
                    <input
                      type="checkbox"
                      checked={checkedCharacters.includes(char)}
                      onChange={() => handleCharacterSelect(char)}
                      className="mx-2 cursor-pointer"
                      id={`character-${char.id}`}
                    />
                    <label
                      htmlFor={`character-${char.id}`}
                      className="cursor-pointer w-full py-1"
                    >
                      {truncateText(
                        characters.find((c) => c.id == char.id)?.name,
                        25
                      )}
                    </label>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        <div className="pl-2 lg:pl-0 lg:mr-3 mb-6 lg:mb-3 relative w-1/2 lg:w-80">
          <div
            onClick={() => toggleCategoriesModal(!showCategoryModal)}
            className="border-2 border-gray-200 py-3 px-4 rounded bg-white flex justify-between cursor-pointer category-modal"
          >
            <div>
              {checkedChildCategories.length > 0 ? (
                <div>
                  {checkedChildCategories.slice(0, displayLimit).map((cate) => (
                    <span
                      key={cate.id}
                      className="bg-gray-200 py-2 px-2 mr-1 rounded-full border-gray-400 text-ellipsis overflow-hidden whitespace-nowrap"
                    >
                      # {truncateText(cate.name, 4)}
                    </span>
                  ))}
                  {checkedChildCategories.length > displayLimit && (
                    <span>...</span>
                  )}
                </div>
              ) : (
                <span className="text-gray opacity-50">
                  カテゴリを選択してください
                </span>
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
            <div className="relative w-full">
              <div className="absolute bg-white left-0 border-2 border-gray-300 p-4 w-full z-50 shadow-md  category-modal-content">
                {childCategories.map((cate) => (
                  <div
                    key={cate.id}
                    className="flex items-center rounded hover:bg-gray-200 cursor-pointer"
                  >
                    <input
                      type="checkbox"
                      checked={checkedChildCategories.includes(cate)}
                      onChange={() => handleCategoriesSelect(cate)}
                      className="mx-2 cursor-pointer"
                      id={`category-${cate.id}`}
                    />
                    <label
                      htmlFor={`category-${cate.id}`}
                      className="cursor-pointer w-full py-1"
                    >
                      {truncateText(
                        childCategories.find((c) => c.id == cate.id)?.name,
                        25
                      )}
                    </label>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        <button className="flex justify-center items-center lg:justify-start bg-green-600 text-white rounded-md font-bold py-3 pl-4 pr-3 lg:mb-6 w-full lg:w-auto hover:opacity-70 duration-200">
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
