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
