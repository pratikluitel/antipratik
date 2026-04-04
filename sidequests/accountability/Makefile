# Regenerate styles.min.css from styles.css
# Usage: make css
css:
	node -e "\
	const fs=require('fs');\
	let c=fs.readFileSync('styles.css','utf8');\
	c=c.replace(/\/\*[\s\S]*?\*\//g,'').replace(/\s*\n\s*/g,'\n').replace(/\n/g,' ').replace(/\s{2,}/g,' ').replace(/\s*([{};:,>~+])\s*/g,'$$1').replace(/;}/g,'}').trim();\
	fs.writeFileSync('styles.min.css',c);\
	console.log('styles.min.css written,',c.length,'bytes');\
	"
