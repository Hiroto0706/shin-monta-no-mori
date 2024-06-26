import CreateCharacter from "@/components/admin/characters/createForm";
import { getServerAccessToken } from "@/utils/accessToken/server";

const CreateIllustrationPage = async () => {
  const accessToken = getServerAccessToken();

  return (
    <>
      <CreateCharacter
        accessToken={accessToken}
      />
    </>
  );
};

export default CreateIllustrationPage;
