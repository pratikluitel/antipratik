const CATEGORIES = {
  person: { label: "Person", color: "#E24B4A", light: "#FCEBEB" },
  org: { label: "Organisation", color: "#378ADD", light: "#E6F1FB" },
  event: { label: "Event", color: "#E4A135", light: "#FAEEDA" },
  place: { label: "Place", color: "#1D9E75", light: "#EAF3DE" },
  legal: { label: "Legal", color: "#7F77DD", light: "#EEEDFE" },
  finding: { label: "Finding", color: "#888780", light: "#F2F1EF" },
  outcome: { label: "Outcome", color: "#D4537E", light: "#FBEAF0" },
  digital: { label: "Digital", color: "#0F6E56", light: "#E1F5EE" },
};

const nodes = [
  // EVENTS
  { id: "bhadra23", label: "Bhadra 23 March", cat: "event", desc: "Sept 8, 2025. Peaceful GenZ march from Maitighar → Parliament. ~50,000+ participants. Police fired live rounds. 17 killed at Parliament perimeter on the day; 12 more died from wounds. Total Bhadra 23 bullet deaths: 29 in Kathmandu valley, 42 nationwide.", section: "§5.1.2, §12.2" },
  { id: "bhadra24", label: "Bhadra 24 Arson", cat: "event", desc: "Sept 9, 2025. Nationwide arson, looting, vandalism. Singha Durbar, Supreme Court, PM/President residences, 196 party offices burned. TOB motorcycle group and political infiltrators identified. 47 more deaths. NPR 84.45 Cr damage (NPC estimate).", section: "§5.1.3, §5.1.4" },
  { id: "ban26", label: "26-Platform Social Media Ban", cat: "event", desc: "Bhadra 19 (Sept 5). Government banned 26 social media platforms. Immediate trigger for GenZ uprising. Lifted same night as Bhadra 23 deaths — too late.", section: "§5.1.1" },
  { id: "curfew", label: "Curfew Order", cat: "event", desc: "CDO Kathmandu issued curfew at 12:30 PM on Bhadra 23. Crowd noise meant it was inaudible to most protesters. Extended multiple times through Bhadra 24.", section: "§12.11" },
  { id: "commission", label: "Inquiry Commission", cat: "event", desc: "Janachbujh Ayog (Inquiry Commission) formed Bhadra 2082/06/05. Began work 06/09. Extended 3 times. Produced 907-page Pratibeden (report). Commission itself: charged the PM who formed it.", section: "§1–3" },
  { id: "pm_resign", label: "KP Oli Resignation", cat: "event", desc: "PM KP Sharma Oli resigned on Bhadra 24, after deaths and nationwide arson. His PA received calls from police during crisis; he held NSC meeting only at 10 PM — after hours of deaths.", section: "§12.8" },
  { id: "discord_vote", label: "Discord PM Vote", cat: "event", desc: "Inside the Youth Against Corruption Discord server, GenZ protesters held a vote on who should be the next PM. Winner: former Chief Justice Sushila Karki (~62%). She was subsequently appointed.", section: "§12.13" },
  { id: "airport_attack", label: "Airport Attack (Thwarted)", cat: "event", desc: "Tribhuvan International Airport attack planned by Discord users (AppleP, Vishwasthapa, Zeus). Nepal Army successfully blocked it.", section: "§12.13" },
  { id: "singha_durbar_burn", label: "Singha Durbar Burned", cat: "event", desc: "Sept 9, ~2:05 PM. Nepal's government secretariat (equivalent of PM's office complex) partially burned. Diwakar Dulal posted 87 Discord messages trying to save the data centre — ignored. Data centre lost.", section: "§5.1.4, §12.13" },
  { id: "global_college_burn", label: "Global College Burned", cat: "event", desc: "Discord user 'Tony' spread false rape rumour about Global College hostel. College burned on Bhadra 24. Police confirmed no rape had occurred.", section: "§12.13, §14.10" },
  { id: "hilton_burn", label: "Hilton Hotel Burned", cat: "event", desc: "Journalist Dilbhushan Pathak falsely claimed Hilton Hotel was owned by Deuba's son Jayavir Singh Deuba on YouTube. Hotel burned. Commission: 'chilled foreign investment.'", section: "§14.10" },
  { id: "prison_collapse", label: "Prison System Collapse", cat: "event", desc: "14,555 prisoners escaped from prisons and police custody during Bhadra 23–24. Still at large: 5,208 (as of report). Includes 1,174 rapists, 549 murderers, 1,037 drug traffickers, 36 weapons offenders.", section: "§12.12" },
  { id: "banke_prison", label: "Banke Prison Collapse", cat: "event", desc: "1,280 prisoners released by mob. 5 deaths inside prison. 459 still missing from Banke alone. Prison had no contingency plan for simultaneous external mob + internal unrest.", section: "§12.12" },
  { id: "prev_andolan", label: "Previous Durga Prasai Andolan", cat: "event", desc: "Referenced by Commission as precedent for the same security failure pattern: police going directly to live ammunition without lathi charge / non-lethal escalation.", section: "§5 (DIGP testimony)" },

  // PEOPLE — accused
  { id: "kp_oli", label: "PM KP Sharma Oli", cat: "person", desc: "Prime Minister. Charged: Reckless + Negligent killing (Muluki Penal Code §181 + §182). Live fire lasted ~4 hours. No Cabinet decision on force. NSC meeting held only at 10 PM after deaths. Refused to answer Commission questions. Said they were 'questions for SEE students.'", section: "§13.1" },
  { id: "lekhak", label: "Home Min. Ramesh Lekhak", cat: "person", desc: "Home Minister. Charged: §181 + §182 (Reckless + Negligent killing). Communicated through his PA rather than direct command to security forces. No corrective order to stop firing.", section: "§13.1" },
  { id: "khapung", label: "IGP Chandra Kuber Khapung", cat: "person", desc: "Inspector General of Nepal Police. Charged: §181 + §182. Highest police authority. 42 calls on CDR with DIGP Bishwa Adhikari during crisis. Issued no stop-firing order for ~4 hours.", section: "§13.1" },
  { id: "dubadi", label: "Home Secretary Dubadi", cat: "person", desc: "Gokarnaman Dubadi. Charged: §182 (Negligent killing). Senior civil servant at Home Ministry. Mobile number 9851322222.", section: "§13.2" },
  { id: "ayal", label: "APF IGP Raju Aryal", cat: "person", desc: "Armed Police Force Inspector General. Charged: §182 (Negligent killing). Issued Victor Control orders during Bhadra 24 including prisoner release order.", section: "§13.2" },
  { id: "hut_raj", label: "NID Chief Hut Raj Thapa", cat: "person", desc: "National Investigation Department Chief. Charged: §182. Failed to provide timely intelligence. CDR shows only 8 calls during 2-day crisis.", section: "§13.2" },
  { id: "rijal", label: "CDO Khaviraj Rijal", cat: "person", desc: "Chief District Officer, Kathmandu. Charged: §182. Issued curfew and gave sequential force orders verbally. Said in testimony that 'even during curfew the PJA can give shoot orders' — claiming he didn't give one.", section: "§13.2" },
  { id: "shah_aigp", label: "AIGP Siddha Vikram Shah", cat: "person", desc: "Nepal Police Operations Chief. Departmental action recommended (Police Act §9(4)). CDR: 69 calls with senior officers during 2-day crisis.", section: "§13.2" },
  { id: "om_rana", label: "DIGP Om Bahadur Rana", cat: "person", desc: "Acting head, Kathmandu Valley Police during Bhadra 23. Departmental action recommended. Victor Control orders issued. CDR: 53 calls.", section: "§13.2" },
  { id: "bishwa", label: "DIGP Bishwa Adhikari", cat: "person", desc: "Kathmandu District Police head. Departmental action recommended. Most active CDR node: 95 total calls. Brig. Gen. Manoj Baidwar was Army commander at Shital Niwas which burned.", section: "§13.2" },
  { id: "deep_samsher", label: "DSP Deep Samsher JBR", cat: "person", desc: "Ranipokhari Valley Police. Departmental action recommended (Police Act §9(4)).", section: "§13.2" },
  { id: "kandel", label: "DSP Rishiram Kandel", cat: "person", desc: "Nepal Police Special Task Force at Parliament. Departmental action recommended.", section: "§13.2" },
  { id: "kaki_igp", label: "IGP Dan Bahadur Kaki", cat: "person", desc: "Current IGP (was AIGP on Bhadra 23–24). Departmental action recommended under Police Act §9(4). Victor Control orders + weapons not secured on Bhadra 24.", section: "§13.2" },
  { id: "baidwar", label: "Brig. Gen. Manoj Baidwar", cat: "person", desc: "Army security commander, Shital Niwas (President's residence) — BURNED. Army Act §105 action recommended.", section: "§13.2" },
  { id: "diwakar_khadka", label: "Col. Diwakar Khadka", cat: "person", desc: "Army commander, Baluwatar (PM residence) — BURNED. Army Act §105.", section: "§13.2" },
  { id: "ganesh_khadka", label: "Col. Ganesh Khadka", cat: "person", desc: "Army commander, Singha Durbar — PARTIALLY BURNED. Army Act §105. Also District Security Committee member.", section: "§13.2" },
  { id: "santosh_dhungel", label: "Lt. Col. Santosh Dhungel", cat: "person", desc: "Army commander, Parliament / BICC — ATTACKED. Army Act §105.", section: "§13.2" },
  { id: "khanal_nid", label: "JD Krishna Prasad Khanal", cat: "person", desc: "NID Valley Operations. Special Service Act Rule 10.2 action for failing to provide timely intelligence.", section: "§13.2" },
  { id: "gachhadar", label: "JD Riben Kumar Gachhadar", cat: "person", desc: "NID Kathmandu District Chief. Same Special Service Act action.", section: "§13.2" },
  { id: "paudel_apf", label: "APF Maj. Gen. Narayan Datta Paudel", cat: "person", desc: "APF HQ Halchok Operations Commander. APF Act §112 action recommended. Called Sudan Gurung on Bhadra 24 to help save APF HQ.", section: "§13.2" },
  { id: "shrestha_apf", label: "APF Maj. Gen. Suresh Kumar Shrestha", cat: "person", desc: "APF Battalion Commander, Kathmandu. APF Act §112.", section: "§13.2" },
  { id: "jeevan_kc", label: "APF SP Jeevan KC", cat: "person", desc: "APF Disaster Rescue, Sinamangal. APF Act §112.", section: "§13.2" },

  // PEOPLE — organizers
  { id: "sudan_gurung", label: "Sudan Gurung", cat: "person", desc: "Primary GenZ organiser, President of Hami Nepal NGO. Filed CDO permit. Spent NPR 50,000 personal funds. Tried to control crowd with megaphone. On Bhadra 24, personally saved both Police HQ and APF HQ at Khapung's request. IT company + Adventures Coach entrepreneur.", section: "§5.2, testimony #126" },
  { id: "anil_baniya", label: "Anil Baniya", cat: "person", desc: "Co-organiser, Hami Nepal. Co-signed CDO permit. Called DIGP Bishwa Adhikari the night before to warn about Discord Molotov discussions. Also visited Cyber Bureau in person.", section: "§5.2" },
  { id: "bablu_gupta", label: "Bablu Gupta", cat: "person", desc: "Co-organiser of Bhadra 23 march. Now a serving government minister in the post-uprising cabinet. Most direct evidence of GenZ demand being met.", section: "§5.2" },
  { id: "raksha_bam", label: "Raksha Bam", cat: "person", desc: "Second main GenZ leadership cluster. Led alongside Sudan Gurung. Present at Maitighar from early morning. Gave speeches from the vehicle.", section: "§5.2" },
  { id: "purushottam", label: "Purushottam Yadav", cat: "person", desc: "Journalism student, unemployed. Filed separate CDO permit with Sabal Gautam. Did not know Gurung beforehand. Hid in college building overnight when shooting started.", section: "§5.2, testimony #132" },
  { id: "sushila_karki", label: "PM Sushila Karki", cat: "person", desc: "Former Chief Justice. Won Discord PM vote. Subsequently appointed Prime Minister after Oli's resignation. GenZ's digital democracy worked — the Discord vote outcome became reality.", section: "§12.13" },

  // PEOPLE — heroes
  { id: "nagata_shah", label: "Nagata Kumari Shah", cat: "person", desc: "Watch shop owner, Koteshwor. Ran to save Police Constable Chij Kumar Kumal being beaten by mob. Crowd poured petrol on her — she caught fire. Still saved his life. Commission recommends STATE AWARD.", section: "§13.2" },
  { id: "shivaram_bada", label: "Shivaram Bada", cat: "person", desc: "Hotel Swarnim & Guest House owner, Koteshwor. Hid 7+ police officers being hunted by mob on Bhadra 24. Gave food, clothes, locked doors. Commission recommends award.", section: "§13.2" },
  { id: "ekta_shah", label: "Ekta Shah", cat: "person", desc: "Shot in the left knee at Parliament south gate, Bhadra 23. Sat on a stretcher and took her MBBS internship exam. Scored 57.5% — passed. Commission personally lobbied Government for her scholarship. Government accepted.", section: "§13.2" },

  // PEOPLE — infiltrators/misinformation
  { id: "dilbhushan", label: "Dilbhushan Pathak", cat: "person", desc: "Journalist/YouTuber ('Tough Talk'). Falsely claimed Hilton Hotel owned by Deuba's son Jayavir Singh Deuba. No fact-check. No correction after hotel burned. Commission interrogated him.", section: "§14.10" },
  { id: "tanka_dhakal", label: "Tanka Dhakal", cat: "person", desc: "Social media influencer. Spread '35 skeletons found at Bhatbhateni' false claim. Named first of five co-spreading influencers in §14.10.", section: "§14.10" },
  { id: "tob_leader", label: "Tenzin Dawa Lama (TOB)", cat: "person", desc: "Founder of TOB (Tibetan Original Blood / 'The Original Brother') motorcycle club. Tattoo artist, Kapan. ~12 members now. TOB identified on CCTV inciting violence at Baneswor. Commission recommends prosecution under Penal Code §35.", section: "§13.2, §5.2, testimony #134" },
  { id: "prasai", label: "Durga Prasai", cat: "person", desc: "Political agitator/businessman. Pro-monarchy, Hindu nationalist circles. APF DIGP's sworn testimony named 'Durga Prasai's interest groups' among Bhadra 24 arsonists. Sunsari witness: 'those who burned were all Durga Prasai's people.' Report refers to 'Prasai's previous andolan' as precedent.", section: "§5, §7.1, §14.10" },
  { id: "discord_tony", label: "Discord 'Tony'", cat: "person", desc: "Spread false rape claim about Global College hostel. College burned on Bhadra 24. Police found nothing. Sexual violence rumour spread faster than any denial could travel.", section: "§12.13" },
  { id: "discord_idke", label: "Discord 'IDKEHUAIM'", cat: "person", desc: "First to post Molotov bomb-making instructions: 'Guys, make Molotov...' at 10:40 PM Sept 7. Triggered 356 Molotov mentions over the next 36 hours.", section: "§12.13" },
  { id: "diwakar_dulal", label: "Diwakar Dulal", cat: "person", desc: "Discord user who posted 'save the data centre' 87 times on Sept 9 as Singha Durbar burned. Only named Discord user who tried to stop violence. Ignored by entire server. Data centre lost.", section: "§12.13" },
  { id: "subhash_bhandari", label: "Subhash Bhandari", cat: "person", desc: "Son of Puskar Bhandari. Operating from Basant Vihar, New Delhi, India. Created fake email, impersonated a Nepal Police officer, filed fraudulent complaint. Referred for criminal investigation.", section: "§29.6" },

  // ORGANISATIONS
  { id: "hami_nepal", label: "Hami Nepal NGO", cat: "org", desc: "Non-profit registered at Narapureshwar, Kathmandu. Founded Bhadra 24, 2077 (2020). Track record: 2072 earthquake hospital management, COVID plasma bank + 270 ICU beds, Turkey earthquake 14 tons relief. Admin of 'Youth Against Corruption' Discord. Medical camp at Maitighar on Bhadra 23. Saved both Police HQ and APF HQ on Bhadra 24.", section: "§5.2" },
  { id: "nepal_police", label: "Nepal Police", cat: "org", desc: "2,776 rubber bullets, 7,873 total rounds fired on Bhadra 23–24. 3 officers killed by protesters. 1,616+ injured. 219 police units damaged. NPR 1.59 Cr damage. ~4 hours of live fire without stop order on Bhadra 23. Refused to provide ammunition accounting to Commission.", section: "§12.1, §12.3" },
  { id: "apf", label: "Armed Police Force (APF)", cat: "org", desc: "APF deployed alongside Nepal Police. Victor Control issued by APF IGP Ayal on Bhadra 24 including prisoner release order. Three APF officers recommended for APF Act §112 action.", section: "§13.2" },
  { id: "nepal_army", label: "Nepal Army", cat: "org", desc: "Deployed to vital installations. Thwarted airport attack. Saved President Karki by evacuation. 4 installation commanders recommended for Army Act §105 action for failing to protect Shital Niwas, Baluwatar, Singha Durbar, Parliament.", section: "§13.2" },
  { id: "nid", label: "National Investigation Dept (NID)", cat: "org", desc: "Intelligence failure. Failed to provide timely warning about Discord violence plans or march escalation risk. Two joint directors recommended for Special Service Act action.", section: "§13.2" },
  { id: "discord_servers", label: "Discord Servers (Youth/Yuwa)", cat: "org", desc: "Two servers: 'Youth Against Corruption' (18,010 messages, Admin: Hami Nepal) and 'Yuwa Hub' (12,144 messages, 2-tier Core Member system). Combined 664,000+ messages in 50 hours. Used for march planning, PM vote (Karki won), Molotov coordination, and real-time arson coordination.", section: "§12.13" },
  { id: "cyabra", label: "Cyabra (Israel)", cat: "org", desc: "Israeli social media analytics firm. Report on Nepal protests: 34% of X profiles were fake or inauthentic. Commission notes methodological concerns with the report but cites it as evidence of coordinated inauthentic behaviour.", section: "§14.10" },
  { id: "cmr", label: "Centre for Media Research", cat: "org", desc: "Research finding cited by Commission: 75% of political/social news in Nepal contains misinformation. 76.17% of that originates from media and social media.", section: "§14.10" },
  { id: "himalayan_airlines", label: "Himalayan Airlines", cat: "org", desc: "Issued formal press denial after Air Hostess falsely claimed PM Oli boarded their flight to Dubai. Only corporate entity in the report to issue such a denial.", section: "§14.10" },

  // PLACES
  { id: "maitighar", label: "Maitighar Mandala", cat: "place", desc: "Starting point of Bhadra 23 march. GenZ gathered from 7–8 AM. Hami Nepal ran medical camp here. QR code pamphlets for Discord server distributed here. Crowd ~40–50 initially; swelled to 50,000+.", section: "§5.1.2, §12.2" },
  { id: "parliament", label: "Parliament / BICC", cat: "place", desc: "New Baneswor. Designated march endpoint (restricted zone perimeter). Crowd broke barricade at ~11:47 AM. Live fire at main gate from 12:37 PM. Army guarding but under-resourced. Lt. Col. Santosh Dhungel's command.", section: "§12.2" },
  { id: "singha_durbar", label: "Singha Durbar", cat: "place", desc: "Nepal's government secretariat complex. Partially burned on Bhadra 24 ~2:05 PM. Col. Ganesh Khadka's Army command. Diwakar Dulal posted 87 Discord pleas to save its data centre — all ignored.", section: "§5.1.4" },
  { id: "baluwatar", label: "Baluwatar (PM Residence)", cat: "place", desc: "Prime Minister's official residence. Burned on Bhadra 24. Col. Diwakar Khadka's Army command.", section: "§5.1.4" },
  { id: "shital_niwas", label: "Shital Niwas (President)", cat: "place", desc: "Presidential residence. Attacked and burned on Bhadra 24. Brig. Gen. Manoj Baidwar's Army command. President evacuated by helicopter at ~3:15 PM.", section: "§5.1.4" },
  { id: "supreme_court", label: "Supreme Court", cat: "place", desc: "Burned on Bhadra 24. Nepal's highest court. Commission notes 148,000+ pending cases nationwide, including 28,000+ at Supreme Court alone.", section: "§5.1.4, §14.3" },
  { id: "civil_hospital", label: "Civil Hospital", cat: "place", desc: "Nearest hospital to Parliament. First injured (lower body wound) arrived 12:07 PM. Overwhelmed with gunshot victims from 12:41 PM. Hami Nepal volunteers formed human chain at gate to prevent both sides entering.", section: "§12.2, §5.2" },
  { id: "baneswor", label: "New Baneswor", cat: "place", desc: "Key intersection. Crowd reached barricade here ~11 AM. TOB motorcycles arrived from both sides of Serviceway at ~11:40 AM. BICC PTZ camera captured critical escalation moments.", section: "§12.2" },

  // DIGITAL
  { id: "discord_tony_claim", label: "False Rape Claim", cat: "digital", desc: "Discord 'Tony' posted false claim rape was occurring at Global College hostel during protest. Spread faster than any denial. Police confirmed nothing happened. College burned on Bhadra 24.", section: "§12.13, §14.10" },
  { id: "molotov_356", label: "Molotov: 356 Mentions", cat: "digital", desc: "Molotov cocktails discussed 356 times across both Discord servers over 50 hours. IDKEHUAIM's first post at 10:40 PM Sept 7 triggered the cascade. CCTV confirmed actual Molotov attacks at Parliament.", section: "§12.13" },
  { id: "deepfakes", label: "AI Deepfake Videos", cat: "digital", desc: "Commission confirmed AI-edited/deepfake videos claiming security force atrocities were present in evidence (§12.5, item 6). Metadata stripped — verification impossible. First confirmed AI disinformation in Nepal's political history.", section: "§12.5" },
  { id: "netakhor", label: "netakhor.vercel.app", cat: "digital", desc: "Web platform identified in security intelligence reports as used for mapping attack sites during Bhadra 24. Commission notes this as evidence of coordinated targeting.", section: "§5.2" },
  { id: "bts_data", label: "BTS / CDR Data", cat: "digital", desc: "Commission obtained actual cell tower (BTS) data and Call Detail Records (CDR) for all senior officials. NCELL data showed 39,212 SIM users at key sites during events. CDR mapped entire crisis communication network.", section: "§12.6" },
  { id: "nepo_baby", label: "#NepoBaby Campaigns", cat: "digital", desc: "Pre-march TikTok/Reddit/X campaigns: #NepoBaby, #NepoBabies, #EnoughIsEnough, #NoNotAgain, #AntiCorruption. Weeks of digital organising before any physical march.", section: "§14" },

  // LEGAL
  { id: "pc181", label: "Muluki Penal Code §181", cat: "legal", desc: "Reckless killing (लापरबाहीपूर्ण काम गरी ज्यान मार्न नहुने). Applied to PM Oli, Home Minister Lekhak, IGP Khapung. Commission cites Duckers v. Lynch, Boswell v. State, Pape v. Time in its legal analysis.", section: "§13.1" },
  { id: "pc182", label: "Muluki Penal Code §182", cat: "legal", desc: "Negligent killing (हेलचेक्र्याई गरी ज्यान मार्न नहुने). Applied to: all three §181 accused PLUS Home Secretary, APF IGP, NID Chief, CDO Kathmandu.", section: "§13.1, §13.2" },
  { id: "police_act9", label: "Police Act §9(4)", cat: "legal", desc: "Departmental action provision. Recommended for: AIGP Shah, DIGP Om Rana, DIGP Bishwa Adhikari, DSP Deep Samsher JBR, DSP Rishiram Kandel, IGP Dan Bahadur Kaki.", section: "§13.2" },
  { id: "army105", label: "Army Act §105", cat: "legal", desc: "Departmental action for vital installation commanders who failed to protect assigned sites. Applied to: 4 Army commanders (Brig. Gen. Baidwar, Col. Diwakar Khadka, Col. Ganesh Khadka, Lt. Col. Santosh Dhungel).", section: "§13.2" },
  { id: "apf112", label: "APF Act §112", cat: "legal", desc: "APF departmental action. Applied to: APF Maj. Gen. Narayan Datta Paudel, APF Maj. Gen. Suresh Kumar Shrestha, APF SP Jeevan KC.", section: "§13.2" },
  { id: "pc35", label: "Penal Code §35", cat: "legal", desc: "Criminal action recommendation for TOB motorcycle group for deliberately inciting violence at New Baneswor.", section: "§13.2" },
  { id: "ssa_rule10", label: "Special Service Act Rule 10.2", cat: "legal", desc: "Action for NID Joint Directors Krishna Prasad Khanal and Riben Kumar Gachhadar for intelligence failure.", section: "§13.2" },

  // FINDINGS
  { id: "no_cabinet_decision", label: "No Cabinet Decision on Force", cat: "finding", desc: "Commission found: Cabinet made NO decision on use of force on Bhadra 23. NSC met but made no army deployment decision. Central Security Committee held no emergency meeting even after deaths began. Civilians effectively silent while 42 were shot.", section: "§12.8, §12.9, §12.10" },
  { id: "command_vacuum", label: "Command Vacuum Finding", cat: "finding", desc: "CDR analysis shows PM's office had near-zero direct contact with field commanders. DIGP Bishwa Adhikari and DIGP Om Rana had 53 calls between just the two of them — frantic lateral coordination inside police chain with civilian command silent.", section: "§12.6.1" },
  { id: "bullet_accounting", label: "No Bullet Accounting", cat: "finding", desc: "Neither Nepal Police nor APF provided Commission with records of rounds issued vs expended vs returned. Commission: 'It was difficult to ascertain who specifically fired.' This is the most critical unresolved finding.", section: "§12.7, §13.2" },
  { id: "aimed_fire", label: "Aimed Fire Evidence", cat: "finding", desc: "Autopsy: 14 head shots + 20 chest shots among 48 confirmed gunshot deaths. International law standard prohibits aimed fire at head/chest in crowd situations. Commission: proportionality and necessity principles violated.", section: "§12.7" },
  { id: "all_bullets_nepal", label: "All Bullets = Nepal Forces", cat: "finding", desc: "Ballistic analysis of 14 samples from deaths + 67 from injuries: all match Nepal security forces' weapons (5.56mm rifle army/APF, 12-bore police, 7.62mm SLR police, 9mm police pistol). No external actor firing confirmed.", section: "§12.7" },
  { id: "bhadra24_organised", label: "Bhadra 24 = Organised, Not GenZ", cat: "finding", desc: "Commission conclusion: Bhadra 24 violence was carried out by criminal elements, opportunistic looters, and politically motivated infiltrators — NOT the GenZ movement. Pattern evidence: CCTV broken first, then arson; pickup trucks with petrol and machetes; simultaneous attacks across Kathmandu.", section: "§Upasamhar" },
  { id: "organiser_goodfaith", label: "Organisers Acted in Good Faith", cat: "finding", desc: "Commission finding: Hami Nepal organisers filed legal permits, gave written commitments, warned police about Discord violence in advance, deployed medical team, and personally tried to stop escalation. Violence was not their intent.", section: "§5.2" },
  { id: "deepfake_first", label: "First AI Disinformation in Nepal", cat: "finding", desc: "Commission confirmed AI-generated deepfake videos in evidence. First confirmed use of AI disinformation in Nepal's political history. Metadata stripped — producers unidentified.", section: "§12.5" },

  // OUTCOMES
  { id: "76_dead", label: "76 Total Deaths", cat: "outcome", desc: "Breakdown: 42 by security bullets, 9 fire/burns, 5 blunt trauma, 3 police killed by protesters, 12 unidentified, 1 foreign national. Age: 35 were 14–29 years old. By location: 54 in Kathmandu, 10 in prisons, remainder across districts.", section: "§12.1, §12.7" },
  { id: "2522_injured", label: "2,522 Injured", cat: "outcome", desc: "1,328 protesters injured (1,325 discharged, 3 still hospitalized at report date). 2,068 security forces injured (all discharged). Government paid NPR 35,000 per injured person in two instalments.", section: "§12.1" },
  { id: "weapons_looted", label: "597 Weapons Still Missing", cat: "outcome", desc: "1,342 weapons + 118,525 rounds looted from police/APF. Recovered: 745 weapons + 28,068 rounds. Still missing: 597 weapons + ~90,457 rounds. Commission: 'active national security threat.'", section: "§12.3" },
  { id: "damage_npc", label: "NPR 84.45 Crore Damage (NPC)", cat: "outcome", desc: "National Planning Commission estimate. Commission's own field estimate: NPR 41.89 Cr (couldn't visit all 77 districts). Includes government buildings, party offices, private businesses, vehicles, infrastructure.", section: "§12.4" },
  { id: "compensation_paid", label: "NPR 15 Lakh per Death", cat: "outcome", desc: "Government paid NPR 1.5 million to each death family. Martyrdom declared for all Bhadra 23 dead. Families' uniform demand in testimony: legal action against the guilty (appears in nearly every testimony).", section: "§5.6" },
  { id: "chemicals_india", label: "Arson Samples Sent to India", cat: "outcome", desc: "Chemical samples from Parliament, Singha Durbar, hotel fires sent to Mirzapur, India for analysis (Nepal has no forensic chemistry lab). Results not received by report date. Sodium/magnesium accelerant hypothesis unconfirmed.", section: "§12.13" },
  { id: "fatf_greylist", label: "FATF Grey-listed (2×)", cat: "outcome", desc: "Nepal grey-listed 2008 and 2025. Commission: every political party treated it with 'indifference.' Three unresolved mega-cases: Fake Bhutanese Refugee Scam, Lalita Niwas land fraud (131 convicted, no AML probe), Bhatbhateni fake invoice scam.", section: "§14" },
  { id: "crypto_emergency", label: "Crypto Seizure Emergency", cat: "outcome", desc: "Commission: Nepal has no government crypto wallet. When suspects arrested, co-conspirators can drain wallets remotely before trial. No MLAT for digital assets. Commission recommends Nepal Rastra Bank establish government crypto wallet 'immediately or national security will be affected.'", section: "§14" },
  { id: "political_recycling", label: "30 PMs in 35 Years", cat: "outcome", desc: "Commission documents: 30 Prime Ministers in 35 years (1990–2025). Zero completed full 5-year term. Koirala and Deuba: 5 terms each. KP Oli: 4 terms. Prachanda: 3 terms. CPI score: 34/100, rank 107/180.", section: "§14.1" },
  { id: "sushila_appointed", label: "Karki Appointed PM", cat: "outcome", desc: "Discord GenZ vote winner Sushila Karki was subsequently actually appointed Prime Minister of Nepal. The digital democracy created inside a gaming platform produced a real-world outcome.", section: "§12.13, post-events" },
  { id: "ekta_scholarship", label: "Ekta Shah Gets Scholarship", cat: "outcome", desc: "Commission lobbied government for Ekta Shah (shot in knee, took MBBS exam on stretcher, scored 57.5%) to receive MBBS scholarship. Government accepted. Commission decision 2082/10/15/5.", section: "§13.2" },
];

