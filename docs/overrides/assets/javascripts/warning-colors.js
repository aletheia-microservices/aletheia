// Colors the parts of the warnings and constraints that Aletheia writes, in code blocks marked with the
// "warnings" class (``` { .txt .warnings title="..." }): labels such as "delete:" or "READ (FOREIGN KEY):"
// in semi-bold, service methods in blue and database operations and fields in green, as in the diagrams,
// constraint keywords such as FOREIGN_KEY or [MANDATORY] in orange, and the "..." of call paths faded.
// Only the display changes: the copy button still copies the plain text.
// document$ is Material's observable for each page load, including instant navigation.
// The script runs in its own scope, so its names do not clash with those of the other scripts.
(() => {
  const escape = (text) => text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");

  // a label at the start of a line, after its indentation and an optional "- ", e.g. "missing cascade #1:"
  const LABEL = /^(\s*)((?:- )?[A-Za-z][A-Za-z0-9 ()_]*?(?: #\d+)?:)(?=\s|$)/;

  // the other parts, in order of priority
  const PARTS = [
    ["w-count", /\[NUM_WARNINGS = \d+\]/],
    ["w-faded", /\.\.\./],
    ["w-constraint", /\bFOREIGN_KEY\b|\bREFERENCES\b|\bPRIMARY KEY\b|\bUNIQUE\b|\[MANDATORY\]|\[T\]/],
    ["w-service", /\b[A-Z][A-Za-z0-9_]*\.[A-Za-z_][A-Za-z0-9_]*\(\)/],
    ["w-db", /\b[a-z][a-z0-9_]*(?:\.[A-Za-z0-9_*[\]]+)*\.[A-Z][A-Za-z0-9_]*\(\)/],
    ["w-db", /\b[a-z][a-z0-9_]*\.[A-Za-z0-9_*[\]]+(?:\.[A-Za-z0-9_*[\]]+)*/],
    ["w-db", /(?<=\{)[^{}]+(?=\})/],
  ];
  const ANY = new RegExp(PARTS.map(([, re]) => `(${re.source})`).join("|"), "g");

  const colorLine = (line) => {
    let html = "";
    const label = line.match(LABEL);
    if (label) {
      html += escape(label[1]) + `<span class="w-label">${escape(label[2])}</span>`;
      line = line.slice(label[0].length);
    }
    let last = 0;
    for (const match of line.matchAll(ANY)) {
      const kind = PARTS[match.slice(1).findIndex((group) => group !== undefined)][0];
      html += escape(line.slice(last, match.index)) + `<span class="${kind}">${escape(match[0])}</span>`;
      last = match.index + match[0].length;
    }
    return html + escape(line.slice(last));
  };

  document$.subscribe(() => {
    for (const code of document.querySelectorAll(".warnings pre > code")) {
      if (code.dataset.colored) continue;
      code.innerHTML = code.textContent.split("\n").map(colorLine).join("\n");
      code.dataset.colored = "true";
    }
  });
})();
