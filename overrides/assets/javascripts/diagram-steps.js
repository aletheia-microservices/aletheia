// Steps through a Mermaid diagram that sits, with an ordered list after it, inside a
// <div class="diagram-steps" markdown> block: the diagram shows one item of the list at a time, and the
// parts that a step adds fade in when the diagram reaches that step. Tag those parts in the diagram:
//   - text in a node label: class='step-N' on an element of the label, e.g. <b class='step-2'>...</b>
//   - an edge and its label: an edge id that starts with stepN_, e.g. a step2_into@-.-> b
//   - a node that keeps the style of its classDef only from step N on: class a from-N, e.g. a tainted
//     value; before step N, it is drawn like a plain node, with no fill and a gray outline
//   - text to box only at step N: class='focus-N' on an element of the label, e.g. <span class='focus-N'>
//     around several lines; mermaid.js makes room for the box, which stays transparent until step N
//   - a node to outline only at step N, in the color of the box: class a focus-N
//   At a step that boxes or outlines something, the rest of the diagram fades, so that it stands out.
// It adds previous, next and play buttons under the diagram, and clicking an item of the list shows its
// step. The steps play once when the diagram first scrolls into view, unless the reader prefers reduced
// motion. Without JavaScript, the diagram shows all the steps at once and the list stays as it is.
// mermaid.js fires "mermaid-rendered" with the shadow root that holds the diagram, which Material keeps
// closed; the styles below go into that shadow root, since the page's stylesheets do not reach it.
(() => {
  const DELAY = 2500; // milliseconds between steps while playing
  const FOCUS = "#b58900"; // Solarized yellow, for the boxes and outlines of focus-N

  // the fading of focus-dim leaves out the focused nodes and text: a node that holds focused text fades
  // only its shape and the rest of its text, which setup() wraps in spans so that each piece can fade
  const STYLE = `
    .step-hidden { opacity: 0 !important; }
    .focus-active { border-color: ${FOCUS} !important; }
    .focus-dim .node:not(.focus-active):not(:has(.focus-active)),
    .focus-dim .node:has(.focus-active) > :not(.label),
    .focus-dim .node:has(.focus-active) .label :not(:has(*)):not(.focus-active *),
    .focus-dim .edgePaths > *, .focus-dim .edgeLabels > * { opacity: 0.3; }
    .node, .node > *, .label *, .edgePaths > *, .edgeLabels > *, [data-id^="step"] {
      transition: opacity 0.4s, fill 0.4s, stroke 0.4s, border-color 0.4s;
    }
    @media (prefers-reduced-motion: reduce) {
      .node, .node > *, .label *, .edgePaths > *, .edgeLabels > *, [data-id^="step"] { transition: none; }
    }
  `;

  // the number N in class="... <prefix>-N ..." or data-id="<prefix>N_...", or 0 if there is none
  const tagOf = (el, prefix) => {
    const match = el.getAttribute("class")?.match(new RegExp(`\\b${prefix}-(\\d+)\\b`))
      || el.getAttribute("data-id")?.match(new RegExp(`^${prefix}(\\d+)_`));
    return match ? Number(match[1]) : 0;
  };

  const tagged = (svg, selector, prefix) => [...svg.querySelectorAll(selector)]
    .map((el) => [el, tagOf(el, prefix)])
    .filter(([, n]) => n > 0);

  const setup = (container, root) => {
    const svg = root.querySelector("svg");
    const items = [...container.querySelectorAll(":scope > ol > li")];
    if (!svg || items.length === 0) return;

    const style = document.createElement("style");
    style.textContent = STYLE;
    root.append(style);

    const parts = tagged(svg, '[class*="step-"], [data-id^="step"]', "step");
    const focused = tagged(svg, '[class*="focus-"]', "focus");
    const focusSteps = new Set(focused.map(([, k]) => k));
    for (const label of new Set(focused.map(([el]) => el.closest(".label")).filter(Boolean))) {
      const walker = document.createTreeWalker(label, NodeFilter.SHOW_TEXT);
      const texts = [];
      while (walker.nextNode()) if (walker.currentNode.textContent.trim()) texts.push(walker.currentNode);
      for (const text of texts) {
        const span = document.createElement("span");
        span.style.color = "inherit"; // Mermaid gives every span of a label the text color, e.g. over a red taint
        text.replaceWith(span);
        span.append(text);
      }
    }
    // Mermaid writes the classDef of a node into the style attribute of its shape, with !important, which no
    // stylesheet can override: a node is drawn plain before its step, or outlined at its focus step, by
    // changing that attribute instead
    const shapes = [...svg.querySelectorAll('.node[class*="from-"], .node[class*="focus-"]')].flatMap((node) =>
      [...node.children]
        .filter((child) => !child.classList.contains("label"))
        .map((shape) => [shape, tagOf(node, "from"), tagOf(node, "focus"), shape.getAttribute("style") ?? ""]));

    container.querySelector(":scope > .diagram-steps__controls")?.remove();
    const controls = document.createElement("div");
    controls.className = "diagram-steps__controls";
    controls.innerHTML = `
      <button type="button" data-go="prev" aria-label="Previous step">◀</button>
      <span aria-live="polite"></span>
      <button type="button" data-go="next" aria-label="Next step">▶</button>
      <button type="button" data-go="play"></button>`;
    root.host.after(controls);
    const [prev, status, next, play] = controls.children;

    let step = 1;
    let timer = null;
    let touched = false; // the reader moved through the steps, so they do not play on their own

    const show = (n) => {
      step = Math.min(items.length, Math.max(1, n));
      for (const [el, k] of parts) el.classList.toggle("step-hidden", k > step);
      for (const [shape, from, focus, original] of shapes) {
        shape.setAttribute("style", original);
        if (from > step) {
          shape.style.setProperty("fill", "transparent", "important");
          shape.style.setProperty("stroke", "rgba(128, 128, 128, 0.6)", "important");
        }
        if (focus === step) {
          shape.style.setProperty("stroke", FOCUS, "important");
          shape.style.setProperty("stroke-width", "2px", "important");
        }
      }
      for (const [el, k] of focused) el.classList.toggle("focus-active", k === step);
      svg.classList.toggle("focus-dim", focusSteps.has(step));
      items.forEach((item, i) => item.classList.toggle("diagram-steps__active", i + 1 === step));
      status.textContent = `Step ${step} of ${items.length}`;
      prev.disabled = step === 1;
      next.disabled = step === items.length;
    };

    const stop = () => {
      clearInterval(timer);
      timer = null;
      play.textContent = "Play";
    };

    const start = () => {
      if (step === items.length) show(1);
      play.textContent = "Pause";
      timer = setInterval(() => {
        if (!container.isConnected) return stop();
        show(step + 1);
        if (step === items.length) stop();
      }, DELAY);
    };

    controls.addEventListener("click", (event) => {
      const go = event.target.closest("button")?.dataset.go;
      if (!go) return;
      touched = true;
      if (go === "play") return timer ? stop() : start();
      stop();
      show(go === "next" ? step + 1 : step - 1);
    });

    items.forEach((item, i) => item.addEventListener("click", () => {
      touched = true;
      stop();
      show(i + 1);
    }));

    stop();
    show(1);

    if (matchMedia("(prefers-reduced-motion: reduce)").matches || !("IntersectionObserver" in window)) return;
    const observer = new IntersectionObserver((entries) => {
      if (!entries.some((entry) => entry.isIntersecting)) return;
      observer.disconnect();
      if (!touched) start();
    }, { threshold: 0.5 });
    observer.observe(root.host);
  };

  document.addEventListener("mermaid-rendered", (event) => {
    const container = event.target.closest(".diagram-steps");
    if (container) setup(container, event.detail.root);
  });
})();
