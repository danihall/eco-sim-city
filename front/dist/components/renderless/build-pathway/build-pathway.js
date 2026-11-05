import { BaseComponent } from "components";
import * as $ from "assets";
import * as _ from "utils";
import * as api from "api";
import { isOccupiedTile, newTilesBuffers, updateBuffers } from "../tiles_selection.js";
const REGEX_EFFECTIVE_ARTIFACT = /Road|PedestrianBicycleWay|Rail/;
/**
 * Renderless component that allows user to launch construction
 * of a transport network artifact. The pathway is immediately constructed.
 */
class SimBuildPathway extends BaseComponent {
    #camera;
    #artifact_id = null;
    #anchor = null;
    #can_tick = true;
    #throttled = 0;
    /**
     * Only listen for "click" if all required elements are in the DOM.
     */
    async onComponentMounted() {
        this.channel.subscribe("menu-state-change", this);
    }
    /**
     * Only activates if menu state correlates to this component action.
     * selected item must be "Road|PedestrianBicycleWay|Rail"
     */
    async onMenuStateChange({ selected_action, selected_item }) {
        const can_build = selected_action === "build-pathway" &&
            selected_item &&
            REGEX_EFFECTIVE_ARTIFACT.test(selected_item);
        if (!can_build) {
            return this.#dispose();
        }
        [this.#camera] = await this.channel.request("sim-city", "camera");
        this.#artifact_id = $.ARTIFACTS_BY_NAME[selected_item];
        this.channel.subscribe("city-pointer-event", this);
        this.ownerDocument.addEventListener("keydown", this);
    }
    onCityPointerEvent(event) {
        this.handleEvent(event);
    }
    /**
     * If #anchor is sound, then it means user can "close" the
     * line that was drawn. So let's update the city state via API.
     * Else start drawing desired path.
     */
    async onClick(event) {
        const { ok, ...position } = _.getPositionFromPointer(event, this.#camera);
        if (!ok || !this.#artifact_id || event.button !== 0) {
            return;
        }
        if (this.#anchor) {
            const path_line = this.#computeLine(position);
            if (!path_line) {
                return;
            }
            const [city_state] = await this.channel.request("sim-city", "state");
            const unoccupied_tiles = [];
            _.pathIterator(path_line).forEach((tile) => {
                const is_occupied = isOccupiedTile({ ...tile, city_state });
                if (!is_occupied) {
                    const { q, r } = tile;
                    unoccupied_tiles.push({ q, r });
                }
            });
            if (unoccupied_tiles.length === 0) {
                return;
            }
            try {
                const payload = {
                    artifact_name: $.ARTIFACTS_BY_ID[this.#artifact_id],
                    coordinates: unoccupied_tiles,
                };
                const updated_tiles = await api.buildPathway(payload);
                this.channel.publish("update-city-state", updated_tiles);
            }
            catch (reason) {
                console.error(reason);
            }
            finally {
                this.channel.publish("ui-reset");
                return (this.#anchor = null);
            }
        }
        console.log("drawing...");
        this.#anchor = position;
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
     * Send a message to the canvas worker to refresh the UI.
     */
    async handlePointerMove(event) {
        if (!this.#anchor) {
            return void (this.#can_tick = true);
        }
        const { ok, ...position } = _.getPositionFromPointer(event, this.#camera);
        if (!ok) {
            return void (this.#can_tick = true);
        }
        const [city_state] = await this.channel.request("sim-city", "state");
        if (!city_state) {
            return void (this.#can_tick = true);
        }
        const path_line = this.#computeLine(position);
        if (!path_line) {
            return void (this.#can_tick = true);
        }
        const typed_arrays = newTilesBuffers(path_line.length);
        const context = {
            artifact_id: this.#artifact_id,
            city_state,
            ...typed_arrays,
        };
        const updateBuffersWithContext = updateBuffers.bind(context);
        _.pathIterator(path_line).forEach(updateBuffersWithContext);
        if (!this.#anchor) {
            return void (this.#can_tick = true);
        }
        this.channel.publish("ui-render", {
            ui: "ui_build_artifact",
            qs_buffer: typed_arrays.qs.buffer,
            rs_buffer: typed_arrays.rs.buffer,
            item_ids_buffer: typed_arrays.item_ids.buffer,
            occupied_buffer: typed_arrays.occupied.buffer,
        });
        this.#can_tick = true;
    }
    /**
     * All elements in the returned array are garanteed to form
     * a continuous line. Includes the origin and the destination.
     */
    #computeLine(destination) {
        const origin = this.#anchor;
        const distance = _.computeDistance(origin, destination);
        if (distance === 0) {
            return;
        }
        return Array.from({ length: distance + 1 }, (item, index) => {
            const current_point = index / distance;
            const q = origin.q + (destination.q - origin.q) * current_point;
            const r = origin.r + (destination.r - origin.r) * current_point;
            return _.cubeRound(q, r);
        });
    }
    /**
     * Remove drawn path.
     */
    onPointerleave() {
        this.channel.publish("ui-reset");
    }
    /**
     * Remove drawn path.
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
        this.#artifact_id &&= null;
        this.#anchor &&= null;
        this.channel.unSubscribe("city-pointer-event", this);
        this.ownerDocument.removeEventListener("keydown", this);
    }
}
window.customElements.define("sim-build-pathway", SimBuildPathway);
export { SimBuildPathway };
