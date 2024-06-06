import { Category } from "@/types/category";
import { Character } from "@/types/character";
import {
  fetchCategories,
} from "@/app/admin/illustrations/page";
import CreateIllustration from "@/components/admin/illustrations/createForm";
import {} from "@/utils/accessToken/accessToken";
import { getServerAccessToken } from "@/utils/accessToken/server";

const CreateIllustrationPage = async () => {
  const accessToken = getServerAccessToken();

  return (
    <>
      {/* <CreateIllustration
        characters={characters}
        categories={categories}
        accessToken={accessToken}
      /> */}
    </>
  );
};

export default CreateIllustrationPage;
