import type { SuperKindId, SUPER_KINDS_BY_ID } from "assets";
import { getCityID } from "utils";

type ZoneSuperKindPayload = {
  super_kind_name: (typeof SUPER_KINDS_BY_ID)[SuperKindId];
  top_left: Coordinates;
  bottom_right: Coordinates;
};

async function zoneSuperKind(payload: ZoneSuperKindPayload): Promise<Tile[]> {
  const method = "PATCH";
  const headers = { "Content-Type": "application/json" };
  const url = `/city/${getCityID()}/zoning/update`;
  const body = JSON.stringify(payload);

  const response = await fetch(url, { method, headers, body });

  if (!response.ok) {
    const reason = await response.text();
    throw new Error(`${response.status}\n${reason}`);
  }

  const updated_tiles = (await response.json()) as Partial<Tile>[];

  if (
    !Array.isArray(updated_tiles) ||
    updated_tiles[0].artifact?.status !== "Step_1_AwaitingAssignment"
  ) {
    throw Error(`
      status should be Step_1_AwaitingAssignment.
      received ${updated_tiles[0].artifact?.status}.
      `);
  }

  return updated_tiles as Tile[];
}

export type { ZoneSuperKindPayload };
export { zoneSuperKind };
