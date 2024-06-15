import axios from "axios";
import { FetchIllustrationsResponse } from "@/types/user/illustration";
import Image from "next/image";
import {
  FetchIllustrationsAPI,
  SearchIllustrationsAPI,
} from "@/api/user/illustration";
import ListIllustrations from "@/components/user/illustrations/listIllustrations";

const fetchIllustrations = async (
  page: number = 0,
  query: string | null
): Promise<FetchIllustrationsResponse> => {
  const isSearch = query;
  try {
    const response = !isSearch
      ? await axios.get(FetchIllustrationsAPI(page))
      : await axios.get(SearchIllustrationsAPI(page, query));

    return response.data;
  } catch (error) {
    console.error(error);
    return { illustrations: [] };
  }
};

const AllIllustrationsPage = async ({
  searchParams,
}: {
  searchParams: {
    p: string;
    q: string;
  };
}) => {
  const page = searchParams.p ? parseInt(searchParams.p, 10) : 0;
  const query = searchParams.q ? searchParams.q : "";
  const fetchIllustrationsRes = await fetchIllustrations(page, query);

  return (
    <>
      <div className="max-w-[1800px] m-auto">
        <h1 className="text-xl font-bold">
          {query != "" ? `『${query}』で検索` : "すべてのイラスト"}
        </h1>

        {fetchIllustrationsRes.illustrations.length > 0 && (
          <ListIllustrations
            illustrations={fetchIllustrationsRes.illustrations}
          />
        )}
      </div>
    </>
  );
};

export default AllIllustrationsPage;
