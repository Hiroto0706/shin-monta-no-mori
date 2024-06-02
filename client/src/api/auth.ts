export const VerifyAPI = () => {
  return process.env.NEXT_PUBLIC_BASE_API + "auth/verify";
};

export const AuthLoginAPI = () => {
  return process.env.NEXT_PUBLIC_BASE_API + "auth/login";
};
