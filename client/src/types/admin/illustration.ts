import { Illustration } from "../illustration";

export interface FetchIllustrationsResponse {
  illustrations: Illustration[];
  total_pages: number;
  total_count: number;
}

export interface GetIllustrationResponse {
  illustration: Illustration | null;
}
