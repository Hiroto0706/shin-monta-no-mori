import { Dispatch, SetStateAction, useEffect, useState } from "react";

interface priorityLevel {
  checkedPriorityLevel: number;
  setCheckedPriorityLevel: Dispatch<SetStateAction<number>>;
  showPriorityLevelModal: boolean;
  togglePriorityLevelModal: (status: boolean) => void;
}

const usePriorityLevel = (priorityLevel: number): priorityLevel => {
  const [checkedPriorityLevel, setCheckedPriorityLevel] =
    useState(priorityLevel);
  const [showPriorityLevelModal, setShowPriorityLevelModal] = useState(false);

  const togglePriorityLevelModal = (status: boolean) => {
    setShowPriorityLevelModal(status);
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        (event.target as HTMLElement).closest(".priority-modal") === null &&
        (event.target as HTMLElement).closest(".priority-modal-content") ===
          null
      ) {
        setShowPriorityLevelModal(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return {
    checkedPriorityLevel,
    setCheckedPriorityLevel,
    showPriorityLevelModal,
    togglePriorityLevelModal,
  };
};

export default usePriorityLevel;
