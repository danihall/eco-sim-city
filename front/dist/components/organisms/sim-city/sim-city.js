import { HEX_HEIGHT, HEX_WIDTH, TILES_COUNT } from "constants";
import { BaseComponent } from "components";
import * as api from "api";
import * as $ from "assets";
import * as _ from "utils";
import css from "./sim-city.css" with { type: "css" };
const pointer_events = [
    "pointerenter",
    "pointerdown",
    "pointermove",
    "pointerup",
    "pointerleave",
    "click",
];
class SimCity extends BaseComponent {
    static styles = [css];
    #camera = { x: 0, y: 0, zoom: 1 };
    #shared_buffers = {
        artifact_ids: new SharedArrayBuffer(TILES_COUNT * 1),
        anchor_qs: new SharedArrayBuffer(TILES_COUNT * 1),
        anchor_rs: new SharedArrayBuffer(TILES_COUNT * 1),
        artifact_statuses: new SharedArrayBuffer(TILES_COUNT * 1),
        zones: new SharedArrayBuffer(TILES_COUNT * 1),
        foot_traffic: new SharedArrayBuffer(TILES_COUNT * Uint32Array.BYTES_PER_ELEMENT),
        car_traffic: new SharedArrayBuffer(TILES_COUNT * Uint32Array.BYTES_PER_ELEMENT),
        bus_traffic: new SharedArrayBuffer(TILES_COUNT * Uint32Array.BYTES_PER_ELEMENT),
        bicycle_traffic: new SharedArrayBuffer(TILES_COUNT * Uint32Array.BYTES_PER_ELEMENT),
        train_traffic: new SharedArrayBuffer(TILES_COUNT * Uint32Array.BYTES_PER_ELEMENT),
        subway_traffic: new SharedArrayBuffer(TILES_COUNT * Uint32Array.BYTES_PER_ELEMENT),
    };
    #state = {
        artifact_ids: new Uint8Array(this.#shared_buffers.artifact_ids),
        anchor_qs: new Uint8Array(this.#shared_buffers.anchor_qs),
        anchor_rs: new Uint8Array(this.#shared_buffers.anchor_rs),
        artifact_statuses: new Uint8Array(this.#shared_buffers.artifact_statuses),
        zones: new Uint8Array(this.#shared_buffers.zones),
        foot_traffic: new Uint32Array(this.#shared_buffers.foot_traffic),
        car_traffic: new Uint32Array(this.#shared_buffers.car_traffic),
        bus_traffic: new Uint32Array(this.#shared_buffers.bus_traffic),
        bicycle_traffic: new Uint32Array(this.#shared_buffers.bicycle_traffic),
        train_traffic: new Uint32Array(this.#shared_buffers.train_traffic),
        subway_traffic: new Uint32Array(this.#shared_buffers.subway_traffic),
    };
    render() {
        this.html `<slot></slot>`;
    }
    async onComponentMounted() {
        try {
            const tiles_bytes = await api.getAllTiles();
            this.#state.artifact_ids.set(tiles_bytes.subarray(0, TILES_COUNT));
            this.#state.anchor_qs.set(tiles_bytes.subarray(TILES_COUNT, TILES_COUNT * 2));
            this.#state.anchor_rs.set(tiles_bytes.subarray(TILES_COUNT * 2, TILES_COUNT * 3));
            this.#state.artifact_statuses.set(tiles_bytes.subarray(TILES_COUNT * 3, TILES_COUNT * 4));
            this.#state.zones.set(tiles_bytes.subarray(TILES_COUNT * 4, TILES_COUNT * 5));
            this.ownerDocument.addEventListener("keydown", this);
            pointer_events.forEach((event_name) => this.addEventListener(event_name, this));
            this.channel.subscribe("update-city-state", this);
            this.channel.subscribe("update-city-zones", this);
        }
        catch (reason) {
            console.error(reason);
        }
    }
    handleEvent(event) {
        if (pointer_events.includes(event.type)) {
            return this.channel.publish("city-pointer-event", event);
        }
        super.handleEvent(event);
    }
    onKeydown(event) {
        let should_move = false;
        switch (event.key) {
            case "ArrowUp":
                this.#camera.y -= HEX_HEIGHT * 2;
                should_move = true;
                break;
            case "ArrowRight":
                this.#camera.x += HEX_WIDTH * 2;
                should_move = true;
                break;
            case "ArrowDown":
                this.#camera.y += HEX_HEIGHT * 2;
                should_move = true;
                break;
            case "ArrowLeft":
                this.#camera.x -= HEX_WIDTH * 2;
                should_move = true;
                break;
        }
        if (should_move) {
            this.channel.publish("grid-render");
            this.channel.publish("render");
        }
    }
    onUpdateCityState(updated_tiles) {
        const [first_tile] = updated_tiles;
        const artifact_id = $.ARTIFACTS_BY_NAME[first_tile.artifact.name];
        const super_kind_id = $.SUPER_KINDS_BY_NAME[first_tile.super_kind];
        const status_id = $.TILE_STATUS_BY_NAME[first_tile.artifact.status];
        const display_props = $.ARTIFACTS_SIZES[first_tile.artifact.name];
        const has_free_size = display_props.hex_height === 0 && display_props.hex_width === 0;
        updated_tiles.forEach((tile) => {
            const anchor_index = _.getTileIndex(tile.anchor.q, tile.anchor.r);
            const tile_index = _.getTileIndex(tile.position.q, tile.position.r);
            this.#state.artifact_ids[tile_index] =
                tile_index === anchor_index || has_free_size ? artifact_id : 0;
            this.state.anchor_qs[tile_index] = tile.anchor.q;
            this.state.anchor_rs[tile_index] = tile.anchor.r;
            this.#state.artifact_statuses[tile_index] = status_id;
            this.#state.zones[tile_index] = super_kind_id;
        });
        this.channel.publish("render");
    }
    onUpdateCityZones(params) {
        const { updated_tiles, super_kind } = params;
        const super_kind_id = $.SUPER_KINDS_BY_NAME[super_kind];
        const status_id = $.TILE_STATUS_BY_NAME[updated_tiles[0].artifact.status];
        updated_tiles.forEach((tile) => {
            const tile_index = _.getTileIndex(tile.position.q, tile.position.r);
            this.#state.artifact_statuses[tile_index] = status_id;
            this.#state.zones[tile_index] = super_kind_id;
        });
        this.channel.publish("render");
    }
    get state() {
        return this.#state;
    }
    get sharedBuffers() {
        return this.#shared_buffers;
    }
    get camera() {
        return this.#camera;
    }
}
window.customElements.define("sim-city", SimCity);
export { SimCity };
