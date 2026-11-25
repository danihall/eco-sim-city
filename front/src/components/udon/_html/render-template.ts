import type { InterpolationHint, TargetSelector } from "./compile_template.js";

import { getInterpolationConfig } from "./get_interpolation_config.js";
import { checkPotentialCorrectSlotUsage } from "./check_slot_usage.js";

import {
  TEXT_OR_ELEMENT,
  ATTR_STD,
  ATTR_BOOL,
  CLASS_BOOL,
  DATA,
  compileTemplate,
} from "./compile_template.js";
import { hintSplitter } from "./hint_splitter.js";
import { registerUpdater } from "./register_updater.js";
import * as errors from "./errors.js";

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
function renderTemplate(component: any, parts: TemplateStringsArray, ...interpolations: unknown[]) {
  const template = templates_cache.getOrInsertComputed(parts, compileTemplate);
  const fragment = document.importNode(template.content, true);
  const tree_walker = document.createTreeWalker(fragment, NodeFilter.SHOW_COMMENT);

  let node = tree_walker.nextNode();

  while (node) {
    const comment = node as Comment;
    // @ts-expect-error
    const hint = comment.nodeValue!.split(hintSplitter) as InterpolationHint;
    if (!hint) {
      node = tree_walker.nextNode();
      continue;
    }

    const { type, name, index } = hint;
    const { callback, initial_data } = getInterpolationConfig(type, interpolations[index]);

    if (callback?.name! in component.update) {
      console.warn(errors.callback_already_exists(callback!.name, component));
    }

    const target_element = fragment.querySelector(`[__i${index}="${index}"]`)!;

    switch (type) {
      case CLASS_BOOL: {
        target_element.classList.toggle(name, initial_data);

        if (callback) {
          (component as any)[`_${callback.name}`] = initial_data;
          component.update[callback.name] = registerUpdater((...params: any[]) => {
            const new_data = callback(...params);

            target_element.classList.toggle(name, new_data);
            (component as any)[`_${callback.name}`] = new_data;
          });
        }

        node = tree_walker.nextNode();
        comment.remove();
        break;
      }
      case ATTR_BOOL: {
        target_element.toggleAttribute(name, initial_data);

        if (callback) {
          (component as any)[`_${callback.name}`] = initial_data;
          component.update[callback.name] = registerUpdater((...params: any[]) => {
            const new_data = callback(...params);

            target_element.toggleAttribute(name, new_data);
            (component as any)[`_${callback.name}`] = new_data;
          });
        }

        node = tree_walker.nextNode();
        comment.remove();
        break;
      }
      case ATTR_STD: {
        _setAttributeValue(name, initial_data, target_element);

        if (callback) {
          (component as any)[`_${callback.name}`] = initial_data;
          component.update[callback.name] = registerUpdater((...params: any[]) => {
            const new_data = callback(...params);

            _setAttributeValue(name, new_data, target_element);
            (component as any)[`_${callback.name}`] = new_data;
          });
        }

        node = tree_walker.nextNode();
        comment.remove();
        break;
      }
      case DATA: {
        if (!target_element) {
          console.error(errors.target_element_not_found(name, component));
          node = tree_walker.nextNode();
          comment.remove();
          break;
        }

        if (typeof initial_data === "function") {
          console.error(errors.func_in_data_prop(initial_data.name, name, component));
          node = tree_walker.nextNode();
          comment.remove();
          break;
        }

        (target_element as any)[name] = initial_data;

        node = tree_walker.nextNode();
        comment.remove();
        break;
      }
      case TEXT_OR_ELEMENT: {
        const ok = checkPotentialCorrectSlotUsage(component, comment, initial_data);
        if (!ok) {
          throw Error(errors.incorrect_slot_usage(component));
        }

        if (typeof initial_data === "string" || typeof initial_data === "number") {
          const text_node = document.createTextNode(initial_data.toString());

          comment.before(text_node);

          if (callback) {
            component[`_${callback.name}`] = initial_data;
            component.update[callback.name] = registerUpdater((...params: any[]) => {
              const new_data = callback(...params);

              text_node.data = new_data;
              component[`_${callback.name}`] = new_data;
            });
          }
        } else {
          const slot =
            comment.parentElement instanceof HTMLSlotElement ? comment.parentElement : null;
          const range = document.createRange();
          const elements_set: Element[] = [];
          let future_inject: Comment | null;

          if (slot) {
            range.setStart(component, component.childNodes.length);
          } else {
            range.setStartBefore(comment);
          }

          if (initial_data) {
            const inserted_elements = _insertElements(initial_data, slot, range);
            elements_set.push(...inserted_elements);
          } else {
            future_inject = document.createComment("sisi");
            range.insertNode(future_inject);
          }

          if (callback) {
            (component as any)[`_${callback.name}`] = initial_data;
            (component as any).update[callback.name] = registerUpdater((...params: any[]) => {
              const new_data = callback(...params);

              _removeElements(range, elements_set, future_inject);
              elements_set.length = 0;

              if (new_data) {
                const inserted_elements = _insertElements(new_data, slot, range);
                elements_set.push(...inserted_elements);
              }

              (component as any)[`_${callback.name}`] = new_data;
            });
          }
        }

        node = tree_walker.nextNode();
        comment.remove();
        break;
      }
      default:
        console.warn(errors.unknown_interpolation(type, component));
        node = tree_walker.nextNode();
    }

    target_element?.removeAttribute(`__i${index}`);
  }

  return fragment;
}

function _insertElements(
  data: DocumentFragment | DocumentFragment[],
  slot: HTMLSlotElement | null,
  range: Range,
) {
  const fragments = Array.isArray(data) ? data : [data];
  const container = document.createDocumentFragment();
  const children = [];

  if (slot) {
    fragments.forEach(_linkToSlot, slot);
    container.append(...fragments);
    children.push(...container.children);
    range.insertNode(container);
  } else {
    container.append(...fragments);
    children.push(...container.children);
    range.insertNode(container);
  }

  return children;
}

function _removeElements(range: Range, elements_set: Element[], future_inject: Comment | null) {
  if (elements_set.length) {
    range.setStartBefore(elements_set.at(0)!);
    range.setEndAfter(elements_set.at(-1)!);
    range.deleteContents();
  } else if (future_inject?.isConnected) {
    range.setStartAfter(future_inject);
    future_inject?.remove();
  }
}

/**
 *
 */
function _linkToSlot(this: HTMLSlotElement, fragment: DocumentFragment) {
  [...fragment.children].forEach((element) => (element.slot = this.name));
}

/**
 *
 */
function _setAttributeValue(name: string, data: string, target: Element) {
  if (name === "class") {
    if (data.length) {
      const classes = data.split(" ") as string[];
      classes.forEach((class_name) => target.classList.add(class_name));
    } else {
      target.classList.remove(...target.classList);
    }
  } else {
    target.setAttribute(name, data);
  }
}

export { renderTemplate };
