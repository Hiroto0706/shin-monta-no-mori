import { Category, ChildCategory } from "@/types/category";
import { Dispatch, SetStateAction, useState } from "react";

interface selectCategories {
  childCategories: ChildCategory[];
  checkedChildCategories: ChildCategory[];
  setCheckedChildCategories: Dispatch<SetStateAction<ChildCategory[]>>;
  showCategoryModal: boolean;
  handleCategoriesSelect: (category: ChildCategory) => void;
  toggleCategoriesModal: (status: boolean) => void;
}

const useSelectCategories = (categories: Category[]): selectCategories => {
  const childCategories = categories.flatMap((c) =>
    c.ChildCategory.map((child) => child)
  );
  const [checkedChildCategories, setCheckedChildCategories] = useState<
    ChildCategory[]
  >([]);
  const [showCategoryModal, setShowCategoryModal] = useState(false);

  const handleCategoriesSelect = (category: ChildCategory) => {
    setCheckedChildCategories((prev) => {
      if (prev.some((cate) => cate.id === category.id)) {
        return prev.filter((cate) => cate.id !== category.id);
      }
      return [...prev, category];
    });
  };

  const toggleCategoriesModal = (status: boolean) => {
    setShowCategoryModal(status);
  };

  return {
    childCategories,
    checkedChildCategories,
    setCheckedChildCategories,
    showCategoryModal,
    handleCategoriesSelect,
    toggleCategoriesModal,
  };
};

export default useSelectCategories;
