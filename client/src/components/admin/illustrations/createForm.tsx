"use client";

import useSelectCategories from "@/hooks/selectCategories";
import useSelectCharacters from "@/hooks/selectCharacters";
import { Category } from "@/types/category";
import { Character } from "@/types/character";
import { truncateText } from "@/utils/text";
import Image from "next/image";
import { useEffect, useState } from "react";

type Props = {
  characters: Character[];
  categories: Category[];
};

const CreateIllustration: React.FC<Props> = ({ characters, categories }) => {
  const displayLimit = 5;
  const displayTextLimit = 50;

  const [title, setTitle] = useState("");
  const [filename, setFilename] = useState("");
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

  const [parentCategoryIDs, setParentCategoryIDs] = useState<number[]>(
    categories.flatMap((c) => c.ParentCategory.id)
  );

  const [checkedParentCategoryIDs, setCheckedParentCategoryIDs] = useState<
    number[]
  >([]);
  const [originalImageFile, setOriginalImageFile] = useState<File | null>(null);
  const [simpleImageFile, setSimpleImageFile] = useState<File | null>(null);
  const [originalImageData, setOriginalImageData] = useState<string | null>(
    null
  );
  const [simpleImageData, setSimpleImageData] = useState<string | null>(null);

  const onFileChange = (
    event: React.ChangeEvent<HTMLInputElement>,
    setImageData: React.Dispatch<React.SetStateAction<string | null>>,
    setFile: React.Dispatch<React.SetStateAction<File | null>>
  ) => {
    const files = event.target.files;
    if (files && files.length > 0) {
      const selectedFile = files[0];
      const reader = new FileReader();

      reader.onload = (e: ProgressEvent<FileReader>) => {
        setImageData(e.target?.result as string);
      };
      setFile(selectedFile);
      reader.readAsDataURL(selectedFile);
    } else {
      setFile(null);
    }
  };

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
    <>
      <div className="max-w-7xl m-auto">
        <h1 className="text-2xl font-bold mb-6">イラストの作成</h1>
        <form className="border-2 border-gray-300 rounded-lg p-12 bg-white">
          <div className="mb-16">
            <label className="text-xl">タイトル</label>
            <input
              className="w-full p-4 border border-gray-200 rounded mt-2"
              type="text"
              placeholder="イラストのタイトルを入力してください"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
            />
          </div>

          <div className="mb-16">
            <label className="text-xl">キャラクター</label>
            <div className="relative">
              <div
                onClick={() => toggleCharactersModal(!showCharacterModal)}
                className="border-2 border-gray-200 mt-4 py-4 px-4 rounded bg-white flex justify-between flex-nowrap cursor-pointer character-modal"
              >
                <div>
                  {checkedCharacters.length > 0 ? (
                    <div>
                      {checkedCharacters.slice(0, displayLimit).map((char) => (
                        <span
                          key={char.id}
                          className="bg-gray-200 py-2 px-4 mr-2 rounded-full  border-gray-400"
                        >
                          # {truncateText(char.name, 10)}
                        </span>
                      ))}
                      {checkedCharacters.length > displayLimit && (
                        <span>...</span>
                      )}
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
                <div className="absolute left-0 bg-white border-2 border-gray-300 p-4 rounded w-full z-50 shadow-md character-modal-content">
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
                          displayTextLimit
                        )}
                      </label>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          <div className="mb-16">
            <label className="text-xl">カテゴリ</label>
            <div className="relative">
              <div
                onClick={() => toggleCategoriesModal(!showCategoryModal)}
                className="border-2 border-gray-200 mt-4 py-4 px-4 rounded bg-white flex justify-between cursor-pointer category-modal"
              >
                <div>
                  {checkedChildCategories.length > 0 ? (
                    <div>
                      {checkedChildCategories
                        .slice(0, displayLimit)
                        .map((cate) => (
                          <span
                            key={cate.id}
                            className="bg-gray-200 py-2 px-4 mr-2 rounded-full  border-gray-400"
                          >
                            # {truncateText(cate.name, 8)}
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
                </div>
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
                <div className="absolute left-0 bg-white border-2 border-gray-300 p-4 rounded w-full z-50 shadow-md category-modal-content">
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
                        id={`character-${cate.id}`}
                      />
                      <label
                        htmlFor={`character-${cate.id}`}
                        className="cursor-pointer w-full py-1"
                      >
                        {truncateText(
                          categories
                            .flatMap(
                              (c) =>
                                c.ChildCategory.find((cc) => cc.id == cate.id)
                                  ?.name || ""
                            )
                            .join(""),
                          30
                        )}
                      </label>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          <div className="mb-16">
            <label className="text-xl">ファイル名</label>
            <input
              className="w-full p-4 border border-gray-200 rounded mt-2"
              type="text"
              placeholder="ファイル名を入力してください"
              value={filename}
              onChange={(e) => setFilename(e.target.value)}
              required
            />
          </div>

          <div className="flex flex-wrap mb-16">
            <div className="mb-6 mr-2 w-1/3 min-w-[350px]">
              <label className="text-xl w-full bg-green-600 text-white py-2 px-4 rounded-full">
                オリジナル
              </label>
              <div className="border-2 p-4 mt-4 bg-gray-200 rounded-lg w-80 h-80 flex justify-center items-center">
                {originalImageData ? (
                  <div className="relative w-full h-full">
                    <Image
                      src={originalImageData}
                      alt="オリジナル画像プレビュー"
                      layout="fill"
                      objectFit="contain"
                      className="absolute inset-0"
                    />
                  </div>
                ) : (
                  <span className="flex justify-center items-center">
                    Upload Image
                  </span>
                )}
              </div>
              <input
                type="file"
                onChange={(e) =>
                  onFileChange(e, setOriginalImageData, setOriginalImageFile)
                }
                className="w-full mt-4"
                required
              />
            </div>

            <div className="mb-6 mr-2 w-1/3 min-w-[350px]">
              <label className="text-xl w-full bg-gray-200 py-2 px-4 rounded-full">
                シンプル
              </label>
              <div className="border-2 p-4 mt-4 bg-gray-200 rounded-lg w-80 h-80 flex justify-center items-center">
                {simpleImageData ? (
                  <div className="relative w-full h-full">
                    <Image
                      src={simpleImageData}
                      alt="シンプル画像プレビュー"
                      layout="fill"
                      objectFit="contain"
                      className="absolute inset-0"
                    />
                  </div>
                ) : (
                  <span className="flex justify-center items-center">
                    Upload Image
                  </span>
                )}
              </div>
              <input
                type="file"
                onChange={(e) =>
                  onFileChange(e, setSimpleImageData, setSimpleImageFile)
                }
                className="w-full mt-4"
                required
              />
            </div>
          </div>
        </form>
      </div>
    </>
  );
};

export default CreateIllustration;
