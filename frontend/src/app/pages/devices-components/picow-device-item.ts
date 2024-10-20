import "./picow-options-button";
import "./picow-power-button";

import { css as CSS, html, LitElement, PropertyValues } from "lit";
import { customElement, property } from "lit/decorators.js";
import { CleanUp, globalStylesToShadowRoot } from "ui";

import { ws, WSEventsDevice } from "../../../lib/websocket";
import { PicowStore } from "../../../types";

/**
 * **Tag**: picow-options-button
 *
 * **Attributes**:
 *  - device: `WSEventsDevice` - [json]
 *  - hide: `boolean`
 */
@customElement("picow-device-item")
export class PicowDeviceItem extends LitElement {
    @property({ type: Object, attribute: "device", reflect: true })
    device?: WSEventsDevice;

    @property({ type: Boolean, attribute: "hide", reflect: true })
    hide: boolean = false;

    private store: PicowStore = document.querySelector(`ui-store`)!;
    private cleanup = new CleanUp();

    static get styles() {
        return CSS`
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

                color: hsl(var(--ui-hsl-destructive));

                transform: translateX(-50%);
            }

            .offline-marker[hide] {
                display: none;
            }
        `;
    }

    protected render() {
        let primary = this.device?.server.name || "";
        let secondary = this.device?.server.addr;
        if (!primary) {
            primary = this.device?.server?.addr || "&nbsp;";
            secondary = "&nbsp;";
        }

        return html`
            <div class="current-color"></div>

            <li
                class="is-card"
                style="cursor: pointer;"
                data-server-addr="${this.device?.server.addr || ""}"
                @click=${() => {
                    // TODO: Open a color picker dialog to select a color
                }}
            >
                <ui-label primary="${primary}" secondary="${secondary}">
                    <ui-flex-grid-row gap="0.25rem" align="center">
                        <ui-flex-grid-item>
                            <picow-power-button
                                device=${this.device}
                            ></picow-power-button>
                        </ui-flex-grid-item>

                        <ui-flex-grid-item>
                            <picow-options-button
                                device=${this.device}
                            ></picow-options-button>
                        </ui-flex-grid-item>
                    </ui-flex-grid-row>
                </ui-label>
            </li>

            <ui-secondary
                class="offline-marker"
                ?hide=${!this.device?.server.online}
            ></ui-secondary>
        `;
    }

    protected firstUpdated(_changedProperties: PropertyValues): void {
        globalStylesToShadowRoot(this.shadowRoot!);
        this.classList.add("no-user-select");

        if (this.device?.color) {
            this.style.setProperty(
                "--current-color",
                `rgb(${this.device.color[0] || 0}, ${
                    this.device.color[1] || 0
                }, ${this.device.color[2] || 0})`
            );
        }

        this.cleanup.add(
            ws.events.addListener("message-device", (data) => {
                if (data.server.addr !== this.device?.server.addr) return;
                this.device = data;

                if (!this.device.color) return;

                // Only update "devicesColor" store if color is not 0
                if (this.device.color.filter((c) => c > 0).length > 0) {
                    this.store.updateData("devicesColor", (data) => {
                        if (!this.device || !this.device?.color) return data;
                        data[this.device.server.addr] = this.device.color;
                        return data;
                    });
                }
            })
        );
    }

    disconnectedCallback(): void {
        super.disconnectedCallback();
        this.cleanup.run();
    }
}
