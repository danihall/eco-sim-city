import type { SuperKindName } from "assets";

declare module "components" {
  interface ChannelEvents {
    "city-pointer-event": PointerEvent;
  }
}

declare module "components" {
  interface ChannelEvents {
    "update-city-state": Tile[];
  }
}

declare module "components" {
  interface ChannelEvents {
    "update-city-zones": { updated_tiles: Tile[]; super_kind: SuperKindName };
  }
}
