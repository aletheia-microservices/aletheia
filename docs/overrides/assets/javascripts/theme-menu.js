// Theme menu: replaces Material's theme button, which moves to the next theme on each click, with a button
// that opens the list of all themes (the palette in mkdocs.yml, with each toggle's name and icon) to pick
// one directly. Picking a theme clicks Material's own hidden input for it, so Material still applies the
// theme, saves it (localStorage "__palette") and restores it on the next page before drawing it.
// Keyboard: Enter, Space or the arrow keys open the list; the arrow keys, Home and End move through it;
// Enter or Space picks a theme; Escape closes the list; Tab leaves it.
(() => {
  const form = document.querySelector('[data-md-component="palette"]');
  const inputs = form ? [...form.querySelectorAll("input")] : [];
  if (!inputs.length || document.querySelector(".theme-menu")) return;

  const schemeOf = (input) => input.getAttribute("data-md-color-scheme");
  const nameOf = (input) => input.getAttribute("aria-label") || schemeOf(input);
  // Material puts each theme's icon in the label right after its input
  const iconOf = (input) => input.nextElementSibling?.querySelector("svg")?.cloneNode(true);

  // Material restores the saved theme by its position in the palette, so when the palette changes (a theme
  // added, removed or moved), point the saved position back at the saved theme. This runs before Material
  // reads it, which it does once the page has loaded.
  try {
    const saved = __md_get("__palette");
    const index = inputs.findIndex((input) => schemeOf(input) === saved?.color?.scheme);
    if (index >= 0 && index !== saved.index) __md_set("__palette", { ...saved, index });
  } catch {}

  const menu = document.createElement("div");
  menu.className = "md-header__option theme-menu";

  const button = document.createElement("button");
  button.type = "button";
  button.className = "md-header__button md-icon theme-menu__button";
  button.setAttribute("aria-haspopup", "true");
  button.setAttribute("aria-expanded", "false");
  button.setAttribute("aria-controls", "theme-menu__list");

  const list = document.createElement("ul");
  list.className = "theme-menu__list";
  list.id = "theme-menu__list";
  list.setAttribute("role", "menu");
  list.setAttribute("aria-label", "Theme");
  list.hidden = true;

  const items = inputs.map((input, index) => {
    const item = document.createElement("button");
    item.type = "button";
    item.tabIndex = -1;
    item.setAttribute("role", "menuitemradio");
    const icon = iconOf(input);
    if (icon) item.append(icon);
    const name = document.createElement("span");
    name.textContent = nameOf(input);
    item.append(name);
    item.addEventListener("click", () => select(index));
    const li = document.createElement("li");
    li.setAttribute("role", "none");
    li.append(item);
    list.append(li);
    return item;
  });

  // the theme on the page, which Material sets on the body
  const current = () =>
    Math.max(0, inputs.findIndex((input) => schemeOf(input) === document.body.getAttribute("data-md-color-scheme")));

  // show the current theme's icon on the button and tick the theme in the list
  const update = () => {
    const index = current();
    button.replaceChildren(iconOf(inputs[index]) || document.createTextNode("Theme"));
    button.title = `Theme: ${nameOf(inputs[index])}`;
    button.setAttribute("aria-label", button.title);
    items.forEach((item, i) => item.setAttribute("aria-checked", String(i === index)));
  };

  const open = (focus) => {
    list.hidden = false;
    button.setAttribute("aria-expanded", "true");
    items[focus].focus();
  };

  const close = (focusButton) => {
    if (list.hidden) return;
    list.hidden = true;
    button.setAttribute("aria-expanded", "false");
    if (focusButton) button.focus();
  };

  // clicking an input that is already checked does nothing, so picking the current theme changes nothing
  const select = (index) => {
    inputs[index].click();
    close(true);
  };

  button.addEventListener("click", () => (list.hidden ? open(current()) : close(false)));
  button.addEventListener("keydown", (event) => {
    if (event.key !== "ArrowDown" && event.key !== "ArrowUp") return;
    event.preventDefault();
    open(event.key === "ArrowDown" ? 0 : items.length - 1);
  });

  list.addEventListener("keydown", (event) => {
    const index = items.indexOf(document.activeElement);
    const next = { ArrowDown: index + 1, ArrowUp: index - 1, Home: 0, End: items.length - 1 }[event.key];
    if (next !== undefined) {
      event.preventDefault();
      items[(next + items.length) % items.length].focus();
    } else if (event.key === "Escape") {
      event.preventDefault();
      close(true);
    } else if (event.key === "Tab") {
      close(false);
    }
  });

  // a click anywhere else closes the list; pointerdown, so that a click on the button itself still toggles it
  document.addEventListener("pointerdown", (event) => {
    if (!menu.contains(event.target)) close(false);
  });

  // Material changes the body's scheme also when it restores the saved theme
  new MutationObserver(update).observe(document.body, { attributes: true, attributeFilter: ["data-md-color-scheme"] });

  menu.append(button, list);
  form.after(menu);
  form.style.display = "none"; // Material's own button; hidden only here, so it stays if this script fails
  update();
})();
