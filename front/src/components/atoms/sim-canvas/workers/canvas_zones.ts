import type { InitData, ReadyMessage, RenderData } from "./_index.js";
import type { SuperKindId, SuperKindName, TileStatusId } from "../../../../assets/_index.js";

import { SUPER_KINDS_BY_ID, TILE_STATUS_BY_ID } from "../../../../assets/_index.js";
import * as _ from "../../../../utils/_index.js";
import * as C from "../../../../constants/_index.js";
import * as helpers from "./helpers/helpers.js";

type WorkerZonesInit = InitData<"artifact_statuses" | "zones", "zones">;

let spritesheet: WorkerZonesInit["sprites"];
let canvas: OffscreenCanvas;
let ctx: OffscreenCanvasRenderingContext2D;
let pixel_ratio: number;
let is_drawing = false;

let artifact_statuses: Uint8Array;
let zones: Uint8Array;

self.addEventListener("message", onMessage);

async function onMessage(event: MessageEvent<WorkerZonesInit | RenderData>) {
  const { data } = event;

  switch (data.type) {
    case "init":
      return init(data).then(self.postMessage);
    case "render":
      if (!is_drawing) {
        is_drawing = true;
        return render(data);
      }
  }
}

async function init(data: WorkerZonesInit): Promise<ReadyMessage> {
  canvas = data.offscreen_canvas;
  ctx = canvas.getContext("2d")!;
  ctx.imageSmoothingEnabled = false;

  artifact_statuses = new Uint8Array(data.shared_buffers.artifact_statuses);
  zones = new Uint8Array(data.shared_buffers.zones);

  spritesheet = data.sprites;
  pixel_ratio = data.pixel_ratio;

  return { ready: true };
}

function render(data: RenderData) {
  helpers.clearCanvas(ctx, data.camera, pixel_ratio);

  for (const { q, r, center } of helpers.visibleTilesIterator(ctx, data.camera)) {
    const tile_index = _.getTileIndex(q, r);
    const zone = getZone(tile_index);

    if (zone) {
      drawZone(zone, center);
    }
  }

  is_drawing = false;
}

function getZone(tile_index: number) {
  const zone_type = SUPER_KINDS_BY_ID[zones[tile_index] as SuperKindId];
  const status = TILE_STATUS_BY_ID[artifact_statuses[tile_index] as TileStatusId];

  if (
    (zone_type === "Industrial" || zone_type === "Residential") &&
    status === "Step_1_AwaitingAssignment"
  ) {
    return zone_type;
  }
}

function drawZone(zone: SuperKindName, center: PixelCoordinates) {
  ctx.drawImage(spritesheet[zone], center.x - C.HEX_WIDTH / 2, center.y - C.HEX_HEIGHT / 2);
}

export type { WorkerZonesInit };
