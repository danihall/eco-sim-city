const updates_tick = new Set();
function registerUpdater(callback) {
    const updater = (...params) => {
        if (updates_tick.has(updater)) {
            return;
        }
        const is_initiator = updates_tick.size === 0;
        updates_tick.add(updater);
        try {
            callback(...params);
        }
        finally {
            if (is_initiator) {
                updates_tick.clear();
            }
        }
    };
    return updater;
}
export { registerUpdater };
