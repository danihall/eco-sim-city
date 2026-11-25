import type { BaseComponentClass } from "../udon.js";

const instances = new Set<BaseComponentClass<any>>();

export { instances };
