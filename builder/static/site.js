// コードブロックに行番号とコピーボタンを付け、URL の #xxx が指す開閉ブロックを開く。
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
  btn.textContent = "[コピー]";
  btn.addEventListener("click", async () => {
    try {
      await copyText(text);
      btn.textContent = "[コピーしました]";
      btn.classList.add("done");
    } catch {
      btn.textContent = "[コピーできませんでした]";
    }
    setTimeout(() => {
      btn.textContent = "[コピー]";
      btn.classList.remove("done");
    }, 1500);
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
