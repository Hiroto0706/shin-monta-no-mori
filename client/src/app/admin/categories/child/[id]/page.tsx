import axios from "axios";
import { getServerAccessToken } from "@/utils/accessToken/server";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { FetchCategoriesAPI, GetChildCategoryAPI } from "@/api/category";
import {
  FetchCategoriesResponse,
  GetChildCategoryResponse,
} from "@/types/category";
import EditChildCategory from "@/components/admin/categories/childCategory/editForm";

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

export const getChildCategory = async (
  id: number,
  accessToken: string | undefined
): Promise<GetChildCategoryResponse> => {
  try {
    const response = await axios.get(GetChildCategoryAPI(id), {
      headers: {
        Authorization: SetBearerToken(accessToken),
      },
    });

    return response.data;
  } catch (error) {
    console.error(error);
    return { child_category: null };
  }
};

const EditChildCategoryPage = async ({
  params,
}: {
  params: { id: number };
}) => {
  const accessToken = getServerAccessToken();
  const fetchCategoriesRes = await fetchCategories(accessToken);
  const getChildCategoryRes = await getChildCategory(params.id, accessToken);

  return (
    <>
      {fetchCategoriesRes.categories.length > 0 &&
        getChildCategoryRes.child_category && (
          <EditChildCategory
            categories={fetchCategoriesRes.categories}
            childCategory={getChildCategoryRes.child_category}
            accessToken={accessToken}
          />
        )}
    </>
  );
};

export default EditChildCategoryPage;
