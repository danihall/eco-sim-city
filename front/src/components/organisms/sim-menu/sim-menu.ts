import type { ArtifactsKey, HeatmapKey, ZoneKey } from "constants";
import { UdonComponent } from "components";
import * as _ from "utils";
import css_global from "../../../styles.css" with { type: "css" };
import css_menu from "./sim-menu.css" with { type: "css" };

type MenuState = {
  selected_action: string | undefined;
  selected_item: ArtifactsKey | ZoneKey | void;
};

declare module "components" {
  interface ComponentTagNameMap {
    "sim-menu": SimMenu;
  }
}

declare module "components" {
  interface ChannelEvents {
    "menu-state-change": MenuState;
  }
}

class SimMenu extends UdonComponent {
  static css = [css_global, css_menu];

  #state: MenuState = {
    selected_action: undefined,
    selected_item: undefined,
  };

  render() {
    return this.html /*html*/ `
      <div data-action="build" >
          <button data-item="BusStop">bus stop</button>
          <button data-item="TrainStation">train station</button>
          <button data-item="SubwayStation">subway station</button>
          <button data-item="MixedForest">mixed forest</button>
          <button data-item="UniversityCampus">university campus</button>
          <button data-item="CoalPlant">centrale à charbon</button>
      </div>

      <div data-action="update-zoning">
          <button data-item="Industrial">zone industrial</button>
          <button data-item="Residential">zone residential</button>
      </div>

      <div data-action="revert-zoning">
          <button>unzone</button>
      </div>

      <div data-action="build-pathway">
        <button data-item="Road">road</button>
        <button data-item="Rail">rail</button>
        <button data-item="PedestrianBicycleWay">pedestrian/bicycle way</button>
      </div>

      <div data-action="destroy">
        <button>destroy</button>
      </div>

      <div data-action="ordinances">
        <button>arrêtés municipaux</button>
      </div>

      <div data-action="query">
          <button data-item="traffic">query traffic</button>
      </div>
      <slot></slot>
      `;
  }

  async onMounted() {
    this.addEventListener("click", this);
  }

  onClick(event: MouseEvent) {
    let action_type;
    let item_type;

    for (const element of event.composedPath()) {
      action_type ??= (element as HTMLElement).dataset?.action;
      item_type ??= (element as HTMLElement).dataset?.item as ArtifactsKey | ZoneKey;

      if (action_type && item_type) {
        break;
      }
    }

    this.#state.selected_action = action_type ?? "";
    this.#state.selected_item = item_type;

    this.channel.publish("menu-state-change", this.state);
  }

  get state(): Readonly<MenuState> {
    return this.#state;
  }
}

window.customElements.define("sim-menu", SimMenu);

export type { MenuState };
export { SimMenu };
