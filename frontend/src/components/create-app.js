import { html, styles } from "ui";
import createAppBar from "./create-app-bar";

/**
 * @typedef {{
 *  element: HTMLDivElement;
 * }} App
 */

/**
 * @returns {App}
 */
export default function () {
    const el = document.createElement("div");

    el.innerHTML = html`
        <ui-theme-handler mode="dark"></ui-theme-handler>
        <ui-store storageprefix="picow:" storage></ui-store>

        <div class="app-bar"></div>
        <div class="drawer"></div>

        <ui-container
            style="${styles({
                width: "100%",
                height: "100%",
            })}"
        >
            <ui-stack-layout></ui-stack-layout>
        </ui-container>

        <ui-alerts></ui-alerts>
    `;

    const appBar = createAppBar();
    el.querySelector(`div.app-bar`).replaceWith(appBar.element);

    // TODO: Create the drawer
    const drawer = createDrawer();
    el.querySelector(`div.drawer`).replaceWith(drawer.element);

    return {
        element: el,
    };
}
