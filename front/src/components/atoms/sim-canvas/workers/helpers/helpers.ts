import type { CameraState } from "components";
import type { Visibility } from "../_index.js";
import * as C from "../../../../../constants/_index.js";
import { getTileIndex } from "../../../../../utils/_index.js";

type VisibleTile = {
  q: number;
  r: number;
  center: PixelCoordinates;
};

function* visibleTilesIterator(
  context: OffscreenCanvasRenderingContext2D,
  camera: CameraState,
): Generator<VisibleTile> {
  const canvas = context.canvas;
  const visibility: Visibility = {
    width: canvas.width / camera.zoom,
    height: canvas.height / camera.zoom,
  };

  const min_r = computeMinR(camera.y);
  const max_r = computeMaxR(camera.y, visibility.height);

  for (let r = min_r; r <= max_r; r++) {
    const min_q = computeMinQ(camera.x, r);
    const max_q = computeMaxQ(camera.x, r, visibility.width);

    for (let q = min_q; q <= max_q; q++) {
      const center = hexToPixel(q, r);

      if (isOutOfView(camera, center, visibility)) {
        continue;
      }

      yield { q, r, center };
    }
  }
}

function computeMinR(camera_y: number) {
  return Math.max(
    C.GRID_MIN,
    Math.floor((camera_y - C.MAX_ARTIFACT_HEIGHT - C.GRID_OFFSET_Y) / C.SPACING_R),
  );
}

function computeMaxR(camera_y: number, visible_height: number) {
  return Math.min(
    C.GRID_MAX,
    Math.ceil((camera_y + visible_height - C.GRID_OFFSET_Y) / C.SPACING_R) + 1,
  );
}

function computeMinQ(camera_x: number, r: number) {
  return Math.max(
    C.GRID_MIN,
    Math.floor((camera_x - C.GRID_OFFSET_X - (C.SPACING_Q / 2) * r) / C.SPACING_Q),
  );
}

function computeMaxQ(camera_x: number, r: number, visible_witdh: number) {
  return Math.min(
    C.GRID_MAX,
    Math.ceil((camera_x + visible_witdh - C.GRID_OFFSET_X - (C.SPACING_Q / 2) * r) / C.SPACING_Q) +
      1,
  );
}

function hexToPixel(q: number, r: number): PixelCoordinates {
  return {
    x: C.SPACING_Q * q + (C.SPACING_Q / 2) * r + C.GRID_OFFSET_X,
    y: C.SPACING_R * r + C.GRID_OFFSET_Y,
  };
}

function isOutOfView(camera: CameraState, center: PixelCoordinates, visibility: Visibility) {
  return (
    center.x < camera.x - C.HEX_WIDTH ||
    center.x > camera.x + visibility.width + C.HEX_WIDTH ||
    center.y < camera.y - C.MAX_ARTIFACT_HEIGHT ||
    center.y > camera.y + visibility.height + C.HEX_HEIGHT
  );
}

function clearCanvas(
  context: OffscreenCanvasRenderingContext2D,
  camera: CameraState,
  pixel_ratio: number,
) {
  const canvas = context.canvas;
  const effective_zoom = camera.zoom * pixel_ratio;

  context.resetTransform();
  context.clearRect(0, 0, canvas.width, canvas.height);
  context.setTransform(
    effective_zoom,
    0,
    0,
    effective_zoom,
    -camera.x * camera.zoom * pixel_ratio,
    -camera.y * camera.zoom * pixel_ratio,
  );
}

export {
  computeMinR,
  computeMaxR,
  computeMinQ,
  computeMaxQ,
  hexToPixel,
  clearCanvas,
  visibleTilesIterator,
  isOutOfView,
};
