/**
 * If interpolation is inside a `<slot>`, we must check if the component uses
 * a `ShadowDOM`. If the data is not null, it must be a `DocumentFragment`
 * or an array of `DocumentFragment`s. Finally, each `DocumentFragment` must have
 * an `HTMLElement` as its first child.
 */
function checkPotentialCorrectSlotUsage(component: any, comment: Comment, data: unknown) {
  if (data === null || data === undefined) {
    return true;
  }
  if (!(comment.parentElement instanceof HTMLSlotElement)) {
    return true;
  }
  if (!component.constructor.useShadowDOM) {
    return false;
  }
  if (data instanceof DocumentFragment && data.firstElementChild) {
    return true;
  }
  if (
    Array.isArray(data) &&
    data.every((item) => item instanceof DocumentFragment && item.firstElementChild)
  ) {
    return true;
  }

  return false;
}

export { checkPotentialCorrectSlotUsage };
