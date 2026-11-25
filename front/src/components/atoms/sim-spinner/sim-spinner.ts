import { UdonComponent } from "components";

import css_spinner from "./sim-spinner.css" with { type: "css" };

declare module "components" {
  interface ComponentTagNameMap {
    "sim-spinner": SimSpinner;
  }
}

class SimSpinner extends UdonComponent {
  static css = [css_spinner];

  render() {
    return this.html /* html */ `
        <svg viewBox="0 0 20 20" width="20" height="20">
            <path pathLength="100" d="M2.5,10a7.5,7.5 0 1,0 15,0a7.5,7.5 0 1,0 -15,0"/>
        </svg>`;
  }
}

window.customElements.define("sim-spinner", SimSpinner);

export { SimSpinner };
