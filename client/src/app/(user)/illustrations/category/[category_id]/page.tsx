import axios from "axios";
import { FetchIllustrationsResponse } from "@/types/user/illustration";
import { FetchIllustrationsByCategoryAPI } from "@/api/user/illustration";
import { GetChildCategoryResponse } from "@/types/user/categories";
import { GetChildCategoryAPI } from "@/api/user/category";
import IllustrationListByCategoryTemplate from "@/components/user/illustrations/list/illustrationListByCategoryTemplate";

const fetchIllustrationsByCategoryID = async (
  category_id: number,
  page: number = 0
): Promise<FetchIllustrationsResponse> => {
  try {
    const response = await axios.get(
      FetchIllustrationsByCategoryAPI(category_id, page),
      {
        headers: {
          "Cache-Control": "no-store",
          "CDN-Cache-Control": "no-store",
          "Vercel-CDN-Cache-Control": "no-store",
        },
      }
    );
    return response.data;
  } catch (error) {
    console.error(error);
    return { illustrations: [] };
  }
};

const getChildCategory = async (
  category_id: number
): Promise<GetChildCategoryResponse> => {
  try {
    const response = await axios.get(GetChildCategoryAPI(category_id), {
      headers: {
        "Cache-Control": "no-store",
        "CDN-Cache-Control": "no-store",
        "Vercel-CDN-Cache-Control": "no-store",
      },
    });
    return response.data;
  } catch (error) {
    console.error(error);
    return { child_category: null };
  }
};

const FetchIllustrationsByCategoryID = async ({
  params,
}: {
  params: { category_id: number };
}) => {
  const fetchIllustrationsByCategoryIDRes =
    await fetchIllustrationsByCategoryID(params.category_id);
  const getChildCategoryRes = await getChildCategory(params.category_id);

  return (
    <>
      <div className="w-full max-w-[1100px]  2xl:max-w-[1600px] m-auto">
        <IllustrationListByCategoryTemplate
          illustrations={fetchIllustrationsByCategoryIDRes.illustrations}
          categoryID={params.category_id}
          childCategory={getChildCategoryRes.child_category}
        />
      </div>
    </>
  );
};

export default FetchIllustrationsByCategoryID;
// export const revalidate = 0;
