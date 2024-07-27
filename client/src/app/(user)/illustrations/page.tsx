import axios from "axios";
import { FetchIllustrationsResponse } from "@/types/user/illustration";
import { FetchIllustrationsAPI } from "@/api/user/illustration";
import ListIllustrations from "@/components/user/illustrations/listIllustrations";
import Breadcrumb from "@/components/common/breadCrumb";
import Link from "next/link";

const fetchIllustrations = async (
  page: number = 0
): Promise<FetchIllustrationsResponse> => {
  try {
    const response = await axios.get(FetchIllustrationsAPI(page), {
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

const AllIllustrationsPage = async () => {
  const fetchIllustrationsRes = await fetchIllustrations();

  return (
    <>
      <div className="w-full max-w-[1100px] 2xl:max-w-[1600px] m-auto">
        <Breadcrumb />
        <h1 className="text-xl font-bold mb-6">すべてのイラスト</h1>

        {fetchIllustrationsRes.illustrations.length > 0 ? (
          <ListIllustrations
            initialIllustrations={fetchIllustrationsRes.illustrations}
            fetchType={{}}
          />
        ) : (
          <div>
            イラストが見つかりませんでした
            <Link
              href="/"
              className="text-sm ml-4 underline border-blue-600 text-blue-600 cursor-pointer hover:text-blue-700 duration-200"
            >
              ホームに戻る
            </Link>
          </div>
        )}
      </div>
    </>
  );
};

export default AllIllustrationsPage;
// export const revalidate = 0;
