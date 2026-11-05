import type { CameraState, CityData } from "components";
import type { SPRITESHEET } from "constants";
import type { WorkerArtifactsInit } from "./canvas_artifacts.js";
import type { WorkerGridInit } from "./canvas_grid.js";
import type { WorkerHeatmapInit } from "./canvas_heatmap.js";
import type { WorkerUiInit } from "./canvas_ui.js";
import type { WorkerZonesInit } from "./canvas_zones.js";

type ReadyMessage = {
  ready: boolean;
};

type Sprites<S extends keyof typeof SPRITESHEET> = {
  [Key in S]: { [SubKey in keyof (typeof SPRITESHEET)[Key]]: ImageBitmap };
}[S];

type InitData<T extends keyof CityData = never, S extends keyof typeof SPRITESHEET = never> = {
  type: "init";
  pixel_ratio: number;
  offscreen_canvas: OffscreenCanvas;
  shared_buffers: Pick<CityData, T>;
  sprites: Sprites<S>;
};

type RenderData = {
  type: "render";
  camera: CameraState;
};

type ResetData = {
  type: "reset";
  camera: CameraState;
};

type Visibility = {
  width: number;
  height: number;
};

export type {
  ReadyMessage,
  Sprites,
  InitData,
  RenderData,
  ResetData,
  Visibility,
  WorkerArtifactsInit,
  WorkerGridInit,
  WorkerHeatmapInit,
  WorkerUiInit,
  WorkerZonesInit,
};
