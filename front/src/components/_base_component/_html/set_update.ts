const updates_tick = new Set<(...args: any[]) => any>();

function setUpdate(callback: (...params: any[]) => void) {
  const updater = (...params: any[]) => {
    if (updates_tick.has(updater)) {
      return;
    }

    const is_initiator = updates_tick.size === 0;

    updates_tick.add(updater);

    try {
      callback(...params);
    } finally {
      if (is_initiator) {
        updates_tick.clear();
      }
    }
  };

  return updater;
}

export { setUpdate };
