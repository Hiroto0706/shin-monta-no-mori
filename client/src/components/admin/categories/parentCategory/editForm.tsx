"use client";

import axios from "axios";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { DeleteCharacterAPI, EditCharacterAPI } from "@/api/admin/character";
import { ParentCategory } from "@/types/category";
import { DeleteParentCategoryAPI, EditParentCategoryAPI } from "@/api/admin/category";

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
          className="border-2 border-gray-300 rounded-lg p-12 bg-white"
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

          <div className="flex flex-wrap mb-16">
            <div className="mb-6 mr-2 w-1/3 min-w-[350px]">
              <div className="border-2 p-4 mt-4 bg-gray-200 rounded-lg w-80 h-80 flex justify-center items-center">
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
