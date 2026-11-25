import { UdonComponent } from "components";
import css from "./sim-button-icon.css" with { type: "css" };
class SimButtonIcon extends UdonComponent {
    static css = [css];
    disabled = (toggle = false) => toggle;
    render() {
        return this.html /* html */ `
      <button :disabled:=${this.disabled}>
        <slot></slot>
      </button>
      `;
    }
    async onMounted() {
        this.addEventListener("click", this);
    }
    onClick() {
        this.update.disabled(!this._disabled);
    }
}
window.customElements.define("sim-button-icon", SimButtonIcon);
export { SimButtonIcon };
