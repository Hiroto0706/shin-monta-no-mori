export const truncateText = (
  text: string | undefined,
  maxTextLength: number = 20
) => {
  return text && text.length > maxTextLength
    ? `${text.slice(0, maxTextLength)}...`
    : text;
};

const calcDate = (dateString: string) => {
  const date = new Date(dateString);

  const year = date.getFullYear();
  const month = ("0" + (date.getMonth() + 1)).slice(-2);
  const day = ("0" + date.getDate()).slice(-2);
  const hours = ("0" + date.getHours()).slice(-2);
  const minutes = ("0" + date.getMinutes()).slice(-2);

  return {
    year,
    month,
    day,
    hours,
    minutes,
  };
};

export const FormatDate = (dateString: string): string => {
  const { year, month, day, hours, minutes } = calcDate(dateString);

  return `${year}/${month}/${day} ${hours}:${minutes}`;
};

export const UpdatedAtFormat = (dateString: string): string => {
  const { year, month, day } = calcDate(dateString);

  return `${year}年${month}月${day}日`;
};
