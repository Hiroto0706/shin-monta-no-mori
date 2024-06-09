import CreateParentCategory from "@/components/admin/categories/parentCategory/createForm";
import { getServerAccessToken } from "@/utils/accessToken/server";

const CreateParentCategoryPage = async () => {
  const accessToken = getServerAccessToken();

  return (
    <>
      <CreateParentCategory accessToken={accessToken} />
    </>
  );
};

export default CreateParentCategoryPage;
