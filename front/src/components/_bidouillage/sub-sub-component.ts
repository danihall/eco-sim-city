import { UdonComponent } from "components";

declare module "components" {
  interface ComponentTagNameMap {
    "sub-sub-component": SubSubComponent;
  }
}

type SubSubReactive = {
  count: <T extends number>(param: T) => T;
};

class SubSubComponent extends UdonComponent<SubSubReactive> {
  declare count: (param: number) => number;

  render() {
    return this.html /* html */ `
      <p style="background-color: orange;">
        ${this.count}
      </p>`;
  }

  async onMounted() {
    this.addEventListener("click", this);
  }

  onClick(event: PointerEvent) {
    event.stopImmediatePropagation();
    this.update.count(this._count + 1);
  }
}

window.customElements.define("sub-sub-component", SubSubComponent);

export { SubSubComponent };
