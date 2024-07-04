import axios from "axios";
import CreateChildCategory from "@/components/admin/categories/childCategory/createForm";
import { getServerAccessToken } from "@/utils/accessToken/server";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { FetchAllCategoriesAPI } from "@/api/admin/category";
import { FetchCategoriesResponse } from "@/types/admin/category";

const fetchCategories = async (
  accessToken: string | undefined
): Promise<FetchCategoriesResponse> => {
  try {
    const response = await axios.get(FetchAllCategoriesAPI(), {
      headers: {
        Authorization: SetBearerToken(accessToken),
      },
    });

    return response.data;
  } catch (error) {
    console.error(error);
    return { categories: [], total_count: 0, total_pages: 0 };
  }
};

const CreateChildCategoryPage = async ({
  searchParams,
}: {
  searchParams: {
    parent_id: number;
  };
}) => {
  const accessToken = getServerAccessToken();
  const fetchCategoriesRes = await fetchCategories(accessToken);

  return (
    <>
      {fetchCategoriesRes.categories.length > 0 && (
        <CreateChildCategory
          parentID={searchParams.parent_id}
          categories={fetchCategoriesRes.categories}
          accessToken={accessToken}
        />
      )}
    </>
  );
};

export default CreateChildCategoryPage;
