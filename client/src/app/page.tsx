import axios from "axios";
import { FetchIllustrationsAPI } from "@/api/user/illustration";
import { FetchIllustrationsResponse } from "@/types/user/illustration";
import Image from "next/image";
import React from "react";

const fetchIllustrations = async (): Promise<FetchIllustrationsResponse> => {
  try {
    const response = await axios.get(FetchIllustrationsAPI());
    return response.data;
  } catch (error) {
    console.error(error);
    return { illustrations: [] };
  }
};

const Home = async () => {
  const fetchIllustrationsRes = await fetchIllustrations();

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
      description: "もんたにイラストを描いてほしい方はこちらへ",
      src: "/top-request.png",
      link: "",
      color: "bg-blue-50",
    },
    {
      name: "お問い合わせ",
      description: "もんたがオリジナルイラストを描かせていただきます",
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
      link: "",
    },
    {
      name: "X (Twitter)",
      src: "/sns/twitter.png",
      link: "",
    },
  ];

  return (
    <>
      <div className="max-w-[1100px] m-auto mt-40 px-12">
        <section className="mb-40">
          <h2 className="text-2xl font-bold mb-6 text-black">新着イラスト</h2>

          <div className="my-12">
            {fetchIllustrationsRes.illustrations.length > 0 && (
              <div className="grid grid-cols-2 gap-x-4 gap-y-8 md:grid-cols-5">
                {fetchIllustrationsRes.illustrations
                  .slice(0, 10)
                  .map((illustration) => (
                    <div
                      key={illustration.Image.id}
                      className="group cursor-pointer"
                    >
                      <div className="mb-2 p-4 border-2 border-gray-200 rounded-xl bg-white">
                        <Image
                          className="group-hover:scale-110 duration-200 image"
                          src={illustration.Image.original_src}
                          alt={illustration.Image.title}
                          fill
                        />
                      </div>
                      <span className="group-hover:text-green-600 group-hover:font-bold duration-200">
                        {illustration.Image.title}
                      </span>
                    </div>
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
          <h2 className="text-2xl font-bold mb-6 text-black">
            もんたのもりとは
          </h2>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-x-20">
            {images.map((image, index) => (
              <div key={index} className="mb-8">
                <div className="border-2 border-gray-200 mb-4 rounded-lg">
                  <Image
                    className="bg-white rounded-lg p-4 image"
                    src={image.src}
                    alt={image.alt}
                    fill
                  />
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
          <h2 className="text-2xl font-bold mb-6 text-black">しようれい</h2>
          <p>もんたの森のイラストを様々な場面で使おう！</p>

          <div className="flex flex-wrap justify-between">
            <Image
              className="my-8 image"
              src="/example-line.svg"
              alt="lineでの使用例"
              fill
            />
            <Image
              className="my-8 image"
              src="/example-slack.svg"
              alt="slackでの使用例"
              fill
            />
            <Image
              className="my-8 image"
              src="/example-slide.svg"
              alt="slideでの使用例"
              fill
            />
            <Image
              className="my-8 image"
              src="/example-icons.svg"
              alt="iconでの使用例"
              fill
            />
          </div>
        </section>

        <section className="mb-40">
          <h2 className="text-2xl font-bold mb-6 text-black">そのほか</h2>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {others.map((other, i) => (
              <a
                key={i}
                href={other.link}
                className="pointer-cursor px-2 pt-2 pb-4 hover:bg-gray-100 duration-200 rounded-lg mb-4"
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
        </section>
      </div>

      <footer>
        <div className="max-w-[1100px] m-auto px-12">
          <div className="border-t border-gray-200 py-12 flex flex-wrap">
            {sns.map((s, index) => (
              <a
                key={index}
                href={s.link}
                className="flex items-center my-2 cursor-pointer block w-40 group"
              >
                <Image
                  className="image mx-1 group-hover:opacity-70"
                  alt={s.name}
                  src={s.src}
                  fill
                />
                <span className="text-lg border-b border-gray-600 group-hover:opacity-70 duration-200">
                  {s.name}
                </span>
              </a>
            ))}
          </div>
          <div className="py-12 border-t border-gray-200">
            &copy;もんたの森 2024
          </div>
        </div>
      </footer>
    </>
  );
};

export default Home;
