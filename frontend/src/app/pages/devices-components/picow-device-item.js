import "./picow-power-button";
import "./picow-options-button";

import { CleanUp, globalStylesToShadowRoot, html } from "ui";
import { deviceEvents } from "../../../lib";

export default class PicowDeviceItem extends HTMLElement {
    /**
     * @param {object} options
     * @param {Device | null} [options.device]
     */
    constructor({ device = null }) {
        super();

        /**
         * @type {PicowStore}
         */
        this.store = document.querySelector(`ui-store`);

        /**
         * @type {Device | null}
         */
        this.device = device;

        this.cleanup = new CleanUp();

        this.picow = {
            root: this,

            /**
             * @param {Device} device
             */
            set(device) {
                this.root.device = device;

                const list = this.root.querySelector(`li.is-card`);
                list.setAttribute("data-server-addr", device.server.addr);

                if (!!device.color) {
                    this.root.style.setProperty(
                        "--current-color",
                        `rgb(${device.color[0] || 0}, ${
                            device.color[1] || 0
                        }, ${device.color[2] || 0})`
                    );
                }

                // ------------ //
                // Update Label //
                // ------------ //

                {
                    let primary = device.server.name || "";
                    let secondary = device.server.addr;
                    if (!primary) {
                        primary = device.server.addr;
                        secondary = "&nbsp;";
                    }

                    /**
                     * @type {import("ui").UILabel}
                     */
                    const label = this.root.querySelector(`ui-label`);
                    label.ui.primary = primary;
                    label.ui.secondary = secondary;
                }

                // ------------------- //
                // Update Power Button //
                // ------------------- //

                {
                    /**
                     * @type {import("./picow-power-button").default}
                     */
                    const power = this.root.querySelector(`picow-power-button`);
                    power.picow.set(device);
                }

                // --------------------- //
                // Update Options Button //
                // --------------------- //

                {
                    // TODO: Pass device data to the options button
                }

                // --------------------- //
                // Update Offline Marker //
                // --------------------- //

                {
                    const marker =
                        this.root.shadowRoot.querySelector(`.offline-marker`);

                    if (device.server.isOffline) marker.removeAttribute("hide");
                    else marker.setAttribute("hide", "");
                }
            },
        };

        this.#render();
    }

    #render() {
        this.classList.add("no-user-select");

        this.attachShadow({ mode: "open" });
        globalStylesToShadowRoot(this.shadowRoot);

        this.shadowRoot.innerHTML = html`
            <style>
                :host {
                    display: block;
                    position: relative;
                    border-radius: var(--ui-radius);
                }

                .current-color {
                    position: absolute;
                    top: var(--ui-spacing);
                    right: var(--ui-spacing);
                    bottom: var(--ui-spacing);
                    left: var(--ui-spacing);

                    border-radius: var(--ui-radius);

                    box-shadow: 0 0 8px 1px var(--current-color, transparent);

                    transition: box-shadow 0.35s linear;
                }

                .offline-marker {
                    position: absolute;
                    top: -0.25rem;
                    left: 50%;

                    color: var(--ui-destructive);

                    transform: translateX(-50%);
                }

                .offline-marker[hide] {
                    display: none;
                }
            </style>

            <div class="current-color"></div>
            <slot></slot>
            <ui-secondary class="offline-marker"></ui-secondary>
        `;

        this.innerHTML = html`
            <li class="is-card" style="cursor: pointer;">
                <ui-label>
                    <ui-flex-grid-row gap="0.25rem">
                        <ui-flex-grid-item>
                            <picow-power-button></picow-power-button>
                        </ui-flex-grid-item>

                        <ui-flex-grid-item>
                            <picow-options-button></picow-options-button>
                        </ui-flex-grid-item>
                    </ui-flex-grid-row>
                </ui-label>
            </li>
        `;

        /**
         * @type {HTMLLIElement}
         */
        const card = this.querySelector("li.is-card");
        card.onclick = async () => {
            // TODO: Open a color picker dialog to select a color
        };

        this.picow.set(this.device);
    }

    connectedCallback() {
        this.cleanup.add(
            deviceEvents.events.on("message", (device) => {
                if (device.server.addr !== this.device.server.addr) return;
                this.picow.set(device);
            })
        );
    }

    disconnectedCallback() {
        this.cleanup.run();
    }
}

console.debug(`Register the "picow-device-item"`);
customElements.define("picow-device-item", PicowDeviceItem);