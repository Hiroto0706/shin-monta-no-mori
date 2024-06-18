import axios from "axios";
import { FetchIllustrationsAPI } from "@/api/user/illustration";
import { FetchIllustrationsResponse } from "@/types/user/illustration";
import Image from "next/image";
import React from "react";
import {
  FetchCategoriesResponse,
  FetchChildCategoriesResponse,
} from "@/types/user/categories";
import {
  FetchCategoriesAllAPI,
  FetchChildCategoriesAPI,
} from "@/api/user/category";
import TopHeader from "@/components/user/top/topHeader";

const fetchIllustrations = async (): Promise<FetchIllustrationsResponse> => {
  try {
    const response = await axios.get(FetchIllustrationsAPI());
    return response.data;
  } catch (error) {
    console.error(error);
    return { illustrations: [] };
  }
};

const fetchChildCategories =
  async (): Promise<FetchChildCategoriesResponse> => {
    try {
      const response = await axios.get(FetchChildCategoriesAPI());
      return response.data;
    } catch (error) {
      console.error(error);
      return { child_categories: [] };
    }
  };

const fetchCategories = async (): Promise<FetchCategoriesResponse> => {
  try {
    const response = await axios.get(FetchCategoriesAllAPI());
    return response.data;
  } catch (error) {
    console.error("カテゴリの取得に失敗しました", error);
    return { categories: [] };
  }
};

