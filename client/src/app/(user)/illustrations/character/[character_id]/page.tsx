import axios from "axios";
import { FetchIllustrationsResponse } from "@/types/user/illustration";
import { FetchIllustrationsByCharacterAPI } from "@/api/user/illustration";
import ListIllustrations from "@/components/user/illustrations/listIllustrations";
import { GetCharacterAPI } from "@/api/user/character";
import { GetCharacterResponse } from "@/types/user/characters";

const fetchIllustrationsByCharacterID = async (
  character_id: number,
  page: number = 0
): Promise<FetchIllustrationsResponse> => {
  try {
    const response = await axios.get(
      FetchIllustrationsByCharacterAPI(character_id, page)
    );
    return response.data;
  } catch (error) {
    console.error(error);
    return { illustrations: [] };
  }
};

const getCharacter = async (
  character_id: number
): Promise<GetCharacterResponse> => {
  try {
    const response = await axios.get(GetCharacterAPI(character_id));
    return response.data;
  } catch (error) {
    console.error(error);
    return { character: null };
  }
};

const FetchIllustrationsByCategoryID = async ({
  params,
}: {
  params: { character_id: number };
}) => {
  const fetchIllustrationsByCategoryIDRes =
    await fetchIllustrationsByCharacterID(params.character_id);
  const getCharacterRes = await getCharacter(params.character_id);

  return (
    <>
      <div className="w-full max-w-[1100px]  2xl:max-w-[1600px] m-auto">
        <h1 className="text-xl font-bold mb-6">
          {getCharacterRes.character != null ? (
            <>{`『${getCharacterRes.character.name}』でキャラクター検索`}</>
          ) : (
            <div>存在しないキャラクターを検索しています</div>
          )}
        </h1>

        {fetchIllustrationsByCategoryIDRes.illustrations.length > 0 &&
        getCharacterRes.character != null ? (
          <ListIllustrations
            initialIllustrations={
              fetchIllustrationsByCategoryIDRes.illustrations
            }
            fetchType={{ characterID: params.character_id }}
          />
        ) : (
          <div>
            イラストが見つかりませんでした
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
