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
import { FetchCharactersResponse } from "@/types/user/characters";
import { FetchAllCharactersAPI } from "@/api/user/character";
import IllustrationCard from "@/components/user/illustrations/illustrationCard";
import Link from "next/link";
import Footer from "@/components/user/common/footer/footer";

export const dynamicParams = false;

const fetchIllustrations = async (): Promise<FetchIllustrationsResponse> => {
  try {
    const response = await axios.get(FetchIllustrationsAPI(), {
      headers: {
        "Cache-Control": "no-store",
        "CDN-Cache-Control": "no-store",
        "Vercel-CDN-Cache-Control": "no-store",
      },
    });
    return response.data;
  } catch (error) {
    console.error(error);
    return { illustrations: [] };
  }
};

const fetchChildCategories =
  async (): Promise<FetchChildCategoriesResponse> => {
    try {
      const response = await axios.get(FetchChildCategoriesAPI(), {
        headers: {
          "Cache-Control": "no-store",
          "CDN-Cache-Control": "no-store",
          "Vercel-CDN-Cache-Control": "no-store",
        },
      });
      return response.data;
    } catch (error) {
      console.error("キャラクターの取得に失敗しました", error);
      return { child_categories: [] };
    }
  };

const fetchCategories = async (): Promise<FetchCategoriesResponse> => {
  try {
    const response = await axios.get(FetchCategoriesAllAPI(), {
      headers: {
        "Cache-Control": "no-store",
        "CDN-Cache-Control": "no-store",
        "Vercel-CDN-Cache-Control": "no-store",
      },
    });
    return response.data;
  } catch (error) {
    console.error("カテゴリの取得に失敗しました", error);
    return { categories: [] };
  }
};

const fetchCharacters = async (): Promise<FetchCharactersResponse> => {
  try {
    const response = await axios.get(FetchAllCharactersAPI(), {
      headers: {
        "Cache-Control": "no-store",
        "CDN-Cache-Control": "no-store",
        "Vercel-CDN-Cache-Control": "no-store",
      },
    });
    return response.data;
  } catch (error) {
    console.error("キャラクターの取得に失敗しました", error);
    return { characters: [] };
  }
};

const Home = async () => {
  const fetchIllustrationsRes = await fetchIllustrations();
  const fetchChildCategoriesRes = await fetchChildCategories();
  const fetchCategoriesRes = await fetchCategories();
  const fetchCharactersRes = await fetchCharacters();

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
      link: "https://forms.gle/gfvuc6GwiURxJNR68",
      color: "bg-blue-50",
    },
    {
      name: "お問い合わせ",
      description: "もんたの森に関する「ちょっとわかんない」ことはこちらへ",
      src: "/top-inquiry.png",
      link: "https://forms.gle/THqHAigzTZa7J9D28",
      color: "bg-yellow-50",
    },
    {
      name: "フォーラム",
      description: "バグの報告、機能リクエストなどはこちらへ",
      src: "/top-forum.png",
      link: "https://forms.gle/i4Fp9Xoeq4fkMEc88",
      color: "bg-red-50",
    },
  ];

  const sns = [
    {
      name: "Instagram",
      src: "/sns/instagram.png",
      link: "https://www.instagram.com/monta_no_mori/",
    },
    {
      name: "X (Twitter)",
      src: "/sns/twitter.png",
      link: "https://x.com/monta_no_mori",
    },
  ];

  return (
    <>
      <TopHeader
        child_categories={fetchChildCategoriesRes.child_categories}
        characters={fetchCharactersRes.characters}
        categories={fetchCategoriesRes.categories}
      />

      <div className="max-w-[1100px] m-auto mt-24 px-4 md:px-12">
        <section className="mb-40">
          <h2 className="text-2xl font-bold mb-6 text-black">新着イラスト</h2>

          <div className="my-12">
            {fetchIllustrationsRes.illustrations.length > 0 && (
              <div className="grid grid-cols-2 gap-x-4 gap-y-8 md:grid-cols-5">
                {fetchIllustrationsRes.illustrations
                  .slice(0, 10)
                  .map((illustration) => (
                    <IllustrationCard
                      key={illustration.Image.id}
                      illustration={illustration}
                    />
                  ))}
              </div>
            )}
          </div>

          <Link
            href="/illustrations"
            className="bg-green-600 py-4 text-white font-bold rounded-xl border-2 border-green-600 flex justify-center hover:bg-white hover:text-green-600 duration-200 cursor-pointer"
          >
            もっとみる
          </Link>
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
              <Link
                key={i}
                href={other.link}
                target="_blank"
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
              </Link>
            ))}
          </div>

          <div className="lg:flex">
            <div className="w-full lg:w-1/2 h-28 md:h-40 lg:h-32 flex items-center justify-center my-2 lg:my-0 lg:mr-4">
              <div className="w-full max-w-[450px] md:max-w-[550px]">
                <Link
                  href="https://store.line.me/stickershop/author/2887587/ja"
                  target="_blank"
                  className="cursor-pointer hover:opacity-70 duration-200"
                >
                  <Image
                    className="image"
                    src="/montanomori-line-widget.svg"
                    alt="もんたの森のLINEはこちら"
                    fill
                  />
                </Link>
              </div>
            </div>
            <div className="w-full lg:w-1/2 h-28 md:h-40 lg:h-32 flex items-center justify-center my-2 lg:my-0 lg:ml-4">
              <div className="w-full max-w-[450px] md:max-w-[550px]">
                <Link
                  href="https://www.instagram.com/monta_no_mori/"
                  target="_blank"
                  className="cursor-pointer hover:opacity-70 duration-200"
                >
                  <Image
                    className="image"
                    src="/montanomori-instagram-widget.svg"
                    alt="もんたの森のInstagramはこちら"
                    fill
                  />
                </Link>
              </div>
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
                      <Link
                        key={cc.id}
                        href={`/illustrations/category/${cc.id}`}
                        className="mr-4 hover:bg-gray-200 duration-200 py-2 px-4 cursor-pointer rounded-full"
                      >
                        # {cc.name}
                      </Link>
                    ))}
                  </div>
                </div>
              ))}
            </>
          )}
        </section>

        <section className="border-t border-gray-200 py-12">
          {fetchCharactersRes.characters.length > 0 && (
            <>
              <div className="flex flex-wrap">
                {fetchCharactersRes.characters.map((character, index) => (
                  <div key={character.id} className="my-2 flex items-center">
                    <Link
                      key={character.id}
                      href={`/illustrations/category/${character.id}`}
                      className="mr-4 hover:bg-gray-200 duration-200 py-2 px-2 cursor-pointer rounded-full flex items-center"
                    >
                      <Image
                        className="border border-gray-200 rounded-full bg-white shadow"
                        src={character.src}
                        alt={character.filename.String}
                        width={36}
                        height={36}
                      />
                      <span className="ml-4">{character.name}</span>
                    </Link>
                  </div>
                ))}
              </div>
            </>
          )}
        </section>
      </div>

      <Footer sns={sns} />
    </>
  );
};

export default Home;
// export const revalidate = 0;
