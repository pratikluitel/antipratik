const svg = d3.select("#graph");
const g = svg.append("g");

const zoom = d3.zoom().scaleExtent([0.15, 4]).on("zoom", e => g.attr("transform", e.transform));
svg.call(zoom);

// Arrow markers
const defs = svg.append("defs");
Object.entries(CATEGORIES).forEach(([key, cat]) => {
  defs.append("marker")
    .attr("id", `arrow-${key}`)
    .attr("viewBox", "0 0 10 10").attr("refX", 18).attr("refY", 5)
    .attr("markerWidth", 6).attr("markerHeight", 6).attr("orient", "auto-start-reverse")
    .append("path").attr("d", "M2 1L8 5L2 9").attr("fill", "none")
    .attr("stroke", cat.color).attr("stroke-width", 1.5)
    .attr("stroke-linecap", "round").attr("stroke-linejoin", "round");
});
