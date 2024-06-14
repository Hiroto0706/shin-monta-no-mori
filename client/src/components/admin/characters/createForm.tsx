"use client";

import { useRouter } from "next/navigation";
import Image from "next/image";
import axios from "axios";
import { useState } from "react";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { CreateCharacterAPI } from "@/api/admin/character";

type Props = {
  accessToken: string | undefined;
};

const CreateCharacter: React.FC<Props> = ({ accessToken }) => {
  const router = useRouter();

  const [name, setName] = useState("");
  const [filename, setFilename] = useState("");

  const [imageFile, setImageFile] = useState<File | null>(null);
  const [imageSrc, setImageSrc] = useState<string | null>(null);

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

  const createIllustration = async (event: React.FormEvent) => {
    event.preventDefault();

    const formData = new FormData();
    formData.append("name", name);
    formData.append("filename", filename);
    if (imageFile) {
      formData.append("image_file", imageFile);
    }

    try {
      const response = await axios.post(CreateCharacterAPI(), formData, {
        headers: {
          Authorization: SetBearerToken(accessToken),
        },
      });

      if (response.status === 200) {
        alert(response.data.message);
        router.push("/admin/characters");
      }
    } catch (error) {
      console.error("キャラクターの作成に失敗しました", error);
      alert("キャラクターの作成に失敗しました");
    }
  };

  return (
    <>
      <div className="max-w-7xl m-auto">
        <h1 className="text-2xl font-bold mb-6">キャラクターの作成</h1>
        <form
          className="border-2 border-gray-300 rounded-lg p-12 bg-white"
          onSubmit={createIllustration}
        >
          <div className="mb-16">
            <label className="text-xl">名前</label>
            <input
              className="w-full p-4 border border-gray-200 rounded mt-2"
              type="text"
              placeholder="キャラクターの名前を入力してください"
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
                required
              />
            </div>
          </div>

          <button className="py-3 bg-green-600 text-white font-bold text-lg rounded-lg w-full hover:bg-white hover:text-green-600 border-2 border-green-600 duration-200">
            キャラクター作成
          </button>
        </form>
      </div>
    </>
  );
};

export default CreateCharacter;
