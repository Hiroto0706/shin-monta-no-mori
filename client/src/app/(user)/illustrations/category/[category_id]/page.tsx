import axios from "axios";
import { FetchIllustrationsResponse } from "@/types/user/illustration";
import { FetchIllustrationsByCategoryAPI } from "@/api/user/illustration";
import ListIllustrations from "@/components/user/illustrations/listIllustrations";
import { GetChildCategoryResponse } from "@/types/user/categories";
import { GetChildCategoryAPI } from "@/api/user/category";

const fetchIllustrationsByCategoryID = async (
  category_id: number
): Promise<FetchIllustrationsResponse> => {
  try {
    const response = await axios.get(
      FetchIllustrationsByCategoryAPI(category_id)
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
    const response = await axios.get(GetChildCategoryAPI(category_id));
    console.log(response);
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
        <h1 className="text-xl font-bold mb-6">
          {getChildCategoryRes.child_category != null ? (
            <>
              {`『${getChildCategoryRes.child_category?.name}』でカテゴリ検索`}
            </>
          ) : (
            <div>存在しないカテゴリを検索しています</div>
          )}
        </h1>

        {fetchIllustrationsByCategoryIDRes.illustrations.length > 0 &&
        getChildCategoryRes.child_category != null ? (
          <ListIllustrations
            illustrations={fetchIllustrationsByCategoryIDRes.illustrations}
          />
        ) : (
          <div>
            イラストが見つかりませんでした{" "}
            <a
              href="/"
              className="text-sm ml-4 underline border-blue-600 text-blue-600 cursor-pointer hover:opacity-70 duration-200"
            >
              ホームに戻る
            </a>
          </div>
        )}
      </div>
    </>
  );
};

export default FetchIllustrationsByCategoryID;
