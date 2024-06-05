import axios from "axios";
import { GetIllustrationResponse } from "@/types/illustration";
import { Category } from "@/types/category";
import { Character } from "@/types/character";
import {
  fetchCategories,
  fetchCharacters,
} from "@/app/admin/illustrations/page";
import { GetIllustrationAPI } from "@/api/illustration";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { getServerAccessToken } from "@/utils/accessToken/server";
import EditIllustration from "@/components/admin/illustrations/editForm";

const getIllustration = async (
  id: number,
  accessToken: string | undefined
): Promise<GetIllustrationResponse> => {
  try {
    const response = await axios.get(GetIllustrationAPI(id), {
      headers: {
        Authorization: SetBearerToken(accessToken),
      },
    });
    return response.data;
  } catch (error) {
    console.error(error);
    return { illustration: null };
  }
};

const EditIllustrationPage = async ({ params }: { params: { id: number } }) => {
  const accessToken = getServerAccessToken();
  const illustrationRes: GetIllustrationResponse = await getIllustration(
    params.id,
    accessToken
  );
  const characters: Character[] = await fetchCharacters(accessToken);
  const categories: Category[] = await fetchCategories(accessToken);

  return (
    <>
      {illustrationRes.illustration && (
        <EditIllustration
          id={params.id}
          illustration={illustrationRes.illustration}
          characters={characters}
          categories={categories}
          accessToken={accessToken}
        />
      )}
    </>
  );
};

export default EditIllustrationPage;
