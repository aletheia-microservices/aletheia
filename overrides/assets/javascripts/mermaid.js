// Material renders ```mermaid blocks through a global `mermaid` object and only loads the library itself
// when that global is missing. This stand-in loads a pinned version on the first diagram and adjusts it:
//   - keeps the colors given to sequence diagram boxes (Material paints every box the same color)
//   - outlines each box in its own color, instead of Mermaid's purple outline and drop shadow
//   - draws participants, lifelines and frames (e.g. par) in the page's gray, instead of Mermaid's purple
//     and Material's accent color
//   - draws the participants of sequence diagrams only at the top, not again at the bottom
//   - makes room for the box that diagram-steps.js draws around text marked focus-N in a label
//   - fires a "mermaid-rendered" event on each diagram once it is on the page, with the shadow root that
//     holds its SVG, which Material keeps closed, so other scripts can change it (e.g. diagram-steps.js)
(() => {
  const SRC = "https://unpkg.com/mermaid@11.17.2/dist/mermaid.min.js";
  let config = {};
  let loaded;

  // !important beats Mermaid's own theme rules, which are scoped to each diagram's id
  const NEUTRAL_CSS = `
    rect.actor, g.actor path, .actor-man circle, .actor-man line {
      stroke: var(--md-default-fg-color--light) !important;
    }
    .actor-line { stroke: var(--md-default-fg-color--lighter) !important; }
    .loopLine, .labelBox { stroke: var(--md-default-fg-color--light) !important; }
    .labelText, .labelText > tspan, .loopText, .loopText > tspan {
      fill: var(--md-default-fg-color--light) !important;
    }
  `;

  // Mermaid sizes each label with these styles, before diagram-steps.js adds its own, so the room for the
  // box goes here; the box stays transparent until diagram-steps.js colors it at its step
  const FOCUS_CSS = `
    :not(.node)[class*="focus-"] {
      display: inline-block;
      padding: 0 4px;
      border: 1.5px solid transparent;
      border-radius: 3px;
    }
  `;

  const adjust = (config) => ({
    ...config,
    themeCSS: (config.themeCSS || "").replace(/rect\.rect\{fill:[^;}]*;/, "rect.rect{") + NEUTRAL_CSS + FOCUS_CSS,
    sequence: { ...config.sequence, mirrorActors: false },
  });

  // a box is drawn as <rect fill="rgba(r, g, b, a)" class="rect">; an inline style beats Mermaid's CSS
  const outlineBoxes = (result) => ({
    ...result,
    svg: result.svg.replace(/<rect\b[^>]*\bclass="rect"[^>]*>/g, (rect) => {
      const rgb = rect.match(/fill="rgba?\((\d+),\s*(\d+),\s*(\d+)/);
      const stroke = rgb ? `rgba(${rgb[1]}, ${rgb[2]}, ${rgb[3]}, 0.5)` : "none";
      return rect.replace("<rect", `<rect style="stroke: ${stroke}; filter: none"`);
    }),
  });

  // Material calls the fn of a render result with the diagram's shadow root, after putting the diagram
  // on the page; the event bubbles up from the diagram, so listeners can tell which diagram it is
  const announce = (result) => ({
    ...result,
    fn: (root) => root.host.dispatchEvent(new CustomEvent("mermaid-rendered", { bubbles: true, detail: { root } })),
  });

  // loading the script replaces this stand-in with the real library in window.mermaid,
  // which Material then calls directly, so its render() is wrapped as well
  const load = () => loaded ??= new Promise((resolve, reject) => {
    const script = document.createElement("script");
    script.src = SRC;
    script.onload = () => {
      const mermaid = window.mermaid;
      const render = mermaid.render.bind(mermaid);
      mermaid.render = (...args) => render(...args).then(outlineBoxes).then(announce);
      mermaid.initialize(adjust(config));
      resolve(mermaid);
    };
    script.onerror = reject;
    document.head.appendChild(script);
  });

  window.mermaid = {
    initialize: (c) => { config = c; },
    render: (...args) => load().then((mermaid) => mermaid.render(...args)),
  };
})();
