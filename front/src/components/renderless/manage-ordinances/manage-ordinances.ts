import type { ChannelEvents, SimButtonSubmit, SimButtonSubmitOptions, SimDialog } from "components";

import { UdonComponent } from "components";
import * as api from "api";
import css_global from "../../../styles.css" with { type: "css" };
import css_form from "./sim-ordinance-form.css" with { type: "css" };
import css_fonts from "../../../fonts.css" with { type: "css" };
import { delay } from "utils";
import * as C from "constants";

const SAVE = "Enregistrer";
const OK = "Ok";
const PROMULGATED = `Arrêtés municipaux promulgués ${C.THUMBS_UP}`;
const SUBMIT_ERROR = "Oups! Un problème est survenu";

/**
 * Renderless component that allows to enact or revoke city ordinances.
 */
class SimManageOrdinances extends UdonComponent {
  #dialog_generic!: SimDialog;

  async onMounted() {
    this.#dialog_generic = await this.channel.request("sim-dialog#generic");
    this.channel.subscribe("menu-state-change");
  }

  // Only activates if user has activated the "ordinances" action from the menu.
  async onMenuStateChange({ selected_action }: ChannelEvents["menu-state-change"]) {
    const should_show_ordinances = selected_action === "ordinances";

    if (should_show_ordinances) {
      try {
        const all_ordinances = (await api.getAllOrdinances()).sort(
          ({ ordinance: a }, { ordinance: b }) => {
            return new Intl.Collator("en").compare(a.name, b.name);
          },
        );

        const form = this.html`
            <sim-ordinance-form ordinances::="${all_ordinances}"></sim-ordinance-form>
          `;

        this.#dialog_generic.update.dynamic_content(form);

        const form_instance = await this.channel.request("sim-ordinance-form");

        await this.#dialog_generic.show();
        form_instance.update.appear(true);
      } catch (reason) {
        console.error(reason);
      }
    }
  }
}

type SimOrdinanceFormUpdaters = {
  appear: (param: boolean) => boolean;
};

/**
 * form component to be used to manipulate city's ordinances.
 */
class SimOrdinanceForm extends UdonComponent<SimOrdinanceFormUpdaters> {
  static css = [css_global, css_fonts, css_form];

  declare ordinances: Ordinance[] | null;

  #form!: HTMLFormElement;
  #initial_formdata!: Set<string>;
  appear = (toggle = false) => toggle;

  render() {
    if (!this.ordinances) {
      return;
    }

    const button_options: SimButtonSubmitOptions = {
      disabled: true,
      action_wording: OK,
    };

    return this.html /* html */ `
      <form class:make-appearance:=${this.appear}>
        <header>
            <h2 class="title-regular--medium">Arrêtés municipaux</h2>
            <span class="text-light--small">coût (&curren;)</span>
        </header>

          ${this.ordinances.map(
            ({ ordinance, is_enacted, total_cost }, index) => this.html /*html*/ `
              <div style:=${`--index: ${index}`}>
                  <input-checkbox
                    name::=${ordinance.name}
                    is_checked::=${is_enacted}
                    label::=${ordinance.label}
                    description::=${ordinance.description}
                  >
                  </input-checkbox>
                  <span class="text-light--small">${total_cost}</span>
              </div>
            `,
          )}

        <footer>
          <sim-button-submit options::=${button_options}></sim-button-submit>
        </footer>
      </form>
      `;
  }

  async onMounted(): Promise<void> {
    this.#form = this.shadowRoot!.querySelector("form")!;
    this.#initial_formdata = new Set(new FormData(this.#form).keys());

    this.#form.addEventListener("change", this);
    this.#form.addEventListener("submit", this);
  }

  async onChange(event: Event) {
    const submit_button = this.shadowRoot!.querySelector<SimButtonSubmit>("sim-button-submit")!;
    const new_formdata = new Set(new FormData(event.currentTarget as HTMLFormElement).keys());
    const effective_changes = this.#initial_formdata.symmetricDifference(new_formdata);

    if (effective_changes.size) {
      submit_button.update.state({ disabled: false, wording: SAVE });
    } else {
      submit_button.update.state({ disabled: true, wording: OK });
    }
  }

  async onSubmit(event: SubmitEvent) {
    event.preventDefault();
    const submit_button = this.shadowRoot!.querySelector<SimButtonSubmit>("sim-button-submit")!;
    const formdata = new FormData(event.currentTarget as HTMLFormElement);

    try {
      submit_button.update.state({ disabled: true, loading: true });
      await api.setOrdinances(formdata);
      await delay(2000);
      this.#initial_formdata = new Set(formdata.keys());
      submit_button.update.state({ loading: false, wording: PROMULGATED });
    } catch (reason) {
      console.error(reason);
      submit_button.update.state({ loading: false, disabled: false, wording: SUBMIT_ERROR });
    }
  }

  onBeforeUnmount() {
    this.#form.removeEventListener("change", this);
    this.#form.removeEventListener("submit", this);
  }
}

declare module "components" {
  interface ComponentTagNameMap {
    "sim-ordinance-form": SimOrdinanceForm;
  }
}

window.customElements.define("sim-manage-ordinances", SimManageOrdinances);
window.customElements.define("sim-ordinance-form", SimOrdinanceForm);

export { SimManageOrdinances, SimOrdinanceForm };
