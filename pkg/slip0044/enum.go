package slip0044

var typesBySymbol = map[string]CoinType{
  "BTC": CoinType_BITCOIN,
  "LTC": CoinType_LITECOIN,
  "DOGE": CoinType_DOGECOIN,
  "RDD": CoinType_REDDCOIN,
  "DASH": CoinType_DASH,
  "PPC": CoinType_PEERCOIN,
  "NMC": CoinType_NAMECOIN,
  "FTC": CoinType_FEATHERCOIN,
  "XCP": CoinType_COUNTERPARTY,
  "BLK": CoinType_BLACKCOIN,
  "NSR": CoinType_NUSHARES,
  "NBT": CoinType_NUBITS,
  "MZC": CoinType_MAZACOIN,
  "VIA": CoinType_VIACOIN,
  "XCH": CoinType_CLEARINGHOUSE,
  "RBY": CoinType_RUBYCOIN,
  "GRS": CoinType_GROESTLCOIN,
  "DGC": CoinType_DIGITALCOIN,
  "CCN": CoinType_CANNACOIN,
  "DGB": CoinType_DIGIBYTE,
  "MONA": CoinType_MONACOIN,
  "CLAM": CoinType_CLAMS,
  "XPM": CoinType_PRIMECOIN,
  "NEOS": CoinType_NEOSCOIN,
  "JBS": CoinType_JUMBUCKS,
  "ZRC": CoinType_ZIFTRCOIN,
  "VTC": CoinType_VERTCOIN,
  "NXT": CoinType_NXT,
  "BURST": CoinType_BURST,
  "MUE": CoinType_MONETARYUNIT,
  "ZOOM": CoinType_ZOOM,
  "VASH": CoinType_VIRTUAL_CASH,
  "CDN": CoinType_CANADA_ECOIN,
  "SDC": CoinType_SHADOWCASH,
  "PKB": CoinType_PARKBYTE,
  "PND": CoinType_PANDACOIN,
  "START": CoinType_STARTCOIN,
  "MOIN": CoinType_MOIN,
  "EXP": CoinType_EXPANSE,
  "EMC2": CoinType_EINSTEINIUM,
  "DCR": CoinType_DECRED,
  "XEM": CoinType_NEM,
  "PART": CoinType_PARTICL,
  "ARG": CoinType_ARGENTUM,
  "SHR": CoinType_SHREEJI,
  "GCR": CoinType_GLOBAL_CURRENCY_RESERVE,
  "NVC": CoinType_NOVACOIN,
  "AC": CoinType_ASIACOIN,
  "BTCD": CoinType_BITCOINDARK,
  "DOPE": CoinType_DOPECOIN,
  "TPC": CoinType_TEMPLECOIN,
  "AIB": CoinType_AIB,
  "EDRC": CoinType_EDRCOIN,
  "SYS": CoinType_SYSCOIN,
  "SLR": CoinType_SOLARCOIN,
  "SMLY": CoinType_SMILEYCOIN,
  "ETH": CoinType_ETHER,
  "ETC": CoinType_ETHER_CLASSIC,
  "PSB": CoinType_PESOBIT,
  "LDCN": CoinType_LANDCOIN,
  "XBC": CoinType_BITCOINPLUS,
  "IOP": CoinType_INTERNET_OF_PEOPLE,
  "NXS": CoinType_NEXUS,
  "INSN": CoinType_INSANECOIN,
  "OK": CoinType_OKCASH,
  "BRIT": CoinType_BRITCOIN,
  "CMP": CoinType_COMPCOIN,
  "CRW": CoinType_CROWN,
  "BELA": CoinType_BELACOIN,
  "ICX": CoinType_ICON,
  "FJC": CoinType_FUJICOIN,
  "MIX": CoinType_MIX,
  "XVG": CoinType_VERGE_CURRENCY,
  "EFL": CoinType_ELECTRONIC_GULDEN,
  "CLUB": CoinType_CLUBCOIN,
  "RICHX": CoinType_RICHCOIN,
  "POT": CoinType_POTCOIN,
  "QRK": CoinType_QUARKCOIN,
  "TRC": CoinType_TERRACOIN,
  "GRC": CoinType_GRIDCOIN,
  "AUR": CoinType_AURORACOIN,
  "IXC": CoinType_IXCOIN,
  "NLG": CoinType_GULDEN,
  "BITB": CoinType_BITBEAN,
  "BTA": CoinType_BATA,
  "XMY": CoinType_MYRIADCOIN,
  "BSD": CoinType_BITSEND,
  "UNO": CoinType_UNOBTANIUM,
  "MTR": CoinType_MASTERTRADER,
  "GB": CoinType_GOLDBLOCKS,
  "SHM": CoinType_SAHAM,
  "CRX": CoinType_CHRONOS,
  "BIQ": CoinType_UBIQUOIN,
  "EVO": CoinType_EVOTION,
  "STO": CoinType_SAVETHEOCEAN,
  "BIGUP": CoinType_BIGUP,
  "GAME": CoinType_GAMECREDITS,
  "DLC": CoinType_DOLLARCOINS,
  "ZYD": CoinType_ZAYEDCOIN,
  "DBIC": CoinType_DUBAICOIN,
  "STRAT": CoinType_STRATIS,
  "SH": CoinType_SHILLING,
  "MARS": CoinType_MARSCOIN,
  "UBQ": CoinType_UBIQ,
  "PTC": CoinType_PESETACOIN,
  "NRO": CoinType_NEUROCOIN,
  "ARK": CoinType_ARK,
  "USC": CoinType_ULTIMATESECURECASHMAIN,
  "THC": CoinType_HEMPCOIN,
  "LINX": CoinType_LINX,
  "ECN": CoinType_ECOIN,
  "DNR": CoinType_DENARIUS,
  "PINK": CoinType_PINKCOIN,
  "ATOM": CoinType_ATOM,
  "PIVX": CoinType_PIVX,
  "FLASH": CoinType_FLASHCOIN,
  "ZEN": CoinType_ZENCASH,
  "PUT": CoinType_PUTINCOIN,
  "ZNY": CoinType_BITZENY,
  "UNIFY": CoinType_UNIFY,
  "XST": CoinType_STEALTHCOIN,
  "BRK": CoinType_BREAKOUT_COIN,
  "VC": CoinType_VCASH,
  "XMR": CoinType_MONERO,
  "VOX": CoinType_VOXELS,
  "NAV": CoinType_NAVCOIN,
  "FCT": CoinType_FACTOM_FACTOIDS,
  "EC": CoinType_FACTOM_ENTRY_CREDITS,
  "ZEC": CoinType_ZCASH,
  "LSK": CoinType_LISK,
  "STEEM": CoinType_STEEM,
  "XZC": CoinType_ZCOIN,
  "RBTC": CoinType_RSK,
  "RPT": CoinType_REALPOINTCOIN,
  "LBC": CoinType_LBRY_CREDITS,
  "KMD": CoinType_KOMODO,
  "BSQ": CoinType_BISQ_TOKEN,
  "RIC": CoinType_RIECOIN,
  "XRP": CoinType_RIPPLE,
  "BCH": CoinType_BITCOIN_CASH,
  "NEBL": CoinType_NEBLIO,
  "ZCL": CoinType_ZCLASSIC,
  "XLM": CoinType_STELLAR_LUMENS,
  "NLC2": CoinType_NOLIMITCOIN2,
  "WHL": CoinType_WHALECOIN,
  "ERC": CoinType_EUROPECOIN,
  "DMD": CoinType_DIAMOND,
  "BTM": CoinType_BYTOM,
  "BIO": CoinType_BIOCOIN,
  "XWCC": CoinType_WHITECOIN_CLASSIC,
  "BTG": CoinType_BITCOIN_GOLD,
  "BTC2X": CoinType_BITCOIN_2X,
  "SSN": CoinType_SUPERSKYNET,
  "TOA": CoinType_TOACOIN,
  "BTX": CoinType_BITCORE,
  "ACC": CoinType_ADCOIN,
  "BCO": CoinType_BRIDGECOIN,
  "ELLA": CoinType_ELLAISM,
  "PIRL": CoinType_PIRL,
  "XNO": CoinType_NANO,
  "VIVO": CoinType_VIVO,
  "FRST": CoinType_FIRSTCOIN,
  "HNC": CoinType_HELLENICCOIN,
  "BUZZ": CoinType_BUZZ,
  "MBRS": CoinType_EMBER,
  "HC": CoinType_HCASH,
  "HTML": CoinType_HTMLCOIN,
  "ODN": CoinType_OBSIDIAN,
  "ONX": CoinType_ONIXCOIN,
  "RVN": CoinType_RAVENCOIN,
  "GBX": CoinType_GOBYTE,
  "BTCZ": CoinType_BITCOINZ,
  "POA": CoinType_POA,
  "NYC": CoinType_NEWYORKCOIN,
  "MXT": CoinType_MARTEXCOIN,
  "WC": CoinType_WINCOIN,
  "MNX": CoinType_MINEXCOIN,
  "BTCP": CoinType_BITCOIN_PRIVATE,
  "MUSIC": CoinType_MUSICOIN,
  "BCA": CoinType_BITCOIN_ATOM,
  "CRAVE": CoinType_CRAVE,
  "STAK": CoinType_STRAKS,
  "WBTC": CoinType_WORLD_BITCOIN,
  "LCH": CoinType_LITECASH,
  "EXCL": CoinType_EXCLUSIVECOIN,
  "LCC": CoinType_LITECOINCASH,
  "XFE": CoinType_FEIRM,
  "EOS": CoinType_EOS,
  "TRX": CoinType_TRON,
  "KOBO": CoinType_KOBOCOIN,
  "HUSH": CoinType_HUSH,
  "BAN": CoinType_BANANO,
  "ETF": CoinType_ETF,
  "OMNI": CoinType_OMNI,
  "BIFI": CoinType_BITCOINFILE,
  "UFO": CoinType_UNIFORM_FISCAL_OBJECT,
  "CNMC": CoinType_CRYPTONODES,
  "BCN": CoinType_BYTECOIN,
  "RIN": CoinType_RINGO,
  "ATP": CoinType_ALAYA,
  "EVT": CoinType_EVERITOKEN,
  "ATN": CoinType_ATN,
  "BIS": CoinType_BISMUTH,
  "NEET": CoinType_NEETCOIN,
  "BOPO": CoinType_BOPOCHAIN,
  "OOT": CoinType_UTRUM,
  "ALIAS": CoinType_ALIAS,
  "MONK": CoinType_MONKEY_PROJECT,
  "BOXY": CoinType_BOXYCOIN,
  "FLO": CoinType_FLO,
  "MEC": CoinType_MEGACOIN,
  "BTDX": CoinType_BITCLOUD,
  "XAX": CoinType_ARTAX,
  "ANON": CoinType_ANON,
  "LTZ": CoinType_LITECOINZ,
  "BITG": CoinType_BITCOIN_GREEN,
  "ICP": CoinType_INTERNET_COMPUTER,
  "SMART": CoinType_SMARTCASH,
  "XUEZ": CoinType_XUEZ,
  "HLM": CoinType_HELIUM,
  "WEB": CoinType_WEBCHAIN,
  "ACM": CoinType_ACTINIUM,
  "NOS": CoinType_NOS_STABLE_COINS,
  "BITC": CoinType_BITCASH,
  "HTH": CoinType_HELP_THE_HOMELESS_COIN,
  "TZC": CoinType_TREZARCOIN,
  "VAR": CoinType_VARDA,
  "IOV": CoinType_IOV,
  "FIO": CoinType_FIO,
  "BSV": CoinType_BITCOINSV,
  "DXN": CoinType_DEXON,
  "QRL": CoinType_QUANTUM_RESISTANT_LEDGER,
  "PCX": CoinType_CHAINX,
  "LOKI": CoinType_LOKI,
  "NIM": CoinType_NIMIQ,
  "SOV": CoinType_SOVEREIGN_COIN,
  "JCT": CoinType_JIBITAL_COIN,
  "SLP": CoinType_SIMPLE_LEDGER_PROTOCOL,
  "EWT": CoinType_ENERGY_WEB,
  "UC": CoinType_ULORD,
  "EXOS": CoinType_EXOS,
  "ECA": CoinType_ELECTRA,
  "SOOM": CoinType_SOOM,
  "XRD": CoinType_REDSTONE,
  "FREE": CoinType_FREECOIN,
  "NPW": CoinType_NEWPOWERCOIN,
  "BST": CoinType_BLOCKSTAMP,
  "NANO": CoinType_BITCOIN_NANO,
  "BTCC": CoinType_BITCOIN_CORE,
  "ZEST": CoinType_ZEST,
  "ABT": CoinType_ARCBLOCK,
  "PION": CoinType_PION,
  "DT3": CoinType_DREAMTEAM3,
  "ZBUX": CoinType_ZBUX,
  "KPL": CoinType_KEPLER,
  "TPAY": CoinType_TOKENPAY,
  "ZILLA": CoinType_CHAINZILLA,
  "ANK": CoinType_ANKER,
  "BCC": CoinType_BCCHAIN,
  "HPB": CoinType_HPB,
  "ONE": CoinType_ONE,
  "SBC": CoinType_SBC,
  "IPC": CoinType_IPCHAIN,
  "DMTC": CoinType_DOMINANTCHAIN,
  "OGC": CoinType_ONEGRAM,
  "SHIT": CoinType_SHITCOIN,
  "ANDES": CoinType_ANDESCOIN,
  "AREPA": CoinType_AREPACOIN,
  "BOLI": CoinType_BOLIVARCOIN,
  "RIL": CoinType_RILCOIN,
  "HTR": CoinType_HATHOR_NETWORK,
  "FCTID": CoinType_FACTOM_ID,
  "BRAVO": CoinType_BRAVO,
  "ALGO": CoinType_ALGORAND,
  "BZX": CoinType_BITCOINZERO,
  "GXX": CoinType_GRAVITYCOIN,
  "HEAT": CoinType_HEAT,
  "XDN": CoinType_DIGITALNOTE,
  "FSN": CoinType_FUSION,
  "CPC": CoinType_CAPRICOIN,
  "BOLD": CoinType_BOLD,
  "IOST": CoinType_IOST,
  "TKEY": CoinType_TKEYCOIN,
  "USE": CoinType_USECHAIN,
  "BCZ": CoinType_BITCOINCZ,
  "IOC": CoinType_IOCOIN,
  "ASF": CoinType_ASOFE,
  "MASS": CoinType_MASS,
  "FAIR": CoinType_FAIRCOIN,
  "NUKO": CoinType_NEKONIUM,
  "GNX": CoinType_GENARO_NETWORK,
  "DIVI": CoinType_DIVI_PROJECT,
  "CMT": CoinType_COMMUNITY,
  "EUNO": CoinType_EUNO,
  "IOTX": CoinType_IOTEX,
  "ONION": CoinType_DEEPONION,
  "8BIT": CoinType_EIGHTBIT,
  "ATC": CoinType_ATOKEN_COIN,
  "BTS": CoinType_BITSHARES,
  "CKB": CoinType_NERVOS_CKB,
  "UGAS": CoinType_ULTRAIN,
  "ADS": CoinType_ADSHARES,
  "ARA": CoinType_AURA,
  "ZIL": CoinType_ZILLIQA,
  "MOAC": CoinType_MOAC,
  "SWTC": CoinType_SWTC,
  "VNSC": CoinType_VNSCOIN,
  "PLUG": CoinType_PLUG,
  "MAN": CoinType_MATRIX_AI_NETWORK,
  "ECC": CoinType_ECCOIN,
  "RPD": CoinType_RAPIDS,
  "RAP": CoinType_RAPTURE,
  "GARD": CoinType_HASHGARD,
  "ZER": CoinType_ZERO,
  "EBST": CoinType_EBOOST,
  "SHARD": CoinType_SHARD,
  "MRX": CoinType_METRIX_COIN,
  "CMM": CoinType_COMMERCIUM,
  "BLOCK": CoinType_BLOCKNET,
  "AUDAX": CoinType_AUDAX,
  "LUNA": CoinType_TERRA,
  "ZPM": CoinType_ZPRIME,
  "KUVA": CoinType_KUVA_UTILITY_NOTE,
  "MEM": CoinType_MEMCOIN,
  "CS": CoinType_CREDITS,
  "SWIFT": CoinType_SWIFTCASH,
  "FIX": CoinType_FIX,
}

