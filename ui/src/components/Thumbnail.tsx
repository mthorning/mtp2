import type { Component } from "solid-js";
import type { ImageData } from "./Gallery";

interface ThumbnailsProps {
  image: ImageData;
  selectImage: () => void;
}

export const Thumbnail: Component<ThumbnailsProps> = ({
  image,
  selectImage,
}) => {
  return (
    <div
      role="button"
      onClick={() => selectImage()}
      class="aspect-square min-w-80"
    >
      <div
        style={{
          "background-image": `url('/api${image.url}')`,
        }}
        class="w-full h-full h-full bg-cover bg-center border-2"
      />
    </div>
  );
};
