// Gives each link in the page content a class by where it leads, which page.css draws as a small arrow after
// some of them: "link-up" or "link-down" (up or down arrow) for a section of the same page above or below
// the link, and "link-external" (diagonal arrow) for anything outside the site, such as GitHub, other
// websites or the paper. A link to another page of this site, or to a section of another page, gets
// "link-docs" and no arrow. Footnote links, links in code (e.g., line numbers), the "edit this page" button
// and image links (e.g., badges on the home page) get no class.
// Each of these links also gets a tooltip that names its target, unless it already has one: the page's
// place in the navigation (e.g., "Getting Started › Tutorial"), the heading of the section it jumps to,
// the GitHub repository and path, the paper, or the website (e.g., "graphviz.org").
// Comparing hosts also works with instant navigation, where Material turns internal links into absolute URLs.
// document$ is Material's observable for each page load, including instant navigation.
// The script runs in its own scope, so its names do not clash with those of the other scripts.
(() => {
  const KINDS = ["up", "down", "docs", "external"];
  const PAPER = "usenix.org/system/files/osdi26-ferreira.pdf";

  // the same page has one key, whether its URL ends in "/", "/index.html" or "#section"
  const pageKey = (url) => {
    let path = new URL(url, location.href).pathname.replace(/index\.html$/, "");
    return path.endsWith("/") ? path : `${path}/`;
  };

  // the pages of the navigation in the left sidebar, each with its name in the navigation, e.g.,
  // "Getting Started › Tutorial" (pages outside a section, such as Home, have only their title)
  const navPages = () => {
    const pages = new Map();
    for (const a of document.querySelectorAll(".md-nav--primary a.md-nav__link[href]")) {
      const key = pageKey(a.href);
      if (pages.has(key)) continue;
      const section = a.closest(".md-nav__item")?.parentElement?.closest(".md-nav__item--nested");
      const sectionName = section?.querySelector(":scope > .md-nav__link")?.textContent.trim();
      const title = a.textContent.trim();
      pages.set(key, sectionName ? `${sectionName} › ${title}` : title);
    }
    return pages;
  };

  const linkKind = (link, here) => {
    if (link.classList.contains("md-content__button") || link.classList.contains("headerlink")) return null;
    if (link.closest("pre, code, .highlight, .highlighttable")) return null; // e.g., the line numbers of code
    if (link.querySelector("img") || link.hostname === "") return null;
    if (link.host !== location.host) return "external";
    const target = pageKey(link.href);
    if (target === here) {
      // a section of this page: up or down, by where it is; footnotes and their back links get no arrow
      if (link.closest(".footnote-ref, .footnote-backref, .footnote")) return null;
      const section = sectionOf(link);
      if (!section) return null;
      return link.compareDocumentPosition(section) & Node.DOCUMENT_POSITION_FOLLOWING ? "down" : "up";
    }
    return "docs";
  };

  // the element that a link to this page jumps to, e.g., the heading of a section, or null
  const sectionOf = (link) => link.hash ? document.getElementById(decodeURIComponent(link.hash.slice(1))) : null;

  // the text of a heading, without the ¶ of its permanent link
  const headingText = (el) => [...el.childNodes].filter((n) => !n.classList?.contains("headerlink")).map((n) => n.textContent).join("").trim();

  const externalName = (link) => {
    if (link.href.includes(PAPER)) return "Aletheia paper (OSDI '26, PDF)";
    const host = link.hostname.replace(/^www\./, "");
    if (host !== "github.com") return host;
    // github.com/{owner}/{repo}, followed by the folder or file, without "tree/{branch}" or "blob/{branch}"
    const [owner, repo, kind, , ...rest] = link.pathname.split("/").filter(Boolean);
    if (!owner || !repo) return host;
    const path = ["tree", "blob"].includes(kind) && rest.length ? `/${rest.join("/")}` : "";
    return `${host}/${owner}/${repo}${path}`;
  };

  document$.subscribe(() => {
    const pages = navPages(), here = pageKey(location.href);
    for (const link of document.querySelectorAll(".md-typeset a[href]")) {
      const kind = linkKind(link, here);
      for (const k of KINDS) link.classList.toggle(`link-${k}`, k === kind);
      if (!kind || (link.title && !link.dataset.autoTitle)) continue;
      if (kind === "up" || kind === "down") {
        const section = sectionOf(link);
        if (!/^H[1-6]$/.test(section.tagName)) continue;
        link.title = headingText(section);
      } else {
        link.title = kind === "external" ? externalName(link) : pages.get(pageKey(link.href)) ?? link.pathname;
      }
      link.dataset.autoTitle = "true";
    }
  });
})();
