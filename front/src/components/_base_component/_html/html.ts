import type { BaseComponent, Props } from "../component.js";
import type { InterpolationHint } from "./compile_template.js";

import { TEXT, ATTR_STD, ATTR_BOOL, CLASS_BOOL, compileTemplate } from "./compile_template.js";
import { splitter } from "./splitter.js";
import { setUpdate } from "./set_update.js";

const templates_cache = new WeakMap<TemplateStringsArray, HTMLTemplateElement>();

/**
 *
 * **rendering**
 *
 * A magic thing about tagged template literals is that
 * we can use the static `parts` of a tagged template as a unique key for a template.
 * This key will be the same for all instances of a same component,
 * so every instances of a given component share the same underlying template.
 * @see https://viperhtml.js.org/hyperhtml/documentation/#introduction-1
 *
 * **reactivity**
 *
 * We can use the comments priorly inserted in a template as "hints"
 * to bind reactive values directly to the DOM via functions and
 * a pre-defined update type (`TEXT`, `ATTR_STD`, `ATTR_BOOL`, `CLASS_BOOL`).
 *
 * So no virtual DOM is needed: every time a function bound to a
 * component instance's prop is called, we just simply execute the function
 * that "remembers" what DOM node we want to target, thanks to the closure.
 */
function renderHTML<P extends Props>(
  component: BaseComponent<P>,
  parts: TemplateStringsArray,
  ...interpolations: ((...args: any[]) => any)[]
) {
  const template = templates_cache.getOrInsertComputed(parts, compileTemplate);
  const fragment = template.content.cloneNode(true);
  const tree_walker = document.createTreeWalker(fragment, NodeFilter.SHOW_COMMENT);

  let node = tree_walker.nextNode();

  while (node) {
    const comment = node as Comment;
    // @ts-expect-error
    const { type, name, index } = comment.nodeValue!.split(splitter) as InterpolationHint;
    const callback = interpolations[index];
    const callback_name = callback.name as keyof P;

    if (callback_name in component.update) {
      console.warn(
        `${callback.name} already exists.
        Previous one registered for "${name}" will be overriden.`,
      );
    }

    const parent_element = comment.parentElement!;
    const data = callback();

    switch (type) {
      case TEXT: {
        const text_node = document.createTextNode(data);
        comment.replaceWith(text_node);

        component.props[callback_name] = data;

        component.update[callback_name] = setUpdate((...params: any[]) => {
          const new_data = callback(...params);

          text_node.data = new_data;
          component.props[callback_name] = new_data;
        });

        node = tree_walker.nextNode();
        break;
      }
      case CLASS_BOOL: {
        parent_element.classList.toggle(name, data);

        component.props[callback_name] = data;

        component.update[callback_name] = setUpdate((...params: any[]) => {
          const new_data = callback(...params);

          parent_element.classList.toggle(name, new_data);
          component.props[callback_name] = new_data;
        });

        node = tree_walker.nextNode();
        comment.remove();
        break;
      }
      case ATTR_BOOL: {
        parent_element.toggleAttribute(name, data);

        component.props[callback_name] = data;

        component.update[callback_name] = setUpdate((...params: any[]) => {
          const new_data = callback(...params);

          parent_element.toggleAttribute(name, new_data);
          component.props[callback_name] = new_data;
        });

        node = tree_walker.nextNode();
        comment.remove();
        break;
      }
      case ATTR_STD: {
        parent_element.setAttribute(name, data);

        component.props[callback_name] = data;

        component.update[callback_name] = setUpdate((...params: any[]) => {
          const new_data = callback(...params);

          parent_element.setAttribute(name, new_data);
          component.props[callback_name] = new_data;
        });

        node = tree_walker.nextNode();
        comment.remove();
        break;
      }
      default:
        console.warn(`Unknown interpolation type: ${type}`);
        node = tree_walker.nextNode();
    }
  }

  component.shadowRoot!.appendChild(fragment);
}

export { renderHTML };
