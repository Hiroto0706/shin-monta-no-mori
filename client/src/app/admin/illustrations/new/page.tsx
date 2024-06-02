import { Category } from "@/types/category";
import { Character } from "@/types/character";
import {
  fetchCategories,
  fetchCharacters,
} from "@/app/admin/illustrations/page";
import CreateIllustration from "@/components/admin/illustrations/createForm";

const CreateIllustrationPage = async () => {
  const characters: Character[] = await fetchCharacters();
  const categories: Category[] = await fetchCategories();
  return (
    <>
      <CreateIllustration characters={characters} categories={categories} />
    </>
  );
};

export default CreateIllustrationPage;
