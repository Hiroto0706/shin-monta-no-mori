"use client";

import Loader from "@/components/common/loader";
import { Category } from "@/types/category";
import { Character } from "@/types/character";
import Image from "next/image";
import Link from "next/link";

type Props = {
  links: {
    id: number;
    href: string;
    icon: string;
    icon_active: string;
    text: string;
  }[];
  selectedLink: number | null;
  characters: Character[] | undefined;
  categories: Category[] | undefined;
};

const SidebarSub: React.FC<Props> = ({
  links,
  selectedLink,
  characters,
  categories,
}) => {
  const selectedLinkObj = links.find((link) => selectedLink === link.id);

  return (
    <>
      <div className="w-72 h-full overflow-y-scroll bg-gray-50 border-r border-gray-200 fixed top-0 left-0 pl-16 pt-16">
        <div className="p-4">
          {selectedLinkObj && (
            <>
              {/* カテゴリーサイドバー */}
              {selectedLinkObj.text === "カテゴリ" && (
                <>
                  {categories != undefined && categories.length >= 0 ? (
                    <>
                      {categories.map((category) => (
                        <div key={category.ParentCategory.id} className="mb-4">
                          <div className="flex items-center my-2">
                            <Image
                              src={category.ParentCategory.src}
                              alt={category.ParentCategory.filename.String}
                              width={24}
                              height={24}
                            />
                            <span className="font-bold text-md text-black ml-2">
                              {category.ParentCategory.name}
                            </span>
                          </div>
                          {category.ChildCategory.map((childCategory) => (
                            <Link
                              key={childCategory.id}
                              href={`/illustrations/category/${childCategory.id}`}
                              className="text-sm py-1 px-2 hover:bg-gray-200 duration-200 rounded-full cursor-pointer block mb-1"
                            >
                              # {childCategory.name}
                            </Link>
                          ))}
                        </div>
                      ))}
                    </>
                  ) : (
                    <>
                      <>
                        <Loader height="h-[86vh]" size={30} />
                      </>
                    </>
                  )}
                </>
              )}

              {/* キャラクターサイドバー */}
              {selectedLinkObj.text === "キャラ" && (
                <>
                  {characters != undefined && characters.length >= 0 ? (
                    <>
                      {characters.map((character) => (
                        <Link
                          key={character.id}
                          href={`/illustrations/character/${character.id}`}
                          className="flex items-center mb-2 hover:bg-gray-200 duration-200 p-1 rounded-full cursor-pointer"
                        >
                          <Image
                            className="border border-gray-200 rounded-full bg-white shadow"
                            src={character.src}
                            alt={character.filename.String}
                            width={36}
                            height={36}
                          />
                          <span className="ml-2 text-sm">{character.name}</span>
                        </Link>
                      ))}
                    </>
                  ) : (
                    <>
                      <Loader height="h-[86vh]" size={30} />
                    </>
                  )}
                </>
              )}
            </>
          )}
        </div>
      </div>
    </>
  );
};

export default SidebarSub;
