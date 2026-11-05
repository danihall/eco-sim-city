import type { InitData, ReadyMessage, RenderData } from "./_index.js";
import type { ArtifactsKey } from "../../../../constants/_index.js";
import type { ArtifactId, ArtifactSizeId } from "../../../../assets/_index.js";

import * as _ from "../../../../utils/_index.js";
import * as C from "../../../../constants/_index.js";
import * as helpers from "./helpers/helpers.js";
import { ARTIFACTS_BY_ID, ARTIFACTS_SIZES } from "../../../../assets/_index.js";

type WorkerArtifactsInit = InitData<"artifact_ids", "artifacts">;

let spritesheet: WorkerArtifactsInit["sprites"];
let canvas: OffscreenCanvas;
let ctx: OffscreenCanvasRenderingContext2D;
let pixel_ratio: number;

let artifact_ids: Uint8Array;

self.addEventListener("message", onMessage);

async function onMessage(event: MessageEvent<WorkerArtifactsInit | RenderData>) {
  const { data } = event;

  switch (data.type) {
    case "init":
      return init(data).then(self.postMessage);
    case "render":
      return render(data);
  }
}

async function init(data: WorkerArtifactsInit): Promise<ReadyMessage> {
  canvas = data.offscreen_canvas;
  ctx = canvas.getContext("2d")!;
  ctx.imageSmoothingEnabled = false;
  pixel_ratio = data.pixel_ratio;

  artifact_ids = new Uint8Array(data.shared_buffers.artifact_ids);
  spritesheet = data.sprites;

  return { ready: true };
}

function render(data: RenderData) {
  helpers.clearCanvas(ctx, data.camera, pixel_ratio);

  for (const { q, r, center } of helpers.visibleTilesIterator(ctx, data.camera)) {
    const tile_index = _.getTileIndex(q, r);
    const artifact_id = artifact_ids[tile_index];

    if (artifact_id > 0) {
      const artifact_name = ARTIFACTS_BY_ID[artifact_ids[tile_index] as ArtifactId];

      ctx.drawImage(
        spritesheet[artifact_name],
        center.x - C.HEX_WIDTH / 2,
        center.y - C.HEX_HEIGHT / 2,
      );
    }
  }
}

export type { WorkerArtifactsInit };
