import type { SuperKindId, TileStatusId, ArtifactId } from "assets";
import type { TilePositionWithIndex } from "utils";
import type { CityData, Concrete } from "components";
import { SUPER_KINDS_BY_ID, TILE_STATUS_BY_ID, ARTIFACTS_BY_ID } from "assets";

type TileOccupancy = Coordinates & { is_occupied: boolean };

type TileSelectionContext = {
  qs: Uint8Array<ArrayBuffer>;
  rs: Uint8Array<ArrayBuffer>;
  item_ids: Uint8Array<ArrayBuffer>;
  occupied: Uint8Array<ArrayBuffer>;
  artifact_id: number;
  city_state: CityData<Concrete>;
};

function isOccupiedTile(tile: TilePositionWithIndex & { city_state: CityData<Concrete> }) {
  const { tile_index, city_state } = tile;
  const super_kind = SUPER_KINDS_BY_ID[city_state.zones[tile_index] as SuperKindId];
  const status = TILE_STATUS_BY_ID[city_state.artifact_statuses[tile_index] as TileStatusId];
  const artifact = ARTIFACTS_BY_ID[city_state.artifact_ids[tile_index] as ArtifactId];

  const is_occupied =
    (super_kind !== "Natural" && status !== "Step_1_AwaitingAssignment") ||
    super_kind === "PublicSystems" ||
    artifact === "DeepWater" ||
    artifact === "ShallowWater";

  return is_occupied;
}

/**
 * Create typedArrays to be transmitted to a canvas Worker.
 */
function newTilesBuffers(capacity: number) {
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
function updateBuffers(this: TileSelectionContext, tile: TilePositionWithIndex, index: number) {
  const { qs, rs, item_ids, occupied, artifact_id, city_state } = this;
  const is_occupied = isOccupiedTile({ ...tile, city_state });

  qs[index] = tile.q;
  rs[index] = tile.r;
  item_ids[index] = artifact_id;
  occupied[index] = Number(is_occupied);
}

export type { TileOccupancy, TileSelectionContext };
export { isOccupiedTile, newTilesBuffers, updateBuffers };
