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
      <div className="w-full max-w-[1100px]  2xl:max-w-[1600px] m-auto">
        {getIllustrationRes.illustration && (
          <DetailImage illustration={getIllustrationRes.illustration} />
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
      </div>
    </>
  );
};

export default IllustrationDetailPage;
