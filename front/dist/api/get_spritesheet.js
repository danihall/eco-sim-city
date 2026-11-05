async function getSpritesheet() {
    const method = "GET";
    const cache = "no-cache";
    const headers = { "Content-Type": "image/png" };
    const url = new URL("../assets/spritesheet.png", import.meta.url);
    const response = await fetch(url, { method, headers, cache });
    if (!response.ok) {
        const reason = await response.text();
        throw new Error(`${response.status}\n${reason}`);
    }
    return response.blob();
}
export { getSpritesheet };
