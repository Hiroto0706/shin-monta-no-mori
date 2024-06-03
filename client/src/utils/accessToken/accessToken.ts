export const SetBearerToken = (accessToken: string | undefined) => {
  return `Bearer ${accessToken}`;
};