const links = [
  // Event chains
  { source: "ban26", target: "bhadra23", label: "triggered" },
  { source: "bhadra23", target: "bhadra24", label: "led to" },
  { source: "bhadra23", target: "commission", label: "caused formation of" },
  { source: "bhadra23", target: "curfew", label: "led to" },
  { source: "bhadra23", target: "pm_resign", label: "eventually caused" },
  { source: "bhadra24", target: "singha_durbar_burn", label: "includes" },
  { source: "bhadra24", target: "hilton_burn", label: "includes" },
  { source: "bhadra24", target: "global_college_burn", label: "includes" },
  { source: "bhadra24", target: "prison_collapse", label: "caused" },
  { source: "bhadra24", target: "weapons_looted", label: "caused" },
  { source: "bhadra24", target: "damage_npc", label: "caused" },
  { source: "discord_servers", target: "bhadra24", label: "coordinated" },
  { source: "discord_servers", target: "discord_vote", label: "hosted" },
  { source: "discord_vote", target: "sushila_appointed", label: "outcome was" },
  { source: "airport_attack", target: "nepal_army", label: "thwarted by" },

  // Person → Event
  { source: "kp_oli", target: "bhadra23", label: "commanded during" },
  { source: "lekhak", target: "bhadra23", label: "commanded during" },
  { source: "khapung", target: "bhadra23", label: "commanded during" },
  { source: "sudan_gurung", target: "bhadra23", label: "organised" },
  { source: "raksha_bam", target: "bhadra23", label: "organised" },
  { source: "anil_baniya", target: "bhadra23", label: "organised" },
  { source: "bablu_gupta", target: "bhadra23", label: "organised" },
  { source: "tob_leader", target: "bhadra23", label: "infiltrated" },
  { source: "prasai", target: "bhadra24", label: "interest groups in" },
  { source: "discord_tony", target: "global_college_burn", label: "caused via false claim" },
  { source: "dilbhushan", target: "hilton_burn", label: "caused via false claim" },
  { source: "discord_idke", target: "molotov_356", label: "started cascade" },
  { source: "diwakar_dulal", target: "singha_durbar_burn", label: "tried to stop (87 times)" },
  { source: "sushila_karki", target: "discord_vote", label: "won" },
  { source: "ekta_shah", target: "ekta_scholarship", label: "received" },
  { source: "nagata_shah", target: "bhadra24", label: "heroism during" },
  { source: "shivaram_bada", target: "bhadra24", label: "heroism during" },

  // Person → Legal charges
  { source: "kp_oli", target: "pc181", label: "charged under" },
  { source: "kp_oli", target: "pc182", label: "charged under" },
  { source: "lekhak", target: "pc181", label: "charged under" },
  { source: "lekhak", target: "pc182", label: "charged under" },
  { source: "khapung", target: "pc181", label: "charged under" },
  { source: "khapung", target: "pc182", label: "charged under" },
  { source: "dubadi", target: "pc182", label: "charged under" },
  { source: "ayal", target: "pc182", label: "charged under" },
  { source: "hut_raj", target: "pc182", label: "charged under" },
  { source: "rijal", target: "pc182", label: "charged under" },
  { source: "shah_aigp", target: "police_act9", label: "action under" },
  { source: "om_rana", target: "police_act9", label: "action under" },
  { source: "bishwa", target: "police_act9", label: "action under" },
  { source: "kaki_igp", target: "police_act9", label: "action under" },
  { source: "baidwar", target: "army105", label: "action under" },
  { source: "diwakar_khadka", target: "army105", label: "action under" },
  { source: "ganesh_khadka", target: "army105", label: "action under" },
  { source: "santosh_dhungel", target: "army105", label: "action under" },
  { source: "paudel_apf", target: "apf112", label: "action under" },
  { source: "shrestha_apf", target: "apf112", label: "action under" },
  { source: "jeevan_kc", target: "apf112", label: "action under" },
  { source: "khanal_nid", target: "ssa_rule10", label: "action under" },
  { source: "gachhadar", target: "ssa_rule10", label: "action under" },
  { source: "tob_leader", target: "pc35", label: "prosecution recommended" },

  // Org → Event
  { source: "hami_nepal", target: "bhadra23", label: "permitted + organised" },
  { source: "hami_nepal", target: "discord_servers", label: "admin of" },
  { source: "nepal_police", target: "bhadra23", label: "deployed at" },
  { source: "apf", target: "bhadra23", label: "deployed at" },
  { source: "nepal_army", target: "bhadra24", label: "deployed at" },
  { source: "nid", target: "bhadra23", label: "intelligence failed for" },

  // Findings
  { source: "commission", target: "no_cabinet_decision", label: "found" },
  { source: "commission", target: "command_vacuum", label: "found" },
  { source: "commission", target: "bullet_accounting", label: "found gap" },
  { source: "commission", target: "aimed_fire", label: "found" },
  { source: "commission", target: "all_bullets_nepal", label: "confirmed" },
  { source: "commission", target: "bhadra24_organised", label: "concluded" },
  { source: "commission", target: "organiser_goodfaith", label: "concluded" },
  { source: "bts_data", target: "command_vacuum", label: "revealed" },
  { source: "aimed_fire", target: "76_dead", label: "caused" },
  { source: "all_bullets_nepal", target: "aimed_fire", label: "supports" },
  { source: "chemicals_india", target: "bhadra24_organised", label: "supports" },

  // Outcomes
  { source: "bhadra23", target: "76_dead", label: "caused" },
  { source: "bhadra23", target: "2522_injured", label: "caused" },
  { source: "kp_oli", target: "no_cabinet_decision", label: "responsible for" },
  { source: "lekhak", target: "no_cabinet_decision", label: "responsible for" },
  { source: "political_recycling", target: "nepo_baby", label: "triggered" },
  { source: "nepo_baby", target: "bhadra23", label: "built momentum for" },
  { source: "fatf_greylist", target: "political_recycling", label: "result of" },
  { source: "commission", target: "compensation_paid", label: "recommended" },
  { source: "commission", target: "ekta_scholarship", label: "secured" },
  { source: "prison_collapse", target: "banke_prison", label: "worst case" },

  // Digital connections
  { source: "discord_servers", target: "nepo_baby", label: "amplified" },
  { source: "discord_servers", target: "molotov_356", label: "contains" },
  { source: "discord_servers", target: "discord_tony_claim", label: "hosted" },
  { source: "discord_tony_claim", target: "global_college_burn", label: "caused" },
  { source: "deepfakes", target: "deepfake_first", label: "is" },
  { source: "netakhor", target: "bhadra24", label: "used to map attacks" },
  { source: "cyabra", target: "discord_servers", label: "analysed" },

  // Place connections
  { source: "maitighar", target: "bhadra23", label: "start of" },
  { source: "baneswor", target: "bhadra23", label: "scene of escalation" },
  { source: "parliament", target: "bhadra23", label: "target of" },
  { source: "singha_durbar", target: "singha_durbar_burn", label: "site of" },
  { source: "baluwatar", target: "pm_resign", label: "residence of" },
  { source: "shital_niwas", target: "bhadra24", label: "attacked in" },
  { source: "supreme_court", target: "bhadra24", label: "burned in" },
  { source: "civil_hospital", target: "76_dead", label: "treated victims of" },

  // Prasai's previous andolan connection
  { source: "prev_andolan", target: "prasai", label: "led by" },
  { source: "prev_andolan", target: "bhadra23", label: "mistake pattern repeated in" },
];

