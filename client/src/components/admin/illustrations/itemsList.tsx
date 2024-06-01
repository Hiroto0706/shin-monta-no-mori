import { Illustration } from "@/types/illustration";

type Props = {
  illustrations: Illustration[];
};

const ItemsList: React.FC<Props> = ({ illustrations }) => {
  return (
    <div className="my-12 w-full bg-white overflow-x-auto rounded-lg border-2 border-gray-200 scrollbar-hide">
      <table className="table-auto min-w-max w-full">
        <thead className="bg-gray-200 border-2 border-gray-200">
          <tr>
            <th className="px-4 py-2">ID</th>
            <th className="px-4 py-2">Title</th>
            <th className="px-4 py-2">Original Source</th>
            <th className="px-4 py-2">Filename</th>
            <th className="px-4 py-2">Created At</th>
            <th className="px-4 py-2">Updated At</th>
            <th className="px-4 py-2">Updated At</th>
            <th className="px-4 py-2">Updated At</th>
          </tr>
        </thead>
        <tbody>
          {illustrations.map((illustration, index) => (
            <tr key={index} className="border-2 border-gray-100">
              <td className="px-8 py-4">{illustration.Image.id}</td>
              <td className="px-8 py-4">{illustration.Image.title}</td>
              <td className="px-8 py-4">{illustration.Image.original_src}</td>
              <td className="px-8 py-4">
                {illustration.Image.original_filename}
              </td>
              <td className="px-8 py-4">{illustration.Image.created_at}</td>
              <td className="px-8 py-4">{illustration.Image.updated_at}</td>
              <td className="px-8 py-4">{illustration.Image.updated_at}</td>
              <td className="px-8 py-4">{illustration.Image.updated_at}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default ItemsList;
