export const Verify = () => {
  return process.env.NEXT_PUBLIC_BASE_API + "auth/verify";
};

export const AuthLogin = () => {
  return process.env.NEXT_PUBLIC_BASE_API + "auth/login";
};
