export const FetchIllustrationsAPI = (page: number) => {
  return (
    process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/list?p=" + page
  );
};

export const GetIllustrationAPI = (id: number) => {
  return (
    process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/" + id
  );
};

export const CreateIllustrationAPI = () => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/create";
};
