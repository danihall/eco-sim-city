import { setCityID } from "utils";
// WIP: for now, hard-code it
setCityID(1);
//     new EventSource("http://localhost:3000/city/1/sse").addEventListener(
//       "message",
//       onServerMessage,
//     );
// function onServerMessage(event: MessageEvent) {
//   const { game_loop_name, body } = JSON.parse(event.data) as SSEPayload;
//   console.log(game_loop_name);
//   // TODO refactor later because back-end will send differently shaped payload.
//   // payload will be shaped like "all-tiles" endpoint with JSON arrays instead of arraybuffers.
//   if (game_loop_name === "artifacts_placing_loop") {
//     console.log(body);
//     body.forEach((item: any) => {
//       const artifact_id = ARTIFACTS_BY_NAME[item.artifact.name];
//       const super_kind_id = SUPER_KINDS_BY_NAME[item.super_kind_type];
//       const status_id = TILE_STATUS_BY_NAME[item.artifact.status];
//       if (!artifact_id || !super_kind_id || !status_id) {
//         console.error(`error in ${game_loop_name}`);
//         return;
//       }
//       const anchor_tile_index = getTileIndex(item.anchor.q, item.anchor.r);
//       item.coordinates.forEach((position: { q: number; r: number }) => {
//         const tile_index = getTileIndex(position.q, position.r);
//         city_state.artifacts_ids[tile_index] = tile_index === anchor_tile_index ? artifact_id : 0;
//         city_state.artifact_statuses[tile_index] = status_id;
//         city_state.zones[tile_index] = super_kind_id;
//       });
//     });
//   }
//   updateRender("tiles");
// }
