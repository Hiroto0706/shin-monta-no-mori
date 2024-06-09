import axios from "axios";
import CreateChildCategory from "@/components/admin/categories/childCategory/createForm";
import { getServerAccessToken } from "@/utils/accessToken/server";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { FetchCategoriesAPI } from "@/api/admin/category";
import { FetchCategoriesResponse } from "@/types/admin/category";

export const fetchCategories = async (
  accessToken: string | undefined
): Promise<FetchCategoriesResponse> => {
  try {
    const response = await axios.get(FetchCategoriesAPI(), {
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
