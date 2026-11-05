import { getTileIndex } from "./conversions.js";
function* chunkIterator(chunk) {
    const min_q = Math.min(chunk.top_left.q, chunk.bottom_right.q);
    const min_r = Math.min(chunk.top_left.r, chunk.bottom_right.r);
    const max_q = Math.max(chunk.top_left.q, chunk.bottom_right.q);
    const max_r = Math.max(chunk.top_left.r, chunk.bottom_right.r);
    for (let q = min_q; q <= max_q; q++) {
        for (let r = min_r; r <= max_r; r++) {
            const tile_index = getTileIndex(q, r);
            yield { tile_index, q, r };
        }
    }
}
function getOrderedChunk(chunk) {
    const min_q = Math.min(chunk.top_left.q, chunk.bottom_right.q);
    const min_r = Math.min(chunk.top_left.r, chunk.bottom_right.r);
    const max_q = Math.max(chunk.top_left.q, chunk.bottom_right.q);
    const max_r = Math.max(chunk.top_left.r, chunk.bottom_right.r);
    return {
        top_left: {
            q: min_q,
            r: min_r,
        },
        bottom_right: {
            q: max_q,
            r: max_r,
        },
    };
}
function getChunkTilesCount(chunk) {
    const width = Math.abs(chunk.bottom_right.q - chunk.top_left.q) + 1;
    const height = Math.abs(chunk.bottom_right.r - chunk.top_left.r) + 1;
    return width * height;
}
function* pathIterator(path) {
    for (let i = 0; i < path.length; i++) {
        const { q, r } = path[i];
        const tile_index = getTileIndex(q, r);
        yield { tile_index, q, r };
    }
}
export { getOrderedChunk, chunkIterator, getChunkTilesCount, pathIterator };
