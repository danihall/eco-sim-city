import * as C from "../../../../../constants/_index.js";
function* visibleTilesIterator(context, camera) {
    const canvas = context.canvas;
    const visibility = {
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
function computeMinR(camera_y) {
    return Math.max(C.GRID_MIN, Math.floor((camera_y - C.MAX_ARTIFACT_HEIGHT - C.GRID_OFFSET_Y) / C.SPACING_R));
}
function computeMaxR(camera_y, visible_height) {
    return Math.min(C.GRID_MAX, Math.ceil((camera_y + visible_height - C.GRID_OFFSET_Y) / C.SPACING_R) + 1);
}
function computeMinQ(camera_x, r) {
    return Math.max(C.GRID_MIN, Math.floor((camera_x - C.GRID_OFFSET_X - (C.SPACING_Q / 2) * r) / C.SPACING_Q));
}
function computeMaxQ(camera_x, r, visible_witdh) {
    return Math.min(C.GRID_MAX, Math.ceil((camera_x + visible_witdh - C.GRID_OFFSET_X - (C.SPACING_Q / 2) * r) / C.SPACING_Q) +
        1);
}
function hexToPixel(q, r) {
    return {
        x: C.SPACING_Q * q + (C.SPACING_Q / 2) * r + C.GRID_OFFSET_X,
        y: C.SPACING_R * r + C.GRID_OFFSET_Y,
    };
}
function isOutOfView(camera, center, visibility) {
    return (center.x < camera.x - C.HEX_WIDTH ||
        center.x > camera.x + visibility.width + C.HEX_WIDTH ||
        center.y < camera.y - C.MAX_ARTIFACT_HEIGHT ||
        center.y > camera.y + visibility.height + C.HEX_HEIGHT);
}
function clearCanvas(context, camera, pixel_ratio) {
    const canvas = context.canvas;
    const effective_zoom = camera.zoom * pixel_ratio;
    context.resetTransform();
    context.clearRect(0, 0, canvas.width, canvas.height);
    context.setTransform(effective_zoom, 0, 0, effective_zoom, -camera.x * camera.zoom * pixel_ratio, -camera.y * camera.zoom * pixel_ratio);
}
export { computeMinR, computeMaxR, computeMinQ, computeMaxQ, hexToPixel, clearCanvas, visibleTilesIterator, isOutOfView, };
