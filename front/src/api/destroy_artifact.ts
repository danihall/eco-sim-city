import type { ArtifactName } from "assets";
import { getCityID } from "utils";

type DestroyArtifactPayload = {
  position: Coordinates;
};

async function destroyArtifact(payload: DestroyArtifactPayload): Promise<Tile[]> {
  const method = "PATCH";
  const headers = { "Content-Type": "application/json" };
  const url = `/city/${getCityID()}/construction/destroy`;
  const body = JSON.stringify(payload);

  const response = await fetch(url, { method, headers, body });

  if (!response.ok) {
    const reason = await response.text();
    throw new Error(`${response.status}\n${reason}`);
  }

  const updated_tiles = (await response.json()) as Partial<Tile>[];
  const [first_tile] = updated_tiles;

  if (
    !Array.isArray(updated_tiles) ||
    first_tile.artifact?.name !== "Wasteland" ||
    first_tile?.artifact?.status !== "Step_1_AwaitingAssignment"
  ) {
    throw TypeError(`
      tile artifact name should be: Wasteland.
      received: ${first_tile.artifact?.name}.
      tile artifact status should be: Step_1_AwaitingAssignment.
      received: ${first_tile?.artifact?.status}.
      `);
  }

  return updated_tiles as Tile[];
}

export type { DestroyArtifactPayload };
export { destroyArtifact };
