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
     * onMounted(){
     *  this.channel.subscribe("app-state-change");
     * }
     *
     * onAppStateChange(data: ChannelEvents["app-state-change"]) {}
     * ```
     */
    subscribe(type) {
        if (!subscribers[type]) {
            subscribers[type] = new Set();
        }
        subscribers[type].add(this.#owner);
    }
    /**
     * Useful if need to manually manage a component's
     * subscriptions and cleanup is needed.
     */
    unSubscribe(type) {
        if (!subscribers[type]) {
            return;
        }
        subscribers[type].delete(this.#owner);
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
     * this.channel.publish("app-state-change", "whatever");
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
    async request(target, ...fields) {
        const [target_tag, target_id] = target.split("#");
        await customElements.whenDefined(target_tag);
        const requestor_key = this.#getEffectiveKey(this.#owner);
        const { target_instance, target_key } = this.#getInstanceAndKey(target_tag, target_id);
        const requestor = `${requestor_key}->${target_key}`;
        const answerer = `${target_key}->${requestor_key}`;
        if (requests.has(answerer)) {
            throw Error(`Circular dependency! "${requestor_key}" must already respond to "${target_key}".
        ${JSON.stringify(Object.fromEntries(requests))}`);
        }
        if (target_instance === this.#owner) {
            throw Error(`"${requestor_key}" cannot request itself.`);
        }
        if (requests.get(requestor) === this.#owner) {
            console.warn(`${requestor_key} already has a pending request to ${target_key}`);
        }
        requests.set(requestor, this.#owner);
        try {
            await target_instance.isReady;
            if (!fields.length) {
                return target_instance;
            }
            return fields.map((item) => target_instance[item]);
        }
        catch (reason) {
            console.error(reason);
        }
        finally {
            requests.delete(requestor);
        }
    }
    /**
     * Throws if no instance was found.
     */
    #getInstanceAndKey(target_tag, target_id) {
        const instance = instances
            .values()
            .find((component) => component.localName === target_tag && (!target_id || component.id === target_id));
        if (!instance) {
            throw Error(`Cannot make request to "${target_tag}". No instance of "<${target_tag}>" exists!`);
        }
        return { target_instance: instance, target_key: this.#getEffectiveKey(instance) };
    }
    /**
     * One might request something from another component by
     * using "my-component" without specifying an `id` and that's okay.
     * This function will get the `id` if it exists to be used as a key.
     * This will allow to distinguish between different instances
     * of the same components.
     */
    #getEffectiveKey(component) {
        return component.id ? `${component.localName}#${component.id}` : component.localName;
    }
    /**
     * Automatically unsubscribe from channel events.
     * Should be called when component is removed (disconnectedCallback).
     */
    dispose() {
        Object.values(subscribers).forEach((subscription_type) => {
            subscription_type.delete(this.#owner);
        });
    }
}
export { Channel };
