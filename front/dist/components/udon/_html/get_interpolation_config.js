import { DATA } from "./compile_template.js";
/**
 * An interpolation that is a `function` is considered reactive.
 * It will be added to component.update object.
 */
function getInterpolationConfig(type, interpolation) {
    const is_reactive = type !== DATA && typeof interpolation === "function";
    if (is_reactive) {
        const callback = interpolation;
        const initial_data = callback();
        return { callback, initial_data };
    }
    else {
        const initial_data = interpolation;
        return { callback: null, initial_data };
    }
}
export { getInterpolationConfig };
