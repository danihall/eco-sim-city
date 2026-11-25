import { UdonComponent } from "components";
import css_component from "./component_base.css" with { type: "css" };
class SimComponent extends UdonComponent {
    static css = [css_component];
    static useShadowDOM = true;
    render() {
        return this.html /* html */ ``;
    }
    async onMounted() { }
    onBeforeUnmount() { }
}
window.customElements.define("sim-component", SimComponent);
export { SimComponent };
