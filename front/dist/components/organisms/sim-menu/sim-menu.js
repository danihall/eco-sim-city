import { BaseComponent } from "components";
import css from "./sim-menu.css" with { type: "css" };
class SimMenu extends BaseComponent {
    static styles = [css];
    #state = {
        selected_action: undefined,
        selected_item: undefined,
    };
    render() {
        this.html `
      <div data-action="build" >
          <button data-item="BusStop">bus stop</button>
          <button data-item="TrainStation">train station</button>
          <button data-item="SubwayStation">subway station</button>
          <button data-item="MixedForest">mixed forest</button>
          <button data-item="UniversityCampus">university campus</button>
          <button data-item="CoalPlant">centrale à charbon</button>
      </div>

      <div data-action="zoning">
          <button data-item="Industrial">zone industrial</button>
          <button data-item="Residential">zone residential</button>
      </div>

      <div data-action="build-pathway">
        <button data-item="Road">road</button>
        <button data-item="Rail">rail</button>
        <button data-item="PedestrianBicycleWay">pedestrian/bicycle way</button>
      </div>

      <div data-action="destroy">
        <button>destroy</button>
      </div>

      <div data-action="query">
          <button data-item="traffic">query traffic</button>
      </div>
      <slot></slot>
      `;
    }
    async onComponentMounted() {
        this.addEventListener("click", this);
    }
    onClick(event) {
        let action_type;
        let item_type;
        for (const element of event.composedPath()) {
            action_type ??= element.dataset?.action;
            item_type ??= element.dataset?.item;
            if (action_type && item_type) {
                break;
            }
        }
        this.#state.selected_action = action_type ?? "";
        this.#state.selected_item = item_type;
        this.channel.publish("menu-state-change", this.state);
    }
    get state() {
        return this.#state;
    }
}
window.customElements.define("sim-menu", SimMenu);
export { SimMenu };
