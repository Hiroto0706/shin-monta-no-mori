"use client";

import axios from "axios";
import Image from "next/image";
import { saveAs } from "file-saver";
import { Illustration } from "@/types/illustration";
import { useState } from "react";
import { CreationTimeFormat } from "@/utils/text";
import Link from "next/link";

type Props = {
  illustration: Illustration;
};

const resizeImageAndCenter = async (
  src: string,
  sizePercentage: number
): Promise<string | null> => {
  const response = await axios.get(src, { responseType: "blob" });
  const blob = response.data;

  const canvas = document.createElement("canvas");
  const ctx = canvas.getContext("2d");
  const image = await createImageBitmap(blob);

  const sizeMultiplier = sizePercentage / 100;
  const resizedWidth = image.width * sizeMultiplier;
  const resizedHeight = image.height * sizeMultiplier;

  // キャンバスのサイズを1200px * 1200pxに固定
  const canvasSize = 1200;
  canvas.width = canvasSize;
  canvas.height = canvasSize;

  // 中央にリサイズした画像を配置するためのオフセットを計算
  const offsetX = (canvasSize - resizedWidth) / 2;
  const offsetY = (canvasSize - resizedHeight) / 2;

  ctx?.clearRect(0, 0, canvas.width, canvas.height);
  ctx?.drawImage(image, offsetX, offsetY, resizedWidth, resizedHeight);

  return new Promise<string | null>((resolve) => {
    canvas.toBlob((newBlob) => {
      if (newBlob) {
        resolve(URL.createObjectURL(newBlob)); // Blob URLを返す
      } else {
        resolve(null);
      }
    }, blob.type);
  });
};

const downloadImage = async (src: string, sizePercentage: number) => {
  try {
    const resizedBlobUrl = await resizeImageAndCenter(src, sizePercentage);
    if (resizedBlobUrl) {
      const fileName = src.substring(src.lastIndexOf("/") + 1);
      const a = document.createElement("a");
      a.href = resizedBlobUrl;
      a.download = fileName;
      a.click();
      URL.revokeObjectURL(resizedBlobUrl);
    }
  } catch (error) {
    console.error("Image download failed", error);
  }
};

const copyImageToClipboard = async (
  src: string,
  sizePercentage: number,
  setIsCopied: React.Dispatch<React.SetStateAction<boolean>>
) => {
  try {
    const resizedBlobUrl = await resizeImageAndCenter(src, sizePercentage);
    if (resizedBlobUrl) {
      const response = await fetch(resizedBlobUrl);
      const blob = await response.blob();
      const clipboardItem = new ClipboardItem({
        [blob.type]: blob,
      });
      await navigator.clipboard.write([clipboardItem]);
      setIsCopied(true);
      setTimeout(() => {
        setIsCopied(false);
      }, 3000); // 3秒後にテキストを戻す
    }
  } catch (err) {
    console.error("Failed to copy on clipboard", err);
  }
};

const DetailImage: React.FC<Props> = ({ illustration }) => {
  const [isSimpleImg, setIsSimpleImg] = useState(false);
  const [isCopied, setIsCopied] = useState(false);
  const [size, setSize] = useState(100);
  const [resizedSrc, setResizedSrc] = useState<string | null>(null);

  const handleSliderChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const newSize = parseInt(e.target.value);
    setSize(newSize);

    resizedImageUrl(newSize, isSimpleImg);
  };

  // resizedImageUrl は画像のリサイズを行い、Blob URLを更新
  const resizedImageUrl = async (newSize: number, srcStatus: boolean) => {
    const src = srcStatus
      ? illustration.Image.simple_src.String
      : illustration.Image.original_src;
    const resizedImageUrl = await resizeImageAndCenter(src, newSize);
    setResizedSrc(resizedImageUrl);
  };

  // toggleImageSrc はoriginal_src と simple_src の表示の切り替えをおこなす
  const toggleImageSrc = (newSrcStatus: boolean) => {
    setIsSimpleImg(newSrcStatus);
    resizedImageUrl(size, newSrcStatus);
  };

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
                {resizedSrc ? (
                  <Image
                    className="absolute inset-0 object-cover w-full h-full"
                    src={resizedSrc}
                    alt={illustration.Image.original_filename}
                    fill
                  />
                ) : (
                  <Image
                    className="absolute inset-0 object-cover w-full h-full"
                    src={illustration.Image.original_src}
                    alt={illustration.Image.original_filename}
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
                  onClick={() => toggleImageSrc(false)}
                >
                  オリジナル
                </div>
                <div
                  className={`p-2 rounded-lg w-1/2 flex justify-center font-bold text-lg cursor-pointer ${
                    isSimpleImg && "bg-green-600 text-white"
                  }`}
                  onClick={() => toggleImageSrc(true)}
                >
                  シンプル
                </div>
              </div>
            )}

            <div className="my-2 flex justify-between items-center resize-bar">
              <label className="w-[20%] min-w-[60px] max-w-[70px] flex justify-between pr-2 text-lg">
                <span>{size}</span>
                <span className="text-gray-400">%</span>
              </label>
              <input
                className="w-[80%] rounded-full"
                type="range"
                min="10"
                max="100"
                step="10"
                value={size}
                onChange={handleSliderChange}
                // sizeは10%が最低値のため -10 している。また、100 / 90 は 90% を 100% に正規化している
                style={{
                  background: `linear-gradient(to right, #17a34a ${
                    (size - 10) * (100 / 90)
                  }%, #E5E7EB ${(size - 10) * (100 / 90)}%)`,
                }}
              />
            </div>
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
                  {CreationTimeFormat(illustration.Image.created_at)}
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
                          : illustration.Image.original_src,
                        size
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
                        size,
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
                {/* FIXME: sp時の画像保存ボタンは画像リサイズしたとしてもデフォルトの画像を返すようになっている */}
                {/* spの時の画像保存ボタン
                <Link
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
                </Link> */}
                <p className="text-sm">
                  ダウンロードボタンをクリックすると、
                  <Link
                    href="/terms-of-service"
                    className="text-blue-600 underline border-blue-600 text-blue-600 cursor-pointer hover:text-blue-700 duration-200"
                  >
                    利用規約
                  </Link>
                  および
                  <Link
                    href="/privacy-policy"
                    className="text-blue-600 underline border-blue-600 text-blue-600 cursor-pointer hover:text-blue-700 duration-200"
                  >
                    プライバシーポリシー
                  </Link>
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
                      <Link
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
                      </Link>
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
                              <Link
                                key={cc.id}
                                href={`/illustrations/category/${cc.id}`}
                                className="my-1 ml-1 mr-2 py-1.5 px-2 cursor-pointer duration-200 hover:bg-gray-200 hover:bg-opacity-70 rounded-full"
                              >
                                # {cc.name}
                              </Link>
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
          <Link
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
          </Link>
          <Link
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
          </Link>
          <Link
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
          </Link>
        </div>
      </section>
    </>
  );
};

export default DetailImage;
