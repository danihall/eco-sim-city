import type { ArtifactName } from "assets";
import { getCityID } from "utils";

type BuildSingleArtifactPayload = {
  anchor: Coordinates;
  artifact_name: ArtifactName;
  width: number;
  height: number;
};

async function buildSingleArtifact(payload: BuildSingleArtifactPayload): Promise<Tile[]> {
  const method = "PATCH";
  const headers = { "Content-Type": "application/json" };
  const url = `/city/${getCityID()}/construction/build`;
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
    (first_tile.super_kind !== "PublicSystems" && first_tile.super_kind !== "Natural") ||
    first_tile?.artifact?.name !== payload.artifact_name ||
    first_tile.artifact.status !== "Step_3_DeliveredOperational" ||
    updated_tiles.every(
      (tile) => tile.anchor?.q !== payload.anchor.q && tile.anchor?.r !== payload.anchor.r,
    )
  ) {
    throw TypeError(`
      tile super_kind should be "PublicSystems" || "Natural".
      received: ${first_tile.super_kind}.
      tile artifact should be: ${payload.artifact_name}.
      received: ${first_tile?.artifact?.name}.
      tile artifact status should be: Step_3_DeliveredOperational.
      received: ${first_tile?.artifact?.status}.
      no tiles correspond to payload.anchor.
      `);
  }

  return updated_tiles as Tile[];
}

export type { BuildSingleArtifactPayload };
export { buildSingleArtifact };
