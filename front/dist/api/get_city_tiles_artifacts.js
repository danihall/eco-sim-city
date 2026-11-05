import { getCityID } from "utils";
async function getAllTiles() {
    const method = "GET";
    const headers = { "Content-Type": "application/octet-stream" };
    const url = `/city/${getCityID()}/all-tiles`;
    const response = await fetch(url, { method, headers });
    if (!response.ok) {
        const reason = await response.text();
        throw new Error(`${response.status}\n${reason}`);
    }
    const buffer = await response.arrayBuffer();
    return new Uint8Array(buffer);
}
export { getAllTiles };