var typesByName = map[string]CoinType{
  "bitcoin": CoinType_BITCOIN,
  "testnet": CoinType_TESTNET,
  "litecoin": CoinType_LITECOIN,
  "dogecoin": CoinType_DOGECOIN,
  "reddcoin": CoinType_REDDCOIN,
  "dash": CoinType_DASH,
  "peercoin": CoinType_PEERCOIN,
  "namecoin": CoinType_NAMECOIN,
  "feathercoin": CoinType_FEATHERCOIN,
  "counterparty": CoinType_COUNTERPARTY,
  "blackcoin": CoinType_BLACKCOIN,
  "nushares": CoinType_NUSHARES,
  "nubits": CoinType_NUBITS,
  "mazacoin": CoinType_MAZACOIN,
  "viacoin": CoinType_VIACOIN,
  "clearinghouse": CoinType_CLEARINGHOUSE,
  "rubycoin": CoinType_RUBYCOIN,
  "groestlcoin": CoinType_GROESTLCOIN,
  "digitalcoin": CoinType_DIGITALCOIN,
  "cannacoin": CoinType_CANNACOIN,
  "digibyte": CoinType_DIGIBYTE,
  "open assets": CoinType_OPEN_ASSETS,
  "monacoin": CoinType_MONACOIN,
  "clams": CoinType_CLAMS,
  "primecoin": CoinType_PRIMECOIN,
  "neoscoin": CoinType_NEOSCOIN,
  "jumbucks": CoinType_JUMBUCKS,
  "ziftrcoin": CoinType_ZIFTRCOIN,
  "vertcoin": CoinType_VERTCOIN,
  "nxt": CoinType_NXT,
  "burst": CoinType_BURST,
  "monetaryunit": CoinType_MONETARYUNIT,
  "zoom": CoinType_ZOOM,
  "virtual cash": CoinType_VIRTUAL_CASH,
  "canada ecoin": CoinType_CANADA_ECOIN,
  "shadowcash": CoinType_SHADOWCASH,
  "parkbyte": CoinType_PARKBYTE,
  "pandacoin": CoinType_PANDACOIN,
  "startcoin": CoinType_STARTCOIN,
  "moin": CoinType_MOIN,
  "expanse": CoinType_EXPANSE,
  "einsteinium": CoinType_EINSTEINIUM,
  "decred": CoinType_DECRED,
  "nem": CoinType_NEM,
  "particl": CoinType_PARTICL,
  "argentum": CoinType_ARGENTUM,
  "libertas": CoinType_LIBERTAS,
  "posw coin": CoinType_POSW_COIN,
  "shreeji": CoinType_SHREEJI,
  "global currency reserve": CoinType_GLOBAL_CURRENCY_RESERVE,
  "novacoin": CoinType_NOVACOIN,
  "asiacoin": CoinType_ASIACOIN,
  "bitcoindark": CoinType_BITCOINDARK,
  "dopecoin": CoinType_DOPECOIN,
  "templecoin": CoinType_TEMPLECOIN,
  "aib": CoinType_AIB,
  "edrcoin": CoinType_EDRCOIN,
  "syscoin": CoinType_SYSCOIN,
  "solarcoin": CoinType_SOLARCOIN,
  "smileycoin": CoinType_SMILEYCOIN,
  "ether": CoinType_ETHER,
  "ether classic": CoinType_ETHER_CLASSIC,
  "pesobit": CoinType_PESOBIT,
  "landcoin": CoinType_LANDCOIN,
  "open chain": CoinType_OPEN_CHAIN,
  "bitcoinplus": CoinType_BITCOINPLUS,
  "internet of people": CoinType_INTERNET_OF_PEOPLE,
  "nexus": CoinType_NEXUS,
  "insanecoin": CoinType_INSANECOIN,
  "okcash": CoinType_OKCASH,
  "britcoin": CoinType_BRITCOIN,
  "compcoin": CoinType_COMPCOIN,
  "crown": CoinType_CROWN,
  "belacoin": CoinType_BELACOIN,
  "icon": CoinType_ICON,
  "fujicoin": CoinType_FUJICOIN,
  "mix": CoinType_MIX,
  "verge currency": CoinType_VERGE_CURRENCY,
  "electronic gulden": CoinType_ELECTRONIC_GULDEN,
  "clubcoin": CoinType_CLUBCOIN,
  "richcoin": CoinType_RICHCOIN,
  "potcoin": CoinType_POTCOIN,
  "quarkcoin": CoinType_QUARKCOIN,
  "terracoin": CoinType_TERRACOIN,
  "gridcoin": CoinType_GRIDCOIN,
  "auroracoin": CoinType_AURORACOIN,
  "ixcoin": CoinType_IXCOIN,
  "gulden": CoinType_GULDEN,
  "bitbean": CoinType_BITBEAN,
  "bata": CoinType_BATA,
  "myriadcoin": CoinType_MYRIADCOIN,
  "bitsend": CoinType_BITSEND,
  "unobtanium": CoinType_UNOBTANIUM,
  "mastertrader": CoinType_MASTERTRADER,
  "goldblocks": CoinType_GOLDBLOCKS,
  "saham": CoinType_SAHAM,
  "chronos": CoinType_CHRONOS,
  "ubiquoin": CoinType_UBIQUOIN,
  "evotion": CoinType_EVOTION,
  "savetheocean": CoinType_SAVETHEOCEAN,
  "bigup": CoinType_BIGUP,
  "gamecredits": CoinType_GAMECREDITS,
  "dollarcoins": CoinType_DOLLARCOINS,
  "zayedcoin": CoinType_ZAYEDCOIN,
  "dubaicoin": CoinType_DUBAICOIN,
  "stratis": CoinType_STRATIS,
  "shilling": CoinType_SHILLING,
  "marscoin": CoinType_MARSCOIN,
  "ubiq": CoinType_UBIQ,
  "pesetacoin": CoinType_PESETACOIN,
  "neurocoin": CoinType_NEUROCOIN,
  "ark": CoinType_ARK,
  "ultimatesecurecashmain": CoinType_ULTIMATESECURECASHMAIN,
  "hempcoin": CoinType_HEMPCOIN,
  "linx": CoinType_LINX,
  "ecoin": CoinType_ECOIN,
  "denarius": CoinType_DENARIUS,
  "pinkcoin": CoinType_PINKCOIN,
  "atom": CoinType_ATOM,
  "pivx": CoinType_PIVX,
  "flashcoin": CoinType_FLASHCOIN,
  "zencash": CoinType_ZENCASH,
  "putincoin": CoinType_PUTINCOIN,
  "bitzeny": CoinType_BITZENY,
  "unify": CoinType_UNIFY,
  "stealthcoin": CoinType_STEALTHCOIN,
  "breakout coin": CoinType_BREAKOUT_COIN,
  "vcash": CoinType_VCASH,
  "monero": CoinType_MONERO,
  "voxels": CoinType_VOXELS,
  "navcoin": CoinType_NAVCOIN,
  "factom factoids": CoinType_FACTOM_FACTOIDS,
  "factom entry credits": CoinType_FACTOM_ENTRY_CREDITS,
  "zcash": CoinType_ZCASH,
  "lisk": CoinType_LISK,
  "steem": CoinType_STEEM,
  "zcoin": CoinType_ZCOIN,
  "rsk": CoinType_RSK,
  "giftblock": CoinType_GIFTBLOCK,
  "realpointcoin": CoinType_REALPOINTCOIN,
  "lbry credits": CoinType_LBRY_CREDITS,
  "komodo": CoinType_KOMODO,
  "bisq token": CoinType_BISQ_TOKEN,
  "riecoin": CoinType_RIECOIN,
  "ripple": CoinType_RIPPLE,
  "bitcoin cash": CoinType_BITCOIN_CASH,
  "neblio": CoinType_NEBLIO,
  "zclassic": CoinType_ZCLASSIC,
  "stellar lumens": CoinType_STELLAR_LUMENS,
  "nolimitcoin2": CoinType_NOLIMITCOIN2,
  "whalecoin": CoinType_WHALECOIN,
  "europecoin": CoinType_EUROPECOIN,
  "diamond": CoinType_DIAMOND,
  "bytom": CoinType_BYTOM,
  "biocoin": CoinType_BIOCOIN,
  "whitecoin classic": CoinType_WHITECOIN_CLASSIC,
  "bitcoin gold": CoinType_BITCOIN_GOLD,
  "bitcoin 2x": CoinType_BITCOIN_2X,
  "superskynet": CoinType_SUPERSKYNET,
  "toacoin": CoinType_TOACOIN,
  "bitcore": CoinType_BITCORE,
  "adcoin": CoinType_ADCOIN,
  "bridgecoin": CoinType_BRIDGECOIN,
  "ellaism": CoinType_ELLAISM,
  "pirl": CoinType_PIRL,
  "nano": CoinType_NANO,
  "vivo": CoinType_VIVO,
  "firstcoin": CoinType_FIRSTCOIN,
  "helleniccoin": CoinType_HELLENICCOIN,
  "buzz": CoinType_BUZZ,
  "ember": CoinType_EMBER,
  "hcash": CoinType_HCASH,
  "htmlcoin": CoinType_HTMLCOIN,
  "obsidian": CoinType_OBSIDIAN,
  "onixcoin": CoinType_ONIXCOIN,
  "ravencoin": CoinType_RAVENCOIN,
  "gobyte": CoinType_GOBYTE,
  "bitcoinz": CoinType_BITCOINZ,
  "poa": CoinType_POA,
  "newyorkcoin": CoinType_NEWYORKCOIN,
  "martexcoin": CoinType_MARTEXCOIN,
  "wincoin": CoinType_WINCOIN,
  "minexcoin": CoinType_MINEXCOIN,
  "bitcoin private": CoinType_BITCOIN_PRIVATE,
  "musicoin": CoinType_MUSICOIN,
  "bitcoin atom": CoinType_BITCOIN_ATOM,
  "crave": CoinType_CRAVE,
  "straks": CoinType_STRAKS,
  "world bitcoin": CoinType_WORLD_BITCOIN,
  "litecash": CoinType_LITECASH,
  "exclusivecoin": CoinType_EXCLUSIVECOIN,
  "lynx": CoinType_LYNX,
  "litecoincash": CoinType_LITECOINCASH,
  "feirm": CoinType_FEIRM,
  "eos": CoinType_EOS,
  "tron": CoinType_TRON,
  "kobocoin": CoinType_KOBOCOIN,
  "hush": CoinType_HUSH,
  "banano": CoinType_BANANO,
  "etf": CoinType_ETF,
  "omni": CoinType_OMNI,
  "bitcoinfile": CoinType_BITCOINFILE,
  "uniform fiscal object": CoinType_UNIFORM_FISCAL_OBJECT,
  "cryptonodes": CoinType_CRYPTONODES,
  "bytecoin": CoinType_BYTECOIN,
  "ringo": CoinType_RINGO,
  "alaya": CoinType_ALAYA,
  "everitoken": CoinType_EVERITOKEN,
  "atn": CoinType_ATN,
  "bismuth": CoinType_BISMUTH,
  "neetcoin": CoinType_NEETCOIN,
  "bopochain": CoinType_BOPOCHAIN,
  "utrum": CoinType_UTRUM,
  "alias": CoinType_ALIAS,
  "monkey project": CoinType_MONKEY_PROJECT,
  "boxycoin": CoinType_BOXYCOIN,
  "flo": CoinType_FLO,
  "megacoin": CoinType_MEGACOIN,
  "bitcloud": CoinType_BITCLOUD,
  "artax": CoinType_ARTAX,
  "anon": CoinType_ANON,
  "litecoinz": CoinType_LITECOINZ,
  "bitcoin green": CoinType_BITCOIN_GREEN,
  "internet computer": CoinType_INTERNET_COMPUTER,
  "smartcash": CoinType_SMARTCASH,
  "xuez": CoinType_XUEZ,
  "helium": CoinType_HELIUM,
  "webchain": CoinType_WEBCHAIN,
  "actinium": CoinType_ACTINIUM,
  "nos stable coins": CoinType_NOS_STABLE_COINS,
  "bitcash": CoinType_BITCASH,
  "help the homeless coin": CoinType_HELP_THE_HOMELESS_COIN,
  "trezarcoin": CoinType_TREZARCOIN,
  "varda": CoinType_VARDA,
  "iov": CoinType_IOV,
  "fio": CoinType_FIO,
  "bitcoinsv": CoinType_BITCOINSV,
  "dexon": CoinType_DEXON,
  "quantum resistant ledger": CoinType_QUANTUM_RESISTANT_LEDGER,
  "chainx": CoinType_CHAINX,
  "loki": CoinType_LOKI,
  "imagewallet": CoinType_IMAGEWALLET,
  "nimiq": CoinType_NIMIQ,
  "sovereign coin": CoinType_SOVEREIGN_COIN,
  "jibital coin": CoinType_JIBITAL_COIN,
  "simple ledger protocol": CoinType_SIMPLE_LEDGER_PROTOCOL,
  "energy web": CoinType_ENERGY_WEB,
  "ulord": CoinType_ULORD,
  "exos": CoinType_EXOS,
  "electra": CoinType_ELECTRA,
  "soom": CoinType_SOOM,
  "redstone": CoinType_REDSTONE,
  "freecoin": CoinType_FREECOIN,
  "newpowercoin": CoinType_NEWPOWERCOIN,
  "blockstamp": CoinType_BLOCKSTAMP,
  "smartholdem": CoinType_SMARTHOLDEM,
  "bitcoin nano": CoinType_BITCOIN_NANO,
  "bitcoin core": CoinType_BITCOIN_CORE,
  "zen protocol": CoinType_ZEN_PROTOCOL,
  "zest": CoinType_ZEST,
  "arcblock": CoinType_ARCBLOCK,
  "pion": CoinType_PION,
  "dreamteam3": CoinType_DREAMTEAM3,
  "zbux": CoinType_ZBUX,
  "kepler": CoinType_KEPLER,
  "tokenpay": CoinType_TOKENPAY,
  "chainzilla": CoinType_CHAINZILLA,
  "anker": CoinType_ANKER,
  "bcchain": CoinType_BCCHAIN,
  "hpb": CoinType_HPB,
  "one": CoinType_ONE,
  "sbc": CoinType_SBC,
  "ipchain": CoinType_IPCHAIN,
  "dominantchain": CoinType_DOMINANTCHAIN,
  "onegram": CoinType_ONEGRAM,
  "shitcoin": CoinType_SHITCOIN,
  "andescoin": CoinType_ANDESCOIN,
  "arepacoin": CoinType_AREPACOIN,
  "bolivarcoin": CoinType_BOLIVARCOIN,
  "rilcoin": CoinType_RILCOIN,
  "hathor network": CoinType_HATHOR_NETWORK,
  "factom id": CoinType_FACTOM_ID,
  "bravo": CoinType_BRAVO,
  "algorand": CoinType_ALGORAND,
  "bitcoinzero": CoinType_BITCOINZERO,
  "gravitycoin": CoinType_GRAVITYCOIN,
  "heat": CoinType_HEAT,
  "digitalnote": CoinType_DIGITALNOTE,
  "fusion": CoinType_FUSION,
  "capricoin": CoinType_CAPRICOIN,
  "bold": CoinType_BOLD,
  "iost": CoinType_IOST,
  "tkeycoin": CoinType_TKEYCOIN,
  "usechain": CoinType_USECHAIN,
  "bitcoincz": CoinType_BITCOINCZ,
  "iocoin": CoinType_IOCOIN,
  "asofe": CoinType_ASOFE,
  "mass": CoinType_MASS,
  "faircoin": CoinType_FAIRCOIN,
  "nekonium": CoinType_NEKONIUM,
  "genaro network": CoinType_GENARO_NETWORK,
  "divi project": CoinType_DIVI_PROJECT,
  "community": CoinType_COMMUNITY,
  "euno": CoinType_EUNO,
  "iotex": CoinType_IOTEX,
  "deeponion": CoinType_DEEPONION,
  "8bit": CoinType_EIGHTBIT,
  "atoken coin": CoinType_ATOKEN_COIN,
  "bitshares": CoinType_BITSHARES,
  "nervos ckb": CoinType_NERVOS_CKB,
  "ultrain": CoinType_ULTRAIN,
  "adshares": CoinType_ADSHARES,
  "aura": CoinType_AURA,
  "zilliqa": CoinType_ZILLIQA,
  "moac": CoinType_MOAC,
  "swtc": CoinType_SWTC,
  "vnscoin": CoinType_VNSCOIN,
  "pl^g": CoinType_PLUG,
  "matrix ai network": CoinType_MATRIX_AI_NETWORK,
  "eccoin": CoinType_ECCOIN,
  "rapids": CoinType_RAPIDS,
  "rapture": CoinType_RAPTURE,
  "hashgard": CoinType_HASHGARD,
  "zero": CoinType_ZERO,
  "eboost": CoinType_EBOOST,
  "shard": CoinType_SHARD,
  "metrix coin": CoinType_METRIX_COIN,
  "commercium": CoinType_COMMERCIUM,
  "blocknet": CoinType_BLOCKNET,
  "audax": CoinType_AUDAX,
  "terra": CoinType_TERRA,
  "zprime": CoinType_ZPRIME,
  "kuva utility note": CoinType_KUVA_UTILITY_NOTE,
  "memcoin": CoinType_MEMCOIN,
  "credits": CoinType_CREDITS,
  "swiftcash": CoinType_SWIFTCASH,
  "fix": CoinType_FIX,
}
