import axios from "axios";
import { FetchCharactersResponse } from "@/types/admin/character";
import { FetchCharactersAPI, SearchCharactersAPI } from "@/api/admin/character";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { getServerAccessToken } from "@/utils/accessToken/server";
import ListCharactersTable from "@/components/admin/characters/listTable";
import Pagination from "@/components/common/pagenation";
import CharactersSearchForm from "@/components/admin/characters/searchForm";

const fetchCharacters = async (
  page: number = 0,
  query: string,
  accessToken: string | undefined
): Promise<FetchCharactersResponse> => {
  const isSearch = query != "";

  try {
    const response = !isSearch
      ? await axios.get(FetchCharactersAPI(page), {
          headers: {
            Authorization: SetBearerToken(accessToken),
          },
        })
      : await axios.get(SearchCharactersAPI(page, query), {
          headers: {
            Authorization: SetBearerToken(accessToken),
          },
        });
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
  const characters = await fetchCharacters(page, query, accessToken);
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

      <CharactersSearchForm />

      {characters.characters.length > 0 ? (
        <ListCharactersTable characters={characters} />
      ) : (
        <p className="mb-6">キャラクターは見つかりませんでした</p>
      )}

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
