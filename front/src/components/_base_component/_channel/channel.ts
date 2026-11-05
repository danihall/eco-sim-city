import type { BaseComponent, ComponentTagNameMap } from "components";

import { instances } from "../_instances/instances.js";

interface ChannelEvents {}

interface Subscriber {
  localName: string;
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
   * onComponentMounted(){
   *  this.subscribe("app-state-change", this);
   * }
   *
   * onAppStateChange(data: ChannelEvents["app-state-change"]) {}
   * ```
   */
  subscribe<K extends keyof ChannelEvents>(type: K, component: Subscriber) {
    if (!subscribers[type]) {
      subscribers[type] = new Set() as Subscribers[K];
    }
    subscribers[type]!.add(component);
  }

  /**
   * Useful if need to manage a component's lifecycle
   * and cleanup is needed.
   */
  unSubscribe<K extends keyof ChannelEvents>(type: K, component: Subscriber) {
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
   *
   */
  request<K extends keyof ComponentTagNameMap, F extends keyof ComponentTagNameMap[K]>(
    target_name: K,
  ): Promise<ComponentTagNameMap[K]>;

  request<K extends keyof ComponentTagNameMap, F extends keyof ComponentTagNameMap[K]>(
    target_name: K,
    ...fields: F[]
  ): Promise<Readonly<ComponentTagNameMap[K][F][]>>;

  async request<K extends keyof ComponentTagNameMap, F extends keyof ComponentTagNameMap[K]>(
    target_name: K,
    ...fields: F[]
  ) {
    await customElements.whenDefined(target_name);

    const local_name = this.#owner.localName;
    const requestor = `${local_name}->${target_name}` as CompositeKey;
    const answerer = `${target_name}->${local_name}` as CompositeKey;

    if (requests.has(answerer)) {
      throw Error(
        `Circular dependency! "${local_name}" must already respond to "${target_name}".
        ${JSON.stringify(Object.fromEntries(requests))}`,
      );
    }

    if (requests.get(requestor) === this.#owner) {
      console.warn(`${local_name} already has a pending request to ${target_name}`);
    }

    requests.set(requestor, this.#owner);

    try {
      const component_instance = instances
        .values()
        .find(
          (component) => component.localName === target_name,
        ) as unknown as ComponentTagNameMap[K];

      if (!component_instance) {
        console.warn(
          `Cannot make request "${requestor}" No instance of "<${target_name}>" exists!`,
        );
        return [];
      }

      await component_instance.isReady;

      if (!fields.length) {
        return component_instance;
      }

      return fields.map((item) => component_instance[item]);
    } catch (reason) {
      console.error(reason);
    } finally {
      requests.delete(requestor);
    }
  }
}

export type { ChannelEvents };
export { Channel };
