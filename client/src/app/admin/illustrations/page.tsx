import axios from "axios";
import { FetchIllustrationsResponse } from "@/types/illustration";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { Character } from "@/types/character";
import { Category } from "@/types/category";
import Pagination from "@/components/common/pagenation";
import SearchBox from "@/components/admin/illustrations/searcBox";
import ListTable from "@/components/admin/illustrations/listTable";
import { FetchIllustrationsAPI } from "@/api/illustration";
import { getServerAccessToken } from "@/utils/accessToken/server";

const fetchIllustrations = async (
  page: number = 0,
  accessToken: string | undefined
): Promise<FetchIllustrationsResponse> => {
  try {
    const response = await axios.get(FetchIllustrationsAPI(page), {
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

export const fetchCharacters = async (
  accessToken: string | undefined
): Promise<Character[]> => {
  try {
    const response = await axios.get(
      process.env.NEXT_PUBLIC_BASE_API + "admin/characters/list",
      {
        headers: {
          Authorization: SetBearerToken(accessToken),
        },
      }
    );
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
    const response = await axios.get(
      process.env.NEXT_PUBLIC_BASE_API + "admin/categories/list",
      {
        headers: {
          Authorization: SetBearerToken(accessToken),
        },
      }
    );
    return response.data.categories;
  } catch (error) {
    console.error(error);
    return [];
  }
};

export default async function IllustrationsListPage({
  searchParams,
}: {
  searchParams: { p: string };
}) {
  const accessToken = getServerAccessToken();
  const page = searchParams.p ? parseInt(searchParams.p, 10) : 0;
  const illustrations: FetchIllustrationsResponse = await fetchIllustrations(
    page,
    accessToken
  );
  const totalCount = illustrations.total_count;
  const totalPages = illustrations.total_pages;
  const characters: Character[] = await fetchCharacters(accessToken);
  const categories: Category[] = await fetchCategories(accessToken);

  return (
    <>
      <a
        href="illustrations/new"
        className="flex items-center bg-white hover:bg-green-600 border-2 border-green-600 text-green-600 hover:text-white rounded-lg py-2 font-bold mb-6 ml-auto w-full lg:w-36 justify-center duration-200"
      >
        + イラスト追加
      </a>

      <SearchBox characters={characters} categories={categories} />

      <ListTable illustrations={illustrations} />

      <Pagination
        currentPage={page}
        count={totalCount}
        totalPages={totalPages}
        path="/admin/illustrations"
      />
    </>
  );
}
