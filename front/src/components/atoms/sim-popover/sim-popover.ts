import { UdonComponent } from "components";
import css_popover from "./sim-popover.css" with { type: "css" };
import css_fonts from "../../../fonts.css" with { type: "css" };

declare module "components" {
  interface ComponentTagNameMap {
    "sim-popover": SimPopover;
  }
}

type SimPopoverOptions = {
  name: string;
  text: string;
  mode?: "click" | "hover";
  positionArea?: string;
};

/**
 * This component is inserted without a Shadow DOM in order to
 * be callable by its anchored trigger via "popovertarget".
 */
class SimPopover extends UdonComponent {
  static css = [css_popover, css_fonts];
  static useShadowDOM = false;

  declare options: SimPopoverOptions;

  #invoker: HTMLElement | null = null;

  render() {
    return this.html /* html */ `
        <div popover id:=${this.options.name + "pop"}>
            <div>
                <p>${this.options.text}</p>
                <div class="tip"></div>
            </div>
        </div>
      `;
  }

  async onMounted() {
    this.#setOpenedState(false);

    if (this.options.mode === "hover") {
      this.#invoker = (this.getRootNode() as ShadowRoot).querySelector(
        `[popovertarget="${this.options.name + "pop"}"]`,
      );

      this.#invoker?.addEventListener("mouseenter", this);
      this.#invoker?.addEventListener("mouseleave", this);
    }

    this.firstElementChild!.addEventListener("toggle", this);
  }

  #setOpenedState(opened: boolean) {
    if (opened) {
      this.__internals.states.add("opened");
    } else {
      this.__internals.states.delete("opened");
    }
  }

  onMouseenter(event: MouseEvent) {
    (this.firstElementChild as HTMLElement).showPopover({ source: this.#invoker! });
  }

  onMouseleave(event: MouseEvent) {
    (this.firstElementChild as HTMLElement).hidePopover();
  }

  onToggle(event: ToggleEvent) {
    this.#setOpenedState(event.newState === "open");
  }

  onBeforeUnmount() {
    this.#invoker?.removeEventListener("mouseenter", this);
    this.#invoker?.removeEventListener("mouseleave", this);
    this.firstElementChild!.removeEventListener("toggle", this);
  }
}

window.customElements.define("sim-popover", SimPopover);

export type { SimPopoverOptions };
export { SimPopover };
