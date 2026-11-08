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
//   - in a sequence diagram, whose messages cannot be tagged: a comment %% step N on its own line, before
//     the messages, notes and blocks (e.g. rect, par) of step N, up to the next such comment; the parts
//     before the first comment belong to step 1, and a block belongs whole to the step where it starts
//     (nested rect blocks are not supported). A sequence diagram shows all its steps at once, and a band
//     in the color of the step's item in the list highlights the parts of the current step
//   At a step that boxes or outlines something, the rest of the diagram fades, so that it stands out.
// It adds previous and next buttons under the diagram, and clicking an item of the list shows its step;
// the diagram only changes step when the reader asks it to. Without JavaScript, the diagram shows all the
// steps at once and the list stays as it is.
// mermaid.js fires "mermaid-rendered" with the shadow root that holds the diagram, which Material keeps
// closed, and the diagram's source; the styles below go into that shadow root, since the page's
// stylesheets do not reach it.
(() => {
  // Solarized yellow, for the boxes and outlines of focus-N, and the gray outline of a from-N node before
  // step N; theme-contrast.css sets stronger colors for the high contrast themes
  const FOCUS = "var(--diagram-focus-color, #b58900)";
  const IDLE = "var(--diagram-idle-color, rgba(128, 128, 128, 0.6))";
  // the background of the band behind what a step adds to a sequence diagram: the background of the
  // step's item in the list (see page.css)
  const BAND = "color-mix(in srgb, var(--md-accent-fg-color) 12%, transparent)";
  const BAND_PADDING = 6; // pixels between the band and what it highlights

  // the fading of focus-dim leaves out the focused nodes and text: a node that holds focused text fades
  // only its shape and the rest of its text, which setup() wraps in spans so that each piece can fade
  const STYLE = `
    .step-hidden { opacity: 0 !important; }
    .focus-active { border-color: ${FOCUS} !important; }
    .focus-dim .node:not(.focus-active):not(:has(.focus-active)),
    .focus-dim .node:has(.focus-active) > :not(.label),
    .focus-dim .node:has(.focus-active) .label :not(:has(*)):not(.focus-active *),
    .focus-dim .edgePaths > *, .focus-dim .edgeLabels > * { opacity: 0.3; }
    .steps-band { fill: ${BAND}; opacity: 0; }
    .steps-band--current { opacity: 1; }
    .node, .node > *, .label *, .edgePaths > *, .edgeLabels > *, [data-id^="step"], [class*="step-"],
    .steps-band {
      transition: opacity 0.4s, fill 0.4s, stroke 0.4s, border-color 0.4s;
    }
    @media (prefers-reduced-motion: reduce) {
      .node, .node > *, .label *, .edgePaths > *, .edgeLabels > *, [data-id^="step"], [class*="step-"],
      .steps-band {
        transition: none;
      }
    }
  `;

  // Mermaid numbers the statements of a sequence diagram that it draws, in the order of the source:
  // messages, notes, and the start, the middle (e.g. and, else) and the end of blocks such as rect and par.
  // A box only groups participants, so neither it nor its end is numbered. Returns the step of each
  // numbered statement, from the %% step N comments, where a block's middle and end take the step of its
  // start, and the step of each rect block, in the order of their ends
  const sequenceSteps = (source) => {
    const steps = [];
    const rects = [];
    const open = []; // the kind and step of each block that has not ended yet
    let step = 1;
    for (const line of source.split("\n").map((l) => l.trim())) {
      const marker = line.match(/^%%\s*step\s+(\d+)$/i);
      if (marker) {
        step = Number(marker[1]);
      } else if (!line || line.startsWith("%%")) {
        continue;
      } else if (/^box\b/i.test(line)) {
        open.push({ kind: "box" });
      } else if (/^(rect|par|par_over|loop|alt|opt|critical|break)\b/i.test(line)) {
        open.push({ kind: line.split(/\s/)[0].toLowerCase(), step });
        steps.push(step);
      } else if (/^(and|else|option)\b/i.test(line)) {
        steps.push(open.at(-1)?.step ?? step);
      } else if (/^end$/i.test(line)) {
        const block = open.pop();
        if (block?.kind === "box") continue;
        steps.push(block?.step ?? step);
        if (block?.kind === "rect") rects.push(block.step);
      } else if (/^(note|activate|deactivate)\b/i.test(line) || /^[^:]*?(-{1,2}(>>|>|x|\))|<<-{1,2}>>)/.test(line)) {
        steps.push(step); // a note, an activation or a message (A->>B: text)
      }
    }
    return { steps, rects };
  };

  // returns the elements of each step of a sequence diagram. Mermaid draws statement K as an element with
  // data-id="iK": a line for a message (a path if it goes back to its sender), whose text it draws just
  // before it, or a group for a note or for the frame of a block, which gets the number of the block's
  // end. Rect blocks have no data-id, and are drawn after the boxes, as <rect class="rect"> too
  const sequenceParts = (svg, source) => {
    const { steps, rects } = sequenceSteps(source);
    const byStep = new Map();
    const tag = (k, elements) => byStep.set(k, [...(byStep.get(k) ?? []), ...elements]);
    let texts = [];
    for (const el of svg.querySelectorAll(".messageText, [data-id]")) {
      if (el.classList.contains("messageText")) {
        texts.push(el);
        continue;
      }
      const number = el.getAttribute("data-id").match(/^i(\d+)$/);
      const k = number && steps[Number(number[1])];
      if (!k) continue;
      if (el.tagName === "g") {
        tag(k, [el]);
      } else {
        tag(k, [el, ...texts]);
        texts = [];
      }
    }
    const boxes = source.split("\n").filter((line) => /^\s*box\b/i.test(line)).length;
    const drawn = [...svg.querySelectorAll("rect.rect")].slice(boxes);
    rects.forEach((k, i) => drawn[i] && tag(k, [drawn[i]]));
    return byStep;
  };

  // a band behind the given elements, behind everything else in the diagram, or null if they have no
  // size yet (e.g. while the diagram is hidden)
  const band = (svg, elements) => {
    const boxes = elements.map((el) => el.getBBox()).filter((box) => box.width || box.height);
    if (boxes.length === 0) return null;
    const left = Math.min(...boxes.map((box) => box.x)) - BAND_PADDING;
    const top = Math.min(...boxes.map((box) => box.y)) - BAND_PADDING;
    const right = Math.max(...boxes.map((box) => box.x + box.width)) + BAND_PADDING;
    const bottom = Math.max(...boxes.map((box) => box.y + box.height)) + BAND_PADDING;
    const rect = document.createElementNS("http://www.w3.org/2000/svg", "rect");
    for (const [name, value] of Object.entries({ x: left, y: top, width: right - left, height: bottom - top, rx: 4 })) {
      rect.setAttribute(name, value);
    }
    rect.classList.add("steps-band");
    svg.prepend(rect);
    return rect;
  };

  // the number N in class="... <prefix>-N ..." or data-id="<prefix>N_...", or 0 if there is none
  const tagOf = (el, prefix) => {
    const match = el.getAttribute("class")?.match(new RegExp(`\\b${prefix}-(\\d+)\\b`))
      || el.getAttribute("data-id")?.match(new RegExp(`^${prefix}(\\d+)_`));
    return match ? Number(match[1]) : 0;
  };

  const tagged = (svg, selector, prefix) => [...svg.querySelectorAll(selector)]
    .map((el) => [el, tagOf(el, prefix)])
    .filter(([, n]) => n > 0);

  const setup = (container, root, source) => {
    const svg = root.querySelector("svg");
    const items = [...container.querySelectorAll(":scope > ol > li")];
    if (!svg || items.length === 0) return;

    const style = document.createElement("style");
    style.textContent = STYLE;
    root.append(style);

    const sequence = svg.getAttribute("aria-roledescription") === "sequence" ? sequenceParts(svg, source) : new Map();
    const bands = new Map(); // the band of each step of a sequence diagram, drawn when first needed

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
      <button type="button" data-go="next" aria-label="Next step">▶</button>`;
    root.host.after(controls);
    const [prev, status, next] = controls.children;

    let step = 1;

    const show = (n) => {
      step = Math.min(items.length, Math.max(1, n));
      for (const [el, k] of parts) el.classList.toggle("step-hidden", k > step);
      for (const [k, elements] of sequence) {
        if (!bands.get(k)) bands.set(k, band(svg, elements));
        bands.get(k)?.classList.toggle("steps-band--current", k === step);
      }
      for (const [shape, from, focus, original] of shapes) {
        shape.setAttribute("style", original);
        if (from > step) {
          shape.style.setProperty("fill", "transparent", "important");
          shape.style.setProperty("stroke", IDLE, "important");
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

    controls.addEventListener("click", (event) => {
      const go = event.target.closest("button")?.dataset.go;
      if (go) show(go === "next" ? step + 1 : step - 1);
    });

    items.forEach((item, i) => item.addEventListener("click", () => show(i + 1)));

    show(1);
  };

  document.addEventListener("mermaid-rendered", (event) => {
    const container = event.target.closest(".diagram-steps");
    if (container) setup(container, event.detail.root, event.detail.source ?? "");
  });
})();
