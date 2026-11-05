import type { InitData, ReadyMessage, RenderData } from "./_index.js";

import * as _ from "./helpers/helpers.js";

type WorkerGridInit = Omit<InitData, "sprites" | "shared_buffers">;

let canvas: OffscreenCanvas;
let ctx: OffscreenCanvasRenderingContext2D;
let pixel_ratio: number;
let is_drawing = false;

self.addEventListener("message", onMessage);

function onMessage(event: MessageEvent<InitData | RenderData>) {
  const { data } = event;

  switch (data.type) {
    case "init":
      return init(data).then(self.postMessage);
    case "render":
      return render(data);
  }
}

async function init(data: InitData): Promise<ReadyMessage> {
  canvas = data.offscreen_canvas;
  ctx = canvas.getContext("2d")!;
  ctx.imageSmoothingEnabled = false;
  ctx.strokeStyle = "#989898";
  ctx.lineWidth = 1;

  pixel_ratio = data.pixel_ratio;

  return { ready: true };
}

function render(data: RenderData) {
  _.clearCanvas(ctx, data.camera, pixel_ratio);

  for (const { center } of _.visibleTilesIterator(ctx, data.camera)) {
    drawHexLines(center);
  }
}

const half_hex = 7.5;
const half_side = 3.5;

function drawHexLines(center: PixelCoordinates) {
  ctx.beginPath();
  ctx.moveTo(center.x, center.y - half_hex);
  ctx.lineTo(center.x + half_hex, center.y - half_side);
  ctx.lineTo(center.x + half_hex, center.y + half_side);
  ctx.lineTo(center.x, center.y + half_hex);
  ctx.lineTo(center.x - half_hex, center.y + half_side);
  ctx.lineTo(center.x - half_hex, center.y - half_side);
  ctx.closePath();

  ctx.stroke();
}

export type { WorkerGridInit };
