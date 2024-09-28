import { html, UIIconButton } from "ui";

export default class PicowOptionsButton extends UIIconButton {
    constructor() {
        super();
        this.picow = {};
        this.#render();
    }

    #render() {}
}

console.debug(`Register the "picow-options-button"`);
customElements.define("picow-options-button", PicowOptionsButton);
