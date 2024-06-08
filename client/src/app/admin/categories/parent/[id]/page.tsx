import axios from "axios";
import { getServerAccessToken } from "@/utils/accessToken/server";
import { GetCategoryResponse } from "@/types/category";
import { GetCategoryAPI } from "@/api/category";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import EditParentCategory from "@/components/admin/categories/parentCategory/editForm";

const getParentCategory = async (
  id: number,
  accessToken: string | undefined
): Promise<GetCategoryResponse> => {
  try {
    const response = await axios.get(GetCategoryAPI(id), {
      headers: {
        Authorization: SetBearerToken(accessToken),
      },
    });
    return response.data;
  } catch (error) {
    console.error(error);
    return { category: null };
  }
};

const EditParentCategoryPage = async ({
  params,
}: {
  params: { id: number };
}) => {
  const accessToken = getServerAccessToken();
  const categoryRes = await getParentCategory(params.id, accessToken);

  return (
    <>
      {categoryRes.category && (
        <EditParentCategory
          id={params.id}
          parentCategory={categoryRes.category.ParentCategory}
          accessToken={accessToken}
        />
      )}
    </>
  );
};

export default EditParentCategoryPage;
