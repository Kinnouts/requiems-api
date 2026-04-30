package words

// definitionEntry holds a single definition for a word.
type definitionEntry struct {
	partOfSpeech string
	definition   string
	example      string
}

// dictionaryEntry holds the full dictionary record for a word.
type dictionaryEntry struct {
	phonetic    string
	definitions []definitionEntry
	synonyms    []string
}

// dictionaryData is the embedded dictionary dataset.
// All keys are stored in lowercase.
var dictionaryData = map[string]dictionaryEntry{
	"ephemeral": {
		phonetic: "/ɪˈfɛm(ə)rəl/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "lasting for a very short time",
				example:      "ephemeral pleasures",
			},
		},
		synonyms: []string{"transient", "fleeting", "momentary", "brief", "short-lived"},
	},
	"serendipity": {
		phonetic: "/ˌsɛr.ənˈdɪp.ɪ.ti/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "noun",
				definition:   "the occurrence and development of events by chance in a happy or beneficial way",
				example:      "a fortunate stroke of serendipity",
			},
		},
		synonyms: []string{"luck", "chance", "fortune", "providence", "happy accident"},
	},
	"melancholy": {
		phonetic: "/ˈmɛl.ən.kɒl.i/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "noun",
				definition:   "a feeling of pensive sadness, typically with no obvious cause",
				example:      "an air of melancholy surrounded him",
			},
			{
				partOfSpeech: "adjective",
				definition:   "having a feeling of melancholy; sad and pensive",
				example:      "she felt a little melancholy",
			},
		},
		synonyms: []string{"sadness", "sorrow", "gloom", "depression", "despondency"},
	},
	"resilience": {
		phonetic: "/rɪˈzɪl.i.əns/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "noun",
				definition:   "the capacity to recover quickly from difficulties; toughness",
				example:      "the resilience of the human spirit",
			},
		},
		synonyms: []string{"toughness", "hardiness", "adaptability", "flexibility", "strength"},
	},
	"eloquent": {
		phonetic: "/ˈɛl.ə.kwənt/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "fluent or persuasive in speaking or writing",
				example:      "an eloquent speech",
			},
		},
		synonyms: []string{"articulate", "expressive", "fluent", "persuasive", "well-spoken"},
	},
	"ambiguous": {
		phonetic: "/æmˈbɪɡ.ju.əs/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "open to more than one interpretation; not having one obvious meaning",
				example:      "an ambiguous statement",
			},
		},
		synonyms: []string{"unclear", "vague", "equivocal", "uncertain", "nebulous"},
	},
	"benevolent": {
		phonetic: "/bɪˈnɛv.ə.lənt/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "well-meaning and kindly",
				example:      "a benevolent smile",
			},
		},
		synonyms: []string{"kind", "generous", "charitable", "philanthropic", "compassionate"},
	},
	"candid": {
		phonetic: "/ˈkæn.dɪd/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "truthful and straightforward; frank",
				example:      "a candid discussion",
			},
		},
		synonyms: []string{"frank", "honest", "open", "forthright", "direct"},
	},
	"diligent": {
		phonetic: "/ˈdɪl.ɪ.dʒənt/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "having or showing care and conscientiousness in one's work or duties",
				example:      "a diligent student",
			},
		},
		synonyms: []string{"hardworking", "assiduous", "industrious", "conscientious", "painstaking"},
	},
	"enigma": {
		phonetic: "/ɪˈnɪɡ.mə/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "noun",
				definition:   "a person or thing that is mysterious or difficult to understand",
				example:      "Mona Lisa's smile is an enigma",
			},
		},
		synonyms: []string{"mystery", "puzzle", "riddle", "conundrum", "paradox"},
	},
	"fortuitous": {
		phonetic: "/fɔːˈtjuː.ɪ.təs/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "happening by accident or chance rather than design",
				example:      "a fortuitous meeting",
			},
		},
		synonyms: []string{"accidental", "chance", "unplanned", "unexpected", "lucky"},
	},
	"gregarious": {
		phonetic: "/ɡrɪˈɡɛər.i.əs/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "fond of company; sociable",
				example:      "she was a gregarious and well-liked teacher",
			},
		},
		synonyms: []string{"sociable", "outgoing", "friendly", "affable", "convivial"},
	},
	"hubris": {
		phonetic: "/ˈhjuː.brɪs/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "noun",
				definition:   "excessive pride or self-confidence",
				example:      "the hubris of the powerful",
			},
		},
		synonyms: []string{"arrogance", "conceit", "pride", "vanity", "egotism"},
	},
	"innocuous": {
		phonetic: "/ɪˈnɒk.ju.əs/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "not harmful or offensive",
				example:      "an innocuous remark",
			},
		},
		synonyms: []string{"harmless", "inoffensive", "safe", "benign", "unobjectionable"},
	},
	"jubilant": {
		phonetic: "/ˈdʒuː.bɪ.lənt/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "feeling or expressing great happiness and triumph",
				example:      "a jubilant crowd",
			},
		},
		synonyms: []string{"elated", "joyful", "triumphant", "exultant", "overjoyed"},
	},
	"keen": {
		phonetic: "/kiːn/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "having or showing eagerness or enthusiasm",
				example:      "a keen gardener",
			},
		},
		synonyms: []string{"eager", "enthusiastic", "avid", "zealous", "ardent"},
	},
	"lucid": {
		phonetic: "/ˈluː.sɪd/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "expressed clearly; easy to understand",
				example:      "a lucid account of the problem",
			},
		},
		synonyms: []string{"clear", "intelligible", "comprehensible", "transparent", "limpid"},
	},
	"meticulous": {
		phonetic: "/mɪˈtɪk.jʊ.ləs/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "showing great attention to detail or correct behaviour",
				example:      "he was meticulous in his record keeping",
			},
		},
		synonyms: []string{"careful", "thorough", "precise", "scrupulous", "conscientious"},
	},
	"nostalgia": {
		phonetic: "/nɒˈstæl.dʒə/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "noun",
				definition:   "a sentimental longing or wistful affection for the past",
				example:      "I was overcome with nostalgia for my school days",
			},
		},
		synonyms: []string{"wistfulness", "longing", "reminiscence", "sentimentality", "homesickness"},
	},
	"opulent": {
		phonetic: "/ˈɒp.jʊ.lənt/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "ostentatiously costly and luxurious",
				example:      "the opulent décor of the hotel",
			},
		},
		synonyms: []string{"luxurious", "lavish", "sumptuous", "rich", "affluent"},
	},
	"pragmatic": {
		phonetic: "/præɡˈmæt.ɪk/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "dealing with things sensibly and realistically in a way that is based on practical rather than theoretical considerations",
				example:      "a pragmatic approach to the problem",
			},
		},
		synonyms: []string{"practical", "realistic", "sensible", "rational", "matter-of-fact"},
	},
	"quintessential": {
		phonetic: "/ˌkwɪn.tɪˈsen.ʃəl/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "representing the most perfect or typical example of a quality or class",
				example:      "the quintessential English gentleman",
			},
		},
		synonyms: []string{"archetypal", "typical", "classic", "exemplary", "definitive"},
	},
	"resilient": {
		phonetic: "/rɪˈzɪl.i.ənt/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "able to withstand or recover quickly from difficult conditions",
				example:      "babies are generally resilient",
			},
		},
		synonyms: []string{"tough", "hardy", "adaptable", "flexible", "strong"},
	},
	"serenity": {
		phonetic: "/sɪˈrɛn.ɪ.ti/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "noun",
				definition:   "the state of being calm, peaceful, and untroubled",
				example:      "an atmosphere of serenity",
			},
		},
		synonyms: []string{"calm", "tranquility", "peace", "stillness", "composure"},
	},
	"tenacious": {
		phonetic: "/tɪˈneɪ.ʃəs/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "tending to keep a firm hold of something; clinging or adhering closely",
				example:      "a tenacious grip",
			},
		},
		synonyms: []string{"persistent", "determined", "resolute", "stubborn", "dogged"},
	},
	"ubiquitous": {
		phonetic: "/juːˈbɪk.wɪ.təs/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "present, appearing, or found everywhere",
				example:      "his ubiquitous influence",
			},
		},
		synonyms: []string{"omnipresent", "pervasive", "universal", "prevalent", "widespread"},
	},
	"verbose": {
		phonetic: "/vɜːˈbəʊs/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "using or expressed in more words than are needed",
				example:      "much academic writing is verbose",
			},
		},
		synonyms: []string{"wordy", "long-winded", "loquacious", "garrulous", "prolix"},
	},
	"whimsical": {
		phonetic: "/ˈwɪm.zɪ.kəl/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "playfully quaint or fanciful, especially in an appealing and amusing way",
				example:      "a whimsical sense of humour",
			},
		},
		synonyms: []string{"fanciful", "playful", "quaint", "quirky", "capricious"},
	},
	"zealous": {
		phonetic: "/ˈzɛl.əs/",
		definitions: []definitionEntry{
			{
				partOfSpeech: "adjective",
				definition:   "having or showing zeal; enthusiastically devoted to a cause",
				example:      "a zealous advocate of free trade",
			},
		},
		synonyms: []string{"fervent", "ardent", "passionate", "devoted", "fanatical"},
	},
}
