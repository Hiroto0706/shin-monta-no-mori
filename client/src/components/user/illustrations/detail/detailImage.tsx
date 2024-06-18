"use client";

import { Illustration } from "@/types/illustration";
import Image from "next/image";
import { useState } from "react";

type Props = {
  illustration: Illustration;
};

const DetailImage: React.FC<Props> = ({ illustration }) => {
  const [isSimpleImg, setIsSimpleImg] = useState(false);

  return (
    <>
      <section className="mb-20">
        <div className="flex flex-col lg:flex-row">
          <div className="w-full max-w-[500px] m-auto lg:m-0 lg:w-2/5 lg:max-w-[400px]">
            <div
              className="relative w-full border-2 border-gray-200 rounded-xl"
              style={{ paddingTop: "100%" }}
            >
              <div className="absolute inset-0 m-8">
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
            <h2 className="text-2xl font-bold mb-4 text-black">
              {illustration.Image.title}
            </h2>

            <div className="mb-8">
              <div className="flex flex-wrap mb-2">
                <button className="bg-green-600 text-white duration-200 font-bold rounded-lg py-2 px-4 hover:bg-green-700">
                  ダウンロード
                </button>
                <button className="bg-gray-200 font-bold rounded-lg py-2 px-4 ml-4 duration-200 hover:bg-gray-300">
                  コピー
                </button>
              </div>
              <p className="text-sm">
                ダウンロードボタンをクリックすると、利用規約及びプライバシーポリシーに同意したものとみなされます。
              </p>
            </div>

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
                        className="flex items-center rounded-full py-1 pl-1 pr-4 cursor-pointer hover:bg-gray-200 duration-200 m-1"
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
                <div className="mb-8">
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
                                className="my-1 ml-1 mr-2 py-1.5 px-2 cursor-pointer duration-200 hover:bg-gray-200 rounded-full"
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

        <div className="my-12 flex justify-center sm:justify-end">
          <a
            href={`http://twitter.com/share?url=https://montanomori.com/illustrations/${illustration.Image.id}&text=${illustration.Image.title}の画像`}
            target="_blank"
            className="pl-1 pr-4 flex items-center border bg-[#1F1F1F] rounded-lg text-white duration-200 cursor-pointer hover:opacity-70"
          >
            <Image
              src="/sns/twitter.png"
              alt="twitterアイコン"
              width={32}
              height={32}
            />
            <span className="text-sm">でシェア</span>
          </a>
          <a
            href={`https://www.facebook.com/share.php?u=https://montanomori.com/illustrations/${illustration.Image.id}`}
            target="_blank"
            className="ml-1 pl-1 pr-4 flex items-center border bg-[#3F50B6] rounded-lg text-white duration-200 cursor-pointer hover:opacity-70"
          >
            <Image
              src="/sns/facebook.png"
              alt="facebookアイコン"
              width={32}
              height={32}
            />
            <span className="text-sm">でシェア</span>
          </a>
          <a
            href={`https://social-plugins.line.me/lineit/share?url=https://montanomori.com/illustrations/${illustration.Image.id}`}
            target="_blank"
            className="ml-1 pl-1 pr-4 flex items-center border bg-[#01C400] rounded-lg text-white duration-200 cursor-pointer hover:opacity-70"
          >
            <Image
              className="p-1"
              src="/sns/line.png"
              alt="lineアイコン"
              width={32}
              height={32}
            />
            <span className="text-sm">でシェア</span>
          </a>
        </div>
      </section>
    </>
  );
};

export default DetailImage;
