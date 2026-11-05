import type { InitData, ReadyMessage, RenderData, ResetData } from "./_index.js";
import type { ArtifactId, SuperKindId } from "../../../../assets/_index.js";

import { ARTIFACTS_BY_ID, SUPER_KINDS_BY_ID } from "../../../../assets/_index.js";
import * as C from "../../../../constants/_index.js";
import * as _ from "./helpers/helpers.js";

type WorkerUiInit = InitData<never, "artifacts" | "zones">;

type UIType = "ui_zones" | "ui_build_artifact" | "ui_destroy_artifact";

type UIData = {
  ui: UIType;
  qs_buffer: ArrayBuffer;
  rs_buffer: ArrayBuffer;
  item_ids_buffer: ArrayBuffer;
  occupied_buffer: ArrayBuffer;
};

type RenderUIData = RenderData & UIData;

let spritesheet: WorkerUiInit["sprites"];
let canvas: OffscreenCanvas;
let ctx: OffscreenCanvasRenderingContext2D;
let pixel_ratio: number;
let is_drawing = false;

self.addEventListener("message", onMessage);

async function onMessage(event: MessageEvent<WorkerUiInit | RenderUIData | ResetData>) {
  const { data } = event;

  switch (data.type) {
    case "init":
      return init(data).then(self.postMessage);
    case "render":
      if (!is_drawing && data.ui) {
        is_drawing = true;
        return render(data);
      }
    case "reset":
      _.clearCanvas(ctx, data.camera, pixel_ratio);
  }
}

async function init(data: WorkerUiInit): Promise<ReadyMessage> {
  canvas = data.offscreen_canvas;
  ctx = canvas.getContext("2d")!;
  ctx.imageSmoothingEnabled = false;
  ctx.fillStyle = "rgba(255, 0, 0, 0.7)";

  spritesheet = data.sprites;
  pixel_ratio = data.pixel_ratio;

  return { ready: true };
}

function render(data: RenderUIData) {
  const qs = new Uint8Array(data.qs_buffer);
  const rs = new Uint8Array(data.rs_buffer);
  const item_ids = new Uint8Array(data.item_ids_buffer);
  const occupied = new Uint8Array(data.occupied_buffer);
  const sprite = _getSpriteKey(data.ui, item_ids[0]);

  _.clearCanvas(ctx, data.camera, pixel_ratio);

  if (data.ui === "ui_destroy_artifact") {
    ctx.filter = "grayscale(1)";
  }

  // For now a lot of artifacts have a sprite of one tile dimension,
  // so it's normal to see only one tile displayed (expect for free-size artifacts).
  for (let index = 0; index <= qs.length; index++) {
    const center = _.hexToPixel(qs[index], rs[index]);

    ctx.drawImage(
      spritesheet[sprite as keyof typeof spritesheet],
      center.x - C.HEX_WIDTH / 2,
      center.y - C.HEX_HEIGHT / 2,
    );

    if (occupied[index] > 0) {
      drawOccupiedHex(center);
    }
  }

  is_drawing = false;
}

function drawOccupiedHex(center: PixelCoordinates) {
  ctx.beginPath();
  ctx.moveTo(center.x, center.y - 8);
  ctx.lineTo(center.x + 8, center.y - 4);
  ctx.lineTo(center.x + 8, center.y + 4);
  ctx.lineTo(center.x, center.y + 8);
  ctx.lineTo(center.x - 8, center.y + 4);
  ctx.lineTo(center.x - 8, center.y - 4);
  ctx.closePath();

  ctx.fill();
}

function _getSpriteKey(ui: UIType, item_id: number) {
  switch (ui) {
    case "ui_zones":
      return SUPER_KINDS_BY_ID[item_id as SuperKindId];
    case "ui_build_artifact":
      return ARTIFACTS_BY_ID[item_id as ArtifactId];
    case "ui_destroy_artifact":
      return ARTIFACTS_BY_ID[item_id as ArtifactId];
  }
}

export type { UIData, RenderUIData, WorkerUiInit };
