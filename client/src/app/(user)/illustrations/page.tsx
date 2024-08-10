import axios from "axios";
import { FetchIllustrationsResponse } from "@/types/user/illustration";
import { FetchIllustrationsAPI } from "@/api/user/illustration";
import IllustrationListTemplate from "@/components/user/illustrations/list/illustrationListTemplate";

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
        <IllustrationListTemplate
          illustrations={fetchIllustrationsRes.illustrations}
        />
      </div>
    </>
  );
};

export default AllIllustrationsPage;
// export const revalidate = 0;
