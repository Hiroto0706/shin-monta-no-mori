"use client";

import { Category } from "@/types/category";
import Image from "next/image";

type Props = {
  categories: Category[];
};

const ListCategoriesTable: React.FC<Props> = ({ categories }) => {
  return (
    <div className="my-12">
      {categories.map((category) => (
        <div
          key={category.ParentCategory.id}
          className="border-2 border-gray-200 rounded-lg p-1 mb-3 bg-white"
        >
          <div className="bg-gray-200 rounded-lg p-4">
            <div className="flex justify-between items-center">
              <div className="flex items-center">
                <Image
                  className="mr-6"
                  src={category.ParentCategory.src}
                  alt={category.ParentCategory.name}
                  width={50}
                  height={50}
                />
                <span className="text-2xl font-bold">
                  {category.ParentCategory.name}
                </span>
              </div>

              <a
                href={`categories/parent/${category.ParentCategory.id}`}
                className="flex items-center justify-center cursor-pointer hover:bg-gray-300 duration-200 py-2 px-4 rounded-lg"
              >
                <Image
                  src="/icon/edit.png"
                  alt="editアイコン"
                  width={20}
                  height={20}
                />
                <span className="ml-1 mt-0.5">編集</span>
              </a>
            </div>
          </div>

          {category.ChildCategory.length > 0 && (
            <div className="pt-6 pb-4 px-4 flex flex-wrap">
              {category.ChildCategory.map((cc) => (
                <a
                  key={cc.id}
                  href={`categories/child/${cc.id}`}
                  className="text-xl mr-4 cursor-pointer py-2 px-4 mb-2 rounded-full hover:bg-gray-200 duration-200"
                >
                  # {cc.name}
                </a>
              ))}
              <a
                href="categories/child/new"
                className="flex items-center text-xl border-2 mb-2 px-4 cursor-pointer rounded-full hover:bg-gray-200 duration-200"
              >
                + 子カテゴリ追加
              </a>
            </div>
          )}
        </div>
      ))}
    </div>
  );
};

export default ListCategoriesTable;
