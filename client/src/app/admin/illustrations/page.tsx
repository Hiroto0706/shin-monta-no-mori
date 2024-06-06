import axios from "axios";
import { FetchIllustrationsResponse } from "@/types/illustration";
import { Character } from "@/types/character";
import { Category } from "@/types/category";
import Pagination from "@/components/common/pagenation";
import SearchForm from "@/components/admin/illustrations/searchForm";
import ListTable from "@/components/admin/illustrations/listTable";
import {
  FetchIllustrationsAPI,
  SearchIllustrationsAPI,
} from "@/api/illustration";
import { FetchAllCharactersAPI } from "@/api/character";
import { FetchCategoriesAPI } from "@/api/category";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { getServerAccessToken } from "@/utils/accessToken/server";

const fetchIllustrations = async (
  page: number = 0,
  accessToken: string | undefined,
  query: string | null,
  characters: string | null,
  categories: string | null
): Promise<FetchIllustrationsResponse> => {
  const isSearch = query || characters || categories;
  try {
    const response = isSearch
      ? await axios.get(
          SearchIllustrationsAPI(page, query, characters, categories),
          {
            headers: {
              Authorization: SetBearerToken(accessToken),
            },
          }
        )
      : await axios.get(FetchIllustrationsAPI(page), {
          headers: {
            Authorization: SetBearerToken(accessToken),
          },
        });
    return response.data;
  } catch (error) {
    console.error(error);
    return { illustrations: [], total_pages: 0, total_count: 0 };
  }
};

export const fetchAllCharacters = async (
  accessToken: string | undefined
): Promise<Character[]> => {
  try {
    const response = await axios.get(FetchAllCharactersAPI(), {
      headers: {
        Authorization: SetBearerToken(accessToken),
      },
    });
    return response.data.characters;
  } catch (error) {
    console.error(error);
    return [];
  }
};

export const fetchCategories = async (
  accessToken: string | undefined
): Promise<Category[]> => {
  try {
    const response = await axios.get(FetchCategoriesAPI(), {
      headers: {
        Authorization: SetBearerToken(accessToken),
      },
    });
    return response.data.categories;
  } catch (error) {
    console.error(error);
    return [];
  }
};

export default async function IllustrationsListPage({
  searchParams,
}: {
  searchParams: {
    p: string;
    q: string;
    characters: string;
    categories: string;
  };
}) {
  const accessToken = getServerAccessToken();
  const page = searchParams.p ? parseInt(searchParams.p, 10) : 0;
  const query = searchParams.q ? searchParams.q : "";
  const charactersParams = searchParams.characters
    ? searchParams.characters
    : "";
  const categoriesParams = searchParams.categories
    ? searchParams.categories
    : "";
  const illustrations: FetchIllustrationsResponse = await fetchIllustrations(
    page,
    accessToken,
    query,
    charactersParams,
    categoriesParams
  );
  const totalCount = illustrations.total_count;
  const totalPages = illustrations.total_pages;
  const characters: Character[] = await fetchAllCharacters(accessToken);
  const categories: Category[] = await fetchCategories(accessToken);

  return (
    <>
      <a
        href="illustrations/new"
        className="flex items-center bg-white hover:bg-green-600 border-2 border-green-600 text-green-600 hover:text-white rounded-lg py-2 font-bold mb-6 ml-auto w-full lg:w-36 justify-center duration-200"
      >
        + イラスト追加
      </a>

      <SearchForm characters={characters} categories={categories} />

      <ListTable illustrations={illustrations} />

      <Pagination
        currentPage={page}
        count={totalCount}
        totalPages={totalPages}
        path="/admin/illustrations"
        query={query}
        charactersParams={charactersParams}
        categoriesParams={categoriesParams}
      />
    </>
  );
}
