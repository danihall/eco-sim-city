import type { Updaters } from "components";

import { UdonComponent } from "components";
import css_global from "../../../styles.css" with { type: "css" };
import css_fonts from "../../../fonts.css" with { type: "css" };
import css_button_submit from "./sim-button-submit.css" with { type: "css" };

declare module "components" {
  interface ComponentTagNameMap {
    "sim-button-submit": SimButtonSubmit;
  }
}

type SimButtonSubmitOptions = {
  form_id?: string;
  disabled?: boolean;
  action_wording?: string;
};

type SimButtonSubmitUpdaters = {
  disabled: (param: boolean) => boolean;
  action_wording: (param: string) => string;
  show_spinner: (param: boolean) => DocumentFragment | null;
};

type State = {
  loading: boolean;
  disabled: boolean;
  wording: string;
};
type SimButtonSubmitComposite = {
  state: (param: Partial<State>) => void;
};

class SimButtonSubmit extends UdonComponent<SimButtonSubmitUpdaters> {
  static css = [css_global, css_fonts, css_button_submit];
  static useShadowDOM = false;

  declare options: SimButtonSubmitOptions | null;

  #easing!: string;

  disabled = (toggle = this.options?.disabled) => toggle;
  action_wording = (wording = this.options?.action_wording ?? "Valider") => {
    // I wrote this to animate the button when updating its text.
    // When all browsers support "allow-keywords" to simply animate between "auto" dimensions
    // with just CSS code, this won't be needed.
    if (this._action_wording !== undefined) {
      const button = this.firstElementChild!;
      const span = button.firstElementChild as HTMLSpanElement;
      const from = this.getBoundingClientRect().width;

      span.getAnimations().forEach(_cancelAnimation);
      button.getAnimations().forEach(_cancelAnimation);
      span.firstChild!.textContent = wording;

      const to = this.getBoundingClientRect().width;

      if (from === to) {
        return wording;
      }

      span.animate([{ opacity: 0 }, { opacity: 1 }], { duration: 600 });
      button.animate([{ width: `${from}px` }, { width: `${to}px` }], {
        duration: 600,
        easing: this.#easing,
      });
    }

    return wording;
  };
  show_spinner = (show = false) => (show ? this.html`<sim-spinner></sim-spinner>` : null);
  // Other updater functions will be added when mounting the component.
  update = {
    state: (param: Partial<State>) => {
      if (Object.hasOwn(param, "loading")) {
        this.update.show_spinner(param.loading!);
      }
      if (Object.hasOwn(param, "wording")) {
        this.update.action_wording(param.wording!);
      }
      if (Object.hasOwn(param, "disabled")) {
        this.update.disabled(param.disabled!);
      }
    },
  } as Updaters<SimButtonSubmitUpdaters & SimButtonSubmitComposite>;

  render() {
    return this.html /* html */ `
      <button
        type="submit"
        class="text-regular"
        :disabled:=${this.disabled}
      >
        <span>${this.action_wording}</span>
        ${this.show_spinner}
      </button>
      `;
  }

  async onMounted() {
    this.#easing = window.getComputedStyle(this).getPropertyValue("--easing") || "ease-in-out";
  }
}

function _cancelAnimation(animation: Animation) {
  animation.cancel();
}

window.customElements.define("sim-button-submit", SimButtonSubmit);

export type { SimButtonSubmitOptions };
export { SimButtonSubmit };
