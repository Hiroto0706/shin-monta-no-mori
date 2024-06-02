"use client";

import Image from "next/image";
import { useState } from "react";

const CreateCharacters = () => {
  const [title, setTitle] = useState("");
  const [filename, setFilename] = useState("");
  const [characters, setCharacters] = useState<number[]>([]);
  const [parentCategories, setParentCategories] = useState<number[]>([]);
  const [childCategories, setChildCategories] = useState<number[]>([]);
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

  return (
    <div className="max-w-7xl m-auto">
      <h1 className="text-2xl font-bold mb-6">イラストの作成</h1>
      <form className="border-2 border-gray-300 rounded-lg p-12">
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

        <div className="flex flex-wrap">
          <div className="mb-6 mr-2 w-1/3 min-w-[350px]">
            <label className="text-xl w-full bg-green-600 text-white py-2 px-4 rounded-full">オリジナル</label>
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
          <label className="text-xl w-full bg-gray-200 py-2 px-4 rounded-full">シンプル</label>
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
  );
};

export default CreateCharacters;
