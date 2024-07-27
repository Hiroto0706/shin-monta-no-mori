"use client";

import { Category } from "@/types/category";
import Image from "next/image";
import Link from "next/link";

type Props = {
  categories: Category[];
};

const ListCategoriesTable: React.FC<Props> = ({ categories }) => {
  return (
    <div className="my-12 text-sm md:text-md">
      {categories.map((category) => (
        <div
          key={category.ParentCategory.id}
          className="border-2 border-gray-200 rounded-lg p-1 mb-3 bg-white"
        >
          <div className="bg-gray-200 rounded-lg py-2 px-4">
            <div className="flex justify-between items-center">
              <div className="flex items-center">
                <Image
                  className="mr-4"
                  src={category.ParentCategory.src}
                  alt={category.ParentCategory.name}
                  width={32}
                  height={32}
                />
                <span className="text-xl font-bold">
                  {category.ParentCategory.name}
                </span>
              </div>

              <Link
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
              </Link>
            </div>
          </div>

          <div className="pt-4 pb-2 px-4 flex flex-wrap">
            {category.ChildCategory.map((cc) => (
              <Link
                key={cc.id}
                href={`categories/child/${cc.id}`}
                className="flex items-center text-lg mr-4 cursor-pointer py-2 px-4 mb-2 rounded-full hover:bg-gray-100 duration-200"
              >
                # {cc.name}
              </Link>
            ))}
            <Link
              href={`categories/child/new?parent_id=${category.ParentCategory.id}`}
              className="flex items-center border-2 mb-2 py-2 px-4 cursor-pointer rounded-full hover:bg-gray-100 duration-200"
            >
              + 子カテゴリ追加
            </Link>
          </div>
        </div>
      ))}
    </div>
  );
};

export default ListCategoriesTable;
