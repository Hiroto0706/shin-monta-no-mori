import axios from "axios";
import CategoriesSearchForm from "@/components/admin/categories/searchForm";
import { Category, FetchCategoriesResponse } from "@/types/category";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { FetchCategoriesAPI, SearchCategoriesAPI } from "@/api/category";
import { getServerAccessToken } from "@/utils/accessToken/server";
import ListCategoriesTable from "@/components/admin/categories/listTable";

export const fetchCategories = async (
  query: string,
  accessToken: string | undefined
): Promise<FetchCategoriesResponse> => {
  const isSearch = query != "";

  try {
    const response = !isSearch
      ? await axios.get(FetchCategoriesAPI(), {
          headers: {
            Authorization: SetBearerToken(accessToken),
          },
        })
      : await axios.get(SearchCategoriesAPI(query), {
          headers: {
            Authorization: SetBearerToken(accessToken),
          },
        });
    return response.data;
  } catch (error) {
    console.error(error);
    return { categories: [] };
  }
};

export default async function CategoriesListPage({
  searchParams,
}: {
  searchParams: {
    q: string;
  };
}) {
  const accessToken = getServerAccessToken();
  const query = searchParams.q ? searchParams.q : "";
  const categoriesRes = await fetchCategories(query, accessToken);

  return (
    <>
      <a
        href="categories/parent/new"
        className="flex items-center bg-white hover:bg-green-600 border-2 border-green-600 text-green-600 hover:text-white rounded-lg py-2 font-bold mb-6 ml-auto w-full lg:w-44 justify-center duration-200"
      >
        + 親カテゴリ追加
      </a>

      <CategoriesSearchForm />

      {categoriesRes.categories && (
        <ListCategoriesTable categories={categoriesRes.categories} />
      )}
    </>
  );
}
