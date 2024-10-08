import axios from "axios";
import { FetchIllustrationsResponse } from "@/types/user/illustration";
import {
  FetchIllustrationsAPI,
  SearchIllustrationsAPI,
} from "@/api/user/illustration";
import IllustrationListBySearchTemplate from "@/components/user/illustrations/list/illustrationListBySearchTemplate";

const fetchIllustrations = async (
  query: string,
  page: number = 0
): Promise<FetchIllustrationsResponse> => {
  try {
    let response;
    if (query != "") {
      response = await axios.get(SearchIllustrationsAPI(page, query), {
        headers: {
          "Cache-Control": "no-store",
          "CDN-Cache-Control": "no-store",
          "Vercel-CDN-Cache-Control": "no-store",
        },
      });
    } else {
      response = await axios.get(FetchIllustrationsAPI(page), {
        headers: {
          "Cache-Control": "no-store",
          "CDN-Cache-Control": "no-store",
          "Vercel-CDN-Cache-Control": "no-store",
        },
      });
    }
    return response.data;
  } catch (error) {
    console.error(error);
    return { illustrations: [] };
  }
};

const SearchIllustrationsPage = async ({
  params,
}: {
  params: { query: string };
}) => {
  const query = params.query ? decodeURIComponent(params.query) : "";
  const fetchIllustrationsRes = await fetchIllustrations(query);

  return (
    <>
      <div className="w-full max-w-[1100px] 2xl:max-w-[1600px] m-auto">
        <IllustrationListBySearchTemplate
          illustrations={fetchIllustrationsRes.illustrations}
          query={query}
        />
      </div>
    </>
  );
};

export default SearchIllustrationsPage;
// export const revalidate = 0;
