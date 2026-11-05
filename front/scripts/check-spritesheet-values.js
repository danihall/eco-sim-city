import * as fs from "node:fs";

const regex_key = /zones|artifacts|heatmap/;

try {
  const file = await import("../dist/constants/spritesheet.js");
  const spritesheet_values = file.SPRITESHEET;

  Object.keys(spritesheet_values).forEach((key) => {
    if (!regex_key.test(key)) {
      throw Error(`${key} is invalid as a sprite type`);
    }
  });

  const config = Object.entries(spritesheet_values).reduce(
    (all_values, current_sprites_types) => {
      const [, sprites] = current_sprites_types;

      return {
        names: [...all_values.names, ...Object.keys(sprites)],
        values: [...all_values.values, ...Object.values(sprites)],
      };
    },
    { names: [], values: [] },
  );

  let base_y = 0;

  config.values.forEach((value, index) => {
    const name = config.names[index];

    if (!Object.hasOwn(value, "x")) {
      throw Error(`missing "x" field in "${name}": ${JSON.stringify(value)}`);
    }
    if (typeof value.x !== "number" || value.x % 16 !== 0) {
      throw Error(`"x" should be a multiple of 16 in "${name}": ${JSON.stringify(value)}`);
    }

    if (!Object.hasOwn(value, "y")) {
      throw Error(`missing "y" field in "${name}": ${JSON.stringify(value)}`);
    }
    if (typeof value.y !== "number" || value.y % 16 !== 0) {
      throw Error(`"y" should be a multiple of 16 in "${name}": ${JSON.stringify(value)}`);
    }
    if (value.y < base_y) {
      throw Error(`
        each subsequent sprite should have a higher "y" value.
        Last "y" value was ${base_y}, "${name}": ${JSON.stringify(value)} has ${value.y}.
        `);
    }

    base_y = value.y;

    if (!Object.hasOwn(value, "width")) {
      throw Error(`missing "width" field in "${name}": ${JSON.stringify(value)}`);
    }
    //TODO apply correct formula (width is not necessarily a multiple of 16)
    if (typeof value.width !== "number" || value.width % 2 !== 0) {
      throw Error(`"width" should be a multiple of 16 in "${name}": ${JSON.stringify(value)}`);
    }

    if (!Object.hasOwn(value, "height")) {
      throw Error(`missing "height" field in "${name}": ${JSON.stringify(value)}`);
    }
    //TODO apply correct formula (height is not necessarily a multiple of 16)
    if (typeof value.height !== "number" || value.height % 2 !== 0) {
      throw Error(`"height" should be a multiple of 16 in "${name}": ${JSON.stringify(value)}`);
    }
  });
} catch (reason) {
  console.error(reason);
}
