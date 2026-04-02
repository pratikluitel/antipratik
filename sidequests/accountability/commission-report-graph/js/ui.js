// Filters
const filtersEl = document.getElementById("filters");
const legendEl = document.getElementById("legend");
Object.entries(CATEGORIES).forEach(([key, cat]) => {
  const btn = document.createElement("button");
  btn.className = "kg-filter-btn"; btn.textContent = cat.label;
  btn.style.color = cat.color; btn.style.borderColor = cat.color + "88";
  btn.onclick = () => toggleCat(key, btn);
  filtersEl.appendChild(btn);

  const leg = document.createElement("div");
  leg.className = "kg-leg";
  leg.innerHTML = `<div class="kg-leg-dot" style="background:${cat.color}"></div>${cat.label} (${nodes.filter(n => n.cat === key).length})`;
  legendEl.appendChild(leg);
});
document.getElementById("count").textContent = `${nodes.length} nodes · ${links.length} connections`;

function toggleCat(key, btn) {
  if (activeCategories.has(key)) { activeCategories.delete(key); btn.classList.add("off"); }
  else { activeCategories.add(key); btn.classList.remove("off"); }
  updateVisibility();
}

function filterSearch(val) {
  searchTerm = val.toLowerCase();
  updateVisibility();
}

function updateVisibility() {
  const visibleIds = new Set(nodes.filter(n => {
    const catOk = activeCategories.has(n.cat);
    const searchOk = !searchTerm || n.label.toLowerCase().includes(searchTerm) || n.desc.toLowerCase().includes(searchTerm);
    return catOk && searchOk;
  }).map(n => n.id));

  nodeSel.attr("opacity", d => visibleIds.has(d.id) ? 1 : 0.07)
    .attr("pointer-events", d => visibleIds.has(d.id) ? "all" : "none");
  linkSel.attr("opacity", d => {
    const sid = typeof d.source === "string" ? d.source : d.source.id;
    const tid = typeof d.target === "string" ? d.target : d.target.id;
    return visibleIds.has(sid) && visibleIds.has(tid) ? 0.5 : 0.04;
  });
}

// Tooltip
const tt = document.getElementById("tooltip");
function showTooltip(e, d) {
  tt.style.display = "block";
  tt.innerHTML = `<div class="kg-tt-type" style="color:${CATEGORIES[d.cat]?.color}">${CATEGORIES[d.cat]?.label}</div><div class="kg-tt-title">${d.label}</div><div class="kg-tt-body">${d.desc.slice(0, 120)}${d.desc.length > 120 ? "…" : ""}</div><div class="kg-tt-links">Click for full details · §${d.section?.split("§")[1] || ""}</div>`;
  positionTooltip(e);
}
function hideTooltip() { tt.style.display = "none"; }
function positionTooltip(e) {
  const x = e.clientX + 14, y = e.clientY - 30;
  tt.style.left = Math.min(x, window.innerWidth - 300) + "px";
  tt.style.top = Math.min(y, window.innerHeight - 150) + "px";
}

// Info panel
function showPanel(d) {
  const panel = document.getElementById("info-panel");
  const content = document.getElementById("panel-content");
  const connectedLinks = links.filter(l => {
    const sid = typeof l.source === "string" ? l.source : l.source.id;
    const tid = typeof l.target === "string" ? l.target : l.target.id;
    return sid === d.id || tid === d.id;
  });
  const connItems = connectedLinks.slice(0, 8).map(l => {
    const sid = typeof l.source === "string" ? l.source : l.source.id;
    const tid = typeof l.target === "string" ? l.target : l.target.id;
    const other = sid === d.id ? tid : sid;
    const otherNode = nodes.find(n => n.id === other);
    const isReverse = sid !== d.id;
    const displayLabel = isReverse && l.labelReverse ? l.labelReverse : l.label;
    const dir = sid === d.id ? "→" : "←";
    return `<div class="conn-item"><span class="conn-arrow">${dir}</span><span>${displayLabel} <strong>${otherNode?.label || other}</strong></span></div>`;
  }).join("");

  content.innerHTML = `
    <div class="type" style="color:${CATEGORIES[d.cat]?.color}">${CATEGORIES[d.cat]?.label}</div>
    <h3>${d.label}</h3>
    <p>${d.desc}</p>
    ${d.section ? `<p style="font-size:10px;color:var(--kg-text3)">${d.section}</p>` : ""}
    ${connectedLinks.length ? `<div class="connections"><h4>${connectedLinks.length} connections</h4>${connItems}${connectedLinks.length > 8 ? `<div style="font-size:10px;color:var(--kg-text3)">+${connectedLinks.length - 8} more</div>` : ""}</div>` : ""}
  `;
  panel.style.display = "block";

  // Highlight connected nodes
  const connIds = new Set([d.id, ...connectedLinks.map(l => {
    const sid = typeof l.source === "string" ? l.source : l.source.id;
    const tid = typeof l.target === "string" ? l.target : l.target.id;
    return sid === d.id ? tid : sid;
  })]);
  nodeSel.attr("opacity", n => connIds.has(n.id) ? 1 : 0.1);
  linkSel.attr("opacity", l => {
    const sid = typeof l.source === "string" ? l.source : l.source.id;
    const tid = typeof l.target === "string" ? l.target : l.target.id;
    return connIds.has(sid) && connIds.has(tid) ? 0.8 : 0.04;
  });
}

function closePanel() {
  document.getElementById("info-panel").style.display = "none";
  updateVisibility();
}

function toggleHelp() {
  const panel = document.getElementById("help-panel");
  const btn = document.getElementById("helpBtn");
  const isOpen = panel.classList.toggle("open");
  btn.classList.toggle("active", isOpen);
}

// Close help panel when clicking the SVG background
svg.on("click.help", () => {
  const panel = document.getElementById("help-panel");
  if (panel.classList.contains("open")) toggleHelp();
});

// Fit Graph to screen
function fitToScreen() {
  const visibleIds = new Set(nodes.filter(n => {
    const catOk = activeCategories.has(n.cat);
    const searchOk = !searchTerm || n.label.toLowerCase().includes(searchTerm) || n.desc.toLowerCase().includes(searchTerm);
    return catOk && searchOk;
  }).map(n => n.id));

  const visibleNodes = nodes.filter(n => visibleIds.has(n.id));

  if (visibleNodes.length === 0) return;

  const minX = d3.min(visibleNodes, d => d.x);
  const maxX = d3.max(visibleNodes, d => d.x);
  const minY = d3.min(visibleNodes, d => d.y);
  const maxY = d3.max(visibleNodes, d => d.y);

  if (minX === undefined || maxX === undefined) return;

  const dx = maxX - minX || 1;
  const dy = maxY - minY || 1;
  const cx = (minX + maxX) / 2;
  const cy = (minY + maxY) / 2;

  // Add padding
  const scale = Math.max(0.15, Math.min(2, 0.85 / Math.max(dx / W, dy / H)));

  const transform = d3.zoomIdentity
    .translate(W / 2 - cx * scale, H / 2 - cy * scale)
    .scale(scale);

  svg.transition().duration(750).call(zoom.transform, transform);
}
