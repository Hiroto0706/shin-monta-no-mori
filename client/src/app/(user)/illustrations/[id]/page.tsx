import axios from "axios";
import Image from "next/image";
import { GetIllustrationAPI } from "@/api/user/illustration";
import { GetIllustrationResponse } from "@/types/user/illustration";
import DetailImage from "@/components/user/illustrations/detail/detailImage";
import RandomIllustrations from "@/components/user/illustrations/detail/RandomIllutrations";
import Breadcrumb from "@/components/common/breadCrumb";
import Link from "next/link";

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
      <div className="w-full max-w-[1100px] m-auto">
        {getIllustrationRes.illustration != null ? (
          <>
            <Breadcrumb
              customString={getIllustrationRes.illustration?.Image.title}
            />

            <DetailImage illustration={getIllustrationRes.illustration} />

            <section className="mb-20">
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

            <section className="mb-20">
              <h2 className="text-xl font-bold text-black">そのほかイラスト</h2>
              <RandomIllustrations exclusion_id={params.id} />
            </section>
          </>
        ) : (
          <div className="my-8">
            <div className="mb-4">お探しのイラストは見つかりませんでした</div>
            <Link
              href="/"
              className="mt-4 underline border-blue-600 text-blue-600 cursor-pointer hover:text-blue-700 duration-200"
            >
              ホームに戻る
            </Link>
          </div>
        )}
      </div>
    </>
  );
};

export default IllustrationDetailPage;
