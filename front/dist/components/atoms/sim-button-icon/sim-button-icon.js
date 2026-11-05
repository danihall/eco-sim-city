import { BaseComponent } from "components";
import css from "./sim-button-icon.css" with { type: "css" };
class SimButtonIcon extends BaseComponent {
    static styles = [css];
    render() {
        this.html `
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
