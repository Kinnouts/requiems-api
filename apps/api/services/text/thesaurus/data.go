package thesaurus

// entry holds synonyms and antonyms for a single word.
type entry struct {
	synonyms []string
	antonyms []string
}

// thesaurusData is the embedded thesaurus dataset.
// All words are stored in lowercase.
var thesaurusData = map[string]entry{
	"happy": {
		synonyms: []string{"joyful", "cheerful", "content", "pleased", "delighted", "glad", "elated", "blissful"},
		antonyms: []string{"sad", "unhappy", "miserable", "sorrowful", "dejected", "gloomy", "melancholy"},
	},
	"sad": {
		synonyms: []string{"unhappy", "sorrowful", "dejected", "miserable", "gloomy", "melancholy", "downcast"},
		antonyms: []string{"happy", "joyful", "cheerful", "elated", "blissful", "content"},
	},
	"big": {
		synonyms: []string{"large", "great", "huge", "enormous", "vast", "immense", "gigantic", "massive"},
		antonyms: []string{"small", "little", "tiny", "miniature", "petite", "minute"},
	},
	"small": {
		synonyms: []string{"little", "tiny", "miniature", "petite", "minute", "compact", "slight"},
		antonyms: []string{"big", "large", "great", "huge", "enormous", "vast", "immense"},
	},
	"fast": {
		synonyms: []string{"quick", "rapid", "swift", "speedy", "hasty", "brisk", "fleet"},
		antonyms: []string{"slow", "sluggish", "leisurely", "unhurried", "gradual"},
	},
	"slow": {
		synonyms: []string{"sluggish", "leisurely", "unhurried", "gradual", "plodding", "languid"},
		antonyms: []string{"fast", "quick", "rapid", "swift", "speedy", "hasty"},
	},
	"good": {
		synonyms: []string{"excellent", "fine", "superb", "splendid", "wonderful", "great", "positive", "favorable"},
		antonyms: []string{"bad", "poor", "terrible", "awful", "dreadful", "horrible", "inferior"},
	},
	"bad": {
		synonyms: []string{"poor", "terrible", "awful", "dreadful", "horrible", "inferior", "unpleasant", "wicked"},
		antonyms: []string{"good", "excellent", "fine", "superb", "splendid", "wonderful", "great"},
	},
	"strong": {
		synonyms: []string{"powerful", "mighty", "sturdy", "robust", "muscular", "vigorous", "tough"},
		antonyms: []string{"weak", "feeble", "frail", "fragile", "delicate", "powerless"},
	},
	"weak": {
		synonyms: []string{"feeble", "frail", "fragile", "delicate", "powerless", "flimsy", "helpless"},
		antonyms: []string{"strong", "powerful", "mighty", "sturdy", "robust", "muscular"},
	},
	"old": {
		synonyms: []string{"aged", "ancient", "elderly", "antique", "archaic", "vintage", "mature"},
		antonyms: []string{"new", "young", "modern", "fresh", "contemporary", "recent"},
	},
	"new": {
		synonyms: []string{"fresh", "modern", "contemporary", "recent", "novel", "current", "latest"},
		antonyms: []string{"old", "aged", "ancient", "antique", "archaic", "vintage", "outdated"},
	},
	"bright": {
		synonyms: []string{"luminous", "radiant", "brilliant", "vivid", "shining", "gleaming", "dazzling"},
		antonyms: []string{"dark", "dim", "dull", "murky", "gloomy", "shadowy"},
	},
	"dark": {
		synonyms: []string{"dim", "murky", "gloomy", "shadowy", "obscure", "dusky", "pitch-black"},
		antonyms: []string{"bright", "luminous", "radiant", "brilliant", "vivid", "shining"},
	},
	"hot": {
		synonyms: []string{"warm", "burning", "scorching", "blazing", "sweltering", "torrid", "fiery"},
		antonyms: []string{"cold", "cool", "chilly", "frigid", "icy", "freezing"},
	},
	"cold": {
		synonyms: []string{"cool", "chilly", "frigid", "icy", "freezing", "frosty", "arctic"},
		antonyms: []string{"hot", "warm", "burning", "scorching", "blazing", "sweltering"},
	},
	"beautiful": {
		synonyms: []string{"attractive", "lovely", "gorgeous", "stunning", "pretty", "elegant", "handsome"},
		antonyms: []string{"ugly", "unattractive", "hideous", "plain", "unsightly"},
	},
	"ugly": {
		synonyms: []string{"unattractive", "hideous", "plain", "unsightly", "repulsive", "grotesque"},
		antonyms: []string{"beautiful", "attractive", "lovely", "gorgeous", "stunning", "pretty"},
	},
	"smart": {
		synonyms: []string{"intelligent", "clever", "bright", "wise", "sharp", "astute", "brilliant"},
		antonyms: []string{"dumb", "stupid", "foolish", "ignorant", "dense", "obtuse"},
	},
	"brave": {
		synonyms: []string{"courageous", "bold", "fearless", "daring", "valiant", "heroic", "intrepid"},
		antonyms: []string{"cowardly", "timid", "fearful", "afraid", "craven", "spineless"},
	},
	"angry": {
		synonyms: []string{"furious", "mad", "enraged", "irate", "livid", "wrathful", "indignant"},
		antonyms: []string{"calm", "peaceful", "content", "pleased", "serene", "tranquil"},
	},
	"calm": {
		synonyms: []string{"peaceful", "serene", "tranquil", "composed", "placid", "quiet", "still"},
		antonyms: []string{"angry", "agitated", "upset", "turbulent", "anxious", "disturbed"},
	},
	"rich": {
		synonyms: []string{"wealthy", "affluent", "prosperous", "well-off", "opulent", "flush"},
		antonyms: []string{"poor", "broke", "destitute", "impoverished", "penniless", "needy"},
	},
	"poor": {
		synonyms: []string{"broke", "destitute", "impoverished", "penniless", "needy", "indigent"},
		antonyms: []string{"rich", "wealthy", "affluent", "prosperous", "well-off", "opulent"},
	},
	"love": {
		synonyms: []string{"adore", "cherish", "treasure", "care for", "like", "fondness", "affection"},
		antonyms: []string{"hate", "detest", "despise", "loathe", "abhor", "dislike"},
	},
	"hate": {
		synonyms: []string{"detest", "despise", "loathe", "abhor", "dislike", "resent"},
		antonyms: []string{"love", "adore", "cherish", "treasure", "like", "appreciate"},
	},
	"begin": {
		synonyms: []string{"start", "commence", "initiate", "launch", "embark", "open", "originate"},
		antonyms: []string{"end", "finish", "stop", "conclude", "terminate", "cease"},
	},
	"end": {
		synonyms: []string{"finish", "stop", "conclude", "terminate", "cease", "close", "complete"},
		antonyms: []string{"begin", "start", "commence", "initiate", "launch", "open"},
	},
	"hard": {
		synonyms: []string{"difficult", "tough", "challenging", "demanding", "arduous", "firm", "solid"},
		antonyms: []string{"easy", "simple", "soft", "effortless", "straightforward", "gentle"},
	},
	"easy": {
		synonyms: []string{"simple", "effortless", "straightforward", "uncomplicated", "basic", "clear"},
		antonyms: []string{"hard", "difficult", "tough", "challenging", "demanding", "arduous"},
	},
	"true": {
		synonyms: []string{"correct", "accurate", "genuine", "real", "valid", "authentic", "factual"},
		antonyms: []string{"false", "incorrect", "wrong", "fake", "untrue", "inaccurate"},
	},
	"false": {
		synonyms: []string{"incorrect", "wrong", "fake", "untrue", "inaccurate", "erroneous", "bogus"},
		antonyms: []string{"true", "correct", "accurate", "genuine", "real", "valid", "authentic"},
	},
	"give": {
		synonyms: []string{"provide", "offer", "grant", "bestow", "donate", "supply", "present"},
		antonyms: []string{"take", "receive", "withdraw", "withhold", "refuse", "deny"},
	},
	"take": {
		synonyms: []string{"grab", "seize", "acquire", "obtain", "receive", "accept", "capture"},
		antonyms: []string{"give", "provide", "offer", "grant", "donate", "return", "relinquish"},
	},
	"open": {
		synonyms: []string{"unlock", "unclose", "expose", "accessible", "available", "clear", "free"},
		antonyms: []string{"close", "shut", "seal", "lock", "block", "restrict"},
	},
	"close": {
		synonyms: []string{"shut", "seal", "lock", "fasten", "block", "restrict", "near"},
		antonyms: []string{"open", "unlock", "expose", "accessible", "available", "distant"},
	},
	"right": {
		synonyms: []string{"correct", "accurate", "proper", "true", "valid", "appropriate", "suitable"},
		antonyms: []string{"wrong", "incorrect", "false", "improper", "unsuitable", "inaccurate"},
	},
	"wrong": {
		synonyms: []string{"incorrect", "false", "improper", "unsuitable", "inaccurate", "erroneous"},
		antonyms: []string{"right", "correct", "accurate", "proper", "true", "valid", "appropriate"},
	},
	"clean": {
		synonyms: []string{"spotless", "tidy", "neat", "immaculate", "pure", "sanitary", "hygienic"},
		antonyms: []string{"dirty", "filthy", "messy", "unclean", "soiled", "grimy", "polluted"},
	},
	"dirty": {
		synonyms: []string{"filthy", "messy", "unclean", "soiled", "grimy", "polluted", "muddy"},
		antonyms: []string{"clean", "spotless", "tidy", "neat", "immaculate", "pure", "sanitary"},
	},
	"wet": {
		synonyms: []string{"damp", "moist", "soggy", "soaked", "drenched", "saturated", "humid"},
		antonyms: []string{"dry", "arid", "parched", "dehydrated"},
	},
	"dry": {
		synonyms: []string{"arid", "parched", "dehydrated", "desiccated", "waterless", "bone-dry"},
		antonyms: []string{"wet", "damp", "moist", "soggy", "soaked", "drenched", "saturated"},
	},
	"full": {
		synonyms: []string{"complete", "packed", "stuffed", "filled", "brimming", "overflowing", "satiated"},
		antonyms: []string{"empty", "vacant", "hollow", "bare", "blank", "void"},
	},
	"empty": {
		synonyms: []string{"vacant", "hollow", "bare", "blank", "void", "depleted", "barren"},
		antonyms: []string{"full", "complete", "packed", "stuffed", "filled", "brimming"},
	},
	"long": {
		synonyms: []string{"lengthy", "extended", "prolonged", "stretched", "elongated", "tall"},
		antonyms: []string{"short", "brief", "compact", "concise", "curtailed"},
	},
	"short": {
		synonyms: []string{"brief", "compact", "concise", "curtailed", "abbreviated", "stubby", "little"},
		antonyms: []string{"long", "lengthy", "extended", "prolonged", "stretched", "tall"},
	},
	"wide": {
		synonyms: []string{"broad", "expansive", "spacious", "roomy", "extensive", "ample"},
		antonyms: []string{"narrow", "thin", "slender", "tight", "cramped", "confined"},
	},
	"narrow": {
		synonyms: []string{"thin", "slender", "tight", "cramped", "confined", "restricted", "slim"},
		antonyms: []string{"wide", "broad", "expansive", "spacious", "roomy", "extensive"},
	},
}
