export interface Character {
  id: number;
  name: string;
  src: string;
  filename: {
    String: string;
    Valid: boolean;
  };
  created_at: string;
  updated_at: string;
}
