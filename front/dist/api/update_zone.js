import { getCityID } from "utils";
async function zoneSuperKind(payload) {
    const method = "PATCH";
    const headers = { "Content-Type": "application/json" };
    const url = `/city/${getCityID()}/zoning/update`;
    const body = JSON.stringify(payload);
    const response = await fetch(url, { method, headers, body });
    if (!response.ok) {
        const reason = await response.text();
        throw new Error(`${response.status}\n${reason}`);
    }
    const updated_tiles = (await response.json());
    if (!Array.isArray(updated_tiles) ||
        updated_tiles[0].artifact?.status !== "Step_1_AwaitingAssignment") {
        throw Error(`
      status should be Step_1_AwaitingAssignment.
      received ${updated_tiles[0].artifact?.status}.
      `);
    }
    return updated_tiles;
}
export { zoneSuperKind };
