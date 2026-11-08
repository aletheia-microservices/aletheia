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
//     (nested rect blocks are not supported). A sequence diagram shows all its steps at once, and fades
//     the parts of the other steps, so that the current step stands out without adding a color of its own
//     to the diagram, whose blocks already have colors with a meaning (e.g. red for the problem)
//   At a step that boxes or outlines something, the rest of the diagram fades, so that it stands out.
// It adds previous and next buttons under the diagram, and an "All steps" button that shows the whole diagram
// at once, with nothing hidden or faded, until the reader picks a step again; clicking an item of the list
// shows its step;
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

  // the fading of focus-dim leaves out the focused nodes and text: a node that holds focused text fades
  // only its shape and the rest of its text, which setup() wraps in spans so that each piece can fade
  const STYLE = `
    .step-hidden { opacity: 0 !important; }
    .focus-active { border-color: ${FOCUS} !important; }
    .focus-dim .node:not(.focus-active):not(:has(.focus-active)),
    .focus-dim .node:has(.focus-active) > :not(.label),
    .focus-dim .node:has(.focus-active) .label :not(:has(*)):not(.focus-active *),
    .focus-dim .edgePaths > *, .focus-dim .edgeLabels > * { opacity: 0.3; }
    .steps-faded { opacity: 0.4; }
    .node, .node > *, .label *, .edgePaths > *, .edgeLabels > *, [data-id^="step"], [class*="step-"],
    [data-id], .messageText, rect.rect {
      transition: opacity 0.4s, fill 0.4s, stroke 0.4s, border-color 0.4s;
    }
    @media (prefers-reduced-motion: reduce) {
      .node, .node > *, .label *, .edgePaths > *, .edgeLabels > *, [data-id^="step"], [class*="step-"],
      [data-id], .messageText, rect.rect {
        transition: none;
      }
    }
  `;

  // Mermaid numbers the statements of a sequence diagram that it draws, in the order of the source:
  // messages, notes, and the start, the middle (e.g. and, else) and the end of blocks such as rect and par.
  // A box only groups participants, so neither it nor its end is numbered. Returns the step of each
  // numbered statement, from the %% step N comments, where a block's middle and end take the step of its
  // start, and the step of each rect block, in the order of their starts
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
        const kind = line.split(/\s/)[0].toLowerCase();
        open.push({ kind, step });
        steps.push(step);
        if (kind === "rect") rects.push(step);
      } else if (/^(and|else|option)\b/i.test(line)) {
        steps.push(open.at(-1)?.step ?? step);
      } else if (/^end$/i.test(line)) {
        const block = open.pop();
        if (block?.kind === "box") continue;
        steps.push(block?.step ?? step);
      } else if (/^(note|activate|deactivate)\b/i.test(line) || /^[^:]*?(-{1,2}(>>|>|x|\))|<<-{1,2}>>)/.test(line)) {
        steps.push(step); // a note, an activation or a message (A->>B: text)
      }
    }
    return { steps, rects };
  };

  // returns the elements of each step of a sequence diagram. Mermaid draws statement K as an element with
  // data-id="iK": a line for a message (a path if it goes back to its sender), whose text it draws just
  // before it, or a group for a note or for the frame of a block, which gets the number of the block's
  // end. Rect blocks have no data-id, and are drawn after the boxes, as <rect class="rect"> too, but in
  // reverse order (the last one first), so they are matched to the rect blocks of the source, which are
  // in the order of their starts, by their position from the top of the diagram
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
    const drawn = [...svg.querySelectorAll("rect.rect")].slice(boxes)
      .sort((a, b) => parseFloat(a.getAttribute("y")) - parseFloat(b.getAttribute("y")));
    rects.forEach((k, i) => drawn[i] && tag(k, [drawn[i]]));
    return byStep;
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
      <button type="button" data-go="all" aria-pressed="false">All steps</button>`;
    root.host.after(controls);
    const [prev, status, next, allButton] = controls.children;

    let step = 1;
    let all = false; // whether the diagram shows all its steps at once

    const show = (n, showAll = false) => {
      all = showAll;
      step = Math.min(items.length, Math.max(1, n));
      const shown = all ? items.length : step; // the last step whose parts are shown
      for (const [el, k] of parts) el.classList.toggle("step-hidden", k > shown);
      for (const [k, elements] of sequence) {
        for (const el of elements) el.classList.toggle("steps-faded", !all && k !== step);
      }
      for (const [shape, from, focus, original] of shapes) {
        shape.setAttribute("style", original);
        if (from > shown) {
          shape.style.setProperty("fill", "transparent", "important");
          shape.style.setProperty("stroke", IDLE, "important");
        }
        if (!all && focus === step) {
          shape.style.setProperty("stroke", FOCUS, "important");
          shape.style.setProperty("stroke-width", "2px", "important");
        }
      }
      for (const [el, k] of focused) el.classList.toggle("focus-active", !all && k === step);
      svg.classList.toggle("focus-dim", !all && focusSteps.has(step));
      items.forEach((item, i) => item.classList.toggle("diagram-steps__active", !all && i + 1 === step));
      status.textContent = all ? "All steps" : `Step ${step} of ${items.length}`;
      prev.disabled = !all && step === 1;
      next.disabled = !all && step === items.length;
      allButton.setAttribute("aria-pressed", String(all));
    };

    // the previous and next buttons step from the step that was shown before "All steps"
    controls.addEventListener("click", (event) => {
      const go = event.target.closest("button")?.dataset.go;
      if (go === "all") show(step, !all);
      else if (go) show(all ? step : go === "next" ? step + 1 : step - 1);
    });

    items.forEach((item, i) => item.addEventListener("click", () => show(i + 1)));

    show(1);
  };

  document.addEventListener("mermaid-rendered", (event) => {
    const container = event.target.closest(".diagram-steps");
    if (container) setup(container, event.detail.root, event.detail.source ?? "");
  });
})();
