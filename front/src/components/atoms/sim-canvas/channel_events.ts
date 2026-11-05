import type { UIData } from "./workers/canvas_ui.js";

declare module "components" {
  interface ChannelEvents {
    "grid-render": never;
  }
}

declare module "components" {
  interface ChannelEvents {
    render: never;
  }
}

declare module "components" {
  interface ChannelEvents {
    "ui-render": UIData;
  }
}

declare module "components" {
  interface ChannelEvents {
    "ui-reset": never;
  }
}

declare module "components" {
  interface ChannelEvents {
    "heatmap-render": never;
  }
}
