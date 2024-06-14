"use client";

import { useRouter } from "next/navigation";
import axios from "axios";
import { useEffect, useState } from "react";
import { SetBearerToken } from "@/utils/accessToken/accessToken";
import { CreateChildCategoryAPI } from "@/api/admin/category";
import { Category } from "@/types/category";
import Image from "next/image";
import { truncateText } from "@/utils/text";

type Props = {
  parentID: number;
  categories: Category[];
  accessToken: string | undefined;
};

const CreateChildCategory: React.FC<Props> = ({
  parentID,
  categories,
  accessToken,
}) => {
  const router = useRouter();
  const [name, setName] = useState("");
  const parentCategories = categories.map((c) => c.ParentCategory);
  const [checkedParentCategoryID, setCheckedParentCategoryID] =
    useState<number>(parentID);
  const [showCategoryModal, setShowCategoryModal] = useState(false);
  const toggleCategoriesModal = (status: boolean) => {
    setShowCategoryModal(status);
  };

  const createIllustration = async (event: React.FormEvent) => {
    event.preventDefault();

    const formData = new FormData();
    formData.append("name", name);
    formData.append("parent_id", checkedParentCategoryID.toString());

    try {
      const response = await axios.post(CreateChildCategoryAPI(), formData, {
        headers: {
          Authorization: SetBearerToken(accessToken),
        },
      });

      if (response.status === 200) {
        alert(response.data.message);
        router.push("/admin/categories");
      }
    } catch (error) {
      console.error("子カテゴリの作成に失敗しました", error);
      alert("子カテゴリの作成に失敗しました");
    }
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        (event.target as HTMLElement).closest(".category-modal") === null &&
        (event.target as HTMLElement).closest(".category-modal-content") ===
          null
      ) {
        toggleCategoriesModal(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return (
    <>
      <div className="max-w-7xl m-auto">
        <h1 className="text-2xl font-bold mb-6">子カテゴリの作成</h1>
        <form
          className="border-2 border-gray-300 rounded-lg p-12 bg-white"
          onSubmit={createIllustration}
        >
          <div className="mb-16">
            <label className="text-xl">名前</label>
            <input
              className="w-full p-4 border border-gray-200 rounded mt-2"
              type="text"
              placeholder="名前を入力してください"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>

          <div className="mb-16">
            <label className="text-xl">親カテゴリ</label>
            <div className="relative">
              <div
                onClick={() => toggleCategoriesModal(!showCategoryModal)}
                className="border-2 border-gray-200 mt-4 py-4 px-4 rounded bg-white flex justify-between cursor-pointer category-modal"
              >
                <span>
                  {
                    parentCategories.find(
                      (c) => c.id == checkedParentCategoryID
                    )?.name
                  }
                </span>
                <Image
                  className={`duration-100 ${
                    !showCategoryModal ? "rotate-90" : "-rotate-90"
                  }`}
                  src="/icon/arrow.png"
                  alt="arrowアイコン"
                  width={20}
                  height={20}
                />
              </div>
              {showCategoryModal && (
                <div className="absolute left-0 bg-white border-2 border-gray-300 p-4 rounded w-full z-10 max-h-60 overflow-y-auto shadow-md category-modal-content">
                  {parentCategories.map((cate) => (
                    <div
                      key={cate.id}
                      className="flex items-center rounded hover:bg-gray-200 cursor-pointer"
                    >
                      <input
                        type="radio"
                        checked={cate.id == checkedParentCategoryID}
                        onChange={() => setCheckedParentCategoryID(cate.id)}
                        className="mx-2 cursor-pointer"
                        id={`character-${cate.id}`}
                      />
                      <label
                        htmlFor={`character-${cate.id}`}
                        className="cursor-pointer w-full py-1"
                      >
                        {truncateText(cate.name, 30)}
                      </label>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          <button className="py-3 bg-green-600 text-white font-bold text-lg rounded-lg w-full hover:bg-white hover:text-green-600 border-2 border-green-600 duration-200">
            子カテゴリ作成
          </button>
        </form>
      </div>
    </>
  );
};

export default CreateChildCategory;
