import { TEXT_OR_ELEMENT, ATTR_STD, ATTR_BOOL, CLASS_BOOL, DATA } from "./compile_template.js";
const REGEX_INTERPOLATION_EFFECTIVE_TYPE = new RegExp(`${[TEXT_OR_ELEMENT, ATTR_STD, ATTR_BOOL, CLASS_BOOL, DATA].join("|")}`);
const hintSplitter = {
    [Symbol.split](string) {
        const [type, name, _index] = string.split(":");
        if (!REGEX_INTERPOLATION_EFFECTIVE_TYPE.test(type)) {
            console.warn(`unknown interpolation in HTML comment: <!--${type}-->.\nLeft comment as is.`);
            return;
        }
        return {
            type,
            name,
            index: Number(_index),
        };
    },
};
export { hintSplitter };
