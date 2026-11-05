import type { MenuState } from "./sim-menu.js";

declare module "components" {
  interface ChannelEvents {
    "menu-state-change": MenuState;
  }
}
