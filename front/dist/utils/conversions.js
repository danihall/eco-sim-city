import * as C from "../constants/_index.js";
function cubeRound(q_frac, r_frac) {
    const x_frac = q_frac;
    const z_frac = r_frac;
    const y_frac = -x_frac - z_frac;
    let x = Math.round(x_frac);
    let y = Math.round(y_frac);
    let z = Math.round(z_frac);
    const x_diff = Math.abs(x - x_frac);
    const y_diff = Math.abs(y - y_frac);
    const z_diff = Math.abs(z - z_frac);
    if (x_diff > y_diff && x_diff > z_diff) {
        x = -y - z;
    }
    else if (y_diff > z_diff) {
        y = -x - z;
    }
    else {
        z = -x - y;
    }
    return { q: x, r: z };
}
function pixelToHex(world_x, world_y) {
    const r_frac = (world_y - C.GRID_OFFSET_Y) / C.SPACING_R;
    const q_frac = (world_x - C.GRID_OFFSET_X - (C.SPACING_Q / 2) * r_frac) / C.SPACING_Q;
    const { q, r } = cubeRound(q_frac, r_frac);
    const ok = q >= C.GRID_MIN && q <= C.GRID_MAX && r >= C.GRID_MIN && r <= C.GRID_MAX;
    return { q, r, ok };
}
function getPositionFromPointer(event, camera) {
    const x = event.clientX / camera.zoom + camera.x;
    const y = event.clientY / camera.zoom + camera.y;
    return pixelToHex(x, y);
}
function computeDistance(origin, destination) {
    const q_diff = origin.q - destination.q;
    const r_diff = origin.r - destination.r;
    return (Math.abs(q_diff) + Math.abs(r_diff) + Math.abs(q_diff + r_diff)) / 2;
}
function getTileIndex(q, r) {
    return q + r * 128;
}
export { cubeRound, pixelToHex, getTileIndex, getPositionFromPointer, computeDistance };
