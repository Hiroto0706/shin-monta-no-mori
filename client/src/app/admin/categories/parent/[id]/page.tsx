import axios from "axios";
import { getServerAccessToken } from "@/utils/accessToken/server";
import { GetParentCategoryResponse } from "@/types/category";
import { GetParentCategoryAPI } from "@/api/category";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import EditParentCategory from "@/components/admin/categories/parentCategory/editForm";

const getParentCategory = async (
  id: number,
  accessToken: string | undefined
): Promise<GetParentCategoryResponse> => {
  try {
    const response = await axios.get(GetParentCategoryAPI(id), {
      headers: {
        Authorization: SetBearerToken(accessToken),
      },
    });
    return response.data;
  } catch (error) {
    console.error(error);
    return { parent_category: null };
  }
};

const EditParentCategoryPage = async ({
  params,
}: {
  params: { id: number };
}) => {
  const accessToken = getServerAccessToken();
  const parentCategoryRes = await getParentCategory(params.id, accessToken);

  return (
    <>
      {parentCategoryRes.parent_category && (
        <EditParentCategory
          id={params.id}
          parentCategory={parentCategoryRes.parent_category}
          accessToken={accessToken}
        />
      )}
    </>
  );
};

export default EditParentCategoryPage;
