import axios from "axios";
import { FetchIllustrationsResponse } from "@/types/user/illustration";
import { FetchIllustrationsByCharacterAPI } from "@/api/user/illustration";
import { GetCharacterAPI } from "@/api/user/character";
import { GetCharacterResponse } from "@/types/user/characters";
import IllustrationListByCharacterTemplate from "@/components/user/illustrations/list/illustrationListByCharacterTemplate";

const fetchIllustrationsByCharacterID = async (
  character_id: number,
  page: number = 0
): Promise<FetchIllustrationsResponse> => {
  try {
    const response = await axios.get(
      FetchIllustrationsByCharacterAPI(character_id, page),
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

const getCharacter = async (
  character_id: number
): Promise<GetCharacterResponse> => {
  try {
    const response = await axios.get(GetCharacterAPI(character_id), {
      headers: {
        "Cache-Control": "no-store",
        "CDN-Cache-Control": "no-store",
        "Vercel-CDN-Cache-Control": "no-store",
      },
    });
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
        <IllustrationListByCharacterTemplate
          illustrations={fetchIllustrationsByCategoryIDRes.illustrations}
          characterID={params.character_id}
          character={getCharacterRes.character}
        />
      </div>
    </>
  );
};

export default FetchIllustrationsByCategoryID;
// export const revalidate = 0;
