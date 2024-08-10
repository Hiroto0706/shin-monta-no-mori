import { create } from "zustand";

interface sidebarStore {
  isShow: boolean;
  toggleIsShow: (state: boolean) => void;
}

const useSidebarStore = create<sidebarStore>((set) => ({
  isShow: true,
  toggleIsShow: (state) => set({ isShow: state }),
}));

export default useSidebarStore;