const Home = async () => {
  const fetchIllustrationsRes = await fetchIllustrations();
  const fetchChildCategoriesRes = await fetchChildCategories();
  const fetchCategoriesRes = await fetchCategories();

  const images = [
    {
      title: "おりじなりてぃ",
      description: "もんたの森の絵は\nクセが強くて他と被らない",
      src: "/top-yahhoi.png",
      alt: "おりじなりてぃのある画像",
    },
    {
      title: "ゆるいでざいん",
      description: "もんたの森ではゆーるい\nデザインを大切にしています",
      src: "/top-normal.png",
      alt: "ゆるいでざいんな画像",
    },
    {
      title: "はんようせいがある",
      description: "使える場面の多いイラストを\n描くことを心がけています",
      src: "/top-sumasenn.png",
      alt: "つかえるばめんがおおい画像",
    },
    {
      title: "くせつよもある",
      description: "にっちな場面でしか使えない\nイラストも描いています",
      src: "/top-nikuway.png",
      alt: "くせつよもあるな画像",
    },
  ];

  const others = [
    {
      name: "イラストの依頼",
      description: "もんたがオリジナルイラストを描かせていただきます",
      src: "/top-request.png",
      link: "",
      color: "bg-blue-50",
    },
    {
      name: "お問い合わせ",
      description: "もんたの森に関する「ちょっとわかんない」ことはこちらへ",
      src: "/top-inquiry.png",
      link: "",
      color: "bg-yellow-50",
    },
    {
      name: "フォーラム",
      description: "バグの報告、機能リクエストなどはこちらへ",
      src: "/top-forum.png",
      link: "",
      color: "bg-red-50",
    },
  ];

  const sns = [
    {
      name: "Instagram",
      src: "/sns/instagram.png",
      link: "https://www.instagram.com/yoshida_mandanda/",
    },
    {
      name: "X (Twitter)",
      src: "/sns/twitter.png",
      link: "",
    },
  ];

  return (
    <>
      <TopHeader child_categories={fetchChildCategoriesRes.child_categories} />

      <div className="max-w-[1100px] m-auto mt-40 px-4 md:px-12">
        <section className="mb-40">
          <h2 className="text-2xl font-bold mb-6 text-black">新着イラスト</h2>

          <div className="my-12">
            {fetchIllustrationsRes.illustrations.length > 0 && (
              <div className="grid grid-cols-2 gap-x-4 gap-y-8 md:grid-cols-5">
                {fetchIllustrationsRes.illustrations
                  .slice(0, 10)
                  .map((illustration) => (
                    <a
                      key={illustration.Image.id}
                      href={`/illustrations/${illustration.Image.id}`}
                      className="group cursor-pointer"
                    >
                      <div
                        className="mb-2 border-2 border-gray-200 rounded-xl bg-white relative w-full overflow-hidden"
                        style={{ paddingTop: "100%" }}
                      >
                        <div className="absolute inset-0 m-4">
                          <Image
                            className="group-hover:scale-110 duration-200 absolute top-0 left-0 w-full h-full object-cover"
                            src={illustration.Image.original_src}
                            alt={illustration.Image.title}
                            fill
                          />
                        </div>
                      </div>
                      <span className="group-hover:text-green-600 group-hover:font-bold duration-200">
                        {illustration.Image.title}
                      </span>
                    </a>
                  ))}
              </div>
            )}
          </div>

          <a
            href="/illustrations"
            className="bg-green-600 py-4 text-white font-bold rounded-xl border-2 border-green-600 flex justify-center hover:bg-white hover:text-green-600 duration-200 cursor-pointer"
          >
            もっとみる
          </a>
        </section>

        <section className="mb-40">
          <div className="mb-6">
            <h2 className="text-2xl font-bold text-black mb-2">
              もんたのもりとは
            </h2>
            <p>
              もんたの森は無料で画像を保存・コピーして使うことのできるフリー画像サイトです
            </p>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-x-8 md:gap-x-16">
            {images.map((image, index) => (
              <div key={index} className="mb-8">
                <div
                  className="border-2 border-gray-200 mb-4 rounded-lg relative w-full"
                  style={{ paddingTop: "100%" }}
                >
                  <div className="absolute inset-0 m-4">
                    <Image
                      className="bg-white rounded-lg absolute top-0 left-0 w-full h-full object-cover"
                      src={image.src}
                      alt={image.alt}
                      fill
                    />
                  </div>
                </div>

                <div>
                  <h4 className="font-bold text-lg mb-2 text-black">
                    {image.title}
                  </h4>
                  <p className="text-sm">
                    {image.description.split("\n").map((line, i) => (
                      <React.Fragment key={i}>{line}</React.Fragment>
                    ))}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </section>

        <section className="mb-40">
          <div className="mb-6">
            <h2 className="text-2xl font-bold mb-2 text-black">しようれい</h2>
            <p>もんたの森のイラストは様々な場面で使うことができます</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
            <div className="relative w-full" style={{ paddingTop: "100%" }}>
              <Image
                className="absolute w-full h-full object-contain"
                src="/example-line.svg"
                alt="lineでの使用例"
                fill
              />
            </div>

            <div className="relative w-full" style={{ paddingTop: "100%" }}>
              <Image
                className="absolute w-full h-full object-contain"
                src="/example-slack.svg"
                alt="slackでの使用例"
                fill
              />
            </div>
            <div className="relative w-full" style={{ paddingTop: "100%" }}>
              <Image
                className="absolute w-full h-full object-contain"
                src="/example-slide.svg"
                alt="slideでの使用例"
                fill
              />
            </div>
            <div className="relative w-full" style={{ paddingTop: "100%" }}>
              <Image
                className="absolute w-full h-full object-contain"
                src="/example-icons.svg"
                alt="iconでの使用例"
                fill
              />
            </div>
          </div>
        </section>

        <section className="mb-40">
          <h2 className="text-2xl font-bold mb-6 text-black">そのほか</h2>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-12">
            {others.map((other, i) => (
              <a
                key={i}
                href={other.link}
                className="pointer-cursor px-2 pt-2 pb-4 hover:bg-gray-100 duration-200 rounded-lg"
              >
                <div
                  className={`flex justify-center border rounded-lg border-gray-200 ${other.color}`}
                >
                  <Image
                    className="py-2"
                    src={other.src}
                    alt={other.name}
                    width={120}
                    height={120}
                  />
                </div>
                <div className="my-4">
                  <h4 className="mb-1 text-black text-md font-bold">
                    {other.name}
                  </h4>
                  <p className="text-sm">{other.description}</p>
                </div>
              </a>
            ))}
          </div>

          <div className="w-full md:flex">
            <div className="w-full md:w-1/2 py-4 md:py-0 md:pr-4">
              <a
                href="https://store.line.me/stickershop/author/2887587/ja"
                target="_blank"
                className="cursor-pointer hover:opacity-70 duration-200 w-full"
              >
                <Image
                  className="image"
                  src="/montanomori-line-widget.svg"
                  alt="もんたの森のLINEはこちら"
                  fill
                />
              </a>
            </div>
            <div className="w-full md:w-1/2 py-4 md:pl-4">
              <a
                href="https://www.instagram.com/yoshida_mandanda/"
                target="_blank"
                className="cursor-pointer hover:opacity-70 duration-200 w-full"
              >
                <Image
                  className="image"
                  src="/montanomori-instagram-widget.svg"
                  alt="もんたの森のInstagramはこちら"
                  fill
                />
              </a>
            </div>
          </div>
        </section>

        <section className="border-t border-gray-200 py-12">
          {fetchCategoriesRes.categories.length > 0 && (
            <>
              {fetchCategoriesRes.categories.map((category, index) => (
                <div
                  key={category.ParentCategory.id}
                  className={`${
                    index == fetchCategoriesRes.categories.length - 1
                      ? "mb-0"
                      : "mb-6"
                  }`}
                >
                  <div className="mb-4 flex items-center">
                    <Image
                      src={category.ParentCategory.src}
                      alt={category.ParentCategory.filename.String}
                      width={24}
                      height={24}
                    />
                    <span className="ml-2 font-bold text-black">
                      {category.ParentCategory.name}
                    </span>
                  </div>

                  <div className="flex flex-wrap">
                    {category.ChildCategory.map((cc) => (
                      <a
                        key={cc.id}
                        href={`/illustrations/category/${cc.id}`}
                        className="mr-4 hover:bg-gray-200 duration-200 py-2 px-4 cursor-pointer rounded-full"
                      >
                        # {cc.name}
                      </a>
                    ))}
                  </div>
                </div>
              ))}
            </>
          )}
        </section>
      </div>

      <footer>
        <div className="max-w-[1100px] m-auto px-4 md:px-12">
          <div className="border-t border-gray-200 py-12 flex flex-wrap">
            {sns.map((s, index) => (
              <a
                key={index}
                href={s.link}
                target="_blank"
                className="flex items-center my-2 cursor-pointer block w-40 group"
              >
                <div className="w-8 h-8 relative">
                  <Image
                    className="absolute w-full h-full mx-1 group-hover:opacity-70"
                    alt={s.name}
                    src={s.src}
                    fill
                  />
                </div>
                <span className="text-lg border-b border-gray-600 group-hover:opacity-70 duration-200 ml-2">
                  {s.name}
                </span>
              </a>
            ))}
          </div>
          <div className="py-12 border-t border-gray-200 flex justify-center">
            &copy;もんたの森 2024
          </div>
        </div>
      </footer>
    </>
  );
};

export default Home;
