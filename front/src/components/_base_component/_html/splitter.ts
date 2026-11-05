import type { InterpolationHint } from "./compile_template.js";
import { TEXT, ATTR_STD, ATTR_BOOL, CLASS_BOOL } from "./compile_template.js";

const REGEX_INTERPOLATION_EFFECTIVE_TYPE = new RegExp(
  `${[TEXT, ATTR_STD, ATTR_BOOL, CLASS_BOOL].join("|")}`,
);

const splitter = {
  [Symbol.split](string: string): InterpolationHint | void {
    const [type, name, index] = string.split(":");

    if (!REGEX_INTERPOLATION_EFFECTIVE_TYPE.test(type)) {
      console.warn(`interpolation of type ${type} is invalid`);
      return;
    }

    return { type, name, index: Number(index) } as InterpolationHint;
  },
};

export { splitter };
