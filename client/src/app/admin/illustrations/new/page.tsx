import { Category } from "@/types/category";
import { Character } from "@/types/character";
import {
  fetchAllCharacters,
  fetchCategories,
} from "@/app/admin/illustrations/page";
import CreateIllustration from "@/components/admin/illustrations/createForm";
import {} from "@/utils/accessToken/accessToken";
import { getServerAccessToken } from "@/utils/accessToken/server";

const CreateIllustrationPage = async () => {
  const accessToken = getServerAccessToken();
  const characters: Character[] = await fetchAllCharacters(accessToken);
  const categories: Category[] = await fetchCategories(accessToken);

  return (
    <>
      <CreateIllustration
        characters={characters}
        categories={categories}
        accessToken={accessToken}
      />
    </>
  );
};

export default CreateIllustrationPage;
