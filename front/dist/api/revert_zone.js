import { getCityID } from "utils";
async function revertZone(payload) {
    const method = "PATCH";
    const headers = { "Content-Type": "application/json" };
    const url = `/city/${getCityID()}/zoning/revert`;
    const body = JSON.stringify(payload);
    const response = await fetch(url, { method, headers, body });
    if (!response.ok) {
        const reason = await response.text();
        throw new Error(`${response.status}\n${reason}`);
    }
    const updated_tiles = (await response.json());
    if (!Array.isArray(updated_tiles) ||
        (updated_tiles.length && updated_tiles[0].artifact?.status !== "Step_3_DeliveredOperational")) {
        throw Error(`
      status should be Step_3_DeliveredOperational.
      received ${updated_tiles[0].artifact?.status}.
      `);
    }
    return updated_tiles;
}
export { revertZone };
