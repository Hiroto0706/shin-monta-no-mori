import axios from "axios";
import { Character, FetchCharactersResponse } from "@/types/character";
import SearchForm from "@/components/admin/characters/searchForm";
import { FetchCharactersAPI } from "@/api/character";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { getServerAccessToken } from "@/utils/accessToken/server";
import ListCharactersTable from "@/components/admin/characters/listTable";
import Pagination from "@/components/common/pagenation";

export const fetchCharacters = async (
  page: number,
  accessToken: string | undefined
): Promise<FetchCharactersResponse> => {
  try {
    const response = await axios.get(FetchCharactersAPI(page), {
      headers: {
        Authorization: SetBearerToken(accessToken),
      },
    });
    console.log(response.data);
    return response.data;
  } catch (error) {
    console.error(error);
    return { characters: [], total_pages: 0, total_count: 0 };
  }
};

export default async function IllustrationsListPage({
  searchParams,
}: {
  searchParams: {
    p: string;
    q: string;
  };
}) {
  const accessToken = getServerAccessToken();
  const page = searchParams.p ? parseInt(searchParams.p, 10) : 0;
  const query = searchParams.q ? searchParams.q : "";
  const characters = await fetchCharacters(page, accessToken);
  const totalCount = characters.total_count;
  const totalPages = characters.total_pages;

  return (
    <>
      <a
        href="characters/new"
        className="flex items-center bg-white hover:bg-green-600 border-2 border-green-600 text-green-600 hover:text-white rounded-lg py-2 font-bold mb-6 ml-auto w-full lg:w-44 justify-center duration-200"
      >
        + キャラクター追加
      </a>

      <SearchForm />

      <ListCharactersTable characters={characters} />

      <Pagination
        currentPage={page}
        count={totalCount}
        totalPages={totalPages}
        path="/admin/characters"
        query={query}
      />
    </>
  );
}
