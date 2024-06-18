import axios from "axios";
import Image from "next/image";
import { GetIllustrationAPI } from "@/api/user/illustration";
import { GetIllustrationResponse } from "@/types/user/illustration";
import DetailImage from "@/components/user/illustrations/detail/detailImage";

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
        {getIllustrationRes.illustration && (
          <DetailImage illustration={getIllustrationRes.illustration} />
        )}

        <section className="mb-20">
          <div className="lg:flex">
            <div className="w-full lg:w-1/2 h-28 md:h-40 lg:h-32 flex items-center justify-center my-2 lg:my-0 lg:mr-4">
              <div className="w-full max-w-[450px] md:max-w-[550px]">
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
            <div className="w-full lg:w-1/2 h-28 md:h-40 lg:h-32 flex items-center justify-center my-2 lg:my-0 lg:ml-4">
              <div className="w-full max-w-[450px] md:max-w-[550px]">
                <a
                  href="https://store.line.me/stickershop/author/2887587/ja"
                  className="cursor-pointer hover:opacity-70 duration-200"
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
          </div>
        </section>

        <section className="mb-20">
          <h3>関連イラスト</h3>
        </section>
      </div>
    </>
  );
};

export default IllustrationDetailPage;
