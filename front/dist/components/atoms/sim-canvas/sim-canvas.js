import { BaseComponent } from "components";
import * as C from "constants";
import * as api from "api";
import css from "./sim-canvas.css" with { type: "css" };
const spritesheet_blob = await api.getSpritesheet();
class SimCanvas extends BaseComponent {
    static styles = [css];
    #city;
    #worker;
    #name = this.id;
    render() {
        this.html `<canvas></canvas>`;
    }
    async onComponentMounted() {
        this.#city = await this.channel.request("sim-city");
        const canvas = this.shadowRoot.querySelector("canvas");
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
                    const init_data = {
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
                    const init_data = {
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
                    const init_data = {
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
                    const init_data = {
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
                    const init_data = {
                        type: "init",
                        pixel_ratio,
                        offscreen_canvas,
                        shared_buffers,
                        sprites,
                    };
                    this.#worker = await this.#initCanvasWorkers(init_data);
                    console.log(this.#name);
                    return (this.channel.subscribe("ui-reset", this),
                        this.channel.subscribe("ui-render", this));
                }
            }
        }
        catch (reason) {
            console.error(reason);
        }
    }
    /**
     * For convenience, let's provide the correct spritesheet to the worker.
     */
    async #generateSprites(...sprite_types) {
        const all_sprites = {};
        for (const sprites of sprite_types) {
            const sprites_generator = this.#newSpritesGenerator(sprites);
            for await (const result of sprites_generator) {
                const image = result;
                const name = sprites_generator.next().value;
                all_sprites[name] = image;
            }
        }
        return all_sprites;
    }
    /**
     * "Cut" in the loaded global spritesheet and
     * generate an *ImgeBitmap* from it.
     */
    *#newSpritesGenerator(sprite_type) {
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
    async #initCanvasWorkers(init_data) {
        const worker_opts = { type: "module" };
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
    onUiRender(data) {
        const message = {
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
    get name() {
        return this.#name;
    }
}
window.customElements.define("sim-canvas", SimCanvas);
export { SimCanvas };
