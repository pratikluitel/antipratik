// Layout boundaries and states
const W = window.innerWidth, H = window.innerHeight - 60;
let activeCategories = new Set(Object.keys(CATEGORIES));
let searchTerm = "";

// Define key nodes that should be larger
const keyNodes = new Set([
  "pm_karki",
  "kp_oli",
  "sudan_gurung",
  "lekhak",
  "hami_nepal",
  "discord_servers",
  "discord_vote",
  "command_vacuum",
  "aimed_fire",
  "total_deaths",
  "karki_pm",
  "pm_resign"
]);

// Define critical event nodes that should have special styling
const criticalEvents = new Set([
  "bhadra23",
  "bhadra24"
]);

// Calculate node degrees for hub identification
const nodeDegrees = {};
nodes.forEach(n => {
  nodeDegrees[n.id] = links.filter(l => {
    const srcId = typeof l.source === "string" ? l.source : l.source.id;
    const tgtId = typeof l.target === "string" ? l.target : l.target.id;
    return srcId === n.id || tgtId === n.id;
  }).length;
});

// Identify hub nodes (top 10 by degree) - these become galaxy centers
const hubNodes = nodes
  .map(n => ({ ...n, degree: nodeDegrees[n.id] }))
  .sort((a, b) => b.degree - a.degree)
  .slice(0, 10)
  .map(n => n.id);

// Initialize hub nodes at distant positions to prevent violent repulsion jitter at load
const peripheryPositions = [
  { x: W * 0.85, y: H / 2 },      // right
  { x: W / 2, y: H * 0.15 },      // top
  { x: W * 0.15, y: H / 2 },      // left
  { x: W / 2, y: H * 0.85 },      // bottom
  { x: W * 0.75, y: H * 0.25 },   // top-right
  { x: W * 0.25, y: H * 0.25 },   // top-left
  { x: W * 0.75, y: H * 0.75 },   // bottom-right
  { x: W * 0.25, y: H * 0.75 },   // bottom-left
  { x: W * 0.85, y: H * 0.85 },   // far bottom-right
  { x: W * 0.15, y: H * 0.15 },   // far top-left
];

hubNodes.forEach((hubId, idx) => {
  const node = nodes.find(n => n.id === hubId);
  if (node && idx < peripheryPositions.length) {
    node.x = peripheryPositions[idx].x;
    node.y = peripheryPositions[idx].y;
  }
});
