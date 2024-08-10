export interface ParentCategory {
  id: number;
  name: string;
  src: string;
  filename: {
    String: string;
    Valid: boolean;
  };
  priority_level: number;
  created_at: string;
  updated_at: string;
}

export interface ChildCategory {
  id: number;
  name: string;
  parent_id: number;
  priority_level: number;
  created_at: string;
  updated_at: string;
}

export interface Category {
  ParentCategory: ParentCategory;
  ChildCategory: ChildCategory[];
}
