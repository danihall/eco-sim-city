import { UdonComponent } from "components";

import css_component from "./component_base.css" with { type: "css" };

declare module "components" {
  interface ComponentTagNameMap {
    "sim-component": SimComponent;
  }
}

type SimComponentReactive = {};

class SimComponent extends UdonComponent<SimComponentReactive> {
  static css = [css_component];
  static useShadowDOM = true;

  render() {
    return this.html /* html */ ``;
  }

  async onMounted() {}

  onBeforeUnmount() {}
}

window.customElements.define("sim-component", SimComponent);

export { SimComponent };
