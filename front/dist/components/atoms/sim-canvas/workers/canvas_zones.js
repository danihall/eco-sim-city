import { SUPER_KINDS_BY_ID, TILE_STATUS_BY_ID } from "../../../../assets/_index.js";
import * as _ from "../../../../utils/_index.js";
import * as C from "../../../../constants/_index.js";
import * as helpers from "./helpers/helpers.js";
let spritesheet;
let canvas;
let ctx;
let pixel_ratio;
let is_drawing = false;
let artifact_statuses;
let zones;
self.addEventListener("message", onMessage);
async function onMessage(event) {
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
async function init(data) {
    canvas = data.offscreen_canvas;
    ctx = canvas.getContext("2d");
    ctx.imageSmoothingEnabled = false;
    artifact_statuses = new Uint8Array(data.shared_buffers.artifact_statuses);
    zones = new Uint8Array(data.shared_buffers.zones);
    spritesheet = data.sprites;
    pixel_ratio = data.pixel_ratio;
    return { ready: true };
}
function render(data) {
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
function getZone(tile_index) {
    const zone_type = SUPER_KINDS_BY_ID[zones[tile_index]];
    const status = TILE_STATUS_BY_ID[artifact_statuses[tile_index]];
    if ((zone_type === "Industrial" || zone_type === "Residential") &&
        status === "Step_1_AwaitingAssignment") {
        return zone_type;
    }
}
function drawZone(zone, center) {
    ctx.drawImage(spritesheet[zone], center.x - C.HEX_WIDTH / 2, center.y - C.HEX_HEIGHT / 2);
}
