// Create selections immediately (before warm-up, so they're ready to render)
const linkSel = g.append("g").selectAll("line").data(links).join("line")
  .attr("stroke-width", 0.8).attr("stroke-opacity", 0.5)
  .attr("stroke", d => {
    const src = nodes.find(n => n.id === (typeof d.source === "string" ? d.source : d.source.id));
    return src ? CATEGORIES[src.cat]?.color || "#888" : "#888";
  })
  .attr("marker-end", d => {
    const src = nodes.find(n => n.id === (typeof d.source === "string" ? d.source : d.source.id));
    return src ? `url(#arrow-${src.cat})` : "url(#arrow-finding)";
  });

const nodeSel = g.append("g").selectAll("g").data(nodes).join("g")
  .attr("class", "node-g")
  .style("cursor", "pointer")
  .style("will-change", "transform") // Offload translation coordinate calculations to GPU
  .call(d3.drag()
    .on("start", (e, d) => { if (!e.active) sim.alphaTarget(0.3).restart(); d.fx = d.x; d.fy = d.y; })
    .on("drag", (e, d) => { d.fx = e.x; d.fy = e.y; })
    .on("end", (e, d) => { if (!e.active) sim.alphaTarget(0); d.fx = null; d.fy = null; }))
  .on("click", (e, d) => showPanel(d))
  .on("mouseenter", (e, d) => showTooltip(e, d))
  .on("mouseleave", () => hideTooltip());

nodeSel.append("circle")
  .attr("r", d => {
    const linkCount = links.filter(l => l.source === d.id || l.target === d.id || (l.source.id === d.id) || (l.target.id === d.id)).length;
    let radius = Math.max(14, Math.min(80, 8 + Math.pow(linkCount, 1.35) * 0.65));
    // Critical events get +18px size boost
    if (criticalEvents.has(d.id)) {
      radius = Math.min(105, radius + 18);
    }
    // Key nodes get +25px size boost
    else if (keyNodes.has(d.id)) {
      radius = Math.min(110, radius + 25);
    }
    return radius;
  })
  .attr("fill", d => {
    if (criticalEvents.has(d.id)) return "#E91E6388"; // Hot pink with transparency
    return CATEGORIES[d.cat]?.color + "cc" || "#88878088";
  })
  .attr("stroke", d => {
    if (criticalEvents.has(d.id)) return "#E91E63"; // Hot pink
    return CATEGORIES[d.cat]?.color || "#888";
  })
  .attr("stroke-width", d => criticalEvents.has(d.id) ? 3 : 1.5);

nodeSel.append("text")
  .attr("text-anchor", "middle").attr("dy", "0.35em")
  .attr("font-size", d => {
    const linkCount = links.filter(l => l.source === d.id || l.target === d.id || (l.source.id === d.id) || (l.target.id === d.id)).length;
    let fontSize = Math.max(12, Math.min(20, 8 + Math.pow(linkCount, 1.35) * 0.35));
    // Key nodes get larger text
    if (keyNodes.has(d.id)) {
      fontSize = Math.min(26, fontSize + 5);
    }
    return fontSize;
  })
  .attr("font-family", "'DM Sans', system-ui, sans-serif")
  .attr("fill", d => CATEGORIES[d.cat]?.color || "#888")
  .attr("font-weight", d => criticalEvents.has(d.id) ? "900" : "500")
  .attr("pointer-events", "none")
  .attr("y", d => {
    const linkCount = links.filter(l => l.source === d.id || l.target === d.id || (l.source.id === d.id) || (l.target.id === d.id)).length;
    let radius = Math.max(14, Math.min(80, 8 + Math.pow(linkCount, 1.35) * 0.65));
    let gap = 18;
    if (criticalEvents.has(d.id)) {
      radius = Math.min(105, radius + 18);
      gap = 24;
    } else if (keyNodes.has(d.id)) {
      radius = Math.min(110, radius + 25);
      gap = 28;
    }
    return radius + gap;
  })
  .each(function(d) {
    const el = d3.select(this);
    const words = d.label.split(/\s+/);
    let line = [];
    const maxLineLen = 14; 
    let tspan = el.append("tspan").attr("x", 0).attr("dy", "0.35em");
    
    for (let i = 0; i < words.length; i++) {
        const word = words[i];
        line.push(word);
        const text = line.join(" ");
        if (text.length > maxLineLen && line.length > 1) {
            line.pop();
            tspan.text(line.join(" "));
            line = [word];
            tspan = el.append("tspan").attr("x", 0).attr("dy", "1.1em").text(word);
        } else {
            tspan.text(text);
        }
    }
  });

// Ensure critical events render visually on top of all other nodes (Z-index)
nodeSel.filter(d => criticalEvents.has(d.id)).raise();

// Render the warm-up state before attaching the tick listener
linkSel
  .attr("x1", d => d.source.x).attr("y1", d => d.source.y)
  .attr("x2", d => d.target.x).attr("y2", d => d.target.y);
nodeSel.attr("transform", d => `translate(${d.x},${d.y})`);

sim.on("tick", () => {
  linkSel
    .attr("x1", d => d.source.x).attr("y1", d => d.source.y)
    .attr("x2", d => d.target.x).attr("y2", d => d.target.y);
  nodeSel.attr("transform", d => `translate(${d.x},${d.y})`);
});

// Restart simulation after tick listener is attached
sim.restart();

// Hide loading overlay after warm-up
setTimeout(() => {
  const loadingEl = document.getElementById('loading');
  if (loadingEl) loadingEl.classList.add('hidden');
}, 150);

// Center view on midpoint between Bhadra 23 and Bhadra 24, with bias toward graph center of mass
setTimeout(() => {
  const bhadra23 = nodes.find(n => n.id === "bhadra23");
  const bhadra24 = nodes.find(n => n.id === "bhadra24");

  if (bhadra23 && bhadra24) {
    // Calculate center of mass of all nodes
    const centerOfMassX = nodes.reduce((sum, n) => sum + n.x, 0) / nodes.length;
    const centerOfMassY = nodes.reduce((sum, n) => sum + n.y, 0) / nodes.length;

    // Midpoint between Bhadra nodes
    const bhadraX = (bhadra23.x + bhadra24.x) / 2;
    const bhadraY = (bhadra23.y + bhadra24.y) / 2;

    // Blend with 60% weight toward center of mass, 40% toward Bhadra midpoint
    const finalCenterX = centerOfMassX * 0.6 + bhadraX * 0.4;
    const finalCenterY = centerOfMassY * 0.6 + bhadraY * 0.4;

    // Apply zoom transform: zoom to responsive scale and center on the blended point
    // On mobile devices (width < 768px), we use a much wider zoom (e.g. 0.35) to fit both nodes
    const k = W < 768 ? 0.35 : 0.8; 
    const x = W / 2 - finalCenterX * k;
    const y = H / 2 - finalCenterY * k;
    const transform = d3.zoomIdentity.translate(x, y).scale(k);
    svg.call(zoom.transform, transform); // Immediate, no transition
  }
}, 50);
