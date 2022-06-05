import { onMount, onCleanup } from "solid-js";
import type { Component } from "solid-js";
import type { ImageData } from "./Gallery";

interface ImageOverlayProps {
  close: () => void;
  next: () => void;
  previous: () => void;
  image: () => ImageData | null;
}

function useKeyDown(handler: (event: KeyboardEvent) => void) {
  onMount(() => {
    document.addEventListener("keydown", handler);
  });

  onCleanup(() => document.removeEventListener("keydown", handler));
}

export const ImageOverlay: Component<ImageOverlayProps> = ({
  close,
  next,
  previous,
  image,
}) => {
  useKeyDown((e: KeyboardEvent) => {
    switch (e.key) {
      case "Escape":
        close();
        break;
      case "ArrowLeft":
        previous();
        break;
      case "ArrowRight":
        next();
        break;
    }
  });

  return (
    <>
      <div
        onClick={close}
        class="fixed top-0 left-0 right-0 bottom-0 flex items-center bg-slate-900 opacity-90"
      ></div>
      <img
        src={`/api${image()?.url}`}
        class="border-4 fixed top-1/2 left-1/2 -translate-y-1/2 -translate-x-1/2 h-3/4"
      />
    </>
  );
};
