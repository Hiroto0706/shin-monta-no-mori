export const FetchIllustrationsAPI = (page: number) => {
  return (
    process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/list?p=" + page
  );
};

export const CreateIllustrationAPI = () => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/create";
};
