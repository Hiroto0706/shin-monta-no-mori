import axios from "axios";
import { GetCharacterResponse } from "@/types/admin/character";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { getServerAccessToken } from "@/utils/accessToken/server";
import EditCharacter from "@/components/admin/characters/editForm";
import { GetCharacterAPI } from "@/api/admin/character";

const getCharacter = async (
  id: number,
  accessToken: string | undefined
): Promise<GetCharacterResponse> => {
  try {
    const response = await axios.get(GetCharacterAPI(id), {
      headers: {
        Authorization: SetBearerToken(accessToken),
      },
    });
    return response.data;
  } catch (error) {
    console.error(error);
    return { character: null };
  }
};

const EditIllustrationPage = async ({ params }: { params: { id: number } }) => {
  const accessToken = getServerAccessToken();
  const characterRes = await getCharacter(params.id, accessToken);

  return (
    <>
      {characterRes.character && (
        <EditCharacter
          id={params.id}
          character={characterRes.character}
          accessToken={accessToken}
        />
      )}
    </>
  );
};

export default EditIllustrationPage;
