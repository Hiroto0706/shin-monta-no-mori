const maxTextLength = 20;
export const truncateText = (text: string | undefined) => {
  return text && text.length > maxTextLength
    ? `${text.slice(0, maxTextLength)}...`
    : text;
};
