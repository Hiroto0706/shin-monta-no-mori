"use client";

import { useRouter } from "next/navigation";
import { formatDate } from "@/utils/text";
import Image from "next/image";
import { FetchCharactersResponse } from "@/types/admin/character";

type Props = {
  characters: FetchCharactersResponse;
};

const ListCharactersTable: React.FC<Props> = ({ characters }) => {
  const router = useRouter();
  return (
    <div className="my-12 w-full bg-white overflow-x-auto rounded-lg border-2 border-gray-200 scrollbar-hide">
      <table className="table-auto min-w-max w-full">
        <thead className="bg-gray-200 border-2 border-gray-200 py-4">
          <tr>
            <th className="px-6 py-4">ID</th>
            <th className="px-6 py-4">タイトル</th>
            <th className="px-6 py-4">イメージ</th>
            <th className="px-6 py-4">ファイル名</th>
            <th className="px-6 py-4">最終更新日時</th>
            <th className="px-6 py-4">作成日時</th>
          </tr>
        </thead>
        <tbody>
          {characters.characters.map((character, index) => (
            <tr
              key={index}
              className="border-2 border-gray-100 cursor-pointer duration-200 hover:bg-gray-50"
              onClick={() => router.push(`characters/edit/${character.id}`)}
            >
              <td className="px-6 py-4">{character.id}</td>
              <td className="px-6 py-4">{character.name}</td>
              <td className="px-6 py-4">
                <div className="flex">
                  <Image
                    className="border-2 rounded-lg border-gray-200 p-2"
                    src={character.src}
                    alt={`${character.name}の画像`}
                    width={90}
                    height={90}
                  />
                </div>
              </td>
              <td className="px-6 py-4">{character.filename.String}</td>
              <td className="px-6 py-4">{formatDate(character.updated_at)}</td>
              <td className="px-6 py-4">{formatDate(character.created_at)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default ListCharactersTable;
