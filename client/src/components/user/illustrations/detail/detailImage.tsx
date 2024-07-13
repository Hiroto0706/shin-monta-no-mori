"use client";

import axios from "axios";
import Image from "next/image";
import { saveAs } from "file-saver";
import { Illustration } from "@/types/illustration";
import { useState } from "react";
import { UpdatedAtFormat } from "@/utils/text";

type Props = {
  illustration: Illustration;
};

export const downloadImage = async (src: string) => {
  try {
    const response = await axios.get(src, {
      responseType: "blob",
    });
    const fileName = src.substring(src.lastIndexOf("/") + 1);
    saveAs(response.data, fileName);
  } catch (error) {
    console.error("Image download failed", error);
  }
};

export const copyImageToClipboard = async (
  src: string,
  setIsCopied: React.Dispatch<React.SetStateAction<boolean>>
) => {
  try {
    const response = await axios.get(src, { responseType: "blob" });
    const blob = response.data;

    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");
    const image = await createImageBitmap(blob);
    canvas.width = image.width;
    canvas.height = image.height;
    if (ctx) {
      ctx.drawImage(image, 0, 0);
    }

    canvas.toBlob(async (newBlob) => {
      if (newBlob) {
        const clipboardItem = new ClipboardItem({ [newBlob.type]: newBlob });
        await navigator.clipboard.write([clipboardItem]);
        setIsCopied(true);
        setTimeout(() => {
          setIsCopied(false);
        }, 3000); // 3秒後にテキストを戻す
      }
    }, blob.type);
  } catch (err) {
    console.error("Failed to copy on clipboard", err);
  }
};

