"use client";

import { Category } from "@/types/category";
import { Character } from "@/types/character";
import Image from "next/image";

type Props = {
  links: {
    id: number;
    href: string;
    icon: string;
    icon_active: string;
    text: string;
  }[];
  selectedLink: number;
  characters: Character[];
  categories: Category[];
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
      <div className="w-60 min-h-screen overflow-y-scroll bg-gray-50 border-r border-gray-200 fixed top-0 left-0 pl-16 pt-16">
        <div className="pt-4 px-2">
          {selectedLinkObj && (
            <>
              {/* カテゴリーサイドバー */}
              {selectedLinkObj.text === "カテゴリ" && (
                <>
                  {categories.map((category) => (
                    <div key={category.ParentCategory.id} className="mb-4">
                      <div className="flex items-center mb-1">
                        <Image
                          src={category.ParentCategory.src}
                          alt={category.ParentCategory.filename.String}
                          width={20}
                          height={20}
                        />
                        <span className="font-bold text-md ml-2">
                          {category.ParentCategory.name}
                        </span>
                      </div>
                      {category.ChildCategory.map((childCategory) => (
                        <a
                          key={childCategory.id}
                          href={`/illustrations/category/${childCategory.id}`}
                          className="text-sm py-1 px-2 hover:bg-gray-200 duration-200 rounded-lg cursor-pointer block"
                        >
                          <span>{childCategory.name}</span>
                        </a>
                      ))}
                    </div>
                  ))}
                </>
              )}

              {/* キャラクターサイドバー */}
              {selectedLinkObj.text === "キャラ" && (
                <>
                  {characters.map((character) => (
                    <div
                      key={character.id}
                      className="flex items-center mb-2 hover:bg-gray-200 duration-200 p-1 rounded-lg cursor-pointer"
                    >
                      <Image
                        className="border-2 border-gray-200 rounded-full bg-white"
                        src={character.src}
                        alt={character.filename.String}
                        width={32}
                        height={32}
                      />
                      <span className="ml-1">{character.name}</span>
                    </div>
                  ))}
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
