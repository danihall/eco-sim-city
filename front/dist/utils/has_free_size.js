import * as $ from "../assets/_index.js";
function hasFreeSize(artifact_id) {
    const artifact_name = $.ARTIFACTS_BY_ID[artifact_id];
    const size_props = $.ARTIFACTS_SIZES[artifact_name];
    return size_props.hex_height === 0 && size_props.hex_width === 0;
}
export { hasFreeSize };
