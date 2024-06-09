import { Illustration } from "../illustration";

export interface FetchIllustrationsResponse {
  illustrations: Illustration[];
}

export interface GetIllustrationResponse {
  illustration: Illustration | null;
}
