import { SUPER_KINDS_BY_ID, TILE_STATUS_BY_ID, ARTIFACTS_BY_ID } from "assets";
function isOccupiedTile(tile) {
    const { tile_index, city_state } = tile;
    const super_kind = SUPER_KINDS_BY_ID[city_state.zones[tile_index]];
    const status = TILE_STATUS_BY_ID[city_state.artifact_statuses[tile_index]];
    const artifact = ARTIFACTS_BY_ID[city_state.artifact_ids[tile_index]];
    const is_occupied = (super_kind !== "Natural" && status !== "Step_1_AwaitingAssignment") ||
        super_kind === "PublicSystems" ||
        artifact === "DeepWater" ||
        artifact === "ShallowWater";
    return is_occupied;
}
/**
 * Create typedArrays to be transmitted to a canvas Worker.
 */
function newTilesBuffers(capacity) {
    const qs_buffer = new ArrayBuffer(capacity);
    const rs_buffer = new ArrayBuffer(capacity);
    const item_ids_buffer = new ArrayBuffer(capacity);
    const occupied_buffer = new ArrayBuffer(capacity);
    return {
        qs: new Uint8Array(qs_buffer),
        rs: new Uint8Array(rs_buffer),
        item_ids: new Uint8Array(item_ids_buffer),
        occupied: new Uint8Array(occupied_buffer),
    };
}
/**
 * To be used in any iteration.
 */
function updateBuffers(tile, index) {
    const { qs, rs, item_ids, occupied, artifact_id, city_state } = this;
    const is_occupied = isOccupiedTile({ ...tile, city_state });
    qs[index] = tile.q;
    rs[index] = tile.r;
    item_ids[index] = artifact_id;
    occupied[index] = Number(is_occupied);
}
export { isOccupiedTile, newTilesBuffers, updateBuffers };
