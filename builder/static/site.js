// コードブロックに行番号とコピーボタンを付け、URL の #xxx が指す開閉ブロックを開く。
const ICONS = {
  copy: '<rect x="9" y="9" width="12" height="12" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>',
  done: '<path d="M4 12l5 5L20 6"/>',
  failed: '<path d="M6 6l12 12M18 6L6 18"/>',
};

function setIcon(btn, name, label) {
  btn.innerHTML =
    '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">' +
    ICONS[name] +
    "</svg>";
  btn.title = label;
  btn.setAttribute("aria-label", label);
  btn.classList.toggle("done", name === "done");
}

document.querySelectorAll(".md pre > code").forEach((code) => {
  const pre = code.parentElement;
  const text = code.textContent.replace(/\n$/, "");

  code.textContent = "";
  for (const line of text.split("\n")) {
    const span = document.createElement("span");
    span.className = "line";
    span.textContent = line;
    code.appendChild(span);
  }

  const wrap = document.createElement("div");
  wrap.className = "code-block";
  pre.replaceWith(wrap);
  wrap.appendChild(pre);

  const btn = document.createElement("button");
  btn.type = "button";
  btn.className = "copy-btn";
  setIcon(btn, "copy", "コピー");
  btn.addEventListener("click", async () => {
    try {
      await copyText(text);
      setIcon(btn, "done", "コピーしました");
    } catch {
      setIcon(btn, "failed", "コピーできませんでした");
    }
    setTimeout(() => setIcon(btn, "copy", "コピー"), 1500);
  });
  wrap.appendChild(btn);
});

async function copyText(text) {
  if (navigator.clipboard && window.isSecureContext) {
    return navigator.clipboard.writeText(text);
  }
  // file:// などで Clipboard API が使えないときの代替
  const ta = document.createElement("textarea");
  ta.value = text;
  ta.style.position = "fixed";
  ta.style.opacity = "0";
  document.body.appendChild(ta);
  ta.select();
  const ok = document.execCommand("copy");
  ta.remove();
  if (!ok) throw new Error("copy failed");
}

function openTarget() {
  const id = decodeURIComponent(location.hash.slice(1));
  const el = id && document.getElementById(id);
  if (el && el.tagName === "DETAILS") el.open = true;
}
openTarget();
window.addEventListener("hashchange", openTarget);

document.querySelectorAll("[data-toggle-all]").forEach((btn) => {
  btn.addEventListener("click", () => {
    const open = btn.dataset.toggleAll === "open";
    document.querySelectorAll("main details").forEach((d) => (d.open = open));
  });
});
