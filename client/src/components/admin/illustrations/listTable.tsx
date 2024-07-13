"use client";

import { useRouter } from "next/navigation";
import { FetchIllustrationsResponse } from "@/types/admin/illustration";
import { FormatDate, truncateText } from "@/utils/text";
import Image from "next/image";

type Props = {
  illustrations: FetchIllustrationsResponse;
};

const ListIllustrationsTable: React.FC<Props> = ({ illustrations }) => {
  const router = useRouter();
  const maxTextLength = 15;
  return (
    <div className="my-12 w-full bg-white overflow-x-auto rounded-lg border-2 border-gray-200 scrollbar-hide">
      <table className="table-auto min-w-max w-full">
        <thead className="bg-gray-200 border-2 border-gray-200 py-4">
          <tr>
            <th className="px-6 py-4">ID</th>
            <th className="px-6 py-4">タイトル</th>
            <th className="px-6 py-4">イメージ</th>
            <th className="px-6 py-4">ファイル名</th>
            <th className="px-6 py-4">キャラクター</th>
            <th className="px-6 py-4">カテゴリー</th>
            <th className="px-6 py-4">最終更新日時</th>
            <th className="px-6 py-4">作成日時</th>
          </tr>
        </thead>
        <tbody>
          {illustrations.illustrations.map((illustration, index) => (
            <tr
              key={index}
              className="border-2 border-gray-100 cursor-pointer duration-200 hover:bg-gray-50"
              onClick={() =>
                router.push(`illustrations/edit/${illustration.Image.id}`)
              }
            >
              <td className="px-6 py-4">{illustration.Image.id}</td>
              <td className="px-6 py-4">{illustration.Image.title}</td>
              <td className="px-6 py-4">
                <div className="flex">
                  <Image
                    className="border-2 rounded-lg border-gray-200 p-2"
                    src={illustration.Image.original_src}
                    alt={`${illustration.Image.title}の画像`}
                    width={90}
                    height={90}
                  />
                  {illustration.Image.simple_src.Valid &&
                    illustration.Image.simple_filename.Valid && (
                      <Image
                        className="ml-4 border-2 rounded-lg border-gray-200 p-2"
                        src={illustration.Image.simple_src.String}
                        alt={`${illustration.Image.title}の画像`}
                        width={90}
                        height={90}
                      />
                    )}
                </div>
              </td>
              <td className="px-6 py-4">
                {illustration.Image.original_filename}
              </td>
              <td className="px-6 py-4">
                <div className="flex">
                  {illustration.Characters != undefined && (
                    <>
                      {illustration.Characters.flatMap((c, index) => (
                        <Image
                          key={index}
                          className={`border-2 rounded-full border-gray-200 ${
                            index > 0 ? "ml-2" : ""
                          }`}
                          src={c.Character.src}
                          alt={`${c.Character.name}の画像`}
                          width={60}
                          height={60}
                        />
                      ))}
                    </>
                  )}
                </div>
              </td>
              <td className="px-6 py-4">
                {" "}
                <div className="flex">
                  {illustration.Categories != undefined && (
                    <>
                      {illustration.Categories.flatMap((category) =>
                        category.ChildCategory.map((c, index) => (
                          <span
                            key={c.id}
                            className="py-2 px-4 bg-gray-100 rounded-full mx-1"
                          >
                            # {truncateText(c.name, maxTextLength)}
                          </span>
                        ))
                      )}
                    </>
                  )}
                </div>
              </td>
              <td className="px-6 py-4">
                {FormatDate(illustration.Image.updated_at)}
              </td>
              <td className="px-6 py-4">
                {FormatDate(illustration.Image.created_at)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default ListIllustrationsTable;
