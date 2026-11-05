// import { BaseComponent } from "../../__base_component/component.js";
// import css from "./sim-button.css" with { type: "css" };

// type SimButtonProps = {
//   class_1: (text: string) => string;
//   disabled: (toggle: boolean) => boolean;
//   toggled: (toggle: boolean) => boolean;
// };

// class SimButton extends BaseComponent<SimButtonProps> {
//   static styles = [css];

//   class_1 = (text = "tutu") => text;
//   disabled = (param = false) => param;
//   toggled = (value = false) => value;

//   render() {
//     this.html`
//       <button :disabled=${this.disabled} class="${this.class_1}" :class:is-toggled=${this.toggled}>
//         <slot></slot>
//       </button>
//       `;
//   }

//   onComponentMounted(): void {
//     this.addEventListener("click", this);
//   }

//   onClick() {
//     this.update.disabled?.(!this.props.disabled);
//     this.update.toggled?.(!this.props.toggled);
//   }
// }

// window.customElements.define("sim-button", SimButton);

// export { SimButton };
