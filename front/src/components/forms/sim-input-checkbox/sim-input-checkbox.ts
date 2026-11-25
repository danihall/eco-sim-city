import type { SimPopoverOptions } from "components";

import { UdonComponent } from "components";
import css_global from "../../../styles.css" with { type: "css" };
import css_fonts from "../../../fonts.css" with { type: "css" };
import css_input from "./sim-input-checkbox.css" with { type: "css" };

declare module "components" {
  interface ComponentTagNameMap {
    "input-checkbox": InputCheckbox;
  }
}

/**
 * This component is inserted without a Shadow DOM in order to
 * natively participate in its parent `<form>`.
 */
class InputCheckbox extends UdonComponent {
  static css = [css_global, css_fonts, css_input];
  static useShadowDOM = false;

  declare name: string;
  declare is_checked: boolean;
  declare label: string;
  declare description: string;

  render() {
    const popover_options: SimPopoverOptions = {
      name: this.name,
      text: this.description,
    };

    const popover = this.description
      ? this.html /*html*/ `
          <button
            type="button"
            class="text-light--small"
            popovertarget:=${this.name + "pop"}>&#65110;</button>
          <sim-popover options::=${popover_options}></sim-popover>
          `
      : null;

    return this.html /* html */ `
        <input
            type="checkbox"
            name:=${this.name}
            id:=${this.name}
            :checked:=${this.is_checked}
          />
        <label class="text-regular" for:=${this.name}>${this.label}</label>

        ${popover}
      `;
  }

  async onMounted() {
    this.#setCheckedState(this.is_checked);
    this.addEventListener("change", this);
  }

  #setCheckedState(checked: boolean) {
    if (checked) {
      this.__internals.states.add("checked");
    } else {
      this.__internals.states.delete("checked");
    }
  }

  onChange(event: Event) {
    this.#setCheckedState((event.target as HTMLInputElement).checked);
  }

  onBeforeUnmount() {
    this.removeEventListener("change", this);
  }
}

window.customElements.define("input-checkbox", InputCheckbox);

export { InputCheckbox };
