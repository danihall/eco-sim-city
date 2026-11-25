import { UdonComponent } from "components";
import css_global from "../../../styles.css" with { type: "css" };
import css_fonts from "../../../fonts.css" with { type: "css" };
import css_button_submit from "./sim-button-submit.css" with { type: "css" };
class SimButtonSubmit extends UdonComponent {
    static css = [css_global, css_fonts, css_button_submit];
    static useShadowDOM = false;
    #easing;
    disabled = (toggle = this.options?.disabled) => toggle;
    action_wording = (wording = this.options?.action_wording ?? "Valider") => {
        // I wrote this to animate the button when updating its text.
        // When all browsers support "allow-keywords" to simply animate between "auto" dimensions
        // with just CSS code, this won't be needed.
        if (this._action_wording !== undefined) {
            const button = this.firstElementChild;
            const span = button.firstElementChild;
            const from = this.getBoundingClientRect().width;
            span.getAnimations().forEach(_cancelAnimation);
            button.getAnimations().forEach(_cancelAnimation);
            span.firstChild.textContent = wording;
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
    show_spinner = (show = false) => (show ? this.html `<sim-spinner></sim-spinner>` : null);
    // Other updater functions will be added when mounting the component.
    update = {
        state: (param) => {
            if (Object.hasOwn(param, "loading")) {
                this.update.show_spinner(param.loading);
            }
            if (Object.hasOwn(param, "wording")) {
                this.update.action_wording(param.wording);
            }
            if (Object.hasOwn(param, "disabled")) {
                this.update.disabled(param.disabled);
            }
        },
    };
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
function _cancelAnimation(animation) {
    animation.cancel();
}
window.customElements.define("sim-button-submit", SimButtonSubmit);
export { SimButtonSubmit };
