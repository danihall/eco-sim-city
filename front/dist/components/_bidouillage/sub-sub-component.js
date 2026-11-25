import { UdonComponent } from "components";
class SubSubComponent extends UdonComponent {
    render() {
        return this.html /* html */ `
      <p style="background-color: orange;">
        ${this.count}
      </p>`;
    }
    async onMounted() {
        this.addEventListener("click", this);
    }
    onClick(event) {
        event.stopImmediatePropagation();
        this.update.count(this._count + 1);
    }
}
window.customElements.define("sub-sub-component", SubSubComponent);
export { SubSubComponent };
