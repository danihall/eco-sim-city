import { Channel } from "./_channel/channel.js";
import { renderHTML } from "./_html/html.js";
import { instances } from "./_instances/instances.js";
/**
 * `BaseComponent` is to be extended by components.
 * It provides better ergonomics to declare a Web Component's HTML
 * via *tagged template literals* as well as providing a simple
 * syntax to update parts of a components without any virtual DOM.
 *
 * Like *React.js* or *Vue.js*, components are not truly reactive,
 * one must trigger an update in an imperative way at some point.
 * Note that you can define updates that themselves trigger
 * other updates from other components, so a **guard system**
 * is in place to prevent circular updates.
 *
 * The whole strategy for rendering and reactivity is heavily inspired by *HyperHTML* and *Solid.js*.
 * My implementation is rudimentary and most certainly NOT the most performant.
 * But it should give you enough to get by and build a complex components system
 * that is sound and remains simple, by staying close to the web primitives.
 * Note that a component can be completely renderless, it can be used to just
 * encapsulate some logic or any data.
 *
 * In addition to the templating & "reactive" ergonomics, the `BaseComponent`
 * gives access to an enhanced pub/sub implementation that, along your run-of-the-mill
 * subscribe/publish/unsubscribe features, provides a `request` method to retrieve
 * data from another component. So you don't have to walk the DOM and get **component-A**
 * to use the data you need from **component-B**.
 *
 * @see {@link render}
 * @see {@link update}
 * @see {@link props}
 * @see {@link channel}
 */
class BaseComponent extends HTMLElement {
    static styles;
    #should_render;
    #is_ready = Promise.withResolvers();
    channel = new Channel(this);
    constructor() {
        super();
        const should_render = Object.hasOwn(this.constructor.prototype, "render");
        if (should_render) {
            const constructor = this.constructor;
            this.attachShadow({ mode: "open" }).adoptedStyleSheets = constructor.styles;
        }
        else {
            this.style.cssText = "display:none;";
        }
        this.#should_render = should_render;
    }
    /**
     * components that want to render html must call
     * this method, using the `html` method.
     * ```
     * render(){
     *  this.html`hello world`;
     * }
     * ```
     */
    render() { }
    /**
     * This method is mandatory for intial rendering.
     * It will bind interpolations to fine-grain
     * updates on the DOM.
     */
    html(parts, ...interpolations) {
        renderHTML(this, parts, ...interpolations);
    }
    /**
     * components have the ability to re-render a dynamic part of their template,
     * by accessing this field and call the function bound to that part.
     * ```
     * class MyComponent extends BaseComponent {
     *  name = (param = "bob") => param;
     *
     *  render = () => this.html`hi my name is ${this.name}`
     *
     *  onComponentMounted() {
     *    this.update.name("joey");
     *  }
     * }
     * ```
     * This will execute a fine-grain update on the DOM.
     */
    update = {};
    /**
     * componentss dynamic parts which are updated on the DOM,
     * have their up-to-date value always available on `this.props`
     * ```
     * class MyComponent extends BaseComponent {
     *  name = (param = "lucy") => param.toUppercase();
     *
     *  render = () => this.html`hi my name is ${this.name}`
     *
     *  onComponentMounted() {
     *    console.log(this.props.name);// "LUCY"
     *    this.update.name("mona");
     *    console.log(this.props.name);// "MONA"
     *  }
     * }
     * ```
     */
    props = {};
    /**
     * Render the component and mark it as ready when `onComponentMounted` is fulfilled.
     */
    async connectedCallback() {
        instances.add(this);
        if (this.#should_render) {
            this.render();
        }
        await this.onComponentMounted();
        this.#is_ready.resolve();
    }
    /**
     * components can add event listeners, start setting up other things when overriding this method.
     */
    async onComponentMounted() { }
    /**
     * ```
     * onComponentMounted(){
     *  this.addEventListener("click", this)
     * }
     *
     * onClick(event) {}
     * ```
     *
     * @see https://webreflection.medium.com/dom-handleevent-a-cross-platform-standard-since-year-2000-5bf17287fd38
     */
    handleEvent(event) {
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
    handleChannelEvent(func_name, data) {
        if (func_name in this) {
            // @ts-expect-error
            this[func_name](data);
        }
        else {
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
export { BaseComponent };
