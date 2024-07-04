import axios from "axios";
import CategoriesSearchForm from "@/components/admin/categories/searchForm";
import { FetchCategoriesResponse } from "@/types/admin/category";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { FetchCategoriesAPI, SearchCategoriesAPI } from "@/api/admin/category";
import { getServerAccessToken } from "@/utils/accessToken/server";
import ListCategoriesTable from "@/components/admin/categories/listTable";
import Pagination from "@/components/common/pagenation";

const fetchCategories = async (
  page: number = 0,
  query: string,
  accessToken: string | undefined
): Promise<FetchCategoriesResponse> => {
  const isSearch = query != "";

  try {
    const response = !isSearch
      ? await axios.get(FetchCategoriesAPI(page), {
          headers: {
            Authorization: SetBearerToken(accessToken),
          },
        })
      : await axios.get(SearchCategoriesAPI(page, query), {
          headers: {
            Authorization: SetBearerToken(accessToken),
          },
        });
    return response.data;
  } catch (error) {
    console.error(error);
    return { categories: [], total_pages: 0, total_count: 0 };
  }
};

export default async function CategoriesListPage({
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
  const categoriesRes = await fetchCategories(page, query, accessToken);
  const totalCount = categoriesRes.total_count;
  const totalPages = categoriesRes.total_pages;

  return (
    <>
      <a
        href="categories/parent/new"
        className="flex items-center bg-white hover:bg-green-600 border-2 border-green-600 text-green-600 hover:text-white rounded-lg py-2 font-bold mb-6 ml-auto w-full lg:w-44 justify-center duration-200"
      >
        + 親カテゴリ追加
      </a>

      <CategoriesSearchForm />

      {categoriesRes.categories.length > 0 ? (
        <ListCategoriesTable categories={categoriesRes.categories} />
      ) : (
        <>カテゴリは見つかりませんでした</>
      )}

      <Pagination
        currentPage={page}
        count={totalCount}
        totalPages={totalPages}
        path="/admin/categories"
        query={query}
      />
    </>
  );
}
