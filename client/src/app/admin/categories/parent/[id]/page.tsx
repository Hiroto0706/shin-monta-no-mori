import axios from "axios";
import { getServerAccessToken } from "@/utils/accessToken/server";

const EditParentCategoryPage = async ({
  params,
}: {
  params: { id: number };
}) => {
  const accessToken = getServerAccessToken();

  return <>parent : id = {params.id}</>;
};

export default EditParentCategoryPage;
