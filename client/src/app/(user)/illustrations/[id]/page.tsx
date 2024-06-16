import axios from "axios";
import Image from "next/image";
import { GetIllustrationAPI } from "@/api/user/illustration";
import { GetIllustrationResponse } from "@/types/user/illustration";

const getIllustration = async (
  id: number
): Promise<GetIllustrationResponse> => {
  try {
    const response = await axios.get(GetIllustrationAPI(id));

    return response.data;
  } catch (error) {
    console.error(error);
    return { illustration: null };
  }
};

const IllustrationDetailPage = async ({
  params,
}: {
  params: { id: number };
}) => {
  const getIllustrationRes = await getIllustration(params.id);

  return (
    <>
      {getIllustrationRes.illustration && (
        <section className="mb-20">
          <div className="flex flex-col lg:flex-row">
            <div className="w-full max-w-[500px] m-auto lg:m-0 lg:w-2/5 lg:max-w-[400px]">
              <div
                className="relative w-full border-2 border-gray-200 rounded-xl"
                style={{ paddingTop: "100%" }}
              >
                <div className="absolute inset-0 m-8">
                  <Image
                    className="absolute inset-0 object-cover w-full h-full"
                    src={getIllustrationRes.illustration.Image.original_src}
                    alt={
                      getIllustrationRes.illustration.Image.original_filename
                    }
                    fill
                  />
                </div>
              </div>

              <div className="mt-2 flex justify-between items-center p-1 rounded bg-gray-100 rounded-lg">
                <div className="p-2 bg-green-600 rounded-lg text-white w-1/2 flex justify-center font-bold text-lg">
                  オリジナル
                </div>
                <div className="p-2 rounded-full w-1/2 flex justify-center font-bold text-lg">
                  シンプル
                </div>
              </div>
            </div>

            <div className="w-full lg:w-3/5 mt-8 lg:mt-0 lg:ml-8">
              <h2 className="text-2xl font-bold mb-4 text-black">
                {getIllustrationRes.illustration.Image.title}
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

              {getIllustrationRes.illustration.Characters.length > 0 && (
                <>
                  <div className="mb-8">
                    <h3 className="text-lg font-bold mb-2 text-black">
                      キャラクター
                    </h3>
                    <div className="flex">
                      {getIllustrationRes.illustration.Characters.map(
                        (c, index) => (
                          <a
                            key={c.Character.id}
                            href={`/illustrations/character/${c.Character.id}`}
                            className={`flex items-center rounded-full py-1 pl-1 pr-4 cursor-pointer hover:bg-gray-200 duration-200 ${
                              index == 0 ? "ml-0" : "ml-2"
                            }`}
                          >
                            <Image
                              className="border border-gray-200 rounded-full bg-white"
                              src={c.Character.src}
                              alt={c.Character.name}
                              width={40}
                              height={40}
                            />
                            <span className="ml-2">{c.Character.name}</span>
                          </a>
                        )
                      )}
                    </div>
                  </div>
                </>
              )}

              {getIllustrationRes.illustration.Categories.length > 0 && (
                <>
                  <div className="mb-8">
                    <h3 className="text-lg font-bold mb-2 text-black">
                      カテゴリ
                    </h3>
                    <div>
                      <>
                        <div className="flex mb-2">
                          {getIllustrationRes.illustration.Categories.map(
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
                        </div>

                        <div className="flex flex-wrap">
                          {getIllustrationRes.illustration.Categories.map(
                            (category) => (
                              <>
                                {category.ChildCategory.map((cc) => (
                                  <a
                                    key={cc.id}
                                    href={`/illustrations/category/${cc.id}`}
                                    className="mr-2 py-1 px-2 cursor-pointer duration-200 hover:bg-gray-200 rounded-full"
                                  >
                                    # {cc.name}
                                  </a>
                                ))}
                              </>
                            )
                          )}
                        </div>
                      </>
                    </div>
                  </div>
                </>
              )}
            </div>
          </div>
        </section>
      )}

      <section className="mb-20">
        <div className="w-full h-32 flex items-center">
          <div className="w-full md:w-3/4 lg:w-1/2 lg:max-w-[550px] lg:min-w-[400px]">
            <a
              href="https://store.line.me/stickershop/author/2887587/ja"
              className="cursor-pointer hover:opacity-70 duration-200"
            >
              <Image
                className="image"
                src="/montanomori-line-widget.svg"
                alt="もんたの森のLINEはこちら"
                fill
              />
            </a>
          </div>
        </div>
      </section>

      <section className="mb-20">
        <h3>関連イラスト</h3>
      </section>
    </>
  );
};

export default IllustrationDetailPage;
