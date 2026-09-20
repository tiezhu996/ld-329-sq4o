import { defineStore } from 'pinia';
import { ref } from 'vue';

export type AppView = 'dashboard' | 'match';

const VIEW_STORAGE_KEY = 'cyskillswap.view';

function loadView(): AppView {
  return localStorage.getItem(VIEW_STORAGE_KEY) === 'match' ? 'match' : 'dashboard';
}

export const useAppStore = defineStore('app', () => {
  const view = ref<AppView>(loadView());

  function setView(next: AppView): void {
    view.value = next;
    localStorage.setItem(VIEW_STORAGE_KEY, next);
  }

  return { view, setView };
});
