"use client";

import axios from "axios";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { DeleteCharacterAPI, EditCharacterAPI } from "@/api/admin/character";
import { ParentCategory } from "@/types/category";
import {
  DeleteParentCategoryAPI,
  EditParentCategoryAPI,
} from "@/api/admin/category";
import usePriorityLevel from "@/hooks/priorityLevel";
import { PriorityLevel } from "@/types/admin/priorityLevel";

type Props = {
  id: number;
  parentCategory: ParentCategory;
  accessToken: string | undefined;
};

const EditParentCategory: React.FC<Props> = ({
  id,
  parentCategory,
  accessToken,
}) => {
  const router = useRouter();

  const [name, setName] = useState(parentCategory.name);
  const [filename, setFilename] = useState(parentCategory.filename.String);

  const [imageFile, setImageFile] = useState<File | null>(null);
  const [imageSrc, setImageSrc] = useState<string | null>(parentCategory.src);

  const {
    checkedPriorityLevel,
    setCheckedPriorityLevel,
    showPriorityLevelModal,
    togglePriorityLevelModal,
  } = usePriorityLevel(parentCategory.priority_level);

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

  const editParentCategory = async (event: React.FormEvent) => {
    event.preventDefault();

    if (filename != parentCategory.filename.String && imageFile == null) {
      alert("新しい画像を設定してください");
      return;
    }

    const formData = new FormData();
    formData.append("name", name);
    formData.append("filename", filename);
    if (imageFile) {
      formData.append("image_file", imageFile);
    }
    formData.append("priority_level", checkedPriorityLevel.toString());

    try {
      const response = await axios.put(EditParentCategoryAPI(id), formData, {
        headers: {
          Authorization: SetBearerToken(accessToken),
        },
      });

      if (response.status === 200) {
        alert(response.data.message);
      }
    } catch (error) {
      console.error("親カテゴリの編集に失敗しました", error);
      alert("親カテゴリの編集に失敗しました");
    }
  };

  const deleteParentCategory = async (id: number) => {
    if (!confirm(`本当に「${name}」を削除してもよろしいですか？`)) {
      return;
    }

    try {
      const response = await axios.delete(DeleteParentCategoryAPI(id), {
        headers: {
          Authorization: SetBearerToken(accessToken),
        },
      });

      if (response.status === 200) {
        alert(response.data.message);
        router.push("/admin/categories");
      }
    } catch (error) {
      console.error("親カテゴリの削除に失敗しました", error);
      alert("親カテゴリの削除に失敗しました");
    }
  };

  return (
    <>
      <div className="max-w-7xl m-auto">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold">親カテゴリの編集</h1>
          <button
            onClick={() => deleteParentCategory(id)}
            className="bg-red-500 text-white py-2 px-4 rounded-lg flex items-center"
          >
            <Image
              src="/icon/trash.png"
              alt="trashアイコン"
              width={20}
              height={20}
            />
            <span className="ml-1">削除</span>
          </button>
        </div>

        <form
          className="border-2 border-gray-300 rounded-lg p-2 md:p-12 bg-white"
          onSubmit={editParentCategory}
        >
          <div className="mb-16">
            <label className="text-xl">名前</label>
            <input
              className="w-full p-4 border border-gray-200 rounded mt-2"
              type="text"
              placeholder="親カテゴリの名前を入力してください"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
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

          <div className="mb-16">
            <div className="mb-6 mr-2 w-1/3 min-w-[300px] md:min-w-[350px]">
              <div className="border-2 p-4 mt-4 bg-gray-200 rounded-lg w-60 h-60 md:w-80 md:h-80 flex justify-center items-center">
                {imageSrc ? (
                  <div className="relative w-full h-full">
                    <Image
                      src={imageSrc}
                      alt="画像プレビュー"
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
                onChange={(e) => onFileChange(e, setImageSrc, setImageFile)}
                className="w-full mt-4"
                required={imageSrc !== "" ? false : true}
              />
            </div>
            <p className="text-sm">※ png形式の画像をアップロードしてください</p>
          </div>

          <div className="mb-16">
            <label className="text-xl">優先度</label>
            <div className="relative">
              <div
                onClick={() =>
                  togglePriorityLevelModal(!showPriorityLevelModal)
                }
                className="border-2 border-gray-200 mt-4 py-4 px-4 rounded bg-white flex justify-between flex-nowrap cursor-pointer priority-modal"
              >
                <div>{PriorityLevel[checkedPriorityLevel]}</div>
                <Image
                  className={`duration-100 ${
                    !showPriorityLevelModal ? "rotate-90" : "-rotate-90"
                  }`}
                  src="/icon/arrow.png"
                  alt="arrowアイコン"
                  width={20}
                  height={20}
                />
              </div>
              {showPriorityLevelModal && (
                <div className="absolute left-0 bg-white border-2 border-gray-300 p-4 rounded w-full max-h-60 overflow-y-auto z-10 shadow-md priority-modal-content">
                  {PriorityLevel.map((level, i) => (
                    <div
                      key={i}
                      className="flex items-center rounded hover:bg-gray-200 cursor-pointer"
                    >
                      <input
                        type="radio"
                        checked={checkedPriorityLevel === i}
                        onChange={() => setCheckedPriorityLevel(i)}
                        className="mx-2 cursor-pointer"
                        id={`character-${i}`}
                      />
                      <label
                        htmlFor={`character-${i}`}
                        className="cursor-pointer w-full py-1"
                      >
                        {level}
                      </label>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          <button className="py-3 bg-green-600 text-white font-bold text-lg rounded-lg w-full hover:bg-white hover:text-green-600 border-2 border-green-600 duration-200">
            親カテゴリ編集
          </button>
        </form>
      </div>
    </>
  );
};

export default EditParentCategory;
