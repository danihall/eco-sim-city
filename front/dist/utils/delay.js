/**
 * `duration` is in ms.
 */
function delay(duration) {
    return new Promise((resolve) => setTimeout(resolve, duration));
}
export { delay };
