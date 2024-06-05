import React from "react";

type Props = {
  currentPage: number;
  totalPages: number;
  count: number;
  path: string;
  query?: string;
  charactersParams?: string;
  categoriesParams?: string;
};

const Pagination = ({
  currentPage,
  totalPages,
  count,
  path,
  query,
  charactersParams,
  categoriesParams,
}: Props) => {
  const limit =
    process.env.NEXT_PUBLIC_IMAGE_FETCH_LIMIT != undefined
      ? Number(process.env.NEXT_PUBLIC_IMAGE_FETCH_LIMIT)
      : 10;
  let startPage = Math.max(0, currentPage - 2);
  let endPage = Math.min(totalPages - 1, currentPage + 2);

  if (currentPage <= 2) {
    endPage = Math.min(4, totalPages - 1);
  } else if (currentPage >= totalPages - 3) {
    startPage = Math.max(0, totalPages - 5);
  }

  const pageNumbers = [];
  for (let i = startPage; i <= endPage; i++) {
    pageNumbers.push(i);
  }

  const buildUrl = (page: number) => {
    const params = new URLSearchParams({ p: page.toString() });
    if (query) params.append("q", query);
    if (charactersParams) params.append("characters", charactersParams);
    if (categoriesParams) params.append("categories", categoriesParams);
    return `${path}?${params.toString()}`;
  };

  return (
    <div className="flex items-center space-x-2">
      <a href={buildUrl(currentPage - 1)} aria-label="Previous Page">
        <button
          className={`px-3 py-1 border rounded ${
            currentPage === 0 || count < limit
              ? "cursor-not-allowed opacity-50"
              : "hover:bg-gray-200"
          }`}
          disabled={currentPage === 0}
        >
          ＜
        </button>
      </a>
      {pageNumbers.map((number) => (
        <a key={number} href={buildUrl(number)}>
          <button
            className={`px-3 py-1 border rounded ${
              currentPage === number
                ? "bg-green-600 text-white"
                : "hover:bg-gray-200"
            }`}
          >
            {number}
          </button>
        </a>
      ))}
      <a href={buildUrl(currentPage + 1)} aria-label="Next Page">
        <button
          className={`px-3 py-1 border rounded ${
            currentPage === totalPages - 1 || count < limit
              ? "cursor-not-allowed opacity-50"
              : "hover:bg-gray-200"
          }`}
          disabled={currentPage === totalPages - 1}
        >
          ＞
        </button>
      </a>
    </div>
  );
};

export default Pagination;
