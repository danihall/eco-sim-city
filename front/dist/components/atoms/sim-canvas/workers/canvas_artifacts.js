import * as _ from "../../../../utils/_index.js";
import * as C from "../../../../constants/_index.js";
import * as helpers from "./helpers/helpers.js";
import { ARTIFACTS_BY_ID } from "../../../../assets/_index.js";
let spritesheet;
let canvas;
let ctx;
let pixel_ratio;
let artifact_ids;
self.addEventListener("message", onMessage);
async function onMessage(event) {
    const { data } = event;
    switch (data.type) {
        case "init":
            return init(data).then(self.postMessage);
        case "render":
            return render(data);
    }
}
async function init(data) {
    canvas = data.offscreen_canvas;
    ctx = canvas.getContext("2d");
    ctx.imageSmoothingEnabled = false;
    pixel_ratio = data.pixel_ratio;
    artifact_ids = new Uint8Array(data.shared_buffers.artifact_ids);
    spritesheet = data.sprites;
    return { ready: true };
}
function render(data) {
    helpers.clearCanvas(ctx, data.camera, pixel_ratio);
    for (const { q, r, center } of helpers.visibleTilesIterator(ctx, data.camera)) {
        const tile_index = _.getTileIndex(q, r);
        const artifact_id = artifact_ids[tile_index];
        if (artifact_id > 0) {
            const artifact_name = ARTIFACTS_BY_ID[artifact_ids[tile_index]];
            ctx.drawImage(spritesheet[artifact_name], center.x - C.HEX_WIDTH / 2, center.y - C.HEX_HEIGHT / 2);
        }
    }
}
