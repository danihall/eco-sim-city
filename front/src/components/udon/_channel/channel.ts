import type { ComponentTagNameMap } from "../udon.js";

import { instances } from "../_instances/instances.js";

interface ChannelEvents {}

interface Subscriber {
  localName: string;
  id: string | null;
  handleChannelEvent<K extends keyof ChannelEvents>(
    func_name: string,
    data?: ChannelEvents[K],
  ): void;
}

type Subscribers = {
  [K in keyof ChannelEvents]?: Set<Subscriber>;
};

type CompositeKey = `${string}->${string}`;
type Requests = Map<CompositeKey, Subscriber>;

// I asked Claude to give me the types lol.
type TagName = keyof ComponentTagNameMap & string;
type Target = TagName | `${TagName}#${string}`;
type TagOf<T extends Target> = T extends `${infer K}#${string}` ? K : T;

const REGEX_FUNC_NAME = /(?:^|-)([a-z])/g;

const subscribers: Subscribers = {};
const requests: Requests = new Map();

class Channel {
  #owner: Subscriber;

  constructor(owner: Subscriber) {
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
  subscribe<K extends keyof ChannelEvents>(type: K) {
    if (!subscribers[type]) {
      subscribers[type] = new Set() as Subscribers[K];
    }
    subscribers[type]!.add(this.#owner);
  }

  /**
   * Useful if need to manually manage a component's
   * subscriptions and cleanup is needed.
   */
  unSubscribe<K extends keyof ChannelEvents>(type: K) {
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
  publish<K extends keyof ChannelEvents>(type: K, data?: ChannelEvents[K]) {
    if (!subscribers[type]) {
      return;
    }

    const func_name = `on${(type as string).replace(REGEX_FUNC_NAME, (_, t) => t.toUpperCase())}`;

    subscribers[type].forEach((component: Subscriber) => {
      component.handleChannelEvent(func_name, data);
    });
  }

  /**
   * Request data from another component. Or request another component instance.
   * Just be sure to use `await this.channel.request(...)`.
   * The promise is rejected if no instance of target component was found.
   */
  request<T extends Target>(target: T): Promise<ComponentTagNameMap[TagOf<T> & TagName]>;

  request<T extends Target, F extends keyof ComponentTagNameMap[TagOf<T> & TagName]>(
    target: T,
    ...fields: F[]
  ): Promise<Readonly<ComponentTagNameMap[TagOf<T> & TagName][F]>[]>;

  async request<T extends Target, F extends keyof ComponentTagNameMap[TagOf<T> & TagName]>(
    target: T,
    ...fields: F[]
  ) {
    const [target_tag, target_id] = target.split("#");

    await customElements.whenDefined(target_tag);

    const requestor_key = this.#getEffectiveKey(this.#owner);
    const { target_instance, target_key } = this.#getInstanceAndKey(target_tag, target_id);

    const requestor = `${requestor_key}->${target_key}` as CompositeKey;
    const answerer = `${target_key}->${requestor_key}` as CompositeKey;

    if (requests.has(answerer)) {
      throw Error(
        `Circular dependency! "${requestor_key}" must already respond to "${target_key}".
        ${JSON.stringify(Object.fromEntries(requests))}`,
      );
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

      return fields.map((item) => target_instance[item as keyof typeof target_instance]);
    } catch (reason) {
      console.error(reason);
    } finally {
      requests.delete(requestor);
    }
  }

  /**
   * Throws if no instance was found.
   */
  #getInstanceAndKey(target_tag: string, target_id: string | void) {
    const instance = instances
      .values()
      .find(
        (component) =>
          component.localName === target_tag && (!target_id || component.id === target_id),
      );

    if (!instance) {
      throw Error(
        `Cannot make request to "${target_tag}". No instance of "<${target_tag}>" exists!`,
      );
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
  #getEffectiveKey(component: Subscriber) {
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

export type { ChannelEvents };
export { Channel };
