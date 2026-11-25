import { UdonComponent } from "components";

declare module "components" {
  interface ComponentTagNameMap {
    "sub-component": SubComponent;
  }
}

class SubComponent extends UdonComponent {
  declare count: (param: unknown) => unknown;

  render() {
    return this.html /* html */ `
      <div style="border:1px solid">
        <sub-sub-component count::=${this.count}></sub-sub-component>
      </div>`;
  }
}

window.customElements.define("sub-component", SubComponent);

export { SubComponent };
