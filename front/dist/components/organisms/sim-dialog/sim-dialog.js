import { UdonComponent } from "components";
import css_global from "../../../styles.css" with { type: "css" };
import css_fonts from "../../../fonts.css" with { type: "css" };
import css_dialog from "./sim-dialog.css" with { type: "css" };
class SimDialog extends UdonComponent {
    static css = [css_global, css_fonts, css_dialog];
    #dialog;
    dynamic_content = (html = null) => html;
    styles = (style_props = null) => Object.entries(style_props ?? {})
        .map(([key, value]) => `${key}:${value}`)
        .join(";");
    render() {
        const dialog_id = this.id ?? "";
        return this.html /*html*/ `
      <dialog id:=${dialog_id} style:=${this.styles}>
        <button commandfor:=${dialog_id} command="close" class="text-light">fermer</button>
        <div>
            <slot>${this.dynamic_content}</slot>
        </div>
      </dialog>
      `;
    }
    async onMounted() {
        this.#dialog = this.shadowRoot.querySelector("dialog");
        this.#dialog.addEventListener("close", this);
        this.#dialog.addEventListener("click", this);
    }
    show() {
        const { promise, resolve } = Promise.withResolvers();
        this.#dialog.addEventListener("animationend", resolve, { once: true });
        this.#dialog.showModal();
        return promise;
    }
    hide() {
        this.#dialog.close();
    }
    onClose() {
        this.update.dynamic_content(null);
        this.update.styles(null);
    }
    onClick(event) {
        if (event.composedPath()[0] === this.#dialog && this.#dialog.open) {
            this.hide();
        }
    }
    onBeforeUnmount() {
        this.#dialog.removeEventListener("close", this);
        this.#dialog.removeEventListener("click", this);
    }
}
window.customElements.define("sim-dialog", SimDialog);
export { SimDialog };
