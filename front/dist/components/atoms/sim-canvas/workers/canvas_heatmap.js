import * as _ from "../../../../utils/_index.js";
import * as C from "../../../../constants/_index.js";
import * as helpers from "./helpers/helpers.js";
let spritesheet;
let canvas;
let ctx;
let pixel_ratio;
let foot_traffic;
let bicycle_traffic;
let car_traffic;
let bus_traffic;
let subway_traffic;
let train_traffic;
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
    foot_traffic = new Uint32Array(data.shared_buffers.foot_traffic);
    bicycle_traffic = new Uint32Array(data.shared_buffers.bicycle_traffic);
    car_traffic = new Uint32Array(data.shared_buffers.car_traffic);
    bus_traffic = new Uint32Array(data.shared_buffers.bus_traffic);
    subway_traffic = new Uint32Array(data.shared_buffers.subway_traffic);
    train_traffic = new Uint32Array(data.shared_buffers.train_traffic);
    spritesheet = data.sprites;
    pixel_ratio = data.pixel_ratio;
    return { ready: true };
}
function render(data) {
    helpers.clearCanvas(ctx, data.camera, pixel_ratio);
    for (const { q, r, center } of helpers.visibleTilesIterator(ctx, data.camera)) {
        const tile_index = _.getTileIndex(q, r);
        drawTraffic(tile_index, center, data.data_type);
    }
}
function drawTraffic(tile_index, center, data_type) {
    const value = _getTrafficValue(tile_index, data_type);
    if (value <= 0) {
        return;
    }
    const congestion = _getCongestion(value);
    ctx.drawImage(spritesheet[congestion], center.x - C.HEX_WIDTH / 2, center.y - C.HEX_HEIGHT / 2);
}
function _getTrafficValue(tile_index, data_type) {
    switch (data_type) {
        case "foot_traffic":
            return foot_traffic[tile_index];
        case "bicycle_traffic":
            return bicycle_traffic[tile_index];
        case "car_traffic":
            return car_traffic[tile_index];
        case "bus_traffic":
            return bus_traffic[tile_index];
        case "subway_traffic":
            return subway_traffic[tile_index];
        case "train_traffic":
            return train_traffic[tile_index];
        default:
            return -1;
    }
}
function _getCongestion(traffic_value) {
    if (traffic_value < 5) {
        return "very_low";
    }
    if (traffic_value < 15) {
        return "low";
    }
    if (traffic_value < 35) {
        return "medium";
    }
    if (traffic_value < 55) {
        return "high";
    }
    return "very_high";
}
