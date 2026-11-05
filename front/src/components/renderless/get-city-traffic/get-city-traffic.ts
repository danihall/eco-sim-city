import { TILES_COUNT } from "constants";
import { SimMenu } from "components";
import { BaseComponent } from "../../_base_component/component.js";

import * as api from "api";

type TrafficData = {
  foot_bytes: Uint32Array<ArrayBuffer>;
  car_bytes: Uint32Array<ArrayBuffer>;
  bus_bytes: Uint32Array<ArrayBuffer>;
  bicycle_bytes: Uint32Array<ArrayBuffer>;
  train_bytes: Uint32Array<ArrayBuffer>;
  subway_bytes: Uint32Array<ArrayBuffer>;
};

class SimGetCityTraffic extends BaseComponent {
  constructor() {
    super();
    const canvas = this.ownerDocument.getElementById("grid");

    if (canvas) {
      this.simMenu.addEventListener("click", this);
    }
  }

  simMenu = this.ownerDocument.querySelector<SimMenu>("sim-menu")!;

  async onClick() {
    const { selected_action, selected_item } = this.simMenu.state;

    if (selected_action !== "query") {
      return;
    }

    // if (selected_item === "traffic") {
    //   try {
    //     const traffic_bytes = await api.getTraffic(1);

    //     const foot_bytes = traffic_bytes.subarray(0, TILES_COUNT);
    //     const car_bytes = traffic_bytes.subarray(TILES_COUNT, TILES_COUNT * 2);
    //     const bus_bytes = traffic_bytes.subarray(TILES_COUNT * 2, TILES_COUNT * 3);
    //     const bicycle_bytes = traffic_bytes.subarray(TILES_COUNT * 3, TILES_COUNT * 4);
    //     const train_bytes = traffic_bytes.subarray(TILES_COUNT * 4, TILES_COUNT * 5);
    //     const subway_bytes = traffic_bytes.subarray(TILES_COUNT * 5, TILES_COUNT * 6);

    //     pubsub.publish("display_traffic", {
    //       foot_bytes,
    //       car_bytes,
    //       bus_bytes,
    //       bicycle_bytes,
    //       train_bytes,
    //       subway_bytes,
    //     });
    //   } catch (reason) {
    //     console.error(reason);
    //   }
    // }
  }
}

window.customElements.define("sim-get-city-traffic", SimGetCityTraffic);

export type { TrafficData };
export { SimGetCityTraffic };
