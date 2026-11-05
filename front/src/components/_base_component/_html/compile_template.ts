type InterpolationHint = {
  type: typeof TEXT | typeof ATTR_STD | typeof ATTR_BOOL | typeof CLASS_BOOL;
  name: string;
  index: number;
};

const TEXT = "text";
const ATTR_STD = "attr_std";
const ATTR_BOOL = "attr_bool";
const CLASS_BOOL = "class_bool";

const REGEX_ATTRIBUTE =
  /:[\w|\-|_]+(?=="$|=$|='$)|:class:[\w|\-|_]+(?=="$|=$|='$)|[\w|\-|_]+(?=="$|=$|='$)/;
const REGEX_TAG_OPENING = /^<\w*/;
const REGEX_TAG_CLOSING = /^(?:"|'|)>/;
const REGEX_TAG_CLOSING_END = />$|>(?="|'|$)/;

/**
 * inserts a HTML comment where future interpolations will be placed.
 * These comments will be the "needle-in-a-haystack" that we can
 * search for, to be replaced by interpolated values unique for each
 * instance of a given component.
 */
function compileTemplate(parts: TemplateStringsArray): HTMLTemplateElement {
  let text = "";
  const all_attributes: InterpolationHint[] = [];

  parts.forEach((part: string, index: number, array: readonly string[]) => {
    if (index < array.length && array.length > 1) {
      const trimmed = part.trim();
      const new_attribute = REGEX_ATTRIBUTE.exec(trimmed);

      if (new_attribute) {
        const attribute_hint = _getAttributeHint(new_attribute[0], index);
        const tag_opening = REGEX_TAG_OPENING.exec(trimmed);

        all_attributes.push(attribute_hint);
        text += tag_opening ? trimmed.slice(0, tag_opening[0].length) : "";
        return;
      }

      const tag_closing = REGEX_TAG_CLOSING.exec(trimmed);
      const tag_closing_from_end = REGEX_TAG_CLOSING_END.exec(trimmed);

      if (tag_closing || tag_closing_from_end) {
        const attr_comments = all_attributes
          .map((attr) => `<!--${attr.type}:${attr.name}:${attr.index}-->`)
          .join("");
        const text_comment = index < array.length - 1 ? `<!--${TEXT}::${index}-->` : "";

        if (tag_closing) {
          text += `>${attr_comments}${text_comment}${trimmed.slice(tag_closing[0].length)}`;
        } else {
          text += `${trimmed}${text_comment}`;
        }

        all_attributes.length = 0;
        return;
      }
      // If no new_attribute or tag_closing was found,
      // this means this part is another class name.
      const sibling_attribute = { ...all_attributes.at(-1)!, index };
      all_attributes.push(sibling_attribute);
      return;
    }

    text += part;
  });

  // console.log("COMPILED TEMPLATE:\n", text);
  return Object.assign(document.createElement("template"), {
    innerHTML: text,
  });
}

function _getAttributeHint(attribute_name: string, index: number): InterpolationHint {
  if (attribute_name.startsWith(":class:")) {
    return {
      type: CLASS_BOOL,
      name: attribute_name.slice(7),
      index,
    };
  }

  if (attribute_name.startsWith(":")) {
    return {
      type: ATTR_BOOL,
      name: attribute_name.slice(1),
      index,
    };
  }

  return {
    type: ATTR_STD,
    name: attribute_name,
    index,
  };
}

export type { InterpolationHint };
export { TEXT, ATTR_STD, ATTR_BOOL, CLASS_BOOL, compileTemplate };
