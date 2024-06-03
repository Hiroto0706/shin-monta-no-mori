export const FetchIllustrationsAPI = (page: number): string => {
  return (
    process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/list?p=" + page
  );
};

export const GetIllustrationAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/" + id;
};

export const CreateIllustrationAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/create";
};

export const EditIllustrationAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/" + id;
};