const DetailImage: React.FC<Props> = ({ illustration }) => {
  const [isSimpleImg, setIsSimpleImg] = useState(false);
  const [isCopied, setIsCopied] = useState<boolean>(false);
  const siteUrl = "https://www.montanomori.com/";

  return (
    <>
      <section className="mb-20">
        <div className="flex flex-col lg:flex-row">
          <div className="w-full max-w-[500px] m-auto lg:m-0 lg:w-2/5 lg:max-w-[400px]">
            <div
              className="relative w-full border-2 border-gray-200 rounded-xl repeat-bg-image bg-white"
              style={{ paddingTop: "100%" }}
            >
              <div className="absolute inset-0 m-4">
                {!isSimpleImg ? (
                  <Image
                    className="absolute inset-0 object-cover w-full h-full"
                    src={illustration.Image.original_src}
                    alt={illustration.Image.original_filename}
                    fill
                  />
                ) : (
                  <Image
                    className="absolute inset-0 object-cover w-full h-full"
                    src={illustration.Image.simple_src.String}
                    alt={illustration.Image.simple_filename.String}
                    fill
                  />
                )}
              </div>
            </div>

            {illustration.Image.simple_src.Valid && (
              <div className="mt-2 flex justify-between items-center p-1 bg-gray-100 rounded-lg">
                <div
                  className={`p-2 rounded-lg w-1/2 flex justify-center font-bold text-lg cursor-pointer ${
                    !isSimpleImg && "bg-green-600 text-white"
                  }`}
                  onClick={() => setIsSimpleImg(false)}
                >
                  オリジナル
                </div>
                <div
                  className={`p-2 rounded-lg w-1/2 flex justify-center font-bold text-lg cursor-pointer ${
                    isSimpleImg && "bg-green-600 text-white"
                  }`}
                  onClick={() => setIsSimpleImg(true)}
                >
                  シンプル
                </div>
              </div>
            )}
          </div>

          <div className="w-full lg:w-3/5 mt-8 lg:mt-0 lg:ml-8">
            <h2 className="mb-4">
              <div className="text-2xl font-bold mb-1 text-black">
                {illustration.Image.title}
              </div>
              <div className="flex items-center">
                <Image
                  className="mr-1"
                  src="/icon/time.png"
                  alt="timeアイコン"
                  width={12}
                  height={12}
                />
                <span className="text-xs text-gray-400">
                  {UpdatedAtFormat(illustration.Image.updated_at)}
                </span>
              </div>
            </h2>

            <>
              <div className="mb-8">
                {/* sp以外の時の画像保存ボタン */}
                <div className="flex flex-wrap mb-2 hidden sm:flex">
                  <button
                    className="bg-green-600 text-white duration-200 font-bold rounded-lg py-2 px-4 hover:bg-green-700 flex items-center"
                    onClick={() =>
                      downloadImage(
                        isSimpleImg
                          ? illustration.Image.simple_src.String
                          : illustration.Image.original_src
                      )
                    }
                  >
                    <Image
                      src="/icon/download.png"
                      alt="ダウンロードアイコン"
                      width={24}
                      height={24}
                    />
                    <span className="ml-1">ダウンロード</span>
                  </button>
                  <button
                    className="bg-gray-200 font-bold rounded-lg py-2 px-4 ml-4 duration-200 hover:bg-gray-300 flex items-center"
                    onClick={() =>
                      copyImageToClipboard(
                        isSimpleImg
                          ? illustration.Image.simple_src.String
                          : illustration.Image.original_src,
                        setIsCopied
                      )
                    }
                  >
                    <Image
                      src={
                        !isCopied ? "/icon/copy.png" : "/icon/copy-success.png"
                      }
                      alt="コピーアイコン"
                      width={24}
                      height={24}
                    />
                    <span className="ml-1">
                      {!isCopied ? (
                        <span>コピー</span>
                      ) : (
                        <span className="text-green-600">コピーしました</span>
                      )}
                    </span>
                  </button>
                </div>
                {/* spの時の画像保存ボタン */}
                <a
                  href={
                    !isSimpleImg
                      ? illustration.Image.original_src
                      : illustration.Image.simple_src.String
                  }
                  className="w-full mb-2 py-2 px-4 border rounded-lg bg-green-600 block sm:hidden text-white font-bold text-lg flex justify-between cursor-pointer duration-200 hover:bg-green-700"
                >
                  <span>ダウンロード</span>
                  <Image
                    src="/icon/download.png"
                    alt="ダウンロードアイコン"
                    width={24}
                    height={24}
                  />
                </a>
                <p className="text-sm">
                  ダウンロードボタンをクリックすると、
                  <a
                    href="/terms-of-service"
                    className="text-blue-600 underline border-blue-600 text-blue-600 cursor-pointer hover:text-blue-700 duration-200"
                  >
                    利用規約
                  </a>
                  に同意したものとみなされます。
                </p>
              </div>
            </>

            {illustration.Characters.length > 0 && (
              <>
                <div className="mb-8">
                  <h3 className="text-lg font-bold mb-2 text-black">
                    キャラクター
                  </h3>
                  <div className="flex flex-wrap">
                    {illustration.Characters.map((c, index) => (
                      <a
                        key={c.Character.id}
                        href={`/illustrations/character/${c.Character.id}`}
                        className="flex items-center rounded-full py-1 pl-1 pr-4 cursor-pointer hover:bg-gray-200 hover:bg-opacity-70 duration-200 m-1"
                      >
                        <Image
                          className="border shadow border-gray-200 rounded-full bg-white"
                          src={c.Character.src}
                          alt={c.Character.name}
                          width={40}
                          height={40}
                        />
                        <span className="ml-2">{c.Character.name}</span>
                      </a>
                    ))}
                  </div>
                </div>
              </>
            )}

            {illustration.Categories.length > 0 && (
              <>
                <div>
                  <h3 className="text-lg font-bold mb-2 text-black">
                    カテゴリ
                  </h3>
                  <div>
                    <>
                      {/* <div className="flex flex-wrap mb-2">
                          {illustration.Categories.map(
                            (category, index) => (
                              <div
                                key={category.ParentCategory.id}
                                className={`border py-1 px-3 rounded-full flex items-center bg-white ${
                                  index == 0 ? "ml-0" : "ml-2"
                                }`}
                              >
                                <Image
                                  src={category.ParentCategory.src}
                                  alt={category.ParentCategory.name}
                                  width={28}
                                  height={28}
                                />
                                <span className="ml-1 font-bold">
                                  {category.ParentCategory.name}
                                </span>
                              </div>
                            )
                          )}
                        </div> */}

                      <div className="flex flex-wrap">
                        {illustration.Categories.map((category) => (
                          <>
                            {category.ChildCategory.map((cc) => (
                              <a
                                key={cc.id}
                                href={`/illustrations/category/${cc.id}`}
                                className="my-1 ml-1 mr-2 py-1.5 px-2 cursor-pointer duration-200 hover:bg-gray-200 hover:bg-opacity-70 rounded-full"
                              >
                                # {cc.name}
                              </a>
                            ))}
                          </>
                        ))}
                      </div>
                    </>
                  </div>
                </div>
              </>
            )}
          </div>
        </div>

        <div className="mt-8 flex justify-start sm:justify-end">
          <a
            href={`http://twitter.com/share?url=${siteUrl}illustrations/${illustration.Image.id}&text=${illustration.Image.title}の画像`}
            target="_blank"
            className="pl-1 pr-2 flex items-center border bg-[#1F1F1F] rounded-lg text-white duration-200 cursor-pointer hover:opacity-70"
          >
            <Image
              src="/sns/twitter.png"
              alt="twitterアイコン"
              width={28}
              height={28}
            />
            <span className="text-sm">でシェア</span>
          </a>
          <a
            href={`https://www.facebook.com/share.php?u=${siteUrl}illustrations/${illustration.Image.id}`}
            target="_blank"
            className="ml-1 pl-1 pr-2 flex items-center border bg-[#3F50B6] rounded-lg text-white duration-200 cursor-pointer hover:opacity-70"
          >
            <Image
              src="/sns/facebook.png"
              alt="facebookアイコン"
              width={28}
              height={28}
            />
            <span className="text-sm">でシェア</span>
          </a>
          <a
            href={`https://social-plugins.line.me/lineit/share?url=${siteUrl}illustrations/${illustration.Image.id}`}
            target="_blank"
            className="ml-1 pl-1 pr-2 flex items-center border bg-[#01C400] rounded-lg text-white duration-200 cursor-pointer hover:opacity-70"
          >
            <Image
              className="p-1"
              src="/sns/line.png"
              alt="lineアイコン"
              width={28}
              height={28}
            />
            <span className="text-sm">でシェア</span>
          </a>
        </div>
      </section>
    </>
  );
};

export default DetailImage;
