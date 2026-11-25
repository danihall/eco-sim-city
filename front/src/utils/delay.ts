/**
 * `duration` is in ms.
 */
function delay(duration: number) {
  return new Promise((resolve) => setTimeout(resolve, duration));
}

export { delay };
