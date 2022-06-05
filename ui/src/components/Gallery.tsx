import { Show, For, createResource, createSignal } from "solid-js";
import type { Component } from "solid-js";
import { Thumbnail, ImageOverlay } from ".";

export interface ImageData {
  name: string;
  url: string;
}

export const Gallery: Component = () => {
  const [images] = createResource<ImageData[]>(async () => {
    const res = await fetch("/api/images");
    return res.json();
  });

  const [selectedIdx, setSelectedIdx] = createSignal<number | null>(null);

  const selectIdx = (idx: () => number) => () => setSelectedIdx(idx());

  const prevNext = (direction: number) => () => {
    const idx = selectedIdx();
    if (idx === null) return;

    const imageList = images();
    if (imageList === undefined) return;

    let nextIdx = idx + direction;
    if (nextIdx < 0) nextIdx = imageList.length - 1;
    if (nextIdx >= imageList.length) nextIdx = 0;

    setSelectedIdx(nextIdx);
  };

  return (
    <>
      <Show when={selectedIdx() !== null}>
        <ImageOverlay
          close={() => setSelectedIdx(null)}
          image={() => images()?.[selectedIdx() as number] ?? null}
          next={prevNext(1)}
          previous={prevNext(-1)}
        />
      </Show>
      <div class="grid gap-8 grid-cols-1 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
        <For each={images()} fallback={<div>Loading...</div>}>
          {(image, idx) => (
            <Thumbnail image={image} selectImage={selectIdx(idx)} />
          )}
        </For>
      </div>
    </>
  );
};
