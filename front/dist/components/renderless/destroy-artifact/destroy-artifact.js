import { BaseComponent } from "components";
import * as _ from "utils";
import * as api from "api";
import { newTilesBuffers, updateBuffers } from "../tiles_selection.js";
/**
 * Renderless component that allows to destroy an artifact.
 */
class SimDestroyArtifact extends BaseComponent {
    #city;
    #can_tick = true;
    #throttled = 0;
    async onComponentMounted() {
        this.channel.subscribe("menu-state-change", this);
    }
    /**
     * Only activates if user has activated the "destroy" action from the menu.
     */
    async onMenuStateChange({ selected_action }) {
        const can_destroy = selected_action === "destroy";
        if (!can_destroy) {
            return this.#dispose();
        }
        this.#city = await this.channel.request("sim-city");
        this.channel.subscribe("city-pointer-event", this);
        this.ownerDocument.addEventListener("keydown", this);
    }
    onCityPointerEvent(event) {
        this.handleEvent(event);
    }
    /**
     * Throttles drawing.
     */
    onPointermove(event) {
        if (this.#can_tick) {
            this.#can_tick = false;
            this.#throttled = window.setTimeout(this.handlePointerMove.bind(this), 0, event);
        }
    }
    /**
     * Compute which artifact to highlight on the UI canvas.
     */
    async handlePointerMove(event) {
        const { ok, q, r } = _.getPositionFromPointer(event, this.#city.camera);
        if (!ok) {
            return void (this.#can_tick = true);
        }
        const tile_index = _.getTileIndex(q, r);
        const anchor_q = this.#city.state.anchor_qs[tile_index];
        const anchor_r = this.#city.state.anchor_rs[tile_index];
        const anchor_tile_index = _.getTileIndex(anchor_q, anchor_r);
        const anchor_artifact_id = this.#city.state.artifact_ids[anchor_tile_index];
        const chunk = {
            top_left: { q: anchor_q, r: anchor_r },
            bottom_right: { q: anchor_q, r: anchor_r },
        };
        const tiles_count = _.getChunkTilesCount(chunk);
        console.log(tiles_count);
        const typed_arrays = newTilesBuffers(tiles_count);
        const context = {
            artifact_id: anchor_artifact_id,
            city_state: this.#city.state,
            ...typed_arrays,
        };
        const updateBuffersWithContext = updateBuffers.bind(context);
        _.chunkIterator(chunk).forEach(updateBuffersWithContext);
        this.channel.publish("ui-render", {
            ui: "ui_destroy_artifact",
            qs_buffer: typed_arrays.qs.buffer,
            rs_buffer: typed_arrays.rs.buffer,
            item_ids_buffer: typed_arrays.item_ids.buffer,
            occupied_buffer: typed_arrays.occupied.buffer,
        });
        this.#can_tick = true;
    }
    /**
     * call API to effectively destroy artifact.
     * It should be replaced by "Wasteland" artfacts.
     */
    async onClick(event) {
        console.log(event);
        const { ok, q, r } = _.getPositionFromPointer(event, this.#city.camera);
        if (!ok) {
            return;
        }
        const tile_index = _.getTileIndex(q, r);
        const anchor_q = this.#city.state.anchor_qs[tile_index];
        const anchor_r = this.#city.state.anchor_rs[tile_index];
        const payload = {
            position: { q: anchor_q, r: anchor_r },
        };
        try {
            const updated_tiles = await api.destroyArtifact(payload);
            this.channel.publish("update-city-state", updated_tiles);
        }
        catch (reason) {
            console.error(reason);
        }
        finally {
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
        this.channel.publish("ui-reset");
    }
    /**
     * Reset all shit because menu state does not
     * correlate to this component action.
     */
    #dispose() {
        window.clearTimeout(this.#throttled);
        this.#can_tick = true;
        this.channel.unSubscribe("city-pointer-event", this);
        this.ownerDocument.removeEventListener("keydown", this);
    }
}
window.customElements.define("sim-destroy-artifact", SimDestroyArtifact);
export { SimDestroyArtifact };
