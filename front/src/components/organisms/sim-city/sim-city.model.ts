type CameraState = {
  x: number;
  y: number;
  zoom: number;
};

type Concrete = true;

type CityData<T = unknown> = {
  artifact_ids: T extends Concrete ? Uint8Array<SharedArrayBuffer> : SharedArrayBuffer;
  anchor_qs: T extends Concrete ? Uint8Array<SharedArrayBuffer> : SharedArrayBuffer;
  anchor_rs: T extends Concrete ? Uint8Array<SharedArrayBuffer> : SharedArrayBuffer;
  artifact_statuses: T extends Concrete ? Uint8Array<SharedArrayBuffer> : SharedArrayBuffer;
  zones: T extends Concrete ? Uint8Array<SharedArrayBuffer> : SharedArrayBuffer;
  foot_traffic: T extends Concrete ? Uint32Array<SharedArrayBuffer> : SharedArrayBuffer;
  car_traffic: T extends Concrete ? Uint32Array<SharedArrayBuffer> : SharedArrayBuffer;
  bus_traffic: T extends Concrete ? Uint32Array<SharedArrayBuffer> : SharedArrayBuffer;
  bicycle_traffic: T extends Concrete ? Uint32Array<SharedArrayBuffer> : SharedArrayBuffer;
  train_traffic: T extends Concrete ? Uint32Array<SharedArrayBuffer> : SharedArrayBuffer;
  subway_traffic: T extends Concrete ? Uint32Array<SharedArrayBuffer> : SharedArrayBuffer;
};

type CanvasWorker = {
  id: "artifacts" | "grid" | "zones" | "heatmap" | "ui";
  worker: Worker;
};

export type { CameraState, CityData, Concrete, CanvasWorker };
