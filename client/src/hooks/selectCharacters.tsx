import { Character } from "@/types/character";
import { Dispatch, SetStateAction, useState } from "react";

interface selectCharacters {
  checkedCharacters: Character[];
  setCheckedCharacters: Dispatch<SetStateAction<Character[]>>;
  showCharacterModal: boolean;
  handleCharacterSelect: (character: Character) => void;
  toggleCharactersModal: (status: boolean) => void;
}

const useSelectCharacters = (): selectCharacters => {
  const [checkedCharacters, setCheckedCharacters] = useState<Character[]>([]);
  const [showCharacterModal, setShowCharacterModal] = useState(false);

  const handleCharacterSelect = (character: Character) => {
    setCheckedCharacters((prev) => {
      if (prev.some((char) => char.id === character.id)) {
        return prev.filter((char) => char.id !== character.id);
      }
      return [...prev, character];
    });
  };

  const toggleCharactersModal = (status: boolean) => {
    setShowCharacterModal(status);
  };

  return {
    checkedCharacters,
    setCheckedCharacters,
    showCharacterModal,
    handleCharacterSelect,
    toggleCharactersModal,
  };
};

export default useSelectCharacters;
