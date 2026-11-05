import type { ArtifactsKey } from "constants";
import type { ArtifactName, SuperKindName, TileStatusName } from "assets";
import type { CameraState, CityData, ChannelEvents } from "components";
import type { Chunk } from "utils";
import type { BuildSingleArtifactPayload } from "api";

import { BaseComponent } from "components";
import { SPRITESHEET } from "constants";
import * as $ from "assets";
import * as _ from "utils";
import * as api from "api";
import { newTilesBuffers, updateBuffers } from "../tiles_selection.js";

/**
 * Renderless component that allows user to launch construction
 * of an artifact. The designated artifact is immediately constructed,
 * contrary to "Industrial" | "Residential" artifacts
 * (which are spawned via SSE by the back-end simulating the city).
 */
class SimBuildArtifact extends BaseComponent {
  #camera!: CameraState;
  #artifact_to_build: ArtifactsKey | null = null;
  #anchor: Coordinates | null = null;
  #size!: { width: number; height: number; is_free: boolean };

  #can_tick = true;
  #throttled = 0;

  async onComponentMounted() {
    this.channel.subscribe("menu-state-change", this);
  }

  /**
   * Only activates if user has activated the "build" action from the menu,
   * by selecting a constructible artifact (only "PublicSystem" artifacts allowed).
   */
  async onMenuStateChange({ selected_action, selected_item }: ChannelEvents["menu-state-change"]) {
    const can_build =
      selected_action === "build" &&
      selected_item &&
      Object.hasOwn(SPRITESHEET.artifacts, selected_item);
    if (!can_build) {
      return this.#dispose();
    }

    this.#artifact_to_build = selected_item as ArtifactsKey;

    const { hex_width, hex_height } = $.ARTIFACTS_SIZES[this.#artifact_to_build!];
    this.#size = {
      width: hex_width,
      height: hex_height,
      is_free: hex_width === 0 && hex_height === 0,
    };
    [this.#camera] = await this.channel.request("sim-city", "camera");
    this.channel.subscribe("city-pointer-event", this);
    this.ownerDocument.addEventListener("keydown", this);
  }

  onCityPointerEvent(event: PointerEvent) {
    this.handleEvent(event);
  }

  /**
   * If an artifact has a free size (width || height is 0),
   * user can "drag" a chunk to indicate how much place the artifact
   * should occupy. So let's start by registering the anchor (origin)
   * of the artifact. Later in "pointermove", we'll register the
   * space occupied by the artifact.
   */
  onPointerdown(event: PointerEvent) {
    if (!this.#size.is_free || event.button !== 0) {
      return;
    }

    const { ok, ...position } = _.getPositionFromPointer(event, this.#camera);
    if (!ok) {
      return;
    }

    this.#anchor = position;
  }

  /**
   * Throttles zone selection logic and drawing.
   */
  onPointermove(event: PointerEvent) {
    if (this.#can_tick) {
      this.#can_tick = false;
      this.#throttled = window.setTimeout(this.handlePointerMove.bind(this), 0, event);
    }
  }

  /**
   * If user is attempting to build a "free-size" item,
   * no need to update artifact's anchor.
   * We then send a message to the canvas worker to refresh the UI.
   */
  async handlePointerMove(event: PointerEvent) {
    const { ok, q, r } = _.getPositionFromPointer(event, this.#camera);
    if (!ok) {
      return void (this.#can_tick = true);
    }

    const is_drawing = event.pressure >= 0.5 && this.#size.is_free;

    if (!is_drawing) {
      this.#anchor = { q, r };
    }

    const [city_state] = await this.channel.request("sim-city", "state");
    if (!city_state) {
      return void (this.#can_tick = true);
    }

    if (this.#size.is_free) {
      const chunk = this.#getEffectiveChunk({ q, r });
      const buffer_capacity = _.getChunkTilesCount(chunk);
      const typed_arrays = newTilesBuffers(buffer_capacity);
      const context = {
        artifact_id: $.ARTIFACTS_BY_NAME[this.#artifact_to_build!],
        city_state,
        ...typed_arrays,
      };
      const updateBuffersWithContext = updateBuffers.bind(context);

      _.chunkIterator(chunk).forEach(updateBuffersWithContext);

      this.channel.publish("ui-render", {
        ui: "ui_build_artifact",
        qs_buffer: typed_arrays.qs.buffer,
        rs_buffer: typed_arrays.rs.buffer,
        item_ids_buffer: typed_arrays.item_ids.buffer,
        occupied_buffer: typed_arrays.occupied.buffer,
      });
    } else {
      const typed_arrays = newTilesBuffers(1);
      const context = {
        artifact_id: $.ARTIFACTS_BY_NAME[this.#artifact_to_build!],
        city_state,
        ...typed_arrays,
      };
      const tile_index = _.getTileIndex(q, r);

      updateBuffers.call(context, { tile_index, q, r }, 0);

      this.channel.publish("ui-render", {
        ui: "ui_build_artifact",
        qs_buffer: typed_arrays.qs.buffer,
        rs_buffer: typed_arrays.rs.buffer,
        item_ids_buffer: typed_arrays.item_ids.buffer,
        occupied_buffer: typed_arrays.occupied.buffer,
      });
    }

    this.#can_tick = true;
  }

  /**
   * If artifact's width or height is zero, this means
   * it is an artifact that has a free_size.
   */
  #getEffectiveChunk(position: Coordinates): Chunk {
    return {
      top_left: this.#anchor!,
      bottom_right: {
        q:
          this.#size.width === 0
            ? position.q
            : Math.max(this.#size.width + this.#anchor!.q - 1, this.#anchor!.q),
        r:
          this.#size.height === 0
            ? position.r
            : Math.max(this.#size.height + this.#anchor!.r - 1, this.#anchor!.r),
      },
    };
  }

  /**
   * Call API to register updated tiles in database.
   * Update new city's state zones, artifact_statuses and tile's artifact,
   * and re-draw new state.
   */
  async onPointerup(event: PointerEvent) {
    if (!this.#anchor || event.button !== 0) {
      return;
    }

    const { ok, ...position } = _.getPositionFromPointer(event, this.#camera);
    if (!ok) {
      return;
    }

    const chunk = _.getOrderedChunk(this.#getEffectiveChunk(position));

    const payload: BuildSingleArtifactPayload = {
      artifact_name: this.#artifact_to_build!,
      anchor: chunk.top_left,
      width: this.#size.width || chunk.bottom_right.q - chunk.top_left.q + 1,
      height: this.#size.height || chunk.bottom_right.r - chunk.top_left.r + 1,
    };

    try {
      const updated_tiles = await api.buildSingleArtifact(payload);

      this.channel.publish("update-city-state", updated_tiles);
    } catch (reason) {
      console.error(reason);
    } finally {
      this.channel.publish("ui-reset");
    }
  }

  /**
   * Remove drawn artifact.
   */
  onPointerleave() {
    this.channel.publish("ui-reset");
  }

  /**
   * Remove drawn artifact.
   */
  onKeydown() {
    this.#anchor &&= null;
    this.channel.publish("ui-reset");
  }

  /**
   * Reset all shit because menu state does not
   * correlate to this component action.
   */
  #dispose() {
    window.clearTimeout(this.#throttled);

    this.#can_tick = true;
    this.#artifact_to_build &&= null;
    this.#anchor &&= null;

    this.channel.unSubscribe("city-pointer-event", this);
    this.ownerDocument.removeEventListener("keydown", this);
  }
}

window.customElements.define("sim-build-artifact", SimBuildArtifact);

export { SimBuildArtifact };
