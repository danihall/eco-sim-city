type TargetSelector = `__i${number}="${number}"`;

type InterpolationHint = {
  type:
    typeof TEXT_OR_ELEMENT | typeof ATTR_STD | typeof ATTR_BOOL | typeof CLASS_BOOL | typeof DATA;
  name: string;
  index: number;
};

const TEXT_OR_ELEMENT = "text_or_el";
const ATTR_STD = "attr_std";
const ATTR_BOOL = "attr_bool";
const CLASS_BOOL = "class_bool";
const DATA = "data";

const REGEX_SPECIAL_ATTR = /[\w|\-|:]+(?=:="$|:=$|:='$)/;
const REGEX_MARK = /\uE000(\d+)\uE001/;

const textOrElementMark = (index: number) => `\uE000${index}\uE001`;

/**
 * inserts a HTML comment where future interpolations will be placed.
 * These comments will be the "needle-in-a-haystack" that we can
 * search for, to be replaced by interpolated values unique for each
 * instance of a given component.
 */
function compileTemplate(chunks: TemplateStringsArray): HTMLTemplateElement {
  const { html, attribute_hints } = _parseHTML(chunks);

  const template = Object.assign(document.createElement("template"), {
    innerHTML: html,
  });

  const comments_walker = document.createTreeWalker(template.content, NodeFilter.SHOW_COMMENT);
  const user_comments: Comment[] = [];
  while (comments_walker.nextNode()) {
    user_comments.push(comments_walker.currentNode as Comment);
  }

  user_comments.forEach((comment) => comment.remove());

  const text_walker = document.createTreeWalker(template.content, NodeFilter.SHOW_TEXT);
  const marked_texts: Text[] = [];
  while (text_walker.nextNode()) {
    if (REGEX_MARK.test(text_walker.currentNode.nodeValue!)) {
      marked_texts.push(text_walker.currentNode as Text);
    }
  }

  marked_texts.forEach((text) => {
    const marks = text.data.split(REGEX_MARK);

    text.replaceWith(
      ...marks.flatMap<string | Comment>((mark, index) => {
        if (index % 2 === 0) {
          return mark ? [mark] : [];
        }
        return [document.createComment(`${TEXT_OR_ELEMENT}::${mark}`)];
      }),
    );
  });

  template.content.prepend(
    ...attribute_hints
      .filter((hint) => template.content.querySelector(`[__i${hint.index}]`))
      .map(({ type, name, index }) => {
        return document.createComment(`${type}:${name}:${index}`);
      }),
  );

  // console.log("COMPILED TEMPLATE:\n", template.innerHTML);
  return template;
}

function _parseHTML(chunks: TemplateStringsArray) {
  let html = "";
  const attribute_hints: InterpolationHint[] = [];

  chunks.forEach((part: string, index: number, array: readonly string[]) => {
    const has_loose_quote = part.startsWith('"') || part.startsWith("'");

    if (index === array.length - 1) {
      html += has_loose_quote ? part.slice(1) : part;
      return;
    }

    const special_attribute = REGEX_SPECIAL_ATTR.exec(part);

    // handle special attribute mark (including class attribute)
    if (special_attribute) {
      const hint = _getHintFromSpecialAttribute(special_attribute[0], index);
      const selector: TargetSelector = `__i${index}="${index}"`;
      const replacement = part.replace(part.slice(special_attribute.index), selector);

      html += has_loose_quote ? replacement.slice(1) : replacement;
      attribute_hints.push(hint);
      return;
    }

    // if no special attribute, this must be a text node or an element
    html += part + textOrElementMark(index);
  });

  return { html, attribute_hints };
}
function _getHintFromSpecialAttribute(name: string, index: number): InterpolationHint {
  //class:any-css-class-name:=
  if (name.startsWith("class:")) {
    return {
      type: CLASS_BOOL,
      name: name.slice(6),
      index,
    };
  }
  // :disabled:=
  if (name.startsWith(":")) {
    return {
      type: ATTR_BOOL,
      name: name.slice(1),
      index,
    };
  }
  // any_json::=
  if (name.endsWith(":")) {
    return {
      type: DATA,
      name: name.slice(0, name.length - 1),
      index,
    };
  }

  // any-attribute-name:=
  return {
    type: ATTR_STD,
    name: name,
    index,
  };
}

export type { InterpolationHint, TargetSelector };
export { TEXT_OR_ELEMENT, ATTR_STD, ATTR_BOOL, CLASS_BOOL, DATA, compileTemplate };
