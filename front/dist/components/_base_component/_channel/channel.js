import { instances } from "../_instances/instances.js";
const REGEX_FUNC_NAME = /(?:^|-)([a-z])/g;
const subscribers = {};
const requests = new Map();
class Channel {
    #owner;
    constructor(owner) {
        this.#owner = owner;
    }
    /**
     * Subscribe to a notification from another component.
     * Uses pretty much the same implicit syntax for method's name
     * than the `handleEvent` interface.
     * ```
     * onComponentMounted(){
     *  this.subscribe("app-state-change", this);
     * }
     *
     * onAppStateChange(data: ChannelEvents["app-state-change"]) {}
     * ```
     */
    subscribe(type, component) {
        if (!subscribers[type]) {
            subscribers[type] = new Set();
        }
        subscribers[type].add(component);
    }
    /**
     * Useful if need to manage a component's lifecycle
     * and cleanup is needed.
     */
    unSubscribe(type, component) {
        if (!subscribers[type]) {
            return;
        }
        subscribers[type].delete(component);
    }
    /**
     * Notify a `ChannelEvents` to all subscribers.
     * ```
     * declare module "components" {
     *  interface ChannelEvents {
     *    "app-state-change": string
     *  }
     * }
     * ```
     * then in the component class:
     * ```
     * this.publish("app-state-change", "whatever");
     * ```
     */
    publish(type, data) {
        if (!subscribers[type]) {
            return;
        }
        const func_name = `on${type.replace(REGEX_FUNC_NAME, (_, t) => t.toUpperCase())}`;
        subscribers[type].forEach((component) => {
            component.handleChannelEvent(func_name, data);
        });
    }
    async request(target_name, ...fields) {
        await customElements.whenDefined(target_name);
        const local_name = this.#owner.localName;
        const requestor = `${local_name}->${target_name}`;
        const answerer = `${target_name}->${local_name}`;
        if (requests.has(answerer)) {
            throw Error(`Circular dependency! "${local_name}" must already respond to "${target_name}".
        ${JSON.stringify(Object.fromEntries(requests))}`);
        }
        if (requests.get(requestor) === this.#owner) {
            console.warn(`${local_name} already has a pending request to ${target_name}`);
        }
        requests.set(requestor, this.#owner);
        try {
            const component_instance = instances
                .values()
                .find((component) => component.localName === target_name);
            if (!component_instance) {
                console.warn(`Cannot make request "${requestor}" No instance of "<${target_name}>" exists!`);
                return [];
            }
            await component_instance.isReady;
            if (!fields.length) {
                return component_instance;
            }
            return fields.map((item) => component_instance[item]);
        }
        catch (reason) {
            console.error(reason);
        }
        finally {
            requests.delete(requestor);
        }
    }
}
export { Channel };
