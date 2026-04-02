// Simulation
const sim = d3.forceSimulation(nodes)
  .force("link", d3.forceLink(links).id(d => d.id).distance(d => {
    const cats = new Set([d.source.cat || "finding", d.target.cat || "finding"]);
    if (cats.has("legal") || cats.has("finding")) return 100;
    if (cats.has("outcome")) return 120;
    return 140;
  }).strength(0.2))
  .force("charge", d3.forceManyBody().strength(-600).distanceMax(600))
  .force("topPlacement", d3.forceY(H * 0.1).strength(d => criticalEvents.has(d.id) ? 0.4 : 0))
  .force("center", d3.forceCenter(W / 2, H / 2).strength(0.15))
  .force("collision", d3.forceCollide(d => {
    const linkCount = links.filter(l => {
      const srcId = typeof l.source === "string" ? l.source : l.source.id;
      const tgtId = typeof l.target === "string" ? l.target : l.target.id;
      return srcId === d.id || tgtId === d.id;
    }).length;
    const nodeRadius = Math.max(8, Math.min(75, 8 + Math.pow(linkCount, 1.5) * 0.6));
    return nodeRadius + 80;
  }))
  .force("hubRepulsion", (alpha) => {
    // Strong repulsion between hub nodes to keep them as distinct galaxies
    // Strength scales with node size - larger nodes push harder
    for (let i = 0; i < hubNodes.length; i++) {
      for (let j = i + 1; j < hubNodes.length; j++) {
        const n1 = nodes.find(n => n.id === hubNodes[i]);
        const n2 = nodes.find(n => n.id === hubNodes[j]);
        if (!n1 || !n2) continue;

        const dx = n2.x - n1.x;
        const dy = n2.y - n1.y;
        const dist = Math.sqrt(dx * dx + dy * dy) || 1;

        // Scale repulsion based on average node size
        const sizeWeight = (nodeDegrees[n1.id] + nodeDegrees[n2.id]) / 2;
        const minDist = 500 + sizeWeight * 18; // Much larger minimum distance for big nodes
        const forceStrength = 0.6 + sizeWeight * 0.06; // Much stronger repulsion force

        if (dist < minDist) {
          // Multiply the applied force by the simulation alpha to keep forces balanced when dragging!
          const force = ((minDist - dist) / dist) * forceStrength * alpha;

          n1.vx -= force * dx;
          n1.vy -= force * dy;
          n2.vx += force * dx;
          n2.vy += force * dy;
        }
      }
    }
  });

// Warm-up phase: run simulation in background to pre-position nodes
for (let i = 0; i < 150; i++) sim.tick();
sim.stop();
