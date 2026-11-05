import { BaseComponent } from "components";

import css from "./sim-button-icon.css" with { type: "css" };

declare module "components" {
  interface ComponentTagNameMap {
    "sim-button-icon": SimButtonIcon;
  }
}

type SimButtonIconProps = {
  class_1: (text: string) => string;
  disabled: (toggle: boolean) => boolean;
  toggled: (toggle: boolean) => boolean;
};

class SimButtonIcon extends BaseComponent<SimButtonIconProps> {
  static styles = [css];

  render() {
    this.html`
      <button class="sim-btn">
        <slot></slot>
      </button>
      `;
  }

  async onComponentMounted() {
    this.addEventListener("click", this);
  }

  onClick() {
    this.update.disabled?.(!this.props.disabled);
    this.update.toggled?.(!this.props.toggled);
  }
}

window.customElements.define("sim-button-icon", SimButtonIcon);

export { SimButtonIcon };
