const callback_already_exists = (callback_name, component) => `"${callback_name}" callback already exists in ${component.constructor.name}.
  Previous one will be overriden.`;
const incorrect_slot_usage = (component) => `Interpolation in ${component.constructor.name} placed in a "slot" must be a DocumentFragment || DocumentFragment[].
Each DocumentFragment must start with a HTMLElement. Also,  must use a shadowDOM for slots to work.`;
const unknown_interpolation = (type, component) => `Unknown interpolation type: ${type} in ${component.constructor.name}`;
const target_element_not_found = (callback_name, component) => `Target element not found to put ${callback_name} in ${component.constructor.name}.`;
const func_in_data_prop = (func_name, prop_name, component) => `No function is allowed as data given to child component:
  ${component.constructor.name} gave "${prop_name}::=${func_name}" but "${func_name}"" is a Function.`;
export { callback_already_exists, incorrect_slot_usage, unknown_interpolation, target_element_not_found, func_in_data_prop, };
