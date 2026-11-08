// Colors the comments of code blocks marked with the "annotated" class by what they point out
// (see page.css): database taints, service taints, or anything else, such as where a value is stored.
// The highlighted line (hl_lines) that holds a comment gets the same color, and any identifier on that
// line named by the block's data-service, data-database or data-other attribute (a space-separated
// list, e.g., data-other="postID message") is colored too. A highlighted line without a comment takes the color of the next highlighted line that
// has one, so a statement that spans several lines (e.g., a struct literal) can be highlighted from its start.
// document$ is Material's observable for each page load, including instant navigation.
const KINDS = { service: "annotation-service", database: "annotation-database", other: "annotation-other" };

document$.subscribe(() => {
  for (const comment of document.querySelectorAll(".annotated .c1")) {
    const text = comment.textContent;
    const kind = text.includes("database taint") ? KINDS.database
      : text.includes("service taint") ? KINDS.service
      : KINDS.other;
    comment.classList.add(kind);
    comment.closest(".hll")?.classList.add(kind);
  }

  for (const block of document.querySelectorAll(".annotated")) {
    let next = null;
    for (const line of [...block.querySelectorAll(".hll")].reverse()) {
      const own = Object.values(KINDS).find((kind) => line.classList.contains(kind));
      if (own) next = own;
      else if (next) line.classList.add(next);
    }

    for (const [attr, kind] of Object.entries(KINDS)) {
      const names = block.dataset[attr]?.split(" ");
      if (!names) continue;
      for (const line of block.querySelectorAll(`.hll.${kind}`)) {
        for (const token of line.querySelectorAll("span:not(.c1)")) {
          if (names.includes(token.textContent)) token.classList.add(kind);
        }
      }
    }
  }
});
