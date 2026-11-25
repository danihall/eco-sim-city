import { UdonComponent } from "components";
class SubComponent extends UdonComponent {
    render() {
        return this.html /* html */ `
      <div style="border:1px solid">
        <sub-sub-component count::=${this.count}></sub-sub-component>
      </div>`;
    }
}
window.customElements.define("sub-component", SubComponent);
export { SubComponent };
