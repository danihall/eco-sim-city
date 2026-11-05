import { BaseComponent } from "components";
import { newTilesBuffers, updateBuffers } from "../tiles_selection.js";
import * as _ from "utils";
import * as $ from "assets";
import * as api from "api";
const REGEX_VALID_ZONE = /Industrial|Residential/;
/**
 * Renderless component that enable user to select tiles
 * and assign a given super_kind to them.
 */
class SimZoneSuperKind extends BaseComponent {
    #camera;
    #can_tick = true;
    #throttled = 0;
    #chunk = {
        top_left: { q: 0, r: 0 },
        bottom_right: { q: 0, r: 0 },
    };
    #is_valid_chunk = false;
    #super_kind = null;
    async onComponentMounted() {
        this.channel.subscribe("menu-state-change", this);
    }
    /**
     * Only activates if user has indeed activated zone selection
     * from the menu. Valid zone must be "Industrial | Residential".
     */
    async onMenuStateChange({ selected_action, selected_item }) {
        const can_zone = selected_action === "zoning" && selected_item && REGEX_VALID_ZONE.test(selected_item);
        if (!can_zone) {
            return this.#dispose();
        }
        [this.#camera] = await this.channel.request("sim-city", "camera");
        this.#super_kind = selected_item;
        this.channel.subscribe("city-pointer-event", this);
        this.ownerDocument.addEventListener("keydown", this);
    }
    onCityPointerEvent(event) {
        this.handleEvent(event);
    }
    /**
     * Only listen for "pointermove", "pointerup" and "keydown"
     * if user has indeed activated zone selection from the menu.
     */
    onPointerdown(event) {
        if (event.button !== 0) {
            return;
        }
        const { q, r, ok } = _.getPositionFromPointer(event, this.#camera);
        if (!ok) {
            return;
        }
        this.#chunk.top_left.q = q;
        this.#chunk.top_left.r = r;
        this.#is_valid_chunk = true;
    }
    /**
     * Throttles zone selection logic and drawing.
     */
    onPointermove(event) {
        if (this.#can_tick) {
            this.#can_tick = false;
            this.#throttled = window.setTimeout(this.handlePointerMove.bind(this), 0, event);
        }
    }
    /**
     * If coordinates are out-of-bounds, cancel everything.
     * Else, only ask to draw if chunk has effectively changed.
     * Drawing instructions are communicated to worker via Transferables ArrayBuffers.
     * Contrary to other canvas workers, *ui* canvas does not need to keep
     * an internal state from SharedArrayBuffers, since it does only concern itself
     * with drawing temporary stuff stemming from user actions.
     */
    async handlePointerMove(event) {
        const { q, r, ok } = _.getPositionFromPointer(event, this.#camera);
        if (!ok || event.pressure !== 0.5 || !this.#is_valid_chunk) {
            this.#is_valid_chunk = false;
            this.#can_tick = true;
            return this.channel.publish("ui-reset");
        }
        if (this.#chunk.bottom_right.q !== q || this.#chunk.bottom_right.r !== r) {
            this.#chunk.bottom_right.q = q;
            this.#chunk.bottom_right.r = r;
            const [city_state] = await this.channel.request("sim-city", "state");
            if (!city_state) {
                return void (this.#can_tick = true);
            }
            const buffer_capacity = _.getChunkTilesCount(this.#chunk);
            const typed_arrays = newTilesBuffers(buffer_capacity);
            const context = {
                artifact_id: $.SUPER_KINDS_BY_NAME[this.#super_kind],
                city_state,
                ...typed_arrays,
            };
            const updateBuffersWithContext = updateBuffers.bind(context);
            _.chunkIterator(this.#chunk).forEach(updateBuffersWithContext);
            this.channel.publish("ui-render", {
                ui: "ui_zones",
                qs_buffer: typed_arrays.qs.buffer,
                rs_buffer: typed_arrays.rs.buffer,
                item_ids_buffer: typed_arrays.item_ids.buffer,
                occupied_buffer: typed_arrays.occupied.buffer,
            });
        }
        this.#can_tick = true;
    }
    /**
     * Call API to register new zone in database.
     * Update new city's state zones and artifact_statuses,
     * and re-draw new state.
     */
    async onPointerup(event) {
        if (!this.#is_valid_chunk || event.button !== 0) {
            return this.channel.publish("ui-reset");
        }
        const { q, r } = _.getPositionFromPointer(event, this.#camera);
        this.#chunk.bottom_right.q = q;
        this.#chunk.bottom_right.r = r;
        const effective_chunk = _.getOrderedChunk(this.#chunk);
        const payload = {
            super_kind_name: this.#super_kind,
            top_left: effective_chunk.top_left,
            bottom_right: effective_chunk.bottom_right,
        };
        try {
            const updated_tiles = await api.zoneSuperKind(payload);
            this.channel.publish("update-city-zones", { updated_tiles, super_kind: this.#super_kind });
        }
        catch (reason) {
            console.error(reason);
        }
        finally {
            this.channel.publish("ui-reset");
        }
    }
    /**
     * Remove drawn zone.
     */
    onPointerleave() {
        this.channel.publish("ui-reset");
    }
    /**
     * Remove drawn zone.
     */
    onKeydown() {
        this.#dispose();
        this.channel.publish("ui-reset");
    }
    /**
     * Reset all shit because menu state does not
     * correlate to this component action.
     */
    #dispose() {
        this.#chunk = {
            top_left: { q: 0, r: 0 },
            bottom_right: { q: 0, r: 0 },
        };
        this.#super_kind = null;
        this.#is_valid_chunk = false;
        this.#can_tick = true;
        window.clearTimeout(this.#throttled);
        this.channel.unSubscribe("city-pointer-event", this);
        this.ownerDocument.removeEventListener("keydown", this);
    }
}
window.customElements.define("sim-zone-super-kind", SimZoneSuperKind);
export { SimZoneSuperKind };
