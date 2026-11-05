import type { ArtifactName } from "assets";
import { getCityID } from "utils";

type BuildPathwayPayload = {
  artifact_name: ArtifactName;
  coordinates: Coordinates[];
};

async function buildPathway(payload: BuildPathwayPayload): Promise<Tile[]> {
  const method = "PATCH";
  const headers = { "Content-Type": "application/json" };
  const url = `/city/${getCityID()}/construction/build-pathway`;
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
    first_tile.super_kind !== "PublicSystems" ||
    first_tile?.artifact?.name !== payload.artifact_name ||
    updated_tiles.length !== payload.coordinates.length ||
    first_tile.artifact.status !== "Step_3_DeliveredOperational"
  ) {
    throw TypeError(`
      tile artifact should be: ${payload.artifact_name}.
      received: ${first_tile?.artifact?.name}.
      updated tiles count should be: ${payload.coordinates.length}.
      received ${updated_tiles.length}.
      tile artifact status should be: Step_3_DeliveredOperational.
      received: ${first_tile?.artifact?.status}.
      `);
  }

  return updated_tiles as Tile[];
}

export type { BuildPathwayPayload };
export { buildPathway };
