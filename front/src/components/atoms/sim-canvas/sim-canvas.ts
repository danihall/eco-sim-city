import type { CanvasWorker, SimCity } from "components";
import type {
  Sprites,
  InitData,
  WorkerArtifactsInit,
  WorkerGridInit,
  WorkerHeatmapInit,
  WorkerUiInit,
  WorkerZonesInit,
} from "./workers/_index.js";
import type { UIData, RenderUIData } from "./workers/canvas_ui.js";

import { BaseComponent } from "components";
import * as C from "constants";
import * as api from "api";
import css from "./sim-canvas.css" with { type: "css" };

declare module "components" {
  interface ComponentTagNameMap {
    "sim-canvas": SimCanvas;
  }
}

const spritesheet_blob = await api.getSpritesheet();

class SimCanvas extends BaseComponent {
  static styles = [css];

  #city!: SimCity;
  #worker!: CanvasWorker["worker"];
  #name: CanvasWorker["id"] = this.id as CanvasWorker["id"];

  render() {
    this.html`<canvas></canvas>`;
  }

  async onComponentMounted() {
    this.#city = await this.channel.request("sim-city");

    const canvas = this.shadowRoot!.querySelector("canvas")!;
    const { width, height } = this.#city.getBoundingClientRect();
    const pixel_ratio = window.devicePixelRatio || 1;

    canvas.width = Math.trunc(width * pixel_ratio);
    canvas.height = Math.trunc(height * pixel_ratio);
    canvas.style.width = `${width}px`;
    canvas.style.height = `${height}px`;

    const offscreen_canvas = canvas.transferControlToOffscreen();
    const shared_buffers = this.#city.sharedBuffers;

    try {
      switch (this.name) {
        case "grid": {
          const init_data: WorkerGridInit = {
            type: "init",
            pixel_ratio,
            offscreen_canvas,
          };

          this.#worker = await this.#initCanvasWorkers(init_data);
          console.log(this.#name);
          this.onGridRender();
          return this.channel.subscribe("grid-render", this);
        }
        case "artifacts": {
          const sprites = await this.#generateSprites("artifacts");
          const init_data: WorkerArtifactsInit = {
            type: "init",
            pixel_ratio,
            offscreen_canvas,
            shared_buffers,
            sprites,
          };

          this.#worker = await this.#initCanvasWorkers(init_data);
          console.log(this.#name);
          this.onRender();
          return this.channel.subscribe("render", this);
        }
        case "zones": {
          const sprites = await this.#generateSprites("zones");
          const init_data: WorkerZonesInit = {
            type: "init",
            pixel_ratio,
            offscreen_canvas,
            shared_buffers,
            sprites,
          };

          this.#worker = await this.#initCanvasWorkers(init_data);
          console.log(this.#name);
          this.onRender();
          return this.channel.subscribe("render", this);
        }
        case "heatmap": {
          const sprites = await this.#generateSprites("heatmap");
          const init_data: WorkerHeatmapInit = {
            type: "init",
            pixel_ratio,
            offscreen_canvas,
            shared_buffers,
            sprites,
          };

          this.#worker = await this.#initCanvasWorkers(init_data);
          console.log(this.#name);
          return this.channel.subscribe("heatmap-render", this);
        }
        case "ui": {
          const sprites = await this.#generateSprites("zones", "artifacts");
          const init_data: WorkerUiInit = {
            type: "init",
            pixel_ratio,
            offscreen_canvas,
            shared_buffers,
            sprites,
          };

          this.#worker = await this.#initCanvasWorkers(init_data);
          console.log(this.#name);
          return (
            this.channel.subscribe("ui-reset", this),
            this.channel.subscribe("ui-render", this)
          );
        }
      }
    } catch (reason) {
      console.error(reason);
    }
  }

  /**
   * For convenience, let's provide the correct spritesheet to the worker.
   */
  async #generateSprites<S extends keyof typeof C.SPRITESHEET>(...sprite_types: S[]) {
    const all_sprites = {} as Sprites<S>;

    for (const sprites of sprite_types) {
      const sprites_generator = this.#newSpritesGenerator(sprites);

      for await (const result of sprites_generator) {
        const image = result as ImageBitmap;
        const name = sprites_generator.next().value as keyof Sprites<S>;

        all_sprites[name] = image;
      }
    }

    return all_sprites;
  }

  /**
   * "Cut" in the loaded global spritesheet and
   * generate an *ImgeBitmap* from it.
   */
  *#newSpritesGenerator<S extends keyof typeof C.SPRITESHEET>(sprite_type: S) {
    const entries = Object.entries(C.SPRITESHEET[sprite_type]);

    for (const sprite_props of entries) {
      const [name, { x, y, width, height }] = sprite_props;

      yield createImageBitmap(spritesheet_blob, x, y, width || C.HEX_WIDTH, height || C.HEX_HEIGHT);
      yield name;
    }
  }

  /**
   * Launch worker with *Transferable* objects if needed.
   */
  async #initCanvasWorkers<
    M extends
      WorkerArtifactsInit | WorkerGridInit | WorkerZonesInit | WorkerHeatmapInit | WorkerUiInit,
  >(init_data: M): Promise<CanvasWorker["worker"]> {
    const worker_opts: WorkerOptions = { type: "module" };
    const url = new URL(`./workers/canvas_${this.#name}.js`, import.meta.url);
    const worker = new Worker(url, worker_opts);

    return new Promise((resolve, reject) => {
      worker.onmessage = (message) => {
        worker.onmessage = null;

        if (message.data?.ready) {
          worker.addEventListener("message", console.log);
          worker.addEventListener("error", console.error);
          resolve(worker);
        }

        reject();
      };

      worker.postMessage(init_data, [
        init_data.offscreen_canvas,
        ...("sprites" in init_data ? Object.values(init_data.sprites) : []),
      ]);
    });
  }

  /**
   * Should clear UI canvas.
   */
  onUiReset() {
    this.#worker.postMessage({ type: "reset", camera: this.#city.camera });
  }

  /**
   * Should render new stuff in UI canvas.
   */
  onUiRender(data: UIData) {
    const message: RenderUIData = {
      type: "render",
      camera: this.#city.camera,
      ...data,
    };

    this.#worker.postMessage(message, [
      data.qs_buffer,
      data.rs_buffer,
      data.item_ids_buffer,
      data.occupied_buffer,
    ]);
  }

  /**
   * Should render grid canvas.
   */
  onGridRender() {
    this.#worker.postMessage({ type: "render", camera: this.#city.camera });
  }

  /**
   * Should render in heatmap canvas.
   */
  onHeatmapRender() {
    this.#worker.postMessage({ type: "render", camera: this.#city.camera });
  }

  /**
   * Should render in artifacts or zones canvas.
   */
  onRender() {
    this.#worker.postMessage({ type: "render", camera: this.#city.camera });
  }

  get name(): Readonly<CanvasWorker["id"] | void> {
    return this.#name;
  }
}

window.customElements.define("sim-canvas", SimCanvas);

export { SimCanvas };
