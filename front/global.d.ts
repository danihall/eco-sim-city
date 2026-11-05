declare module "*.css" {
  const stylesheet: CSSStyleSheet;
  export default stylesheet;
}

declare type gameLoop = "artifacts_placing_loop" | "artifacts_completion_loop";

declare type SSEPayload = {
  game_loop_name: gameLoop;
  body: any;
};

declare type Coordinates = {
  q: number;
  r: number;
};

declare type PixelCoordinates = {
  x: number;
  y: number;
};

declare type ArtifactInstance = {
  label: string;
  name: ArtifactName;
  status: TileStatusName;
};

declare type Tile = {
  artifact: ArtifactInstance;
  position: Coordinates;
  anchor: Coordinates;
  super_kind: SuperKindName;
};

declare type Sprite = {
  x: number;
  y: number;
  width: number;
  height: number;
};