// Layout
const W = window.innerWidth, H = window.innerHeight - 60;
let activeCategories = new Set(Object.keys(CATEGORIES));
let searchTerm = "";

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

// Simulation
const sim = d3.forceSimulation(nodes)
  .force("link", d3.forceLink(links).id(d => d.id).distance(d => {
    const cats = new Set([d.source.cat || "finding", d.target.cat || "finding"]);
    if (cats.has("legal") || cats.has("finding")) return 80;
    if (cats.has("outcome")) return 100;
    return 120;
  }).strength(0.3))
  .force("charge", d3.forceManyBody().strength(-300).distanceMax(400))
  .force("center", d3.forceCenter(W / 2, H / 2))
  .force("collision", d3.forceCollide(28));

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
    return Math.max(8, Math.min(22, 8 + linkCount * 1.2));
  })
  .attr("fill", d => CATEGORIES[d.cat]?.color + "cc" || "#88878088")
  .attr("stroke", d => CATEGORIES[d.cat]?.color || "#888")
  .attr("stroke-width", 1.5);

nodeSel.append("text")
  .attr("text-anchor", "middle").attr("dy", "0.35em")
  .attr("font-size", 9).attr("font-family", "-apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif")
  .attr("fill", d => CATEGORIES[d.cat]?.color || "#888")
  .attr("font-weight", "500")
  .attr("pointer-events", "none")
  .attr("y", d => {
    const linkCount = links.filter(l => l.source === d.id || l.target === d.id || (l.source.id === d.id) || (l.target.id === d.id)).length;
    return Math.max(8, Math.min(22, 8 + linkCount * 1.2)) + 11;
  })
  .text(d => d.label.length > 18 ? d.label.slice(0, 16) + "…" : d.label);

sim.on("tick", () => {
  linkSel
    .attr("x1", d => d.source.x).attr("y1", d => d.source.y)
    .attr("x2", d => d.target.x).attr("y2", d => d.target.y);
  nodeSel.attr("transform", d => `translate(${d.x},${d.y})`);
});

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
    const dir = sid === d.id ? "→" : "←";
    return `<div class="conn-item"><span class="conn-arrow">${dir}</span><span>${l.label} <strong>${otherNode?.label || other}</strong></span></div>`;
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
