import { UdonComponent } from "components";
class TestComponent extends UdonComponent {
    class_name = (name = "test-container") => name;
    slotted_content = (content = this.html `<div class="is-slotted">give me bacon</div>`) => content;
    content_1 = (content = ["a"]) => content?.map((item) => this.html `<div>element content 1 ${item}</div>`);
    content_2 = (content = this.html `<div style="background-color: lightblue">element content 2</div>`) => content;
    coutnsisi = (count = 0) => count;
    render() {
        return this.html /* html */ `
      <div style="border: 2px solid lightblue;" class:=${this.class_name}>
          <slot>
            ${this.slotted_content}
            ${this.content_1}
            ${this.content_2}
          </slot>

          <!--<section>
            <div>
            </div>
          </section>-->
      </div>
      `;
    }
    async onMounted() {
        this.update.slotted_content(this.html `<span>give me more bacon</span>`);
        this.update.content_1(["d", "e", "f", "g"]);
        this.update.content_2(this.html `<div style="background-color: orange">element content 2</div>`);
    }
}
window.customElements.define("test-component", TestComponent);
export { TestComponent };
