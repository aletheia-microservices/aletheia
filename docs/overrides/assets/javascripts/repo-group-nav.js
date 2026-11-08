(() => {
  let observer;

  document$.subscribe(() => {
    observer?.disconnect();

    const nav = document.querySelector(".repo-group-nav");
    const header = document.querySelector(".md-header");
    if (!nav || !header) return;

    const article = nav.closest(".md-typeset");
    const updateOffsets = () => {
      article.style.setProperty("--repo-group-header-height", `${header.offsetHeight}px`);
      article.style.setProperty("--repo-group-nav-height", `${nav.offsetHeight}px`);
    };

    updateOffsets();
    observer = new ResizeObserver(updateOffsets);
    observer.observe(header);
    observer.observe(nav);
  });
})();
