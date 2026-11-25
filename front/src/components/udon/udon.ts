import type { ChannelEvents } from "./_channel/channel.js";

import { Channel } from "./_channel/channel.js";
import { renderTemplate } from "./_html/render-template.js";
import { instances } from "./_instances/instances.js";

interface ComponentTagNameMap {}

type ReactiveFields = {
  [key: string]: ((...params: any[]) => unknown) | unknown;
};

// I asked Claude to give me these types lol.
type Callbacks<R> = {
  [K in keyof R as R[K] extends (...args: any[]) => any ? K : never]: R[K];
};
type Updaters<R> = {
  [K in keyof Callbacks<R>]: Callbacks<R>[K] extends (...args: infer A) => any
    ? (...args: A) => void
    : never;
};
type ReactiveFieldsValues<P> = string extends keyof P
  ? {}
  : {
      readonly [K in keyof Callbacks<P> as `_${K & string}`]: Callbacks<P>[K] extends (
        ...args: any[]
      ) => infer R
        ? R
        : never;
    };

type UdonComponent<R extends ReactiveFields = {}> = BaseComponentClass<R> & ReactiveFieldsValues<R>;

/**
 * @see README
 */
class BaseComponentClass<R extends ReactiveFields = {}> extends HTMLElement {
  static css: CSSStyleSheet[];
  static useShadowDOM = true;

  declare __internals: ElementInternals;

  #should_render: boolean;
  #is_rendered = false;
  #is_ready = Promise.withResolvers<void>();

  channel = new Channel(this);

  constructor() {
    super();

    const should_render = Object.hasOwn(this.constructor.prototype, "render");

    if (should_render) {
      const constructor = this.constructor as typeof UdonComponent;

      if (constructor.useShadowDOM) {
        this.attachShadow({
          mode: "open",
        }).adoptedStyleSheets = constructor.css ?? [];
      }
    } else {
      this.style.cssText = "display:none;";
    }

    this.__internals = this.attachInternals();
    this.#should_render = should_render;
  }

  /**
   * components that want to render html must call
   * this method, using the `html` method.
   * ```
   * render() {
   *  return this.html`hello world`;
   * }
   * ```
   */
  render(): DocumentFragment | void {}

  /**
   * This method is mandatory for intial rendering.
   * It will bind interpolations to fine-grain
   * updates on the DOM.
   */
  html(parts: TemplateStringsArray, ...interpolations: unknown[]) {
    return renderTemplate(this, parts, ...interpolations);
  }

  /**
   * components have the ability to re-render a dynamic part of their template,
   * by accessing this field and call the function bound to that part.
   * ```
   * class MyComponent extends BaseComponent {
   *  name = (param = "bob") => param;
   *
   *  render() {
   *    return this.html`hi my name is ${this.name}`
   *  }
   *
   *  onMounted() {
   *    this.update.name("joey");
   *  }
   * }
   * ```
   * This will execute a fine-grain update on the DOM.
   */
  update = {} as Updaters<R>;

  /**
   * Render the component and mark it as ready when `onMounted` is fulfilled.
   */
  async connectedCallback() {
    instances.add(this);

    if (this.#should_render && !this.#is_rendered) {
      const constructor = this.constructor as typeof UdonComponent;

      if (!constructor.useShadowDOM) {
        const current_css = (this.getRootNode() as ShadowRoot).adoptedStyleSheets;

        constructor.css.forEach((stylesheet) => {
          if (!current_css.includes(stylesheet)) {
            (this.getRootNode() as ShadowRoot).adoptedStyleSheets.push(stylesheet);
          }
        });
      }

      const fragment = this.render();

      if (fragment) {
        (this.shadowRoot ?? this).appendChild(fragment);
      }

      this.#is_rendered = true;
    }

    await this.onMounted();
    this.#is_ready.resolve();
  }

  /**
   * components can add event listeners, start setting up other things when overriding this method.
   */
  async onMounted() {}

  /**
   * components can do some clean-up before being taken out of the DOM.
   */
  onBeforeUnmount() {}

  /**
   * Automatically unsubscribe from Channel events and
   * remove component from known instances.
   * Event listeners are NOT automatically removed,
   * this must be done in `onBeforeUnmount()` in the component.
   */
  disconnectedCallback() {
    this.onBeforeUnmount();
    this.channel.dispose();
    instances.delete(this);
  }

  /**
   * ```
   * onMounted(){
   *  this.addEventListener("click", this)
   * }
   *
   * onClick(event) {}
   * ```
   *
   * @see https://webreflection.medium.com/dom-handleevent-a-cross-platform-standard-since-year-2000-5bf17287fd38
   */
  handleEvent(event: Event) {
    const handler = `on${event.type[0].toUpperCase() + event.type.slice(1)}`;
    // @ts-expect-error
    this[handler]?.(event);
  }

  /**
   * Used for components that subscribed to a specific `ChannelEvent`.\
   * A Component that subscribed to a specific event must have a
   * method's name derived from the name of the event:
   * ```
   * "app-state-change" -> onAppStateChange();
   * ```
   */
  handleChannelEvent<K extends keyof ChannelEvents>(func_name: string, data: ChannelEvents[K]) {
    if (func_name in this) {
      // @ts-expect-error
      this[func_name](data);
    } else {
      console.error(`function "${func_name}" doesn't exist on ${this.constructor.name}`);
    }
  }

  /**
   * Can be used to await for a another component resolution.\
   * @see {@link connectedCallback}
   */
  get isReady() {
    return this.#is_ready.promise;
  }
}

const UdonComponent = BaseComponentClass as unknown as (new <
  R extends ReactiveFields = {},
>() => UdonComponent<R>) &
  Pick<typeof BaseComponentClass, keyof typeof BaseComponentClass>;

export type { ChannelEvents, Updaters, ComponentTagNameMap, ReactiveFields, BaseComponentClass };
export { UdonComponent };
