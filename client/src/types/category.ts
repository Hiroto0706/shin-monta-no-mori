export interface ParentCategory {
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

export interface ChildCategory {
  id: number;
  name: string;
  parent_id: number;
  created_at: string;
  updated_at: string;
}

export interface Category {
  ParentCategory: ParentCategory;
  ChildCategory: ChildCategory[];
}