package musicxml

import (
	"fmt"
	"encoding/xml"
)

// From xsd simple type with enumerate restriction "above-below"
type AboveBelow string

const (
	Above AboveBelow = "above"
	Below AboveBelow = "below"
)

// From xsd simple type with enumerate restriction "accidental-value"
type AccidentalValue string

const (
	AccidentalSharp              AccidentalValue = "sharp"
	AccidentalNatural            AccidentalValue = "natural"
	AccidentalFlat               AccidentalValue = "flat"
	AccidentalDoubleSharp        AccidentalValue = "double-sharp"
	AccidentalSharpSharp         AccidentalValue = "sharp-sharp"
	AccidentalFlatFlat           AccidentalValue = "flat-flat"
	AccidentalNaturalSharp       AccidentalValue = "natural-sharp"
	AccidentalNaturalFlat        AccidentalValue = "natural-flat"
	AccidentalQuarterFlat        AccidentalValue = "quarter-flat"
	AccidentalQuarterSharp       AccidentalValue = "quarter-sharp"
	AccidentalThreeQuartersFlat  AccidentalValue = "three-quarters-flat"
	AccidentalThreeQuartersSharp AccidentalValue = "three-quarters-sharp"
	AccidentalSharpDown          AccidentalValue = "sharp-down"
	AccidentalSharpUp            AccidentalValue = "sharp-up"
	AccidentalNaturalDown        AccidentalValue = "natural-down"
	AccidentalNaturalUp          AccidentalValue = "natural-up"
	AccidentalFlatDown           AccidentalValue = "flat-down"
	AccidentalFlatUp             AccidentalValue = "flat-up"
	AccidentalDoubleSharpDown    AccidentalValue = "double-sharp-down"
	AccidentalDoubleSharpUp      AccidentalValue = "double-sharp-up"
	AccidentalFlatFlatDown       AccidentalValue = "flat-flat-down"
	AccidentalFlatFlatUp         AccidentalValue = "flat-flat-up"
	AccidentalArrowDown          AccidentalValue = "arrow-down"
	AccidentalArrowUp            AccidentalValue = "arrow-up"
	AccidentalTripleSharp        AccidentalValue = "triple-sharp"
	AccidentalTripleFlat         AccidentalValue = "triple-flat"
	AccidentalSlashQuarterSharp  AccidentalValue = "slash-quarter-sharp"
	AccidentalSlashSharp         AccidentalValue = "slash-sharp"
	AccidentalSlashFlat          AccidentalValue = "slash-flat"
	AccidentalDoubleSlashFlat    AccidentalValue = "double-slash-flat"
	AccidentalSharp1             AccidentalValue = "sharp-1"
	AccidentalSharp2             AccidentalValue = "sharp-2"
	AccidentalSharp3             AccidentalValue = "sharp-3"
	AccidentalSharp5             AccidentalValue = "sharp-5"
	AccidentalFlat1              AccidentalValue = "flat-1"
	AccidentalFlat2              AccidentalValue = "flat-2"
	AccidentalFlat3              AccidentalValue = "flat-3"
	AccidentalFlat4              AccidentalValue = "flat-4"
	AccidentalSori               AccidentalValue = "sori"
	AccidentalKoron              AccidentalValue = "koron"
	AccidentalOther              AccidentalValue = "other"
)

// From xsd simple type with enumerate restriction "arrow-direction"
type ArrowDirection string

const (
	ArrowDirectionLeft               ArrowDirection = "left"
	ArrowDirectionUp                 ArrowDirection = "up"
	ArrowDirectionRight              ArrowDirection = "right"
	ArrowDirectionDown               ArrowDirection = "down"
	ArrowDirectionNorthwest          ArrowDirection = "northwest"
	ArrowDirectionNortheast          ArrowDirection = "northeast"
	ArrowDirectionSoutheast          ArrowDirection = "southeast"
	ArrowDirectionSouthwest          ArrowDirection = "southwest"
	ArrowDirectionLeftRight          ArrowDirection = "left right"
	ArrowDirectionUpDown             ArrowDirection = "up down"
	ArrowDirectionNorthwestSoutheast ArrowDirection = "northwest southeast"
	ArrowDirectionNortheastSouthwest ArrowDirection = "northeast southwest"
	ArrowDirectionOther              ArrowDirection = "other"
)

// From xsd simple type with enumerate restriction "arrow-style"
type ArrowStyle string

const (
	ArrowStyleSingle   ArrowStyle = "single"
	ArrowStyleDouble   ArrowStyle = "double"
	ArrowStyleFilled   ArrowStyle = "filled"
	ArrowStyleHollow   ArrowStyle = "hollow"
	ArrowStylePaired   ArrowStyle = "paired"
	ArrowStyleCombined ArrowStyle = "combined"
	ArrowStyleOther    ArrowStyle = "other"
)

// From xsd simple type with enumerate restriction "backward-forward"
type BackwardForward string

const (
	BackwardForwardBackward BackwardForward = "backward"
	BackwardForwardForward  BackwardForward = "forward"
)

// From xsd simple type with enumerate restriction "bar-style"
type BarStyle string

const (
	BarStyleRegular    BarStyle = "regular"
	BarStyleDotted     BarStyle = "dotted"
	BarStyleDashed     BarStyle = "dashed"
	BarStyleHeavy      BarStyle = "heavy"
	BarStyleLightLight BarStyle = "light-light"
	BarStyleLightHeavy BarStyle = "light-heavy"
	BarStyleHeavyLight BarStyle = "heavy-light"
	BarStyleHeavyHeavy BarStyle = "heavy-heavy"
	BarStyleTick       BarStyle = "tick"
	BarStyleShort      BarStyle = "short"
	BarStyleNone       BarStyle = "none"
)

// From xsd simple type with enumerate restriction "beam-value"
type BeamValue string

const (
	BeamBegin        BeamValue = "begin"
	BeamContinue     BeamValue = "continue"
	BeamEnd          BeamValue = "end"
	BeamForwardHook  BeamValue = "forward hook"
	BeamBackwardHook BeamValue = "backward hook"
)

// From xsd simple type with enumerate restriction "beater-value"
type BeaterValue string

const (
	BeaterBow                 BeaterValue = "bow"
	BeaterChimeHammer         BeaterValue = "chime hammer"
	BeaterCoin                BeaterValue = "coin"
	BeaterDrumStick           BeaterValue = "drum stick"
	BeaterFinger              BeaterValue = "finger"
	BeaterFingernail          BeaterValue = "fingernail"
	BeaterFist                BeaterValue = "fist"
	BeaterGuiroScraper        BeaterValue = "guiro scraper"
	BeaterHammer              BeaterValue = "hammer"
	BeaterHand                BeaterValue = "hand"
	BeaterJazzStick           BeaterValue = "jazz stick"
	BeaterKnittingNeedle      BeaterValue = "knitting needle"
	BeaterMetalHammer         BeaterValue = "metal hammer"
	BeaterSlideBrushOnGong    BeaterValue = "slide brush on gong"
	BeaterSnareStick          BeaterValue = "snare stick"
	BeaterSpoonMallet         BeaterValue = "spoon mallet"
	BeaterSuperball           BeaterValue = "superball"
	BeaterTriangleBeater      BeaterValue = "triangle beater"
	BeaterTriangleBeaterPlain BeaterValue = "triangle beater plain"
	BeaterWireBrush           BeaterValue = "wire brush"
)

// From xsd simple type with enumerate restriction "bend-shape"
type BendShape string

const (
	BendShapeAngled BendShape = "angled"
	BendShapeCurved BendShape = "curved"
)

// From xsd simple type with enumerate restriction "breath-mark-value"
type BreathMarkValue string

const (
	BreathMarkEmpty   BreathMarkValue = ""
	BreathMarkComma   BreathMarkValue = "comma"
	BreathMarkTick    BreathMarkValue = "tick"
	BreathMarkUpbow   BreathMarkValue = "upbow"
	BreathMarkSalzedo BreathMarkValue = "salzedo"
)

// From xsd simple type with enumerate restriction "caesura-value"
type CaesuraValue string

const (
	CaesuraNormal CaesuraValue = "normal"
	CaesuraThick  CaesuraValue = "thick"
	CaesuraShort  CaesuraValue = "short"
	CaesuraCurved CaesuraValue = "curved"
	CaesuraSingle CaesuraValue = "single"
	CaesuraEmpty  CaesuraValue = ""
)

// From xsd simple type with enumerate restriction "cancel-location"
type CancelLocation string

const (
	CancelLocationLeft          CancelLocation = "left"
	CancelLocationRight         CancelLocation = "right"
	CancelLocationBeforeBarline CancelLocation = "before-barline"
)

// From xsd simple type with enumerate restriction "circular-arrow"
type CircularArrow string

const (
	CircularArrowClockwise     CircularArrow = "clockwise"
	CircularArrowAnticlockwise CircularArrow = "anticlockwise"
)

// From xsd simple type with enumerate restriction "clef-sign"
type ClefSign string

const (
	ClefSignG          ClefSign = "G"
	ClefSignF          ClefSign = "F"
	ClefSignC          ClefSign = "C"
	ClefSignPercussion ClefSign = "percussion"
	ClefSignTAB        ClefSign = "TAB"
	ClefSignJianpu     ClefSign = "jianpu"
	ClefSignNone       ClefSign = "none"
)

// From xsd simple type with enumerate restriction "color"
type Color string

const ()

// From xsd simple type with enumerate restriction "comma-separated-text"
type CommaSeparatedText string

const ()

// From xsd simple type with enumerate restriction "degree-symbol-value"
type DegreeSymbolValue string

const (
	DegreeSymbolMajor          DegreeSymbolValue = "major"
	DegreeSymbolMinor          DegreeSymbolValue = "minor"
	DegreeSymbolAugmented      DegreeSymbolValue = "augmented"
	DegreeSymbolDiminished     DegreeSymbolValue = "diminished"
	DegreeSymbolHalfDiminished DegreeSymbolValue = "half-diminished"
)

// From xsd simple type with enumerate restriction "degree-type-value"
type DegreeTypeValue string

const (
	DegreeTypeAdd      DegreeTypeValue = "add"
	DegreeTypeAlter    DegreeTypeValue = "alter"
	DegreeTypeSubtract DegreeTypeValue = "subtract"
)

// From xsd simple type with enumerate restriction "distance-type"
type DistanceType string

const ()

// From xsd simple type with enumerate restriction "effect-value"
type EffectValue string

const (
	EffectAnvil         EffectValue = "anvil"
	EffectAutoHorn      EffectValue = "auto horn"
	EffectBirdWhistle   EffectValue = "bird whistle"
	EffectCannon        EffectValue = "cannon"
	EffectDuckCall      EffectValue = "duck call"
	EffectGunShot       EffectValue = "gun shot"
	EffectKlaxonHorn    EffectValue = "klaxon horn"
	EffectLionsRoar     EffectValue = "lions roar"
	EffectLotusFlute    EffectValue = "lotus flute"
	EffectMegaphone     EffectValue = "megaphone"
	EffectPoliceWhistle EffectValue = "police whistle"
	EffectSiren         EffectValue = "siren"
	EffectSlideWhistle  EffectValue = "slide whistle"
	EffectThunderSheet  EffectValue = "thunder sheet"
	EffectWindMachine   EffectValue = "wind machine"
	EffectWindWhistle   EffectValue = "wind whistle"
)

// From xsd simple type with enumerate restriction "enclosure-shape"
type EnclosureShape string

const (
	EnclosureShapeRectangle       EnclosureShape = "rectangle"
	EnclosureShapeSquare          EnclosureShape = "square"
	EnclosureShapeOval            EnclosureShape = "oval"
	EnclosureShapeCircle          EnclosureShape = "circle"
	EnclosureShapeBracket         EnclosureShape = "bracket"
	EnclosureShapeInvertedBracket EnclosureShape = "inverted-bracket"
	EnclosureShapeTriangle        EnclosureShape = "triangle"
	EnclosureShapeDiamond         EnclosureShape = "diamond"
	EnclosureShapePentagon        EnclosureShape = "pentagon"
	EnclosureShapeHexagon         EnclosureShape = "hexagon"
	EnclosureShapeHeptagon        EnclosureShape = "heptagon"
	EnclosureShapeOctagon         EnclosureShape = "octagon"
	EnclosureShapeNonagon         EnclosureShape = "nonagon"
	EnclosureShapeDecagon         EnclosureShape = "decagon"
	EnclosureShapeNone            EnclosureShape = "none"
)

// From xsd simple type with enumerate restriction "ending-number"
type EndingNumber string

const ()

// From xsd simple type with enumerate restriction "fan"
type Fan string

const (
	FanAccel Fan = "accel"
	FanRit   Fan = "rit"
	FanNone  Fan = "none"
)

// From xsd simple type with enumerate restriction "fermata-shape"
type FermataShape string

const (
	FermataShapeNormal       FermataShape = "normal"
	FermataShapeAngled       FermataShape = "angled"
	FermataShapeSquare       FermataShape = "square"
	FermataShapeDoubleAngled FermataShape = "double-angled"
	FermataShapeDoubleSquare FermataShape = "double-square"
	FermataShapeDoubleDot    FermataShape = "double-dot"
	FermataShapeHalfCurve    FermataShape = "half-curve"
	FermataShapeCurlew       FermataShape = "curlew"
	FermataShapeEmpty        FermataShape = ""
)

// From xsd simple type with enumerate restriction "font-style"
type FontStyle string

const (
	FontStyleNormal FontStyle = "normal"
	FontStyleItalic FontStyle = "italic"
)

// From xsd simple type with enumerate restriction "font-weight"
type FontWeight string

const (
	FontWeightNormal FontWeight = "normal"
	FontWeightBold   FontWeight = "bold"
)

// From xsd simple type with enumerate restriction "glass-value"
type GlassValue string

const (
	GlassGlassHarmonica GlassValue = "glass harmonica"
	GlassGlassHarp      GlassValue = "glass harp"
	GlassWindChimes     GlassValue = "wind chimes"
)

// From xsd simple type with enumerate restriction "glyph-type"
type GlyphType string

const ()

// From xsd simple type with enumerate restriction "group-barline-value"
type GroupBarlineValue string

const (
	GroupBarlineYes          GroupBarlineValue = "yes"
	GroupBarlineNo           GroupBarlineValue = "no"
	GroupBarlineMensurstrich GroupBarlineValue = "Mensurstrich"
)

// From xsd simple type with enumerate restriction "group-symbol-value"
type GroupSymbolValue string

const (
	GroupSymbolNone    GroupSymbolValue = "none"
	GroupSymbolBrace   GroupSymbolValue = "brace"
	GroupSymbolLine    GroupSymbolValue = "line"
	GroupSymbolBracket GroupSymbolValue = "bracket"
	GroupSymbolSquare  GroupSymbolValue = "square"
)

// From xsd simple type with enumerate restriction "handbell-value"
type HandbellValue string

const (
	HandbellBelltree        HandbellValue = "belltree"
	HandbellDamp            HandbellValue = "damp"
	HandbellEcho            HandbellValue = "echo"
	HandbellGyro            HandbellValue = "gyro"
	HandbellHandMartellato  HandbellValue = "hand martellato"
	HandbellMalletLift      HandbellValue = "mallet lift"
	HandbellMalletTable     HandbellValue = "mallet table"
	HandbellMartellato      HandbellValue = "martellato"
	HandbellMartellatoLift  HandbellValue = "martellato lift"
	HandbellMutedMartellato HandbellValue = "muted martellato"
	HandbellPluckLift       HandbellValue = "pluck lift"
	HandbellSwing           HandbellValue = "swing"
)

// From xsd simple type with enumerate restriction "harmon-closed-location"
type HarmonClosedLocation string

const (
	HarmonClosedLocationRight  HarmonClosedLocation = "right"
	HarmonClosedLocationBottom HarmonClosedLocation = "bottom"
	HarmonClosedLocationLeft   HarmonClosedLocation = "left"
	HarmonClosedLocationTop    HarmonClosedLocation = "top"
)

// From xsd simple type with enumerate restriction "harmon-closed-value"
type HarmonClosedValue string

const (
	HarmonClosedYes  HarmonClosedValue = "yes"
	HarmonClosedNo   HarmonClosedValue = "no"
	HarmonClosedHalf HarmonClosedValue = "half"
)

// From xsd simple type with enumerate restriction "harmony-arrangement"
type HarmonyArrangement string

const (
	HarmonyArrangementVertical   HarmonyArrangement = "vertical"
	HarmonyArrangementHorizontal HarmonyArrangement = "horizontal"
	HarmonyArrangementDiagonal   HarmonyArrangement = "diagonal"
)

// From xsd simple type with enumerate restriction "harmony-type"
type HarmonyType string

const (
	HarmonyTypeExplicit  HarmonyType = "explicit"
	HarmonyTypeImplied   HarmonyType = "implied"
	HarmonyTypeAlternate HarmonyType = "alternate"
)

// From xsd simple type with enumerate restriction "hole-closed-location"
type HoleClosedLocation string

const (
	HoleClosedLocationRight  HoleClosedLocation = "right"
	HoleClosedLocationBottom HoleClosedLocation = "bottom"
	HoleClosedLocationLeft   HoleClosedLocation = "left"
	HoleClosedLocationTop    HoleClosedLocation = "top"
)

// From xsd simple type with enumerate restriction "hole-closed-value"
type HoleClosedValue string

const (
	HoleClosedYes  HoleClosedValue = "yes"
	HoleClosedNo   HoleClosedValue = "no"
	HoleClosedHalf HoleClosedValue = "half"
)

// From xsd simple type with enumerate restriction "kind-value"
type KindValue string

const (
	KindMajor             KindValue = "major"
	KindMinor             KindValue = "minor"
	KindAugmented         KindValue = "augmented"
	KindDiminished        KindValue = "diminished"
	KindDominant          KindValue = "dominant"
	KindMajorSeventh      KindValue = "major-seventh"
	KindMinorSeventh      KindValue = "minor-seventh"
	KindDiminishedSeventh KindValue = "diminished-seventh"
	KindAugmentedSeventh  KindValue = "augmented-seventh"
	KindHalfDiminished    KindValue = "half-diminished"
	KindMajorMinor        KindValue = "major-minor"
	KindMajorSixth        KindValue = "major-sixth"
	KindMinorSixth        KindValue = "minor-sixth"
	KindDominantNinth     KindValue = "dominant-ninth"
	KindMajorNinth        KindValue = "major-ninth"
	KindMinorNinth        KindValue = "minor-ninth"
	KindDominant11th      KindValue = "dominant-11th"
	KindMajor11th         KindValue = "major-11th"
	KindMinor11th         KindValue = "minor-11th"
	KindDominant13th      KindValue = "dominant-13th"
	KindMajor13th         KindValue = "major-13th"
	KindMinor13th         KindValue = "minor-13th"
	KindSuspendedSecond   KindValue = "suspended-second"
	KindSuspendedFourth   KindValue = "suspended-fourth"
	KindNeapolitan        KindValue = "Neapolitan"
	KindItalian           KindValue = "Italian"
	KindFrench            KindValue = "French"
	KindGerman            KindValue = "German"
	KindPedal             KindValue = "pedal"
	KindPower             KindValue = "power"
	KindTristan           KindValue = "Tristan"
	KindOther             KindValue = "other"
	KindNone              KindValue = "none"
)

// From xsd simple type with enumerate restriction "left-center-right"
type LeftCenterRight string

const (
	LeftCenterRightLeft   LeftCenterRight = "left"
	LeftCenterRightCenter LeftCenterRight = "center"
	LeftCenterRightRight  LeftCenterRight = "right"
)

// From xsd simple type with enumerate restriction "left-right"
type LeftRight string

const (
	LeftRightLeft  LeftRight = "left"
	LeftRightRight LeftRight = "right"
)

// From xsd simple type with enumerate restriction "line-end"
type LineEnd string

const (
	LineEndUp    LineEnd = "up"
	LineEndDown  LineEnd = "down"
	LineEndBoth  LineEnd = "both"
	LineEndArrow LineEnd = "arrow"
	LineEndNone  LineEnd = "none"
)

// From xsd simple type with enumerate restriction "line-length"
type LineLength string

const (
	LineLengthShort  LineLength = "short"
	LineLengthMedium LineLength = "medium"
	LineLengthLong   LineLength = "long"
)

// From xsd simple type with enumerate restriction "line-shape"
type LineShape string

const (
	LineShapeStraight LineShape = "straight"
	LineShapeCurved   LineShape = "curved"
)

// From xsd simple type with enumerate restriction "line-type"
type LineType string

const (
	LineTypeSolid  LineType = "solid"
	LineTypeDashed LineType = "dashed"
	LineTypeDotted LineType = "dotted"
	LineTypeWavy   LineType = "wavy"
)

// From xsd simple type with enumerate restriction "line-width-type"
type LineWidthType string

const ()

// From xsd simple type with enumerate restriction "margin-type"
type MarginType string

const (
	MarginTypeOdd  MarginType = "odd"
	MarginTypeEven MarginType = "even"
	MarginTypeBoth MarginType = "both"
)

// From xsd simple type with enumerate restriction "measure-numbering-value"
type MeasureNumberingValue string

const (
	MeasureNumberingNone    MeasureNumberingValue = "none"
	MeasureNumberingMeasure MeasureNumberingValue = "measure"
	MeasureNumberingSystem  MeasureNumberingValue = "system"
)

// From xsd simple type with enumerate restriction "measure-text"
type MeasureText string

const ()

// From xsd simple type with enumerate restriction "membrane-value"
type MembraneValue string

const (
	MembraneBassDrum           MembraneValue = "bass drum"
	MembraneBassDrumOnSide     MembraneValue = "bass drum on side"
	MembraneBongos             MembraneValue = "bongos"
	MembraneChineseTomtom      MembraneValue = "Chinese tomtom"
	MembraneCongaDrum          MembraneValue = "conga drum"
	MembraneCuica              MembraneValue = "cuica"
	MembraneGobletDrum         MembraneValue = "goblet drum"
	MembraneIndoAmericanTomtom MembraneValue = "Indo-American tomtom"
	MembraneJapaneseTomtom     MembraneValue = "Japanese tomtom"
	MembraneMilitaryDrum       MembraneValue = "military drum"
	MembraneSnareDrum          MembraneValue = "snare drum"
	MembraneSnareDrumSnaresOff MembraneValue = "snare drum snares off"
	MembraneTabla              MembraneValue = "tabla"
	MembraneTambourine         MembraneValue = "tambourine"
	MembraneTenorDrum          MembraneValue = "tenor drum"
	MembraneTimbales           MembraneValue = "timbales"
	MembraneTomtom             MembraneValue = "tomtom"
)

// From xsd simple type with enumerate restriction "metal-value"
type MetalValue string

const (
	MetalAgogo            MetalValue = "agogo"
	MetalAlmglocken       MetalValue = "almglocken"
	MetalBell             MetalValue = "bell"
	MetalBellPlate        MetalValue = "bell plate"
	MetalBellTree         MetalValue = "bell tree"
	MetalBrakeDrum        MetalValue = "brake drum"
	MetalCencerro         MetalValue = "cencerro"
	MetalChainRattle      MetalValue = "chain rattle"
	MetalChineseCymbal    MetalValue = "Chinese cymbal"
	MetalCowbell          MetalValue = "cowbell"
	MetalCrashCymbals     MetalValue = "crash cymbals"
	MetalCrotale          MetalValue = "crotale"
	MetalCymbalTongs      MetalValue = "cymbal tongs"
	MetalDomedGong        MetalValue = "domed gong"
	MetalFingerCymbals    MetalValue = "finger cymbals"
	MetalFlexatone        MetalValue = "flexatone"
	MetalGong             MetalValue = "gong"
	MetalHiHat            MetalValue = "hi-hat"
	MetalHighHatCymbals   MetalValue = "high-hat cymbals"
	MetalHandbell         MetalValue = "handbell"
	MetalJawHarp          MetalValue = "jaw harp"
	MetalJingleBells      MetalValue = "jingle bells"
	MetalMusicalSaw       MetalValue = "musical saw"
	MetalShellBells       MetalValue = "shell bells"
	MetalSistrum          MetalValue = "sistrum"
	MetalSizzleCymbal     MetalValue = "sizzle cymbal"
	MetalSleighBells      MetalValue = "sleigh bells"
	MetalSuspendedCymbal  MetalValue = "suspended cymbal"
	MetalTamTam           MetalValue = "tam tam"
	MetalTamTamWithBeater MetalValue = "tam tam with beater"
	MetalTriangle         MetalValue = "triangle"
	MetalVietnameseHat    MetalValue = "Vietnamese hat"
)

// From xsd simple type with enumerate restriction "mode"
type Mode string

const ()

// From xsd simple type with enumerate restriction "mute"
type Mute string

const (
	MuteOn           Mute = "on"
	MuteOff          Mute = "off"
	MuteStraight     Mute = "straight"
	MuteCup          Mute = "cup"
	MuteHarmonNoStem Mute = "harmon-no-stem"
	MuteHarmonStem   Mute = "harmon-stem"
	MuteBucket       Mute = "bucket"
	MutePlunger      Mute = "plunger"
	MuteHat          Mute = "hat"
	MuteSolotone     Mute = "solotone"
	MutePractice     Mute = "practice"
	MuteStopMute     Mute = "stop-mute"
	MuteStopHand     Mute = "stop-hand"
	MuteEcho         Mute = "echo"
	MutePalm         Mute = "palm"
)

// From xsd simple type with enumerate restriction "note-size-type"
type NoteSizeType string

const (
	NoteSizeTypeCue      NoteSizeType = "cue"
	NoteSizeTypeGrace    NoteSizeType = "grace"
	NoteSizeTypeGraceCue NoteSizeType = "grace-cue"
	NoteSizeTypeLarge    NoteSizeType = "large"
)

// From xsd simple type with enumerate restriction "note-type-value"
type NoteTypeValue string

const (
	NoteType1024th  NoteTypeValue = "1024th"
	NoteType512th   NoteTypeValue = "512th"
	NoteType256th   NoteTypeValue = "256th"
	NoteType128th   NoteTypeValue = "128th"
	NoteType64th    NoteTypeValue = "64th"
	NoteType32nd    NoteTypeValue = "32nd"
	NoteType16th    NoteTypeValue = "16th"
	NoteTypeEighth  NoteTypeValue = "eighth"
	NoteTypeQuarter NoteTypeValue = "quarter"
	NoteTypeHalf    NoteTypeValue = "half"
	NoteTypeWhole   NoteTypeValue = "whole"
	NoteTypeBreve   NoteTypeValue = "breve"
	NoteTypeLong    NoteTypeValue = "long"
	NoteTypeMaxima  NoteTypeValue = "maxima"
)

// From xsd simple type with enumerate restriction "notehead-value"
type NoteheadValue string

const (
	NoteheadSlash            NoteheadValue = "slash"
	NoteheadTriangle         NoteheadValue = "triangle"
	NoteheadDiamond          NoteheadValue = "diamond"
	NoteheadSquare           NoteheadValue = "square"
	NoteheadCross            NoteheadValue = "cross"
	NoteheadX                NoteheadValue = "x"
	NoteheadCircleX          NoteheadValue = "circle-x"
	NoteheadInvertedTriangle NoteheadValue = "inverted triangle"
	NoteheadArrowDown        NoteheadValue = "arrow down"
	NoteheadArrowUp          NoteheadValue = "arrow up"
	NoteheadCircled          NoteheadValue = "circled"
	NoteheadSlashed          NoteheadValue = "slashed"
	NoteheadBackSlashed      NoteheadValue = "back slashed"
	NoteheadNormal           NoteheadValue = "normal"
	NoteheadCluster          NoteheadValue = "cluster"
	NoteheadCircleDot        NoteheadValue = "circle dot"
	NoteheadLeftTriangle     NoteheadValue = "left triangle"
	NoteheadRectangle        NoteheadValue = "rectangle"
	NoteheadNone             NoteheadValue = "none"
	NoteheadDo               NoteheadValue = "do"
	NoteheadRe               NoteheadValue = "re"
	NoteheadMi               NoteheadValue = "mi"
	NoteheadFa               NoteheadValue = "fa"
	NoteheadFaUp             NoteheadValue = "fa up"
	NoteheadSo               NoteheadValue = "so"
	NoteheadLa               NoteheadValue = "la"
	NoteheadTi               NoteheadValue = "ti"
	NoteheadOther            NoteheadValue = "other"
)

// From xsd simple type with enumerate restriction "numeral-mode"
type NumeralMode string

const (
	NumeralModeMajor         NumeralMode = "major"
	NumeralModeMinor         NumeralMode = "minor"
	NumeralModeNaturalMinor  NumeralMode = "natural minor"
	NumeralModeMelodicMinor  NumeralMode = "melodic minor"
	NumeralModeHarmonicMinor NumeralMode = "harmonic minor"
)

// From xsd simple type with enumerate restriction "on-off"
type OnOff string

const (
	OnOffOn  OnOff = "on"
	OnOffOff OnOff = "off"
)

// From xsd simple type with enumerate restriction "over-under"
type OverUnder string

const (
	OverUnderOver  OverUnder = "over"
	OverUnderUnder OverUnder = "under"
)

// From xsd simple type with enumerate restriction "pedal-type"
type PedalType string

const (
	PedalTypeStart       PedalType = "start"
	PedalTypeStop        PedalType = "stop"
	PedalTypeSostenuto   PedalType = "sostenuto"
	PedalTypeChange      PedalType = "change"
	PedalTypeContinue    PedalType = "continue"
	PedalTypeDiscontinue PedalType = "discontinue"
	PedalTypeResume      PedalType = "resume"
)

// From xsd simple type with enumerate restriction "pitched-value"
type PitchedValue string

const (
	PitchedCelesta       PitchedValue = "celesta"
	PitchedChimes        PitchedValue = "chimes"
	PitchedGlockenspiel  PitchedValue = "glockenspiel"
	PitchedLithophone    PitchedValue = "lithophone"
	PitchedMallet        PitchedValue = "mallet"
	PitchedMarimba       PitchedValue = "marimba"
	PitchedSteelDrums    PitchedValue = "steel drums"
	PitchedTubaphone     PitchedValue = "tubaphone"
	PitchedTubularChimes PitchedValue = "tubular chimes"
	PitchedVibraphone    PitchedValue = "vibraphone"
	PitchedXylophone     PitchedValue = "xylophone"
)

// From xsd simple type with enumerate restriction "principal-voice-symbol"
type PrincipalVoiceSymbol string

const (
	PrincipalVoiceSymbolHauptstimme PrincipalVoiceSymbol = "Hauptstimme"
	PrincipalVoiceSymbolNebenstimme PrincipalVoiceSymbol = "Nebenstimme"
	PrincipalVoiceSymbolPlain       PrincipalVoiceSymbol = "plain"
	PrincipalVoiceSymbolNone        PrincipalVoiceSymbol = "none"
)

// From xsd simple type with enumerate restriction "right-left-middle"
type RightLeftMiddle string

const (
	RightLeftMiddleRight  RightLeftMiddle = "right"
	RightLeftMiddleLeft   RightLeftMiddle = "left"
	RightLeftMiddleMiddle RightLeftMiddle = "middle"
)

// From xsd simple type with enumerate restriction "semi-pitched"
type SemiPitched string

const (
	SemiPitchedHigh       SemiPitched = "high"
	SemiPitchedMediumHigh SemiPitched = "medium-high"
	SemiPitchedMedium     SemiPitched = "medium"
	SemiPitchedMediumLow  SemiPitched = "medium-low"
	SemiPitchedLow        SemiPitched = "low"
	SemiPitchedVeryLow    SemiPitched = "very-low"
)

// From xsd simple type with enumerate restriction "show-frets"
type ShowFrets string

const (
	ShowFretsNumbers ShowFrets = "numbers"
	ShowFretsLetters ShowFrets = "letters"
)

// From xsd simple type with enumerate restriction "show-tuplet"
type ShowTuplet string

const (
	ShowTupletActual ShowTuplet = "actual"
	ShowTupletBoth   ShowTuplet = "both"
	ShowTupletNone   ShowTuplet = "none"
)

// From xsd simple type with enumerate restriction "staff-divide-symbol"
type StaffDivideSymbol string

const (
	StaffDivideSymbolDown   StaffDivideSymbol = "down"
	StaffDivideSymbolUp     StaffDivideSymbol = "up"
	StaffDivideSymbolUpDown StaffDivideSymbol = "up-down"
)

// From xsd simple type with enumerate restriction "staff-type"
type StaffType string

const (
	StaffTypeOssia     StaffType = "ossia"
	StaffTypeEditorial StaffType = "editorial"
	StaffTypeCue       StaffType = "cue"
	StaffTypeAlternate StaffType = "alternate"
	StaffTypeRegular   StaffType = "regular"
)

// From xsd simple type with enumerate restriction "start-note"
type StartNote string

const (
	StartNoteUpper StartNote = "upper"
	StartNoteMain  StartNote = "main"
	StartNoteBelow StartNote = "below"
)

// From xsd simple type with enumerate restriction "start-stop"
type StartStop string

const (
	StartStopStart StartStop = "start"
	StartStopStop  StartStop = "stop"
)

// From xsd simple type with enumerate restriction "start-stop-continue"
type StartStopContinue string

const (
	StartStopContinueStart    StartStopContinue = "start"
	StartStopContinueStop     StartStopContinue = "stop"
	StartStopContinueContinue StartStopContinue = "continue"
)

// From xsd simple type with enumerate restriction "start-stop-discontinue"
type StartStopDiscontinue string

const (
	StartStopDiscontinueStart       StartStopDiscontinue = "start"
	StartStopDiscontinueStop        StartStopDiscontinue = "stop"
	StartStopDiscontinueDiscontinue StartStopDiscontinue = "discontinue"
)

// From xsd simple type with enumerate restriction "start-stop-single"
type StartStopSingle string

const (
	StartStopSingleStart  StartStopSingle = "start"
	StartStopSingleStop   StartStopSingle = "stop"
	StartStopSingleSingle StartStopSingle = "single"
)

// From xsd simple type with enumerate restriction "stem-value"
type StemValue string

const (
	StemDown   StemValue = "down"
	StemUp     StemValue = "up"
	StemDouble StemValue = "double"
	StemNone   StemValue = "none"
)

// From xsd simple type with enumerate restriction "step"
type Step string

const (
	StepA Step = "A"
	StepB Step = "B"
	StepC Step = "C"
	StepD Step = "D"
	StepE Step = "E"
	StepF Step = "F"
	StepG Step = "G"
)

// From xsd simple type with enumerate restriction "stick-location"
type StickLocation string

const (
	StickLocationCenter     StickLocation = "center"
	StickLocationRim        StickLocation = "rim"
	StickLocationCymbalBell StickLocation = "cymbal bell"
	StickLocationCymbalEdge StickLocation = "cymbal edge"
)

// From xsd simple type with enumerate restriction "stick-material"
type StickMaterial string

const (
	StickMaterialSoft   StickMaterial = "soft"
	StickMaterialMedium StickMaterial = "medium"
	StickMaterialHard   StickMaterial = "hard"
	StickMaterialShaded StickMaterial = "shaded"
	StickMaterialX      StickMaterial = "x"
)

// From xsd simple type with enumerate restriction "stick-type"
type StickType string

const (
	StickTypeBassDrum       StickType = "bass drum"
	StickTypeDoubleBassDrum StickType = "double bass drum"
	StickTypeGlockenspiel   StickType = "glockenspiel"
	StickTypeGum            StickType = "gum"
	StickTypeHammer         StickType = "hammer"
	StickTypeSuperball      StickType = "superball"
	StickTypeTimpani        StickType = "timpani"
	StickTypeWound          StickType = "wound"
	StickTypeXylophone      StickType = "xylophone"
	StickTypeYarn           StickType = "yarn"
)

// From xsd simple type with enumerate restriction "syllabic"
type Syllabic string

const (
	SyllabicSingle Syllabic = "single"
	SyllabicBegin  Syllabic = "begin"
	SyllabicEnd    Syllabic = "end"
	SyllabicMiddle Syllabic = "middle"
)

// From xsd simple type with enumerate restriction "symbol-size"
type SymbolSize string

const (
	SymbolSizeFull     SymbolSize = "full"
	SymbolSizeCue      SymbolSize = "cue"
	SymbolSizeGraceCue SymbolSize = "grace-cue"
	SymbolSizeLarge    SymbolSize = "large"
)

// From xsd simple type with enumerate restriction "sync-type"
type SyncType string

const (
	SyncTypeNone        SyncType = "none"
	SyncTypeTempo       SyncType = "tempo"
	SyncTypeMostlyTempo SyncType = "mostly-tempo"
	SyncTypeMostlyEvent SyncType = "mostly-event"
	SyncTypeEvent       SyncType = "event"
	SyncTypeAlwaysEvent SyncType = "always-event"
)

// From xsd simple type with enumerate restriction "system-relation-number"
type SystemRelationNumber string

const (
	SystemRelationNumberOnlyTop    SystemRelationNumber = "only-top"
	SystemRelationNumberOnlyBottom SystemRelationNumber = "only-bottom"
	SystemRelationNumberAlsoTop    SystemRelationNumber = "also-top"
	SystemRelationNumberAlsoBottom SystemRelationNumber = "also-bottom"
	SystemRelationNumberNone       SystemRelationNumber = "none"
)

// From xsd simple type with enumerate restriction "tap-hand"
type TapHand string

const (
	TapHandLeft  TapHand = "left"
	TapHandRight TapHand = "right"
)

// From xsd simple type with enumerate restriction "text-direction"
type TextDirection string

const (
	TextDirectionLtr TextDirection = "ltr"
	TextDirectionRtl TextDirection = "rtl"
	TextDirectionLro TextDirection = "lro"
	TextDirectionRlo TextDirection = "rlo"
)

// From xsd simple type with enumerate restriction "tied-type"
type TiedType string

const (
	TiedTypeStart    TiedType = "start"
	TiedTypeStop     TiedType = "stop"
	TiedTypeContinue TiedType = "continue"
	TiedTypeLetRing  TiedType = "let-ring"
)

// From xsd simple type with enumerate restriction "time-only"
type TimeOnly string

const ()

// From xsd simple type with enumerate restriction "time-relation"
type TimeRelation string

const (
	TimeRelationParentheses TimeRelation = "parentheses"
	TimeRelationBracket     TimeRelation = "bracket"
	TimeRelationEquals      TimeRelation = "equals"
	TimeRelationSlash       TimeRelation = "slash"
	TimeRelationSpace       TimeRelation = "space"
	TimeRelationHyphen      TimeRelation = "hyphen"
)

// From xsd simple type with enumerate restriction "time-separator"
type TimeSeparator string

const (
	TimeSeparatorNone       TimeSeparator = "none"
	TimeSeparatorHorizontal TimeSeparator = "horizontal"
	TimeSeparatorDiagonal   TimeSeparator = "diagonal"
	TimeSeparatorVertical   TimeSeparator = "vertical"
	TimeSeparatorAdjacent   TimeSeparator = "adjacent"
)

// From xsd simple type with enumerate restriction "time-symbol"
type TimeSymbol string

const (
	TimeSymbolCommon       TimeSymbol = "common"
	TimeSymbolCut          TimeSymbol = "cut"
	TimeSymbolSingleNumber TimeSymbol = "single-number"
	TimeSymbolNote         TimeSymbol = "note"
	TimeSymbolDottedNote   TimeSymbol = "dotted-note"
	TimeSymbolNormal       TimeSymbol = "normal"
)

// From xsd simple type with enumerate restriction "tip-direction"
type TipDirection string

const (
	TipDirectionUp        TipDirection = "up"
	TipDirectionDown      TipDirection = "down"
	TipDirectionLeft      TipDirection = "left"
	TipDirectionRight     TipDirection = "right"
	TipDirectionNorthwest TipDirection = "northwest"
	TipDirectionNortheast TipDirection = "northeast"
	TipDirectionSoutheast TipDirection = "southeast"
	TipDirectionSouthwest TipDirection = "southwest"
)

// From xsd simple type with enumerate restriction "top-bottom"
type TopBottom string

const (
	TopBottomTop    TopBottom = "top"
	TopBottomBottom TopBottom = "bottom"
)

// From xsd simple type with enumerate restriction "tremolo-type"
type TremoloType string

const (
	TremoloTypeStart      TremoloType = "start"
	TremoloTypeStop       TremoloType = "stop"
	TremoloTypeSingle     TremoloType = "single"
	TremoloTypeUnmeasured TremoloType = "unmeasured"
)

// From xsd simple type with enumerate restriction "trill-step"
type TrillStep string

const (
	TrillStepWhole  TrillStep = "whole"
	TrillStepHalf   TrillStep = "half"
	TrillStepUnison TrillStep = "unison"
)

// From xsd simple type with enumerate restriction "two-note-turn"
type TwoNoteTurn string

const (
	TwoNoteTurnWhole TwoNoteTurn = "whole"
	TwoNoteTurnHalf  TwoNoteTurn = "half"
	TwoNoteTurnNone  TwoNoteTurn = "none"
)

// From xsd simple type with enumerate restriction "up-down"
type UpDown string

const (
	UpDownUp   UpDown = "up"
	UpDownDown UpDown = "down"
)

// From xsd simple type with enumerate restriction "up-down-stop-continue"
type UpDownStopContinue string

const (
	UpDownStopContinueUp       UpDownStopContinue = "up"
	UpDownStopContinueDown     UpDownStopContinue = "down"
	UpDownStopContinueStop     UpDownStopContinue = "stop"
	UpDownStopContinueContinue UpDownStopContinue = "continue"
)

// From xsd simple type with enumerate restriction "upright-inverted"
type UprightInverted string

const (
	UprightInvertedUpright  UprightInverted = "upright"
	UprightInvertedInverted UprightInverted = "inverted"
)

// From xsd simple type with enumerate restriction "valign"
type Valign string

const (
	ValignTop      Valign = "top"
	ValignMiddle   Valign = "middle"
	ValignBottom   Valign = "bottom"
	ValignBaseline Valign = "baseline"
)

// From xsd simple type with enumerate restriction "valign-image"
type ValignImage string

const (
	ValignImageTop    ValignImage = "top"
	ValignImageMiddle ValignImage = "middle"
	ValignImageBottom ValignImage = "bottom"
)

// From xsd simple type with enumerate restriction "wedge-type"
type WedgeType string

const (
	WedgeTypeCrescendo  WedgeType = "crescendo"
	WedgeTypeDiminuendo WedgeType = "diminuendo"
	WedgeTypeStop       WedgeType = "stop"
	WedgeTypeContinue   WedgeType = "continue"
)

// From xsd simple type with enumerate restriction "winged"
type Winged string

const (
	WingedNone           Winged = "none"
	WingedStraight       Winged = "straight"
	WingedCurved         Winged = "curved"
	WingedDoubleStraight Winged = "double-straight"
	WingedDoubleCurved   Winged = "double-curved"
)

// From xsd simple type with enumerate restriction "wood-value"
type WoodValue string

const (
	WoodBambooScraper       WoodValue = "bamboo scraper"
	WoodBoardClapper        WoodValue = "board clapper"
	WoodCabasa              WoodValue = "cabasa"
	WoodCastanets           WoodValue = "castanets"
	WoodCastanetsWithHandle WoodValue = "castanets with handle"
	WoodClaves              WoodValue = "claves"
	WoodFootballRattle      WoodValue = "football rattle"
	WoodGuiro               WoodValue = "guiro"
	WoodLogDrum             WoodValue = "log drum"
	WoodMaraca              WoodValue = "maraca"
	WoodMaracas             WoodValue = "maracas"
	WoodQuijada             WoodValue = "quijada"
	WoodRainstick           WoodValue = "rainstick"
	WoodRatchet             WoodValue = "ratchet"
	WoodRecoReco            WoodValue = "reco-reco"
	WoodSandpaperBlocks     WoodValue = "sandpaper blocks"
	WoodSlitDrum            WoodValue = "slit drum"
	WoodTempleBlock         WoodValue = "temple block"
	WoodVibraslap           WoodValue = "vibraslap"
	WoodWhip                WoodValue = "whip"
	WoodWoodBlock           WoodValue = "wood block"
)

// From xsd simple type with enumerate restriction "yes-no"
type YesNo string

const (
	Yes YesNo = "yes"
	No  YesNo = "no"
)

// insertion point for gongstructs declarations

// AttributeGroupBendSound UnNamed source named attribute group "bend-sound"
type AttributeGroupBendSound struct {
	Accelerate YesNo  `xml:"accelerate,attr,omitempty"`
	Beats      string `xml:"beats,attr,omitempty"`
	FirstBeat  string `xml:"first-beat,attr,omitempty"`
	LastBeat   string `xml:"last-beat,attr,omitempty"`
}

// AttributeGroupBezier UnNamed source named attribute group "bezier"
type AttributeGroupBezier struct {
	BezierX       string `xml:"bezier-x,attr,omitempty"`
	BezierY       string `xml:"bezier-y,attr,omitempty"`
	BezierX2      string `xml:"bezier-x2,attr,omitempty"`
	BezierY2      string `xml:"bezier-y2,attr,omitempty"`
	BezierOffset  string `xml:"bezier-offset,attr,omitempty"`
	BezierOffset2 string `xml:"bezier-offset2,attr,omitempty"`
}

// AttributeGroupColor UnNamed source named attribute group "color"
type AttributeGroupColor struct {
	Color string `xml:"color,attr,omitempty"`
}

// AttributeGroupDashedFormatting UnNamed source named attribute group "dashed-formatting"
type AttributeGroupDashedFormatting struct {
	DashLength  string `xml:"dash-length,attr,omitempty"`
	SpaceLength string `xml:"space-length,attr,omitempty"`
}

// AttributeGroupDirective UnNamed source named attribute group "directive"
type AttributeGroupDirective struct {
	Directive YesNo `xml:"directive,attr,omitempty"`
}

// AttributeGroupDocumentAttributes UnNamed source named attribute group "document-attributes"
type AttributeGroupDocumentAttributes struct {
	Version string `xml:"version,attr,omitempty"`
}

// AttributeGroupEnclosure UnNamed source named attribute group "enclosure"
type AttributeGroupEnclosure struct {
	Enclosure string `xml:"enclosure,attr,omitempty"`
}

// AttributeGroupFont UnNamed source named attribute group "font"
type AttributeGroupFont struct {
	FontFamily string `xml:"font-family,attr,omitempty"`
	FontStyle  string `xml:"font-style,attr,omitempty"`
	FontSize   string `xml:"font-size,attr,omitempty"`
	FontWeight string `xml:"font-weight,attr,omitempty"`
}

// AttributeGroupHalign UnNamed source named attribute group "halign"
type AttributeGroupHalign struct {
	Halign string `xml:"halign,attr,omitempty"`
}

// AttributeGroupJustify UnNamed source named attribute group "justify"
type AttributeGroupJustify struct {
	Justify LeftCenterRight `xml:"justify,attr,omitempty"`
}

// AttributeGroupLetterSpacing UnNamed source named attribute group "letter-spacing"
type AttributeGroupLetterSpacing struct {
	LetterSpacing string `xml:"letter-spacing,attr,omitempty"`
}

// AttributeGroupLevelDisplay UnNamed source named attribute group "level-display"
type AttributeGroupLevelDisplay struct {
	Parentheses YesNo      `xml:"parentheses,attr,omitempty"`
	Bracket     YesNo      `xml:"bracket,attr,omitempty"`
	Size        SymbolSize `xml:"size,attr,omitempty"`
}

// AttributeGroupLineHeight UnNamed source named attribute group "line-height"
type AttributeGroupLineHeight struct {
	LineHeight string `xml:"line-height,attr,omitempty"`
}

// AttributeGroupLineLength UnNamed source named attribute group "line-length"
type AttributeGroupLineLength struct {
	LineLength string `xml:"line-length,attr,omitempty"`
}

// AttributeGroupLineShape UnNamed source named attribute group "line-shape"
type AttributeGroupLineShape struct {
	LineShape string `xml:"line-shape,attr,omitempty"`
}

// AttributeGroupLineType UnNamed source named attribute group "line-type"
type AttributeGroupLineType struct {
	LineType string `xml:"line-type,attr,omitempty"`
}

// AttributeGroupOptionalUniqueId UnNamed source named attribute group "optional-unique-id"
type AttributeGroupOptionalUniqueId struct {
	Id string `xml:"id,attr,omitempty"`
}

// AttributeGroupOrientation UnNamed source named attribute group "orientation"
type AttributeGroupOrientation struct {
	Orientation string `xml:"orientation,attr,omitempty"`
}

// AttributeGroupPlacement UnNamed source named attribute group "placement"
type AttributeGroupPlacement struct {
	Placement string `xml:"placement,attr,omitempty"`
}

// AttributeGroupPosition UnNamed source named attribute group "position"
type AttributeGroupPosition struct {
	DefaultX  string `xml:"default-x,attr,omitempty"`
	DefaultY  string `xml:"default-y,attr,omitempty"`
	RelativeX string `xml:"relative-x,attr,omitempty"`
	RelativeY string `xml:"relative-y,attr,omitempty"`
}

// AttributeGroupPrintObject UnNamed source named attribute group "print-object"
type AttributeGroupPrintObject struct {
	PrintObject YesNo `xml:"print-object,attr,omitempty"`
}

// AttributeGroupPrintSpacing UnNamed source named attribute group "print-spacing"
type AttributeGroupPrintSpacing struct {
	PrintSpacing YesNo `xml:"print-spacing,attr,omitempty"`
}

// AttributeGroupPrintStyle UnNamed source named attribute group "print-style"
type AttributeGroupPrintStyle struct {
	AttributeGroupPosition
	AttributeGroupFont
	AttributeGroupColor
}

// AttributeGroupPrintStyleAlign UnNamed source named attribute group "print-style-align"
type AttributeGroupPrintStyleAlign struct {
	AttributeGroupPrintStyle
	AttributeGroupHalign
	AttributeGroupValign
}

// AttributeGroupPrintout UnNamed source named attribute group "printout"
type AttributeGroupPrintout struct {
	PrintDot   YesNo `xml:"print-dot,attr,omitempty"`
	PrintLyric YesNo `xml:"print-lyric,attr,omitempty"`
	AttributeGroupPrintObject
	AttributeGroupPrintSpacing
}

// AttributeGroupSmufl UnNamed source named attribute group "smufl"
type AttributeGroupSmufl struct {
	Smufl string `xml:"smufl,attr,omitempty"`
}

// AttributeGroupSystemRelation UnNamed source named attribute group "system-relation"
type AttributeGroupSystemRelation struct {
	System SystemRelationNumber `xml:"system,attr,omitempty"`
}

// AttributeGroupSymbolFormatting UnNamed source named attribute group "symbol-formatting"
type AttributeGroupSymbolFormatting struct {
	AttributeGroupJustify
	AttributeGroupPrintStyleAlign
	AttributeGroupTextDecoration
	AttributeGroupTextRotation
	AttributeGroupLetterSpacing
	AttributeGroupLineHeight
	AttributeGroupTextDirection
	AttributeGroupEnclosure
}

// AttributeGroupTextDecoration UnNamed source named attribute group "text-decoration"
type AttributeGroupTextDecoration struct {
	Underline   int `xml:"underline,attr,omitempty"`
	Overline    int `xml:"overline,attr,omitempty"`
	LineThrough int `xml:"line-through,attr,omitempty"`
}

// AttributeGroupTextDirection UnNamed source named attribute group "text-direction"
type AttributeGroupTextDirection struct {
	Dir string `xml:"dir,attr,omitempty"`
}

// AttributeGroupTextFormatting UnNamed source named attribute group "text-formatting"
type AttributeGroupTextFormatting struct {
	Lang  string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
	Space string `xml:"http://www.w3.org/XML/1998/namespace space,attr,omitempty"`
	AttributeGroupJustify
	AttributeGroupPrintStyleAlign
	AttributeGroupTextDecoration
	AttributeGroupTextRotation
	AttributeGroupLetterSpacing
	AttributeGroupLineHeight
	AttributeGroupTextDirection
	AttributeGroupEnclosure
}

// AttributeGroupTextRotation UnNamed source named attribute group "text-rotation"
type AttributeGroupTextRotation struct {
	Rotation string `xml:"rotation,attr,omitempty"`
}

// AttributeGroupTrillSound UnNamed source named attribute group "trill-sound"
type AttributeGroupTrillSound struct {
	StartNote   string `xml:"start-note,attr,omitempty"`
	TrillStep   string `xml:"trill-step,attr,omitempty"`
	TwoNoteTurn string `xml:"two-note-turn,attr,omitempty"`
	Accelerate  YesNo  `xml:"accelerate,attr,omitempty"`
	Beats       string `xml:"beats,attr,omitempty"`
	SecondBeat  string `xml:"second-beat,attr,omitempty"`
	LastBeat    string `xml:"last-beat,attr,omitempty"`
}

// AttributeGroupValign UnNamed source named attribute group "valign"
type AttributeGroupValign struct {
	Valign string `xml:"valign,attr,omitempty"`
}

// AttributeGroupValignImage UnNamed source named attribute group "valign-image"
type AttributeGroupValignImage struct {
	Valign string `xml:"valign,attr,omitempty"`
}

// AttributeGroupXPosition UnNamed source named attribute group "x-position"
type AttributeGroupXPosition struct {
	DefaultX  string `xml:"default-x,attr,omitempty"`
	DefaultY  string `xml:"default-y,attr,omitempty"`
	RelativeX string `xml:"relative-x,attr,omitempty"`
	RelativeY string `xml:"relative-y,attr,omitempty"`
}

// AttributeGroupYPosition UnNamed source named attribute group "y-position"
type AttributeGroupYPosition struct {
	DefaultX  string `xml:"default-x,attr,omitempty"`
	DefaultY  string `xml:"default-y,attr,omitempty"`
	RelativeX string `xml:"relative-x,attr,omitempty"`
	RelativeY string `xml:"relative-y,attr,omitempty"`
}

// AttributeGroupImageAttributes UnNamed source named attribute group "image-attributes"
type AttributeGroupImageAttributes struct {
	Source string `xml:"source,attr,omitempty"`
	Type   string `xml:"type,attr,omitempty"`
	Height string `xml:"height,attr,omitempty"`
	Width  string `xml:"width,attr,omitempty"`
	AttributeGroupPosition
	AttributeGroupHalign
	AttributeGroupValignImage
}

// AttributeGroupPrintAttributes UnNamed source named attribute group "print-attributes"
type AttributeGroupPrintAttributes struct {
	StaffSpacing string `xml:"staff-spacing,attr,omitempty"`
	NewSystem    YesNo  `xml:"new-system,attr,omitempty"`
	NewPage      YesNo  `xml:"new-page,attr,omitempty"`
	BlankPage    int    `xml:"blank-page,attr,omitempty"`
	PageNumber   string `xml:"page-number,attr,omitempty"`
}

// AttributeGroupElementPosition UnNamed source named attribute group "element-position"
type AttributeGroupElementPosition struct {
	Element  string `xml:"element,attr,omitempty"`
	Position int    `xml:"position,attr,omitempty"`
}

// AttributeGroupLinkAttributes UnNamed source named attribute group "link-attributes"
type AttributeGroupLinkAttributes struct {
	Href    string `xml:"http://www.w3.org/1999/xlink href,attr,omitempty"`
	Type    string `xml:"http://www.w3.org/1999/xlink type,attr,omitempty"`
	Role    string `xml:"http://www.w3.org/1999/xlink role,attr,omitempty"`
	Title   string `xml:"http://www.w3.org/1999/xlink title,attr,omitempty"`
	Show    string `xml:"http://www.w3.org/1999/xlink show,attr,omitempty"`
	Actuate string `xml:"http://www.w3.org/1999/xlink actuate,attr,omitempty"`
}

// AttributeGroupGroupNameText UnNamed source named attribute group "group-name-text"
type AttributeGroupGroupNameText struct {
	AttributeGroupPrintStyle
	AttributeGroupJustify
}

// AttributeGroupMeasureAttributes UnNamed source named attribute group "measure-attributes"
type AttributeGroupMeasureAttributes struct {
	Number         string `xml:"number,attr,omitempty"`
	Text           string `xml:"text,attr,omitempty"`
	Implicit       YesNo  `xml:"implicit,attr,omitempty"`
	NonControlling YesNo  `xml:"non-controlling,attr,omitempty"`
	Width          string `xml:"width,attr,omitempty"`
	AttributeGroupOptionalUniqueId
}

// AttributeGroupPartAttributes UnNamed source named attribute group "part-attributes"
type AttributeGroupPartAttributes struct {
	Id string `xml:"id,attr,omitempty"`
}

// AttributeGroupPartNameText UnNamed source named attribute group "part-name-text"
type AttributeGroupPartNameText struct {
	AttributeGroupPrintStyle
	AttributeGroupPrintObject
	AttributeGroupJustify
}

// AccidentalText Named source named complex type "accidental-text"
// The accidental-text type represents an element with an accidental
// value and text-formatting attributes.
type AccidentalText struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`
	AttributeGroupTextFormatting

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText AccidentalValue `xml:",chardata"`
}

// Coda Named source named complex type "coda"
// The coda type is the visual indicator of a coda sign. The exact glyph
// can be specified with the smufl attribute. A sound element is also needed to guide
// playback applications reliably.
type Coda struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`
	AttributeGroupPrintStyleAlign
	AttributeGroupOptionalUniqueId
}

// Dynamics Named source named complex type "dynamics"
// Dynamics can be associated either with a note or a general musical
// direction. To avoid inconsistencies between and amongst the letter abbreviations for
// dynamics (what is sf vs. sfz, standing alone or with a trailing dynamic that is not
// always piano), we use the actual letters as the names of these dynamic elements. The
// other-dynamics element allows other dynamic marks that are not covered here.
// Dynamics elements may also be combined to create marks not covered by a single
// element, such as sfmp. These letter dynamic symbols are separated from crescendo,
// decrescendo, and wedge indications. Dynamic representation is inconsistent in
// scores. Many things are assumed by the composer and left out, such as returns to
// original dynamics. The MusicXML format captures what is in the score, but does not
// try to be optimal for analysis or synthesis of dynamics. The placement attribute is
// used when the dynamics are associated with a note. It is ignored when the dynamics
// are associated with a direction. In that case the direction element's placement
// attribute is used instead.
type Dynamics struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyleAlign
	AttributeGroupPlacement
	AttributeGroupTextDecoration
	AttributeGroupEnclosure
	AttributeGroupOptionalUniqueId
	P             string       `xml:"p,omitempty"`
	Pp            string       `xml:"pp,omitempty"`
	Ppp           string       `xml:"ppp,omitempty"`
	Pppp          string       `xml:"pppp,omitempty"`
	Ppppp         string       `xml:"ppppp,omitempty"`
	Pppppp        string       `xml:"pppppp,omitempty"`
	F             string       `xml:"f,omitempty"`
	Ff            string       `xml:"ff,omitempty"`
	Fff           string       `xml:"fff,omitempty"`
	Ffff          string       `xml:"ffff,omitempty"`
	Fffff         string       `xml:"fffff,omitempty"`
	Ffffff        string       `xml:"ffffff,omitempty"`
	Mp            string       `xml:"mp,omitempty"`
	Mf            string       `xml:"mf,omitempty"`
	Sf            string       `xml:"sf,omitempty"`
	Sfp           string       `xml:"sfp,omitempty"`
	Sfpp          string       `xml:"sfpp,omitempty"`
	Fp            string       `xml:"fp,omitempty"`
	Rf            string       `xml:"rf,omitempty"`
	Rfz           string       `xml:"rfz,omitempty"`
	Sfz           string       `xml:"sfz,omitempty"`
	Sffz          string       `xml:"sffz,omitempty"`
	Fz            string       `xml:"fz,omitempty"`
	N             string       `xml:"n,omitempty"`
	Pf            string       `xml:"pf,omitempty"`
	Sfzp          string       `xml:"sfzp,omitempty"`
	OtherDynamics []*OtherText `xml:"other-dynamics,omitempty"`
	InnerXML      string       `xml:",innerxml"`
}

// Empty Named source named complex type "empty"
// The empty type represents an empty element with no attributes.
type Empty struct {
	Name string `xml:"-"`
}

// EmptyPlacement Named source named complex type "empty-placement"
// The empty-placement type represents an empty element with print-style
// and placement attributes.
type EmptyPlacement struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement
}

// EmptyPlacementSmufl Named source named complex type "empty-placement-smufl"
// The empty-placement-smufl type represents an empty element with
// print-style, placement, and smufl attributes.
type EmptyPlacementSmufl struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement
	AttributeGroupSmufl
}

// EmptyPrintStyle Named source named complex type "empty-print-style"
// The empty-print-style type represents an empty element with
// print-style attribute group.
type EmptyPrintStyle struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
}

// EmptyPrintStyleAlign Named source named complex type "empty-print-style-align"
// The empty-print-style-align type represents an empty element with
// print-style-align attribute group.
type EmptyPrintStyleAlign struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyleAlign
}

// EmptyPrintStyleAlignId Named source named complex type "empty-print-style-align-id"
// The empty-print-style-align-id type represents an empty element with
// print-style-align and optional-unique-id attribute groups.
type EmptyPrintStyleAlignId struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyleAlign
	AttributeGroupOptionalUniqueId
}

// EmptyPrintObjectStyleAlign Named source named complex type "empty-print-object-style-align"
// The empty-print-style-align-object type represents an empty element
// with print-object and print-style-align attribute groups.
type EmptyPrintObjectStyleAlign struct {
	Name string `xml:"-"`
	AttributeGroupPrintObject
	AttributeGroupPrintStyleAlign
}

// EmptyTrillSound Named source named complex type "empty-trill-sound"
// The empty-trill-sound type represents an empty element with
// print-style, placement, and trill-sound attributes.
type EmptyTrillSound struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement
	AttributeGroupTrillSound
}

// HorizontalTurn Named source named complex type "horizontal-turn"
// The horizontal-turn type represents turn elements that are horizontal
// rather than vertical. These are empty elements with print-style, placement,
// trill-sound, and slash attributes. If the slash attribute is yes, then a vertical
// line is used to slash the turn. It is no if not specified.
type HorizontalTurn struct {
	Name  string `xml:"-"`
	Slash YesNo  `xml:"slash,attr,omitempty"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement
	AttributeGroupTrillSound
}

// Fermata Named source named complex type "fermata"
// The fermata text content represents the shape of the fermata sign. An
// empty fermata element represents a normal fermata. The fermata type is upright if
// not specified.
type Fermata struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`
	AttributeGroupPrintStyle
	AttributeGroupOptionalUniqueId

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Fingering Named source named complex type "fingering"
// Fingering is typically indicated 1,2,3,4,5. Multiple fingerings may be
// given, typically to substitute fingerings in the middle of a note. The substitution
// and alternate values are "no" if the attribute is not present. For guitar and other
// fretted instruments, the fingering element represents the fretting finger; the pluck
// element represents the plucking finger.
type Fingering struct {
	Name         string `xml:"-"`
	Substitution YesNo  `xml:"substitution,attr,omitempty"`
	Alternate    YesNo  `xml:"alternate,attr,omitempty"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// FormattedSymbol Named source named complex type "formatted-symbol"
// The formatted-symbol type represents a SMuFL musical symbol element
// with formatting attributes.
type FormattedSymbol struct {
	Name string `xml:"-"`
	AttributeGroupSymbolFormatting

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// FormattedSymbolId Named source named complex type "formatted-symbol-id"
// The formatted-symbol-id type represents a SMuFL musical symbol element
// with formatting and id attributes.
type FormattedSymbolId struct {
	Name string `xml:"-"`
	AttributeGroupSymbolFormatting
	AttributeGroupOptionalUniqueId

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// FormattedText Named source named complex type "formatted-text"
// The formatted-text type represents a text element with text-formatting
// attributes.
type FormattedText struct {
	Name string `xml:"-"`
	AttributeGroupTextFormatting

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// FormattedTextId Named source named complex type "formatted-text-id"
// The formatted-text-id type represents a text element with
// text-formatting and id attributes.
type FormattedTextId struct {
	Name string `xml:"-"`
	AttributeGroupTextFormatting
	AttributeGroupOptionalUniqueId

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Fret Named source named complex type "fret"
// The fret element is used with tablature notation and chord diagrams.
// Fret numbers start with 0 for an open string and 1 for the first fret.
type Fret struct {
	Name string `xml:"-"`
	AttributeGroupFont
	AttributeGroupColor

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText int `xml:",chardata"`
}

// Level Named source named complex type "level"
// The level type is used to specify editorial information for different
// MusicXML elements. The content contains identifying and/or descriptive text about
// the editorial status of the parent element. If the reference attribute is yes, this
// indicates editorial information that is for display only and should not affect
// playback. For instance, a modern edition of older music may set reference="yes" on
// the attributes containing the music's original clef, key, and time signature. It is
// no if not specified. The type attribute indicates whether the editorial information
// applies to the start of a series of symbols, the end of a series of symbols, or a
// single symbol. It is single if not specified for compatibility with earlier MusicXML
// versions.
type Level struct {
	Name      string          `xml:"-"`
	Reference YesNo           `xml:"reference,attr,omitempty"`
	Type      StartStopSingle `xml:"type,attr,omitempty"`
	AttributeGroupLevelDisplay

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// MidiDevice Named source named complex type "midi-device"
// The midi-device type corresponds to the DeviceName meta event in
// Standard MIDI Files. The optional port attribute is a number from 1 to 16 that can
// be used with the unofficial MIDI 1.0 port (or cable) meta event. Unlike the
// DeviceName meta event, there can be multiple midi-device elements per MusicXML part.
// The optional id attribute refers to the score-instrument assigned to this device. If
// missing, the device assignment affects all score-instrument elements in the
// score-part.
type MidiDevice struct {
	Name string `xml:"-"`
	Port int    `xml:"port,attr,omitempty"`
	Id   string `xml:"id,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// MidiInstrument Named source named complex type "midi-instrument"
// The midi-instrument type defines MIDI 1.0 instrument playback. The
// midi-instrument element can be a part of either the score-instrument element at the
// start of a part, or the sound element within a part. The id attribute refers to the
// score-instrument affected by the change.
type MidiInstrument struct {
	Name          string `xml:"-"`
	Id            string `xml:"id,attr,omitempty"`
	MidiChannel   int    `xml:"midi-channel,omitempty"`
	MidiName      string `xml:"midi-name,omitempty"`
	MidiBank      int    `xml:"midi-bank,omitempty"`
	MidiProgram   int    `xml:"midi-program,omitempty"`
	MidiUnpitched int    `xml:"midi-unpitched,omitempty"`
	Volume        string `xml:"volume,omitempty"`
	Pan           string `xml:"pan,omitempty"`
	Elevation     string `xml:"elevation,omitempty"`
}

// NameDisplay Named source named complex type "name-display"
// The name-display type is used for exact formatting of multi-font text
// in part and group names to the left of the system. The print-object attribute can be
// used to determine what, if anything, is printed at the start of each system.
// Enclosure for the display-text element is none by default. Language for the
// display-text element is Italian ("it") by default.
type NameDisplay struct {
	Name string `xml:"-"`
	AttributeGroupPrintObject
	DisplayText    []*FormattedText  `xml:"display-text,omitempty"`
	AccidentalText []*AccidentalText `xml:"accidental-text,omitempty"`
}

// OtherPlay Named source named complex type "other-play"
// The other-play element represents other types of playback. The
// required type attribute indicates the type of playback to which the element content
// applies.
type OtherPlay struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Play Named source named complex type "play"
// The play type specifies playback techniques to be used in conjunction
// with the instrument-sound element. When used as part of a sound element, it applies
// to all notes going forward in score order. In multi-instrument parts, the affected
// instrument should be specified using the id attribute. When used as part of a note
// element, it applies to the current note only.
type Play struct {
	Name        string       `xml:"-"`
	Id          string       `xml:"id,attr,omitempty"`
	Ipa         string       `xml:"ipa,omitempty"`
	Mute        string       `xml:"mute,omitempty"`
	SemiPitched string       `xml:"semi-pitched,omitempty"`
	OtherPlay   []*OtherPlay `xml:"other-play,omitempty"`
}

// Segno Named source named complex type "segno"
// The segno type is the visual indicator of a segno sign. The exact
// glyph can be specified with the smufl attribute. A sound element is also needed to
// guide playback applications reliably.
type Segno struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`
	AttributeGroupPrintStyleAlign
	AttributeGroupOptionalUniqueId
}

// StringType Named source named complex type "string-type"
// The string type is used with tablature notation, regular notation
// (where it is often circled), and chord diagrams. String numbers start with 1 for the
// highest pitched full-length string.
type StringType struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText int `xml:",chardata"`
}

// TypedText Named source named complex type "typed-text"
// The typed-text type represents a text element with a type attribute.
type TypedText struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// WavyLine Named source named complex type "wavy-line"
// Wavy lines are one way to indicate trills and vibrato. When used with
// a barline element, they should always have type="continue" set. The smufl attribute
// specifies a particular wavy line glyph from the SMuFL Multi-segment lines range.
type WavyLine struct {
	Name   string            `xml:"-"`
	Type   StartStopContinue `xml:"type,attr,omitempty"`
	Number int               `xml:"number,attr,omitempty"`
	Smufl  string            `xml:"smufl,attr,omitempty"`
	AttributeGroupPosition
	AttributeGroupPlacement
	AttributeGroupColor
	AttributeGroupTrillSound
}

// Attributes Named source named complex type "attributes"
// The attributes element contains musical information that typically
// changes on measure boundaries. This includes key and time signatures, clefs,
// transpositions, and staving. When attributes are changed mid-measure, it affects the
// music in score order, not in MusicXML document order.
type Attributes struct {
	Name string `xml:"-"`
	GroupEditorial
	Divisions    int             `xml:"divisions,omitempty"`
	Key          []*Key          `xml:"key,omitempty"`
	Time         []*Time         `xml:"time,omitempty"`
	Staves       int             `xml:"staves,omitempty"`
	PartSymbol   *PartSymbol     `xml:"part-symbol,omitempty"`
	Instruments  int             `xml:"instruments,omitempty"`
	Clef         []*Clef         `xml:"clef,omitempty"`
	StaffDetails []*StaffDetails `xml:"staff-details,omitempty"`
	Transpose    []*Transpose    `xml:"transpose,omitempty"`
	ForPart      []*ForPart      `xml:"for-part,omitempty"`
	Directive    []*ADirective   `xml:"directive,omitempty"`
	MeasureStyle []*MeasureStyle `xml:"measure-style,omitempty"`
}

// ADirective Named source within outer element "directive"
type ADirective struct {
	Name string `xml:"-"`
	Lang string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
	AttributeGroupPrintStyle

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// BeatRepeat Named source named complex type "beat-repeat"
// The beat-repeat type is used to indicate that a single beat (but
// possibly many notes) is repeated. The slashes attribute specifies the number of
// slashes to use in the symbol. The use-dots attribute indicates whether or not to use
// dots as well (for instance, with mixed rhythm patterns). The value for slashes is 1
// and the value for use-dots is no if not specified. The stop type indicates the first
// beat where the repeats are no longer displayed. Both the start and stop of the beat
// being repeated should be specified unless the repeats are displayed through the end
// of the part. The beat-repeat element specifies a notation style for repetitions. The
// actual music being repeated needs to be repeated within the MusicXML file. This
// element specifies the notation that indicates the repeat.
type BeatRepeat struct {
	Name    string    `xml:"-"`
	Type    StartStop `xml:"type,attr,omitempty"`
	Slashes int       `xml:"slashes,attr,omitempty"`
	UseDots YesNo     `xml:"use-dots,attr,omitempty"`
	GroupSlash
}

// Cancel Named source named complex type "cancel"
// A cancel element indicates that the old key signature should be
// cancelled before the new one appears. This will always happen when changing to C
// major or A minor and need not be specified then. The cancel value matches the fifths
// value of the cancelled key signature (e.g., a cancel of -2 will provide an explicit
// cancellation for changing from B flat major to F major). The optional location
// attribute indicates where the cancellation appears relative to the new key
// signature.
type Cancel struct {
	Name     string `xml:"-"`
	Location string `xml:"location,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText int `xml:",chardata"`
}

// Clef Named source named complex type "clef"
// Clefs are represented by a combination of sign, line, and
// clef-octave-change elements. The optional number attribute refers to staff numbers
// within the part. A value of 1 is assumed if not present. Sometimes clefs are added
// to the staff in non-standard line positions, either to indicate cue passages, or
// when there are multiple clefs present simultaneously on one staff. In this
// situation, the additional attribute is set to "yes" and the line value is ignored.
// The size attribute is used for clefs where the additional attribute is "yes". It is
// typically used to indicate cue clefs. Sometimes clefs at the start of a measure need
// to appear after the barline rather than before, as for cues or for use after a
// repeated section. The after-barline attribute is set to "yes" in this situation. The
// attribute is ignored for mid-measure clefs. Clefs appear at the start of each system
// unless the print-object attribute has been set to "no" or the additional attribute
// has been set to "yes".
type Clef struct {
	Name         string `xml:"-"`
	Number       int    `xml:"number,attr,omitempty"`
	Additional   YesNo  `xml:"additional,attr,omitempty"`
	Size         string `xml:"size,attr,omitempty"`
	AfterBarline YesNo  `xml:"after-barline,attr,omitempty"`
	AttributeGroupPrintStyle
	AttributeGroupPrintObject
	AttributeGroupOptionalUniqueId
	GroupClef
}

// Double Named source named complex type "double"
// The double type indicates that the music is doubled one octave from
// what is currently written. If the above attribute is set to yes, the doubling is one
// octave above what is written, as for mixed flute / piccolo parts in band literature.
// Otherwise the doubling is one octave below what is written, as for mixed cello /
// bass parts in orchestral literature.
type Double struct {
	Name  string `xml:"-"`
	Above YesNo  `xml:"above,attr,omitempty"`
}

// ForPart Named source named complex type "for-part"
// The for-part type is used in a concert score to indicate the
// transposition for a transposed part created from that score. It is only used in
// score files that contain a concert-score element in the defaults. This allows
// concert scores with transposed parts to be represented in a single uncompressed
// MusicXML file. The optional number attribute refers to staff numbers, from top to
// bottom on the system. If absent, the child elements apply to all staves in the
// created part.
type ForPart struct {
	Name   string `xml:"-"`
	Number int    `xml:"number,attr,omitempty"`
	AttributeGroupOptionalUniqueId
	PartClef      *PartClef      `xml:"part-clef,omitempty"`
	PartTranspose *PartTranspose `xml:"part-transpose,omitempty"`
}

// Interchangeable Named source named complex type "interchangeable"
// The interchangeable type is used to represent the second in a pair of
// interchangeable dual time signatures, such as the 6/8 in 3/4 (6/8). A separate
// symbol attribute value is available compared to the time element's symbol attribute,
// which applies to the first of the dual time signatures.
type Interchangeable struct {
	Name         string `xml:"-"`
	Symbol       string `xml:"symbol,attr,omitempty"`
	Separator    string `xml:"separator,attr,omitempty"`
	TimeRelation string `xml:"time-relation,omitempty"`
	GroupTimeSignature
}

// Key Named source named complex type "key"
// The key type represents a key signature. Both traditional and
// non-traditional key signatures are supported. The optional number attribute refers
// to staff numbers. If absent, the key signature applies to all staves in the part.
// Key signatures appear at the start of each system unless the print-object attribute
// has been set to "no".
type Key struct {
	Name   string `xml:"-"`
	Number int    `xml:"number,attr,omitempty"`
	AttributeGroupPrintStyle
	AttributeGroupPrintObject
	AttributeGroupOptionalUniqueId
	GroupTraditionalKey
	GroupNonTraditionalKey
	KeyOctave []*KeyOctave `xml:"key-octave,omitempty"`
}

// KeyAccidental Named source named complex type "key-accidental"
// The key-accidental type indicates the accidental to be displayed in a
// non-traditional key signature, represented in the same manner as the accidental type
// without the formatting attributes.
type KeyAccidental struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText AccidentalValue `xml:",chardata"`
}

// KeyOctave Named source named complex type "key-octave"
// The key-octave type specifies in which octave an element of a key
// signature appears. The content specifies the octave value using the same values as
// the display-octave element. The number attribute is a positive integer that refers
// to the key signature element in left-to-right order. If the cancel attribute is set
// to yes, then this number refers to the canceling key signature specified by the
// cancel element in the parent key element. The cancel attribute cannot be set to yes
// if there is no corresponding cancel element within the parent key element. It is no
// by default.
type KeyOctave struct {
	Name   string `xml:"-"`
	Number int    `xml:"number,attr,omitempty"`
	Cancel YesNo  `xml:"cancel,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText int `xml:",chardata"`
}

// LineDetail Named source named complex type "line-detail"
// If the staff-lines element is present, the appearance of each line may
// be individually specified with a line-detail type. Staff lines are numbered from
// bottom to top. The print-object attribute allows lines to be hidden within a staff.
// This is used in special situations such as a widely-spaced percussion staff where a
// note placed below the higher line is distinct from a note placed above the lower
// line. Hidden staff lines are included when specifying clef lines and determining
// display-step / display-octave values, but are not counted as lines for the purposes
// of the system-layout and staff-layout elements.
type LineDetail struct {
	Name  string `xml:"-"`
	Line  int    `xml:"line,attr,omitempty"`
	Width string `xml:"width,attr,omitempty"`
	AttributeGroupColor
	AttributeGroupLineType
	AttributeGroupPrintObject
}

// MeasureRepeat Named source named complex type "measure-repeat"
// The measure-repeat type is used for both single and multiple measure
// repeats. The text of the element indicates the number of measures to be repeated in
// a single pattern. The slashes attribute specifies the number of slashes to use in
// the repeat sign. It is 1 if not specified. The text of the element is ignored when
// the type is stop. The stop type indicates the first measure where the repeats are no
// longer displayed. Both the start and the stop of the measure-repeat should be
// specified unless the repeats are displayed through the end of the part. The
// measure-repeat element specifies a notation style for repetitions. The actual music
// being repeated needs to be repeated within each measure of the MusicXML file. This
// element specifies the notation that indicates the repeat.
type MeasureRepeat struct {
	Name    string    `xml:"-"`
	Type    StartStop `xml:"type,attr,omitempty"`
	Slashes int       `xml:"slashes,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// MeasureStyle Named source named complex type "measure-style"
// A measure-style indicates a special way to print partial to multiple
// measures within a part. This includes multiple rests over several measures, repeats
// of beats, single, or multiple measures, and use of slash notation. The multiple-rest
// and measure-repeat elements indicate the number of measures covered in the element
// content. The beat-repeat and slash elements can cover partial measures. All but the
// multiple-rest element use a type attribute to indicate starting and stopping the use
// of the style. The optional number attribute specifies the staff number from top to
// bottom on the system, as with clef.
type MeasureStyle struct {
	Name   string `xml:"-"`
	Number int    `xml:"number,attr,omitempty"`
	AttributeGroupFont
	AttributeGroupColor
	AttributeGroupOptionalUniqueId
	MultipleRest  *MultipleRest  `xml:"multiple-rest,omitempty"`
	MeasureRepeat *MeasureRepeat `xml:"measure-repeat,omitempty"`
	BeatRepeat    *BeatRepeat    `xml:"beat-repeat,omitempty"`
	Slash         *Slash         `xml:"slash,omitempty"`
}

// MultipleRest Named source named complex type "multiple-rest"
// The text of the multiple-rest type indicates the number of measures in
// the multiple rest. Multiple rests may use the 1-bar / 2-bar / 4-bar rest symbols, or
// a single shape. The use-symbols attribute indicates which to use; it is no if not
// specified.
type MultipleRest struct {
	Name       string `xml:"-"`
	UseSymbols YesNo  `xml:"use-symbols,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText int `xml:",chardata"`
}

// PartClef Named source named complex type "part-clef"
// The child elements of the part-clef type have the same meaning as for
// the clef type. However that meaning applies to a transposed part created from the
// existing score file.
type PartClef struct {
	Name string `xml:"-"`
	GroupClef
}

// PartSymbol Named source named complex type "part-symbol"
// The part-symbol type indicates how a symbol for a multi-staff part is
// indicated in the score; brace is the default value. The top-staff and bottom-staff
// attributes are used when the brace does not extend across the entire part. For
// example, in a 3-staff organ part, the top-staff will typically be 1 for the right
// hand, while the bottom-staff will typically be 2 for the left hand. Staff 3 for the
// pedals is usually outside the brace. By default, the presence of a part-symbol
// element that does not extend across the entire part also indicates a corresponding
// change in the common barlines within a part.
type PartSymbol struct {
	Name        string `xml:"-"`
	TopStaff    int    `xml:"top-staff,attr,omitempty"`
	BottomStaff int    `xml:"bottom-staff,attr,omitempty"`
	AttributeGroupPosition
	AttributeGroupColor

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText GroupSymbolValue `xml:",chardata"`
}

// PartTranspose Named source named complex type "part-transpose"
// The child elements of the part-transpose type have the same meaning as
// for the transpose type. However that meaning applies to a transposed part created
// from the existing score file.
type PartTranspose struct {
	Name string `xml:"-"`
	GroupTranspose
}

// Slash Named source named complex type "slash"
// The slash type is used to indicate that slash notation is to be used.
// If the slash is on every beat, use-stems is no (the default). To indicate rhythms
// but not pitches, use-stems is set to yes. The type attribute indicates whether this
// is the start or stop of a slash notation style. The use-dots attribute works as for
// the beat-repeat element, and only has effect if use-stems is no.
type Slash struct {
	Name     string    `xml:"-"`
	Type     StartStop `xml:"type,attr,omitempty"`
	UseDots  YesNo     `xml:"use-dots,attr,omitempty"`
	UseStems YesNo     `xml:"use-stems,attr,omitempty"`
	GroupSlash
}

// StaffDetails Named source named complex type "staff-details"
// The staff-details element is used to indicate different types of
// staves. The optional number attribute specifies the staff number from top to bottom
// on the system, as with clef. The print-object attribute is used to indicate when a
// staff is not printed in a part, usually in large scores where empty parts are
// omitted. It is yes by default. If print-spacing is yes while print-object is no, the
// score is printed in cutaway format where vertical space is left for the empty part.
type StaffDetails struct {
	Name      string `xml:"-"`
	Number    int    `xml:"number,attr,omitempty"`
	ShowFrets string `xml:"show-frets,attr,omitempty"`
	AttributeGroupPrintObject
	AttributeGroupPrintSpacing
	StaffType   string         `xml:"staff-type,omitempty"`
	StaffLines  int            `xml:"staff-lines,omitempty"`
	LineDetail  []*LineDetail  `xml:"line-detail,omitempty"`
	StaffTuning []*StaffTuning `xml:"staff-tuning,omitempty"`
	Capo        int            `xml:"capo,omitempty"`
	StaffSize   *StaffSize     `xml:"staff-size,omitempty"`
}

// StaffSize Named source named complex type "staff-size"
// The staff-size element indicates how large a staff space is on this
// staff, expressed as a percentage of the work's default scaling. Values less than 100
// make the staff space smaller while values over 100 make the staff space larger. A
// staff-type of cue, ossia, or editorial implies a staff-size of less than 100, but
// the exact value is implementation-dependent unless specified here. Staff size
// affects staff height only, not the relationship of the staff to the left and right
// margins. In some cases, a staff-size different than 100 also scales the notation on
// the staff, such as with a cue staff. In other cases, such as percussion staves, the
// lines may be more widely spaced without scaling the notation on the staff. The
// scaling attribute allows these two cases to be distinguished. It specifies the
// percentage scaling that applies to the notation. Values less that 100 make the
// notation smaller while values over 100 make the notation larger. The staff-size
// content and scaling attribute are both non-negative decimal values.
type StaffSize struct {
	Name    string `xml:"-"`
	Scaling string `xml:"scaling,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// StaffTuning Named source named complex type "staff-tuning"
// The staff-tuning type specifies the open, non-capo tuning of the lines
// on a tablature staff.
type StaffTuning struct {
	Name string `xml:"-"`
	Line int    `xml:"line,attr,omitempty"`
	GroupTuning
}

// Time Named source named complex type "time"
// Time signatures are represented by the beats element for the numerator
// and the beat-type element for the denominator. The symbol attribute is used to
// indicate common and cut time symbols as well as a single number display. Multiple
// pairs of beat and beat-type elements are used for composite time signatures with
// multiple denominators, such as 2/4 + 3/8. A composite such as 3+2/8 requires only
// one beat/beat-type pair. The print-object attribute allows a time signature to be
// specified but not printed, as is the case for excerpts from the middle of a score.
// The value is "yes" if not present. The optional number attribute refers to staff
// numbers within the part. If absent, the time signature applies to all staves in the
// part.
type Time struct {
	Name      string        `xml:"-"`
	Number    int           `xml:"number,attr,omitempty"`
	Symbol    TimeSymbol    `xml:"symbol,attr,omitempty"`
	Separator TimeSeparator `xml:"separator,attr,omitempty"`
	AttributeGroupPrintStyleAlign
	AttributeGroupPrintObject
	AttributeGroupOptionalUniqueId
	GroupTimeSignature
	Interchangeable *Interchangeable `xml:"interchangeable,omitempty"`
	SenzaMisura     string           `xml:"senza-misura,omitempty"`
}

// Transpose Named source named complex type "transpose"
// The transpose type represents what must be added to a written pitch to
// get a correct sounding pitch. The optional number attribute refers to staff numbers,
// from top to bottom on the system. If absent, the transposition applies to all staves
// in the part. Per-staff transposition is most often used in parts that represent
// multiple instruments.
type Transpose struct {
	Name   string `xml:"-"`
	Number int    `xml:"number,attr,omitempty"`
	AttributeGroupOptionalUniqueId
	GroupTranspose
}

// BarStyleColor Named source named complex type "bar-style-color"
// The bar-style-color type contains barline style and color information.
type BarStyleColor struct {
	Name string `xml:"-"`
	AttributeGroupColor

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Barline Named source named complex type "barline"
// If a barline is other than a normal single barline, it should be
// represented by a barline type that describes it. This includes information about
// repeats and multiple endings, as well as line style. Barline data is on the same
// level as the other musical data in a score - a child of a measure in a partwise
// score, or a part in a timewise score. This allows for barlines within measures, as
// in dotted barlines that subdivide measures in complex meters. The two fermata
// elements allow for fermatas on both sides of the barline (the lower one inverted).
// Barlines have a location attribute to make it easier to process barlines
// independently of the other musical data in a score. It is often easier to set up
// measures separately from entering notes. The location attribute must match where the
// barline element occurs within the rest of the musical data in the score. If location
// is left, it should be the first element in the measure, aside from the print,
// bookmark, and link elements. If location is right, it should be the last element,
// again with the possible exception of the print, bookmark, and link elements. If no
// location is specified, the right barline is the default. The segno, coda, and
// divisions attributes work the same way as in the sound element. They are used for
// playback when barline elements contain segno or coda child elements.
type Barline struct {
	Name      string `xml:"-"`
	Location  string `xml:"location,attr,omitempty"`
	Segno     string `xml:"segno,attr,omitempty"`
	Coda      string `xml:"coda,attr,omitempty"`
	Divisions int `xml:"divisions,attr,omitempty"`
	AttributeGroupOptionalUniqueId
	BarStyle *BarStyleColor `xml:"bar-style,omitempty"`
	GroupEditorial
	WavyLine *WavyLine `xml:"wavy-line,omitempty"`
	Segno1   *Segno    `xml:"segno,omitempty"`
	Coda1    *Coda     `xml:"coda,omitempty"`
	Fermata  *Fermata  `xml:"fermata,omitempty"`
	Ending   *Ending   `xml:"ending,omitempty"`
	Repeat   *Repeat   `xml:"repeat,omitempty"`
}

// Ending Named source named complex type "ending"
// The ending type represents multiple (e.g. first and second) endings.
// Typically, the start type is associated with the left barline of the first measure
// in an ending. The stop and discontinue types are associated with the right barline
// of the last measure in an ending. Stop is used when the ending mark concludes with a
// downward jog, as is typical for first endings. Discontinue is used when there is no
// downward jog, as is typical for second endings that do not conclude a piece. The
// length of the jog can be specified using the end-length attribute. The text-x and
// text-y attributes are offsets that specify where the baseline of the start of the
// ending text appears, relative to the start of the ending line. The number attribute
// indicates which times the ending is played, similar to the time-only attribute used
// by other elements. While this often represents the numeric values for what is under
// the ending line, it can also indicate whether an ending is played during a larger
// dal segno or da capo repeat. Single endings such as "1" or comma-separated multiple
// endings such as "1,2" may be used. The ending element text is used when the text
// displayed in the ending is different than what appears in the number attribute. The
// print-object attribute is used to indicate when an ending is present but not
// printed, as is often the case for many parts in a full score.
type Ending struct {
	Name      string `xml:"-"`
	Number    string `xml:"number,attr,omitempty"`
	Type      string `xml:"type,attr,omitempty"`
	EndLength string `xml:"end-length,attr,omitempty"`
	TextX     string `xml:"text-x,attr,omitempty"`
	TextY     string `xml:"text-y,attr,omitempty"`
	AttributeGroupPrintObject
	AttributeGroupPrintStyle
	AttributeGroupSystemRelation

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Repeat Named source named complex type "repeat"
// The repeat type represents repeat marks. The start of the repeat has a
// forward direction while the end of the repeat has a backward direction. The times
// and after-jump attributes are only used with backward repeats that are not part of
// an ending. The times attribute indicates the number of times the repeated section is
// played. The after-jump attribute indicates if the repeats are played after a jump
// due to a da capo or dal segno.
type Repeat struct {
	Name      string `xml:"-"`
	Direction string `xml:"direction,attr,omitempty"`
	Times     int    `xml:"times,attr,omitempty"`
	AfterJump YesNo  `xml:"after-jump,attr,omitempty"`
	Winged    string `xml:"winged,attr,omitempty"`
}

// Accord Named source named complex type "accord"
// The accord type represents the tuning of a single string in the
// scordatura element. It uses the same group of elements as the staff-tuning element.
// Strings are numbered from high to low.
type Accord struct {
	Name   string `xml:"-"`
	String int    `xml:"string,attr,omitempty"`
	GroupTuning
}

// AccordionRegistration Named source named complex type "accordion-registration"
// The accordion-registration type is used for accordion registration
// symbols. These are circular symbols divided horizontally into high, middle, and low
// sections that correspond to 4', 8', and 16' pipes. Each accordion-high,
// accordion-middle, and accordion-low element represents the presence of one or more
// dots in the registration diagram. An accordion-registration element needs to have at
// least one of the child elements present.
type AccordionRegistration struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyleAlign
	AttributeGroupOptionalUniqueId
	AccordionHigh   string `xml:"accordion-high,omitempty"`
	AccordionMiddle int    `xml:"accordion-middle,omitempty"`
	AccordionLow    string `xml:"accordion-low,omitempty"`
}

// Barre Named source named complex type "barre"
// The barre element indicates placing a finger over multiple strings on
// a single fret. The type is "start" for the lowest pitched string (e.g., the string
// with the highest MusicXML number) and is "stop" for the highest pitched string.
type Barre struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`
	AttributeGroupColor
}

// Bass Named source named complex type "bass"
// The bass type is used to indicate a bass note in popular music chord
// symbols, e.g. G/C. It is generally not used in functional harmony, as inversion is
// generally not used in pop chord symbols. As with root, it is divided into step and
// alter elements, similar to pitches. The arrangement attribute specifies where the
// bass is displayed relative to what precedes it.
type Bass struct {
	Name          string        `xml:"-"`
	Arrangement   string        `xml:"arrangement,attr,omitempty"`
	BassSeparator *StyleText    `xml:"bass-separator,omitempty"`
	BassStep      *BassStep     `xml:"bass-step,omitempty"`
	BassAlter     *HarmonyAlter `xml:"bass-alter,omitempty"`
}

// HarmonyAlter Named source named complex type "harmony-alter"
// The harmony-alter type represents the chromatic alteration of the
// root, numeral, or bass of the current harmony-chord group within the harmony
// element. In some chord styles, the text of the preceding element may include
// alteration information. In that case, the print-object attribute of this type can be
// set to no. The location attribute indicates whether the alteration should appear to
// the left or the right of the preceding element. Its default value varies by element.
type HarmonyAlter struct {
	Name     string    `xml:"-"`
	Location LeftRight `xml:"location,attr,omitempty"`
	AttributeGroupPrintObject
	AttributeGroupPrintStyle

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// BassStep Named source named complex type "bass-step"
// The bass-step type represents the pitch step of the bass of the
// current chord within the harmony element. The text attribute indicates how the bass
// should appear in a score if not using the element contents.
type BassStep struct {
	Name string `xml:"-"`
	Text string `xml:"text,attr,omitempty"`
	AttributeGroupPrintStyle

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Beater Named source named complex type "beater"
// The beater type represents pictograms for beaters, mallets, and sticks
// that do not have different materials represented in the pictogram.
type Beater struct {
	Name string `xml:"-"`
	Tip  string `xml:"tip,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// BeatUnitTied Named source named complex type "beat-unit-tied"
// The beat-unit-tied type indicates a beat-unit within a metronome mark
// that is tied to the preceding beat-unit. This allows two or more tied notes to be
// associated with a per-minute value in a metronome mark, whereas the metronome-tied
// element is restricted to metric relationship marks.
type BeatUnitTied struct {
	Name string `xml:"-"`
	GroupBeatUnit
}

// Bracket Named source named complex type "bracket"
// Brackets are combined with words in a variety of modern directions.
// The line-end attribute specifies if there is a jog up or down (or both), an arrow,
// or nothing at the start or end of the bracket. If the line-end is up or down, the
// length of the jog can be specified using the end-length attribute. The line-type is
// solid if not specified.
type Bracket struct {
	Name      string `xml:"-"`
	Type      string `xml:"type,attr,omitempty"`
	Number    int    `xml:"number,attr,omitempty"`
	LineEnd   string `xml:"line-end,attr,omitempty"`
	EndLength string `xml:"end-length,attr,omitempty"`
	AttributeGroupLineType
	AttributeGroupDashedFormatting
	AttributeGroupPosition
	AttributeGroupColor
	AttributeGroupOptionalUniqueId
}

// Dashes Named source named complex type "dashes"
// The dashes type represents dashes, used for instance with cresc. and
// dim. marks.
type Dashes struct {
	Name   string            `xml:"-"`
	Type   StartStopContinue `xml:"type,attr,omitempty"`
	Number int               `xml:"number,attr,omitempty"`
	AttributeGroupDashedFormatting
	AttributeGroupPosition
	AttributeGroupColor
	AttributeGroupOptionalUniqueId
}

// Degree Named source named complex type "degree"
// The degree type is used to add, alter, or subtract individual notes in
// the chord. The print-object attribute can be used to keep the degree from printing
// separately when it has already taken into account in the text attribute of the kind
// element. The degree-value and degree-type text attributes specify how the value and
// type of the degree should be displayed. A harmony of kind "other" can be spelled
// explicitly by using a series of degree elements together with a root.
type Degree struct {
	Name string `xml:"-"`
	AttributeGroupPrintObject
	DegreeValue *DegreeValue `xml:"degree-value,omitempty"`
	DegreeAlter *DegreeAlter `xml:"degree-alter,omitempty"`
	DegreeType  *DegreeType  `xml:"degree-type,omitempty"`
}

// DegreeAlter Named source named complex type "degree-alter"
// The degree-alter type represents the chromatic alteration for the
// current degree. If the degree-type value is alter or subtract, the degree-alter
// value is relative to the degree already in the chord based on its kind element. If
// the degree-type value is add, the degree-alter is relative to a dominant chord
// (major and perfect intervals except for a minor seventh). The plus-minus attribute
// is used to indicate if plus and minus symbols should be used instead of sharp and
// flat symbols to display the degree alteration. It is no if not specified.
type DegreeAlter struct {
	Name      string `xml:"-"`
	PlusMinus YesNo  `xml:"plus-minus,attr,omitempty"`
	AttributeGroupPrintStyle

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// DegreeType Named source named complex type "degree-type"
// The degree-type type indicates if this degree is an addition,
// alteration, or subtraction relative to the kind of the current chord. The value of
// the degree-type element affects the interpretation of the value of the degree-alter
// element. The text attribute specifies how the type of the degree should be
// displayed.
type DegreeType struct {
	Name string `xml:"-"`
	Text string `xml:"text,attr,omitempty"`
	AttributeGroupPrintStyle

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// DegreeValue Named source named complex type "degree-value"
// The content of the degree-value type is a number indicating the degree
// of the chord (1 for the root, 3 for third, etc). The text attribute specifies how
// the value of the degree should be displayed. The symbol attribute indicates that a
// symbol should be used in specifying the degree. If the symbol attribute is present,
// the value of the text attribute follows the symbol.
type DegreeValue struct {
	Name   string `xml:"-"`
	Symbol string `xml:"symbol,attr,omitempty"`
	Text   string `xml:"text,attr,omitempty"`
	AttributeGroupPrintStyle

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText int `xml:",chardata"`
}

// Direction Named source named complex type "direction"
// A direction is a musical indication that is not necessarily attached
// to a specific note. Two or more may be combined to indicate words followed by the
// start of a dashed line, the end of a wedge followed by dynamics, etc. For
// applications where a specific direction is indeed attached to a specific note, the
// direction element can be associated with the first note element that follows it in
// score order that is not in a different voice. By default, a series of direction-type
// elements and a series of child elements of a direction-type within a single
// direction element follow one another in sequence visually. For a series of
// direction-type children, non-positional formatting attributes are carried over from
// the previous element by default.
type Direction struct {
	Name string `xml:"-"`
	AttributeGroupPlacement
	AttributeGroupDirective
	AttributeGroupSystemRelation
	AttributeGroupOptionalUniqueId
	DirectionType []*DirectionType `xml:"direction-type,omitempty"`
	Offset        *Offset          `xml:"offset,omitempty"`
	GroupEditorialVoiceDirection
	GroupStaff
	Sound     *Sound     `xml:"sound,omitempty"`
	Listening *Listening `xml:"listening,omitempty"`
}

// DirectionType Named source named complex type "direction-type"
// Textual direction types may have more than 1 component due to multiple
// fonts. The dynamics element may also be used in the notations element. Attribute
// groups related to print suggestions apply to the individual direction-type, not to
// the overall direction.
type DirectionType struct {
	Name string `xml:"-"`
	AttributeGroupOptionalUniqueId
	Rehearsal             []*FormattedTextId      `xml:"rehearsal,omitempty"`
	Segno                 []*Segno                `xml:"segno,omitempty"`
	Coda                  []*Coda                 `xml:"coda,omitempty"`
	Words                 []*FormattedTextId      `xml:"words,omitempty"`
	Symbol                []*FormattedSymbolId    `xml:"symbol,omitempty"`
	Wedge                 *Wedge                  `xml:"wedge,omitempty"`
	Dynamics              []*Dynamics             `xml:"dynamics,omitempty"`
	Dashes                *Dashes                 `xml:"dashes,omitempty"`
	Bracket               *Bracket                `xml:"bracket,omitempty"`
	Pedal                 *Pedal                  `xml:"pedal,omitempty"`
	Metronome             *Metronome              `xml:"metronome,omitempty"`
	OctaveShift           *OctaveShift            `xml:"octave-shift,omitempty"`
	HarpPedals            *HarpPedals             `xml:"harp-pedals,omitempty"`
	Damp                  *EmptyPrintStyleAlignId `xml:"damp,omitempty"`
	DampAll               *EmptyPrintStyleAlignId `xml:"damp-all,omitempty"`
	Eyeglasses            *EmptyPrintStyleAlignId `xml:"eyeglasses,omitempty"`
	StringMute            *StringMute             `xml:"string-mute,omitempty"`
	Scordatura            *Scordatura             `xml:"scordatura,omitempty"`
	Image                 *Image                  `xml:"image,omitempty"`
	PrincipalVoice        *PrincipalVoice         `xml:"principal-voice,omitempty"`
	Percussion            []*Percussion           `xml:"percussion,omitempty"`
	AccordionRegistration *AccordionRegistration  `xml:"accordion-registration,omitempty"`
	StaffDivide           *StaffDivide            `xml:"staff-divide,omitempty"`
	OtherDirection        *OtherDirection         `xml:"other-direction,omitempty"`
}

// Effect Named source named complex type "effect"
// The effect type represents pictograms for sound effect percussion
// instruments. The smufl attribute is used to distinguish different SMuFL stylistic
// alternates.
type Effect struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Feature Named source named complex type "feature"
// The feature type is a part of the grouping element used for musical
// analysis. The type attribute represents the type of the feature and the element
// content represents its value. This type is flexible to allow for different analyses.
type Feature struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// FirstFret Named source named complex type "first-fret"
// The first-fret type indicates which fret is shown in the top space of
// the frame; it is fret 1 if the element is not present. The optional text attribute
// indicates how this is represented in the fret diagram, while the location attribute
// indicates whether the text appears to the left or right of the frame.
type FirstFret struct {
	Name     string `xml:"-"`
	Text     string `xml:"text,attr,omitempty"`
	Location string `xml:"location,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText int `xml:",chardata"`
}

// Frame Named source named complex type "frame"
// The frame type represents a frame or fretboard diagram used together
// with a chord symbol. The representation is based on the NIFF guitar grid with
// additional information. The frame type's unplayed attribute indicates what to
// display above a string that has no associated frame-note element. Typical values are
// x and the empty string. If the attribute is not present, the display of the unplayed
// string is application-defined.
type Frame struct {
	Name     string `xml:"-"`
	Height   string `xml:"height,attr,omitempty"`
	Width    string `xml:"width,attr,omitempty"`
	Unplayed string `xml:"unplayed,attr,omitempty"`
	AttributeGroupPosition
	AttributeGroupColor
	AttributeGroupHalign
	AttributeGroupValignImage
	AttributeGroupOptionalUniqueId
	FrameStrings int          `xml:"frame-strings,omitempty"`
	FrameFrets   int          `xml:"frame-frets,omitempty"`
	FirstFret    *FirstFret   `xml:"first-fret,omitempty"`
	FrameNote    []*FrameNote `xml:"frame-note,omitempty"`
}

// FrameNote Named source named complex type "frame-note"
// The frame-note type represents each note included in the frame. An
// open string will have a fret value of 0, while a muted string will not be associated
// with a frame-note element.
type FrameNote struct {
	Name      string      `xml:"-"`
	String    *StringType `xml:"string,omitempty"`
	Fret      *Fret       `xml:"fret,omitempty"`
	Fingering *Fingering  `xml:"fingering,omitempty"`
	Barre     *Barre      `xml:"barre,omitempty"`
}

// Glass Named source named complex type "glass"
// The glass type represents pictograms for glass percussion instruments.
// The smufl attribute is used to distinguish different SMuFL glyphs for wind chimes in
// the Chimes pictograms range, including those made of materials other than glass.
type Glass struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Grouping Named source named complex type "grouping"
// The grouping type is used for musical analysis. When the type
// attribute is "start" or "single", it usually contains one or more feature elements.
// The number attribute is used for distinguishing between overlapping and hierarchical
// groupings. The member-of attribute allows for easy distinguishing of what grouping
// elements are in what hierarchy. Feature elements contained within a "stop" type of
// grouping may be ignored. This element is flexible to allow for different types of
// analyses. Future versions of the MusicXML format may add elements that can represent
// more standardized categories of analysis data, allowing for easier data sharing.
type Grouping struct {
	Name     string `xml:"-"`
	Type     string `xml:"type,attr,omitempty"`
	Number   string `xml:"number,attr,omitempty"`
	MemberOf string `xml:"member-of,attr,omitempty"`
	AttributeGroupOptionalUniqueId
	Feature []*Feature `xml:"feature,omitempty"`
}

// Harmony Named source named complex type "harmony"
// The harmony type represents harmony analysis, including chord symbols
// in popular music as well as functional harmony analysis in classical music. If there
// are alternate harmonies possible, this can be specified using multiple harmony
// elements differentiated by type. Explicit harmonies have all note present in the
// music; implied have some notes missing but implied; alternate represents alternate
// analyses. The print-object attribute controls whether or not anything is printed due
// to the harmony element. The print-frame attribute controls printing of a frame or
// fretboard diagram. The print-style attribute group sets the default for the harmony,
// but individual elements can override this with their own print-style values. The
// arrangement attribute specifies how multiple harmony-chord groups are arranged
// relative to each other. Harmony-chords with vertical arrangement are separated by
// horizontal lines. Harmony-chords with diagonal or horizontal arrangement are
// separated by diagonal lines or slashes.
type Harmony struct {
	Name        string             `xml:"-"`
	Type        string             `xml:"type,attr,omitempty"`
	PrintFrame  YesNo              `xml:"print-frame,attr,omitempty"`
	Arrangement HarmonyArrangement `xml:"arrangement,attr,omitempty"`
	AttributeGroupPrintObject
	AttributeGroupPrintStyle
	AttributeGroupPlacement
	AttributeGroupSystemRelation
	AttributeGroupOptionalUniqueId
	GroupHarmonyChord
	Frame  *Frame  `xml:"frame,omitempty"`
	Offset *Offset `xml:"offset,omitempty"`
	GroupEditorial
	GroupStaff
}

// HarpPedals Named source named complex type "harp-pedals"
// The harp-pedals type is used to create harp pedal diagrams. The
// pedal-step and pedal-alter elements use the same values as the step and alter
// elements. For easiest reading, the pedal-tuning elements should follow standard harp
// pedal order, with pedal-step values of D, C, B, E, F, G, and A.
type HarpPedals struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyleAlign
	AttributeGroupOptionalUniqueId
	PedalTuning []*PedalTuning `xml:"pedal-tuning,omitempty"`
}

// Image Named source named complex type "image"
// The image type is used to include graphical images in a score.
type Image struct {
	Name string `xml:"-"`
	AttributeGroupImageAttributes
	AttributeGroupOptionalUniqueId
}

// InstrumentChange Named source named complex type "instrument-change"
// The instrument-change element type represents a change to the virtual
// instrument sound for a given score-instrument. The id attribute refers to the
// score-instrument affected by the change. All instrument-change child elements can
// also be initially specified within the score-instrument element.
type InstrumentChange struct {
	Name string `xml:"-"`
	Id   string `xml:"id,attr,omitempty"`
	GroupVirtualInstrumentData
}

// Inversion Named source named complex type "inversion"
// The inversion type represents harmony inversions. The value is a
// number indicating which inversion is used: 0 for root position, 1 for first
// inversion, etc. The text attribute indicates how the inversion should be displayed
// in a score.
type Inversion struct {
	Name string `xml:"-"`
	Text string `xml:"text,attr,omitempty"`
	AttributeGroupPrintStyle

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText int `xml:",chardata"`
}

// Kind Named source named complex type "kind"
// Kind indicates the type of chord. Degree elements can then add,
// subtract, or alter from these starting points The attributes are used to indicate
// the formatting of the symbol. Since the kind element is the constant in all the
// harmony-chord groups that can make up a polychord, many formatting attributes are
// here. The use-symbols attribute is yes if the kind should be represented when
// possible with harmony symbols rather than letters and numbers. These symbols
// include: major: a triangle, like Unicode 25B3 minor: -, like Unicode 002D augmented:
// +, like Unicode 002B diminished: °, like Unicode 00B0 half-diminished: ø, like
// Unicode 00F8 For the major-minor kind, only the minor symbol is used when
// use-symbols is yes. The major symbol is set using the symbol attribute in the
// degree-value element. The corresponding degree-alter value will usually be 0 in this
// case. The text attribute describes how the kind should be spelled in a score. If
// use-symbols is yes, the value of the text attribute follows the symbol. The
// stack-degrees attribute is yes if the degree elements should be stacked above each
// other. The parentheses-degrees attribute is yes if all the degrees should be in
// parentheses. The bracket-degrees attribute is yes if all the degrees should be in a
// bracket. If not specified, these values are implementation-specific. The alignment
// attributes are for the entire harmony-chord group of which this kind element is a
// part. The text attribute may use strings such as "13sus" that refer to both the kind
// and one or more degree elements. In this case, the corresponding degree elements
// should have the print-object attribute set to "no" to keep redundant alterations
// from being displayed.
type Kind struct {
	Name               string `xml:"-"`
	UseSymbols         YesNo  `xml:"use-symbols,attr,omitempty"`
	Text               string `xml:"text,attr,omitempty"`
	StackDegrees       YesNo  `xml:"stack-degrees,attr,omitempty"`
	ParenthesesDegrees YesNo  `xml:"parentheses-degrees,attr,omitempty"`
	BracketDegrees     YesNo  `xml:"bracket-degrees,attr,omitempty"`
	AttributeGroupPrintStyle
	AttributeGroupHalign
	AttributeGroupValign

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Listening Named source named complex type "listening"
// The listen and listening types, new in Version 4.0, specify different
// ways that a score following or machine listening application can interact with a
// performer. The listening type handles interactions that change the state of the
// listening application from the specified point in the performance onward. If
// multiple child elements of the same type are present, they should have distinct
// player and/or time-only attributes. The offset element is used to indicate that the
// listening change takes place offset from the current score position. If the
// listening element is a child of a direction element, the listening offset element
// overrides the direction offset element if both elements are present. Note that the
// offset reflects the intended musical position for the change in state. It should not
// be used to compensate for latency issues in particular hardware configurations.
type Listening struct {
	Name           string            `xml:"-"`
	Sync           []*Sync           `xml:"sync,omitempty"`
	OtherListening []*OtherListening `xml:"other-listening,omitempty"`
	Offset         *Offset           `xml:"offset,omitempty"`
}

// MeasureNumbering Named source named complex type "measure-numbering"
// The measure-numbering type describes how frequently measure numbers
// are displayed on this part. The text attribute from the measure element is used for
// display, or the number attribute if the text attribute is not present. Measures with
// an implicit attribute set to "yes" never display a measure number, regardless of the
// measure-numbering setting. The optional staff attribute refers to staff numbers
// within the part, from top to bottom on the system. It indicates which staff is used
// as the reference point for vertical positioning. A value of 1 is assumed if not
// present. The optional multiple-rest-always and multiple-rest-range attributes
// describe how measure numbers are shown on multiple rests when the measure-numbering
// value is not set to none. The multiple-rest-always attribute is set to yes when the
// measure number should always be shown, even if the multiple rest starts midway
// through a system when measure numbering is set to system level. The
// multiple-rest-range attribute is set to yes when measure numbers on multiple rests
// display the range of numbers for the first and last measure, rather than just the
// number of the first measure.
type MeasureNumbering struct {
	Name               string `xml:"-"`
	System             string `xml:"system,attr,omitempty"`
	Staff              int    `xml:"staff,attr,omitempty"`
	MultipleRestAlways YesNo  `xml:"multiple-rest-always,attr,omitempty"`
	MultipleRestRange  YesNo  `xml:"multiple-rest-range,attr,omitempty"`
	AttributeGroupPrintStyleAlign

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Membrane Named source named complex type "membrane"
// The membrane type represents pictograms for membrane percussion
// instruments. The smufl attribute is used to distinguish different SMuFL stylistic
// alternates.
type Membrane struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Metal Named source named complex type "metal"
// The metal type represents pictograms for metal percussion instruments.
// The smufl attribute is used to distinguish different SMuFL stylistic alternates.
type Metal struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Metronome Named source named complex type "metronome"
// The metronome type represents metronome marks and other metric
// relationships. The beat-unit group and per-minute element specify regular metronome
// marks. The metronome-note and metronome-relation elements allow for the
// specification of metric modulations and other metric relationships, such as swing
// tempo marks where two eighths are equated to a quarter note / eighth note triplet.
// Tied notes can be represented in both types of metronome marks by using the
// beat-unit-tied and metronome-tied elements. The parentheses attribute indicates
// whether or not to put the metronome mark in parentheses; its value is no if not
// specified. The print-object attribute is set to no in cases where the metronome
// element represents a relationship or range that is not displayed in the music
// notation.
type Metronome struct {
	Name        string `xml:"-"`
	Parentheses YesNo  `xml:"parentheses,attr,omitempty"`
	AttributeGroupPrintStyleAlign
	AttributeGroupPrintObject
	AttributeGroupJustify
	AttributeGroupOptionalUniqueId
	GroupBeatUnit
	PerMinute         *PerMinute       `xml:"per-minute,omitempty"`
	BeatUnitTied      []*BeatUnitTied  `xml:"beat-unit-tied,omitempty"`
	MetronomeArrows   string           `xml:"metronome-arrows,omitempty"`
	MetronomeRelation string           `xml:"metronome-relation,omitempty"`
	MetronomeNote     []*MetronomeNote `xml:"metronome-note,omitempty"`
}

// MetronomeBeam Named source named complex type "metronome-beam"
// The metronome-beam type works like the beam type in defining metric
// relationships, but does not include all the attributes available in the beam type.
type MetronomeBeam struct {
	Name   string `xml:"-"`
	Number int    `xml:"number,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText BeamValue `xml:",chardata"`
}

// MetronomeNote Named source named complex type "metronome-note"
// The metronome-note type defines the appearance of a note within a
// metric relationship mark.
type MetronomeNote struct {
	Name            string           `xml:"-"`
	MetronomeType   string           `xml:"metronome-type,omitempty"`
	MetronomeDot    string           `xml:"metronome-dot,omitempty"`
	MetronomeBeam   []*MetronomeBeam `xml:"metronome-beam,omitempty"`
	MetronomeTied   *MetronomeTied   `xml:"metronome-tied,omitempty"`
	MetronomeTuplet *MetronomeTuplet `xml:"metronome-tuplet,omitempty"`
}

// MetronomeTied Named source named complex type "metronome-tied"
// The metronome-tied indicates the presence of a tie within a metric
// relationship mark. As with the tied element, both the start and stop of the tie
// should be specified, in this case within separate metronome-note elements.
type MetronomeTied struct {
	Name string    `xml:"-"`
	Type StartStop `xml:"type,attr,omitempty"`
}

// MetronomeTuplet Named source named complex type "metronome-tuplet"
// The metronome-tuplet type uses the same element structure as the
// time-modification element along with some attributes from the tuplet element.
type MetronomeTuplet struct {
	Name string `xml:"-"`
}

// Numeral Named source named complex type "numeral"
// The numeral type represents the Roman numeral or Nashville number part
// of a harmony. It requires that the key be specified in the encoding, either with a
// key or numeral-key element.
type Numeral struct {
	Name         string        `xml:"-"`
	NumeralRoot  *NumeralRoot  `xml:"numeral-root,omitempty"`
	NumeralAlter *HarmonyAlter `xml:"numeral-alter,omitempty"`
	NumeralKey   *NumeralKey   `xml:"numeral-key,omitempty"`
}

// NumeralKey Named source named complex type "numeral-key"
// The numeral-key type is used when the key for the numeral is different
// than the key specified by the key signature. The numeral-fifths element specifies
// the key in the same way as the fifths element. The numeral-mode element specifies
// the mode similar to the mode element, but with a restricted set of values
type NumeralKey struct {
	Name string `xml:"-"`
	AttributeGroupPrintObject
	NumeralFifths int    `xml:"numeral-fifths,omitempty"`
	NumeralMode   string `xml:"numeral-mode,omitempty"`
}

// NumeralRoot Named source named complex type "numeral-root"
// The numeral-root type represents the Roman numeral or Nashville number
// as a positive integer from 1 to 7. The text attribute indicates how the numeral
// should appear in the score. A numeral-root value of 5 with a kind of major would
// have a text attribute of "V" if displayed as a Roman numeral, and "5" if displayed
// as a Nashville number. If the text attribute is not specified, the display is
// application-dependent.
type NumeralRoot struct {
	Name string `xml:"-"`
	Text string `xml:"text,attr,omitempty"`
	AttributeGroupPrintStyle

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText int `xml:",chardata"`
}

// OctaveShift Named source named complex type "octave-shift"
// The octave shift type indicates where notes are shifted up or down
// from their true pitched values because of printing difficulty. Thus a treble clef
// line noted with 8va will be indicated with an octave-shift down from the pitch data
// indicated in the notes. A size of 8 indicates one octave; a size of 15 indicates two
// octaves.
type OctaveShift struct {
	Name   string `xml:"-"`
	Type   string `xml:"type,attr,omitempty"`
	Number int    `xml:"number,attr,omitempty"`
	Size   int    `xml:"size,attr,omitempty"`
	AttributeGroupDashedFormatting
	AttributeGroupPrintStyle
	AttributeGroupOptionalUniqueId
}

// Offset Named source named complex type "offset"
// An offset is represented in terms of divisions, and indicates where
// the direction will appear relative to the current musical location. The current
// musical location is always within the current measure, even at the end of a measure.
// The offset affects the visual appearance of the direction. If the sound attribute is
// "yes", then the offset affects playback and listening too. If the sound attribute is
// "no", then any sound or listening associated with the direction takes effect at the
// current location. The sound attribute is "no" by default for compatibility with
// earlier versions of the MusicXML format. If an element within a direction includes a
// default-x attribute, the offset value will be ignored when determining the
// appearance of that element.
type Offset struct {
	Name  string `xml:"-"`
	Sound YesNo  `xml:"sound,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// OtherDirection Named source named complex type "other-direction"
// The other-direction type is used to define any direction symbols not
// yet in the MusicXML format. The smufl attribute can be used to specify a particular
// direction symbol, allowing application interoperability without requiring every
// SMuFL glyph to have a MusicXML element equivalent. Using the other-direction type
// without the smufl attribute allows for extended representation, though without
// application interoperability.
type OtherDirection struct {
	Name string `xml:"-"`
	AttributeGroupPrintObject
	AttributeGroupPrintStyleAlign
	AttributeGroupSmufl
	AttributeGroupOptionalUniqueId

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// OtherListening Named source named complex type "other-listening"
// The other-listening type represents other types of listening control
// and interaction. The required type attribute indicates the type of listening to
// which the element content applies. The optional player and time-only attributes
// restrict the element to apply to a single player or set of times through a repeated
// section, respectively.
type OtherListening struct {
	Name     string   `xml:"-"`
	Type     string   `xml:"type,attr,omitempty"`
	Player   string   `xml:"player,attr,omitempty"`
	TimeOnly TimeOnly `xml:"time-only,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Pedal Named source named complex type "pedal"
// The pedal type represents piano pedal marks, including damper and
// sostenuto pedal marks. The line attribute is yes if pedal lines are used. The sign
// attribute is yes if Ped, Sost, and * signs are used. For compatibility with older
// versions, the sign attribute is yes by default if the line attribute is no, and is
// no by default if the line attribute is yes. If the sign attribute is set to yes and
// the type is start or sostenuto, the abbreviated attribute is yes if the short P and
// S signs are used, and no if the full Ped and Sost signs are used. It is no by
// default. Otherwise the abbreviated attribute is ignored. The alignment attributes
// are ignored if the sign attribute is no.
type Pedal struct {
	Name        string `xml:"-"`
	Type        string `xml:"type,attr,omitempty"`
	Number      int    `xml:"number,attr,omitempty"`
	Line        YesNo  `xml:"line,attr,omitempty"`
	Sign        YesNo  `xml:"sign,attr,omitempty"`
	Abbreviated YesNo  `xml:"abbreviated,attr,omitempty"`
	AttributeGroupPrintStyleAlign
	AttributeGroupOptionalUniqueId
}

// PedalTuning Named source named complex type "pedal-tuning"
// The pedal-tuning type specifies the tuning of a single harp pedal.
type PedalTuning struct {
	Name       string `xml:"-"`
	PedalStep  Step   `xml:"pedal-step,omitempty"`
	PedalAlter string `xml:"pedal-alter,omitempty"`
}

// PerMinute Named source named complex type "per-minute"
// The per-minute type can be a number, or a text description including
// numbers. If a font is specified, it overrides the font specified for the overall
// metronome element. This allows separate specification of a music font for the
// beat-unit and a text font for the numeric value, in cases where a single metronome
// font is not used.
type PerMinute struct {
	Name string `xml:"-"`
	AttributeGroupFont

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Percussion Named source named complex type "percussion"
// The percussion element is used to define percussion pictogram symbols.
// Definitions for these symbols can be found in Kurt Stone's "Music Notation in the
// Twentieth Century" on pages 206-212 and 223. Some values are added to these based on
// how usage has evolved in the 30 years since Stone's book was published.
type Percussion struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyleAlign
	AttributeGroupEnclosure
	AttributeGroupOptionalUniqueId
	Glass           *Glass     `xml:"glass,omitempty"`
	Metal           *Metal     `xml:"metal,omitempty"`
	Wood            *Wood      `xml:"wood,omitempty"`
	Pitched         *Pitched   `xml:"pitched,omitempty"`
	Membrane        *Membrane  `xml:"membrane,omitempty"`
	Effect          *Effect    `xml:"effect,omitempty"`
	Timpani         *Timpani   `xml:"timpani,omitempty"`
	Beater          *Beater    `xml:"beater,omitempty"`
	Stick           *Stick     `xml:"stick,omitempty"`
	StickLocation   string     `xml:"stick-location,omitempty"`
	OtherPercussion *OtherText `xml:"other-percussion,omitempty"`
}

// Pitched Named source named complex type "pitched"
// The pitched-value type represents pictograms for pitched percussion
// instruments. The smufl attribute is used to distinguish different SMuFL glyphs for a
// particular pictogram within the Tuned mallet percussion pictograms range.
type Pitched struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// PrincipalVoice Named source named complex type "principal-voice"
// The principal-voice type represents principal and secondary voices in
// a score, either for analysis or for square bracket symbols that appear in a score.
// The element content is used for analysis and may be any text value. The symbol
// attribute indicates the type of symbol used. When used for analysis separate from
// any printed score markings, it should be set to none. Otherwise if the type is stop
// it should be set to plain.
type PrincipalVoice struct {
	Name   string    `xml:"-"`
	Type   StartStop `xml:"type,attr,omitempty"`
	Symbol string    `xml:"symbol,attr,omitempty"`
	AttributeGroupPrintStyleAlign
	AttributeGroupOptionalUniqueId

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Print Named source named complex type "print"
// The print type contains general printing parameters, including layout
// elements. The part-name-display and part-abbreviation-display elements may also be
// used here to change how a part name or abbreviation is displayed over the course of
// a piece. They take effect when the current measure or a succeeding measure starts a
// new system. Layout group elements in a print element only apply to the current page,
// system, or staff. Music that follows continues to take the default values from the
// layout determined by the defaults element.
type Print struct {
	Name string `xml:"-"`
	AttributeGroupPrintAttributes
	AttributeGroupOptionalUniqueId
	GroupLayout
	MeasureLayout           *MeasureLayout    `xml:"measure-layout,omitempty"`
	MeasureNumbering        *MeasureNumbering `xml:"measure-numbering,omitempty"`
	PartNameDisplay         *NameDisplay      `xml:"part-name-display,omitempty"`
	PartAbbreviationDisplay *NameDisplay      `xml:"part-abbreviation-display,omitempty"`
}

// Root Named source named complex type "root"
// The root type indicates a pitch like C, D, E vs. a scale degree like
// 1, 2, 3. It is used with chord symbols in popular music. The root element has a
// root-step and optional root-alter element similar to the step and alter elements,
// but renamed to distinguish the different musical meanings.
type Root struct {
	Name      string        `xml:"-"`
	RootStep  *RootStep     `xml:"root-step,omitempty"`
	RootAlter *HarmonyAlter `xml:"root-alter,omitempty"`
}

// RootStep Named source named complex type "root-step"
// The root-step type represents the pitch step of the root of the
// current chord within the harmony element. The text attribute indicates how the root
// should appear in a score if not using the element contents.
type RootStep struct {
	Name string `xml:"-"`
	Text string `xml:"text,attr,omitempty"`
	AttributeGroupPrintStyle

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText Step `xml:",chardata"`
}

// Scordatura Named source named complex type "scordatura"
// Scordatura string tunings are represented by a series of accord
// elements, similar to the staff-tuning elements. Strings are numbered from high to
// low.
type Scordatura struct {
	Name string `xml:"-"`
	AttributeGroupOptionalUniqueId
	Accord []*Accord `xml:"accord,omitempty"`
}

// Sound Named source named complex type "sound"
// The sound element contains general playback parameters. They can stand
// alone within a part/measure, or be a component element within a direction. Tempo is
// expressed in quarter notes per minute. If 0, the sound-generating program should
// prompt the user at the time of compiling a sound (MIDI) file. Dynamics (or MIDI
// velocity) are expressed as a percentage of the default forte value (90 for MIDI
// 1.0). Dacapo indicates to go back to the beginning of the movement. When used it
// always has the value "yes". Segno and dalsegno are used for backwards jumps to a
// segno sign; coda and tocoda are used for forward jumps to a coda sign. If there are
// multiple jumps, the value of these parameters can be used to name and distinguish
// them. If segno or coda is used, the divisions attribute can also be used to indicate
// the number of divisions per quarter note. Otherwise sound and MIDI generating
// programs may have to recompute this. By default, a dalsegno or dacapo attribute
// indicates that the jump should occur the first time through, while a tocoda
// attribute indicates the jump should occur the second time through. The time that
// jumps occur can be changed by using the time-only attribute. The forward-repeat
// attribute indicates that a forward repeat sign is implied but not displayed. It is
// used for example in two-part forms with repeats, such as a minuet and trio where no
// repeat is displayed at the start of the trio. This usually occurs after a barline.
// When used it always has the value of "yes". The fine attribute follows the final
// note or rest in a movement with a da capo or dal segno direction. If numeric, the
// value represents the actual duration of the final note or rest, which can be
// ambiguous in written notation and different among parts and voices. The value may
// also be "yes" to indicate no change to the final duration. If the sound element
// applies only particular times through a repeat, the time-only attribute indicates
// which times to apply the sound element. Pizzicato in a sound element effects all
// following notes. Yes indicates pizzicato, no indicates arco. The pan and elevation
// attributes are deprecated in Version 2.0. The pan and elevation elements in the
// midi-instrument element should be used instead. The meaning of the pan and elevation
// attributes is the same as for the pan and elevation elements. If both are present,
// the mid-instrument elements take priority. The damper-pedal, soft-pedal, and
// sostenuto-pedal attributes effect playback of the three common piano pedals and
// their MIDI controller equivalents. The yes value indicates the pedal is depressed;
// no indicates the pedal is released. A numeric value from 0 to 100 may also be used
// for half pedaling. This value is the percentage that the pedal is depressed. A value
// of 0 is equivalent to no, and a value of 100 is equivalent to yes. Instrument
// changes, MIDI devices, MIDI instruments, and playback techniques are changed using
// the instrument-change, midi-device, midi-instrument, and play elements. When there
// are multiple instances of these elements, they should be grouped together by
// instrument using the id attribute values. The offset element is used to indicate
// that the sound takes place offset from the current score position. If the sound
// element is a child of a direction element, the sound offset element overrides the
// direction offset element if both elements are present. Note that the offset reflects
// the intended musical position for the change in sound. It should not be used to
// compensate for latency issues in particular hardware configurations.
type Sound struct {
	Name           string   `xml:"-"`
	Tempo          string   `xml:"tempo,attr,omitempty"`
	Dynamics       string   `xml:"dynamics,attr,omitempty"`
	Dacapo         YesNo    `xml:"dacapo,attr,omitempty"`
	Segno          string   `xml:"segno,attr,omitempty"`
	Dalsegno       string   `xml:"dalsegno,attr,omitempty"`
	Coda           string   `xml:"coda,attr,omitempty"`
	Tocoda         string   `xml:"tocoda,attr,omitempty"`
	Divisions      int   `xml:"divisions,attr,omitempty"`
	ForwardRepeat  YesNo    `xml:"forward-repeat,attr,omitempty"`
	Fine           string   `xml:"fine,attr,omitempty"`
	TimeOnly       TimeOnly `xml:"time-only,attr,omitempty"`
	Pizzicato      YesNo    `xml:"pizzicato,attr,omitempty"`
	Pan            string   `xml:"pan,attr,omitempty"`
	Elevation      string   `xml:"elevation,attr,omitempty"`
	DamperPedal    string   `xml:"damper-pedal,attr,omitempty"`
	SoftPedal      string   `xml:"soft-pedal,attr,omitempty"`
	SostenutoPedal string   `xml:"sostenuto-pedal,attr,omitempty"`
	AttributeGroupOptionalUniqueId
	InstrumentChange []*InstrumentChange `xml:"instrument-change,omitempty"`
	MidiDevice       []*MidiDevice       `xml:"midi-device,omitempty"`
	MidiInstrument   []*MidiInstrument   `xml:"midi-instrument,omitempty"`
	Play             []*Play             `xml:"play,omitempty"`
	Swing            *Swing              `xml:"swing,omitempty"`
	Offset           *Offset             `xml:"offset,omitempty"`
}

// StaffDivide Named source named complex type "staff-divide"
// The staff-divide element represents the staff division arrow symbols
// found at SMuFL code points U+E00B, U+E00C, and U+E00D.
type StaffDivide struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`
	AttributeGroupPrintStyleAlign
	AttributeGroupOptionalUniqueId
}

// Stick Named source named complex type "stick"
// The stick type represents pictograms where the material of the stick,
// mallet, or beater is included.The parentheses and dashed-circle attributes indicate
// the presence of these marks around the round beater part of a pictogram. Values for
// these attributes are "no" if not present.
type Stick struct {
	Name          string       `xml:"-"`
	Tip           TipDirection `xml:"tip,attr,omitempty"`
	Parentheses   YesNo        `xml:"parentheses,attr,omitempty"`
	DashedCircle  YesNo        `xml:"dashed-circle,attr,omitempty"`
	StickType     string       `xml:"stick-type,omitempty"`
	StickMaterial string       `xml:"stick-material,omitempty"`
}

// StringMute Named source named complex type "string-mute"
// The string-mute type represents string mute on and mute off symbols.
type StringMute struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`
	AttributeGroupPrintStyleAlign
	AttributeGroupOptionalUniqueId
}

// Swing Named source named complex type "swing"
// The swing element specifies whether or not to use swing playback,
// where consecutive on-beat / off-beat eighth or 16th notes are played with unequal
// nominal durations. The straight element specifies that no swing is present, so
// consecutive notes have equal durations. The first and second elements are positive
// integers that specify the ratio between durations of consecutive notes. For example,
// a first element with a value of 2 and a second element with a value of 1 applied to
// eighth notes specifies a quarter note / eighth note tuplet playback, where the first
// note is twice as long as the second note. Ratios should be specified with the
// smallest integers possible. For example, a ratio of 6 to 4 should be specified as 3
// to 2 instead. The optional swing-type element specifies the note type, either eighth
// or 16th, to which the ratio is applied. The value is eighth if this element is not
// present. The optional swing-style element is a string describing the style of swing
// used. The swing element has no effect for playback of grace notes, notes where a
// type element is not present, and notes where the specified duration is different
// than the nominal value associated with the specified type. If a swung note has
// attack and release attributes, those values modify the swung playback.
type Swing struct {
	Name       string        `xml:"-"`
	Straight   string        `xml:"straight,omitempty"`
	First      int           `xml:"first,omitempty"`
	Second     int           `xml:"second,omitempty"`
	SwingType  NoteTypeValue `xml:"swing-type,omitempty"`
	SwingStyle string        `xml:"swing-style,omitempty"`
}

// Sync Named source named complex type "sync"
// The sync type specifies the style that a score following application
// should use the synchronize an accompaniment with a performer. If this type is not
// included in a score, default synchronization depends on the application. The
// optional latency attribute specifies a time in milliseconds that the listening
// application should expect from the performer. The optional player and time-only
// attributes restrict the element to apply to a single player or set of times through
// a repeated section, respectively.
type Sync struct {
	Name     string   `xml:"-"`
	Type     string   `xml:"type,attr,omitempty"`
	Latency  int      `xml:"latency,attr,omitempty"`
	Player   string   `xml:"player,attr,omitempty"`
	TimeOnly TimeOnly `xml:"time-only,attr,omitempty"`
}

// Timpani Named source named complex type "timpani"
// The timpani type represents the timpani pictogram. The smufl attribute
// is used to distinguish different SMuFL stylistic alternates.
type Timpani struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`
}

// Wedge Named source named complex type "wedge"
// The wedge type represents crescendo and diminuendo wedge symbols. The
// type attribute is crescendo for the start of a wedge that is closed at the left
// side, and diminuendo for the start of a wedge that is closed on the right side.
// Spread values are measured in tenths; those at the start of a crescendo wedge or end
// of a diminuendo wedge are ignored. The niente attribute is yes if a circle appears
// at the point of the wedge, indicating a crescendo from nothing or diminuendo to
// nothing. It is no by default, and used only when the type is crescendo, or the type
// is stop for a wedge that began with a diminuendo type. The line-type is solid if not
// specified.
type Wedge struct {
	Name   string `xml:"-"`
	Type   string `xml:"type,attr,omitempty"`
	Number int    `xml:"number,attr,omitempty"`
	Spread string `xml:"spread,attr,omitempty"`
	Niente YesNo  `xml:"niente,attr,omitempty"`
	AttributeGroupLineType
	AttributeGroupDashedFormatting
	AttributeGroupPosition
	AttributeGroupColor
	AttributeGroupOptionalUniqueId
}

// Wood Named source named complex type "wood"
// The wood type represents pictograms for wood percussion instruments.
// The smufl attribute is used to distinguish different SMuFL stylistic alternates.
type Wood struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Encoding Named source named complex type "encoding"
// The encoding element contains information about who did the digital
// encoding, when, with what software, and in what aspects. Standard type values for
// the encoder element are music, words, and arrangement, but other types may be used.
// The type attribute is only needed when there are multiple encoder elements.
type Encoding struct {
	Name                string       `xml:"-"`
	Encoder             []*TypedText `xml:"encoder,omitempty"`
	Software            string       `xml:"software,omitempty"`
	EncodingDescription string       `xml:"encoding-description,omitempty"`
	Supports            []*Supports  `xml:"supports,omitempty"`
}

// Identification Named source named complex type "identification"
// Identification contains basic metadata about the score. It includes
// information that may apply at a score-wide, movement-wide, or part-wide level. The
// creator, rights, source, and relation elements are based on Dublin Core.
type Identification struct {
	Name          string         `xml:"-"`
	Creator       []*TypedText   `xml:"creator,omitempty"`
	Rights        []*TypedText   `xml:"rights,omitempty"`
	Encoding      *Encoding      `xml:"encoding,omitempty"`
	Source        string         `xml:"source,omitempty"`
	Relation      []*TypedText   `xml:"relation,omitempty"`
	Miscellaneous *Miscellaneous `xml:"miscellaneous,omitempty"`
}

// Miscellaneous Named source named complex type "miscellaneous"
// If a program has other metadata not yet supported in the MusicXML
// format, it can go in the miscellaneous element. The miscellaneous type puts each
// separate part of metadata into its own miscellaneous-field type.
type Miscellaneous struct {
	Name               string                `xml:"-"`
	MiscellaneousField []*MiscellaneousField `xml:"miscellaneous-field,omitempty"`
}

// MiscellaneousField Named source named complex type "miscellaneous-field"
// If a program has other metadata not yet supported in the MusicXML
// format, each type of metadata can go in a miscellaneous-field element. The required
// name attribute indicates the type of metadata the element content represents.
type MiscellaneousField struct {
	Name    string `xml:"-"`
	NameXSD string `xml:"name,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Supports Named source named complex type "supports"
// The supports type indicates if a MusicXML encoding supports a
// particular MusicXML element. This is recommended for elements like beam, stem, and
// accidental, where the absence of an element is ambiguous if you do not know if the
// encoding supports that element. For Version 2.0, the supports element is expanded to
// allow programs to indicate support for particular attributes or particular values.
// This lets applications communicate, for example, that all system and/or page breaks
// are contained in the MusicXML file.
type Supports struct {
	Name      string `xml:"-"`
	Type      YesNo  `xml:"type,attr,omitempty"`
	Element   string `xml:"element,attr,omitempty"`
	Attribute string `xml:"attribute,attr,omitempty"`
	Value     string `xml:"value,attr,omitempty"`
}

// Appearance Named source named complex type "appearance"
// The appearance type controls general graphical settings for the
// music's final form appearance on a printed page of display. This includes support
// for line widths, definitions for note sizes, and standard distances between notation
// elements, plus an extension element for other aspects of appearance.
type Appearance struct {
	Name            string             `xml:"-"`
	LineWidth       []*LineWidth       `xml:"line-width,omitempty"`
	NoteSize        []*NoteSize        `xml:"note-size,omitempty"`
	Distance        []*Distance        `xml:"distance,omitempty"`
	Glyph           []*Glyph           `xml:"glyph,omitempty"`
	OtherAppearance []*OtherAppearance `xml:"other-appearance,omitempty"`
}

// Distance Named source named complex type "distance"
// The distance element represents standard distances between notation
// elements in tenths. The type attribute defines what type of distance is being
// defined. Valid values include hyphen (for hyphens in lyrics) and beam.
type Distance struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Glyph Named source named complex type "glyph"
// The glyph element represents what SMuFL glyph should be used for
// different variations of symbols that are semantically identical. The type attribute
// specifies what type of glyph is being defined. The element value specifies what
// SMuFL glyph to use, including recommended stylistic alternates. The SMuFL glyph name
// should match the type. For instance, a type of quarter-rest would use values
// restQuarter, restQuarterOld, or restQuarterZ. A type of g-clef-ottava-bassa would
// use values gClef8vb, gClef8vbOld, or gClef8vbCClef. A type of octave-shift-up-8
// would use values ottava, ottavaBassa, ottavaBassaBa, ottavaBassaVb, or octaveBassa.
type Glyph struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// LineWidth Named source named complex type "line-width"
// The line-width type indicates the width of a line type in tenths. The
// type attribute defines what type of line is being defined. Values include beam,
// bracket, dashes, enclosure, ending, extend, heavy barline, leger, light barline,
// octave shift, pedal, slur middle, slur tip, staff, stem, tie middle, tie tip, tuplet
// bracket, and wedge. The text content is expressed in tenths.
type LineWidth struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// MeasureLayout Named source named complex type "measure-layout"
// The measure-layout type includes the horizontal distance from the
// previous measure. It applies to the current measure only.
type MeasureLayout struct {
	Name            string `xml:"-"`
	MeasureDistance string `xml:"measure-distance,omitempty"`
}

// NoteSize Named source named complex type "note-size"
// The note-size type indicates the percentage of the regular note size
// to use for notes with a cue and large size as defined in the type element. The grace
// type is used for notes of cue size that that include a grace element. The cue type
// is used for all other notes with cue size, whether defined explicitly or implicitly
// via a cue element. The large type is used for notes of large size. The text content
// represent the numeric percentage. A value of 100 would be identical to the size of a
// regular note as defined by the music font.
type NoteSize struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// OtherAppearance Named source named complex type "other-appearance"
// The other-appearance type is used to define any graphical settings not
// yet in the current version of the MusicXML format. This allows extended
// representation, though without application interoperability.
type OtherAppearance struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// PageLayout Named source named complex type "page-layout"
// Page layout can be defined both in score-wide defaults and in the
// print element. Page margins are specified either for both even and odd pages, or via
// separate odd and even page number values. The type is not needed when used as part
// of a print element. If omitted when used in the defaults element, "both" is the
// default. If no page-layout element is present in the defaults element, default page
// layout values are chosen by the application. When used in the print element, the
// page-layout element affects the appearance of the current page only. All other pages
// use the default values as determined by the defaults element. If any child elements
// are missing from the page-layout element in a print element, the values determined
// by the defaults element are used there as well.
type PageLayout struct {
	Name        string       `xml:"-"`
	PageHeight  string       `xml:"page-height,omitempty"`
	PageWidth   string       `xml:"page-width,omitempty"`
	PageMargins *PageMargins `xml:"page-margins,omitempty"`
}

// PageMargins Named source named complex type "page-margins"
// Page margins are specified either for both even and odd pages, or via
// separate odd and even page number values. The type attribute is not needed when used
// as part of a print element. If omitted when the page-margins type is used in the
// defaults element, "both" is the default value.
type PageMargins struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`
	GroupAllMargins
}

// Scaling Named source named complex type "scaling"
// Margins, page sizes, and distances are all measured in tenths to keep
// MusicXML data in a consistent coordinate system as much as possible. The translation
// to absolute units is done with the scaling type, which specifies how many
// millimeters are equal to how many tenths. For a staff height of 7 mm, millimeters
// would be set to 7 while tenths is set to 40. The ability to set a formula rather
// than a single scaling factor helps avoid roundoff errors.
type Scaling struct {
	Name        string `xml:"-"`
	Millimeters string `xml:"millimeters,omitempty"`
	Tenths      string `xml:"tenths,omitempty"`
}

// StaffLayout Named source named complex type "staff-layout"
// Staff layout includes the vertical distance from the bottom line of
// the previous staff in this system to the top line of the staff specified by the
// number attribute. The optional number attribute refers to staff numbers within the
// part, from top to bottom on the system. A value of 1 is used if not present. When
// used in the defaults element, the values apply to all systems in all parts. When
// used in the print element, the values apply to the current system only. This value
// is ignored for the first staff in a system.
type StaffLayout struct {
	Name          string `xml:"-"`
	Number        int    `xml:"number,attr,omitempty"`
	StaffDistance string `xml:"staff-distance,omitempty"`
}

// SystemDividers Named source named complex type "system-dividers"
// The system-dividers element indicates the presence or absence of
// system dividers (also known as system separation marks) between systems displayed on
// the same page. Dividers on the left and right side of the page are controlled by the
// left-divider and right-divider elements respectively. The default vertical position
// is half the system-distance value from the top of the system that is below the
// divider. The default horizontal position is the left and right system margin,
// respectively. When used in the print element, the system-dividers element affects
// the dividers that would appear between the current system and the previous system.
type SystemDividers struct {
	Name         string                      `xml:"-"`
	LeftDivider  *EmptyPrintObjectStyleAlign `xml:"left-divider,omitempty"`
	RightDivider *EmptyPrintObjectStyleAlign `xml:"right-divider,omitempty"`
}

// SystemLayout Named source named complex type "system-layout"
// A system is a group of staves that are read and played simultaneously.
// System layout includes left and right margins and the vertical distance from the
// previous system. The system distance is measured from the bottom line of the
// previous system to the top line of the current system. It is ignored for the first
// system on a page. The top system distance is measured from the page's top margin to
// the top line of the first system. It is ignored for all but the first system on a
// page. Sometimes the sum of measure widths in a system may not equal the system width
// specified by the layout elements due to roundoff or other errors. The behavior when
// reading MusicXML files in these cases is application-dependent. For instance,
// applications may find that the system layout data is more reliable than the sum of
// the measure widths, and adjust the measure widths accordingly. When used in the
// defaults element, the system-layout element defines a default appearance for all
// systems in the score. If no system-layout element is present in the defaults
// element, default system layout values are chosen by the application. When used in
// the print element, the system-layout element affects the appearance of the current
// system only. All other systems use the default values as determined by the defaults
// element. If any child elements are missing from the system-layout element in a print
// element, the values determined by the defaults element are used there as well. This
// type of system-layout element need only be read from or written to the first visible
// part in the score.
type SystemLayout struct {
	Name              string          `xml:"-"`
	SystemMargins     *SystemMargins  `xml:"system-margins,omitempty"`
	SystemDistance    string          `xml:"system-distance,omitempty"`
	TopSystemDistance string          `xml:"top-system-distance,omitempty"`
	SystemDividers    *SystemDividers `xml:"system-dividers,omitempty"`
}

// SystemMargins Named source named complex type "system-margins"
// System margins are relative to the page margins. Positive values
// indent and negative values reduce the margin size.
type SystemMargins struct {
	Name string `xml:"-"`
	GroupLeftRightMargins
}

// Bookmark Named source named complex type "bookmark"
// The bookmark type serves as a well-defined target for an incoming
// simple XLink.
type Bookmark struct {
	Name    string `xml:"-"`
	Id      string `xml:"id,attr,omitempty"`
	NameXSD string `xml:"name,attr,omitempty"`
	AttributeGroupElementPosition
}

// Link Named source named complex type "link"
// The link type serves as an outgoing simple XLink. If a relative link
// is used within a document that is part of a compressed MusicXML file, the link is
// relative to the root folder of the zip file.
type Link struct {
	Name    string `xml:"-"`
	NameXSD string `xml:"name,attr,omitempty"`
	AttributeGroupLinkAttributes
	AttributeGroupElementPosition
	AttributeGroupPosition
}

// Accidental Named source named complex type "accidental"
// The accidental type represents actual notated accidentals. Editorial
// and cautionary indications are indicated by attributes. Values for these attributes
// are "no" if not present. Specific graphic display such as parentheses, brackets, and
// size are controlled by the level-display attribute group.
type Accidental struct {
	Name       string `xml:"-"`
	Cautionary string `xml:"cautionary,attr,omitempty"`
	Editorial  YesNo  `xml:"editorial,attr,omitempty"`
	Smufl      string `xml:"smufl,attr,omitempty"`
	AttributeGroupLevelDisplay
	AttributeGroupPrintStyle

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// AccidentalMark Named source named complex type "accidental-mark"
// An accidental-mark can be used as a separate notation or as part of an
// ornament. When used in an ornament, position and placement are relative to the
// ornament, not relative to the note.
type AccidentalMark struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`
	AttributeGroupLevelDisplay
	AttributeGroupPrintStyle
	AttributeGroupPlacement
	AttributeGroupOptionalUniqueId

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText AccidentalValue `xml:",chardata"`
}

// Arpeggiate Named source named complex type "arpeggiate"
// The arpeggiate type indicates that this note is part of an arpeggiated
// chord. The number attribute can be used to distinguish between two simultaneous
// chords arpeggiated separately (different numbers) or together (same number). The
// direction attribute is used if there is an arrow on the arpeggio sign. By default,
// arpeggios go from the lowest to highest note. The length of the sign can be
// determined from the position attributes for the arpeggiate elements used with the
// top and bottom notes of the arpeggiated chord. If the unbroken attribute is set to
// yes, it indicates that the arpeggio continues onto another staff within the part.
// This serves as a hint to applications and is not required for cross-staff arpeggios.
type Arpeggiate struct {
	Name      string `xml:"-"`
	Number    int    `xml:"number,attr,omitempty"`
	Direction string `xml:"direction,attr,omitempty"`
	Unbroken  YesNo  `xml:"unbroken,attr,omitempty"`
	AttributeGroupPosition
	AttributeGroupPlacement
	AttributeGroupColor
	AttributeGroupOptionalUniqueId
}

// Articulations Named source named complex type "articulations"
// Articulations and accents are grouped together here.
type Articulations struct {
	Name string `xml:"-"`
	AttributeGroupOptionalUniqueId
	Accent            []*EmptyPlacement     `xml:"accent,omitempty"`
	StrongAccent      []*StrongAccent       `xml:"strong-accent,omitempty"`
	Staccato          []*EmptyPlacement     `xml:"staccato,omitempty"`
	Tenuto            []*EmptyPlacement     `xml:"tenuto,omitempty"`
	DetachedLegato    []*EmptyPlacement     `xml:"detached-legato,omitempty"`
	Staccatissimo     []*EmptyPlacement     `xml:"staccatissimo,omitempty"`
	Spiccato          []*EmptyPlacement     `xml:"spiccato,omitempty"`
	Scoop             []*EmptyLine          `xml:"scoop,omitempty"`
	Plop              []*EmptyLine          `xml:"plop,omitempty"`
	Doit              []*EmptyLine          `xml:"doit,omitempty"`
	Falloff           []*EmptyLine          `xml:"falloff,omitempty"`
	BreathMark        []*BreathMark         `xml:"breath-mark,omitempty"`
	Caesura           []*Caesura            `xml:"caesura,omitempty"`
	Stress            []*EmptyPlacement     `xml:"stress,omitempty"`
	Unstress          []*EmptyPlacement     `xml:"unstress,omitempty"`
	SoftAccent        []*EmptyPlacement     `xml:"soft-accent,omitempty"`
	OtherArticulation []*OtherPlacementText `xml:"other-articulation,omitempty"`
}

// Arrow Named source named complex type "arrow"
// The arrow element represents an arrow used for a musical technical
// indication. It can represent both Unicode and SMuFL arrows. The presence of an
// arrowhead element indicates that only the arrowhead is displayed, not the arrow
// stem. The smufl attribute distinguishes different SMuFL glyphs that have an arrow
// appearance such as arrowBlackUp, guitarStrumUp, or handbellsSwingUp. The specified
// glyph should match the descriptive representation.
type Arrow struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement
	AttributeGroupSmufl
	ArrowDirection string `xml:"arrow-direction,omitempty"`
	ArrowStyle     string `xml:"arrow-style,omitempty"`
	Arrowhead      string `xml:"arrowhead,omitempty"`
	CircularArrow  string `xml:"circular-arrow,omitempty"`
}

// Assess Named source named complex type "assess"
// By default, an assessment application should assess all notes without
// a cue child element, and not assess any note with a cue child element. The assess
// type allows this default assessment to be overridden for individual notes. The
// optional player and time-only attributes restrict the type to apply to a single
// player or set of times through a repeated section, respectively. If missing, the
// type applies to all players or all times through the repeated section, respectively.
// The player attribute references the id attribute of a player element defined within
// the matching score-part.
type Assess struct {
	Name     string `xml:"-"`
	Type     YesNo  `xml:"type,attr,omitempty"`
	Player   string `xml:"player,attr,omitempty"`
	TimeOnly string `xml:"time-only,attr,omitempty"`
}

// Backup Named source named complex type "backup"
// The backup and forward elements are required to coordinate multiple
// voices in one part, including music on multiple staves. The backup type is generally
// used to move between voices and staves. Thus the backup element does not include
// voice or staff elements. Duration values should always be positive, and should not
// cross measure boundaries or mid-measure changes in the divisions value.
type Backup struct {
	Name string `xml:"-"`
	GroupDuration
	GroupEditorial
}

// Beam Named source named complex type "beam"
// Beam values include begin, continue, end, forward hook, and backward
// hook. Up to eight concurrent beams are available to cover up to 1024th notes. Each
// beam in a note is represented with a separate beam element, starting with the eighth
// note beam using a number attribute of 1. Note that the beam number does not
// distinguish sets of beams that overlap, as it does for slur and other elements.
// Beaming groups are distinguished by being in different voices and/or the presence or
// absence of grace and cue elements. Beams that have a begin value can also have a fan
// attribute to indicate accelerandos and ritardandos using fanned beams. The fan
// attribute may also be used with a continue value if the fanning direction changes on
// that note. The value is "none" if not specified. The repeater attribute has been
// deprecated in MusicXML 3.0. Formerly used for tremolos, it needs to be specified
// with a "yes" value for each beam using it.
type Beam struct {
	Name     string `xml:"-"`
	Number   int    `xml:"number,attr,omitempty"`
	Repeater YesNo  `xml:"repeater,attr,omitempty"`
	Fan      string `xml:"fan,attr,omitempty"`
	AttributeGroupColor
	AttributeGroupOptionalUniqueId

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Bend Named source named complex type "bend"
// The bend type is used in guitar notation and tablature. A single note
// with a bend and release will contain two bend elements: the first to represent the
// bend and the second to represent the release. The shape attribute distinguishes
// between the angled bend symbols commonly used in standard notation and the curved
// bend symbols commonly used in both tablature and standard notation.
type Bend struct {
	Name  string `xml:"-"`
	Shape string `xml:"shape,attr,omitempty"`
	AttributeGroupPrintStyle
	AttributeGroupBendSound
	BendAlter string         `xml:"bend-alter,omitempty"`
	PreBend   string         `xml:"pre-bend,omitempty"`
	Release   *Release       `xml:"release,omitempty"`
	WithBar   *PlacementText `xml:"with-bar,omitempty"`
}

// BreathMark Named source named complex type "breath-mark"
// The breath-mark element indicates a place to take a breath.
type BreathMark struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Caesura Named source named complex type "caesura"
// The caesura element indicates a slight pause. It is notated using a
// "railroad tracks" symbol or other variations specified in the element content.
type Caesura struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Elision Named source named complex type "elision"
// The elision type represents an elision between lyric syllables. The
// text content specifies the symbol used to display the elision. Common values are a
// no-break space (Unicode 00A0), an underscore (Unicode 005F), or an undertie (Unicode
// 203F). If the text content is empty, the smufl attribute is used to specify the
// symbol to use. Its value is a SMuFL canonical glyph name that starts with lyrics.
// The SMuFL attribute is ignored if the elision glyph is already specified by the text
// content. If neither text content nor a smufl attribute are present, the elision
// glyph is application-specific.
type Elision struct {
	Name  string `xml:"-"`
	Smufl string `xml:"smufl,attr,omitempty"`
	AttributeGroupFont
	AttributeGroupColor

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// EmptyLine Named source named complex type "empty-line"
// The empty-line type represents an empty element with line-shape,
// line-type, line-length, dashed-formatting, print-style and placement attributes.
type EmptyLine struct {
	Name string `xml:"-"`
	AttributeGroupLineShape
	AttributeGroupLineType
	AttributeGroupLineLength
	AttributeGroupDashedFormatting
	AttributeGroupPrintStyle
	AttributeGroupPlacement
}

// Extend Named source named complex type "extend"
// The extend type represents lyric word extension / melisma lines as
// well as figured bass extensions. The optional type and position attributes are added
// in Version 3.0 to provide better formatting control.
type Extend struct {
	Name string            `xml:"-"`
	Type StartStopContinue `xml:"type,attr,omitempty"`
	AttributeGroupPosition
	AttributeGroupColor
}

// Figure Named source named complex type "figure"
// The figure type represents a single figure within a figured-bass
// element.
type Figure struct {
	Name         string     `xml:"-"`
	Prefix       *StyleText `xml:"prefix,omitempty"`
	FigureNumber *StyleText `xml:"figure-number,omitempty"`
	Suffix       *StyleText `xml:"suffix,omitempty"`
	Extend       *Extend    `xml:"extend,omitempty"`
	GroupEditorial
}

// FiguredBass Named source named complex type "figured-bass"
// The figured-bass element represents figured bass notation. Figured
// bass elements take their position from the first regular note (not a grace note or
// chord note) that follows in score order. The optional duration element is used to
// indicate changes of figures under a note. Figures are ordered from top to bottom.
// The value of parentheses is "no" if not present.
type FiguredBass struct {
	Name        string `xml:"-"`
	Parentheses YesNo  `xml:"parentheses,attr,omitempty"`
	AttributeGroupPrintStyleAlign
	AttributeGroupPlacement
	AttributeGroupPrintout
	AttributeGroupOptionalUniqueId
	Figure []*Figure `xml:"figure,omitempty"`
	GroupDuration
	GroupEditorial
}

// Forward Named source named complex type "forward"
// The backup and forward elements are required to coordinate multiple
// voices in one part, including music on multiple staves. The forward element is
// generally used within voices and staves. Duration values should always be positive,
// and should not cross measure boundaries or mid-measure changes in the divisions
// value.
type Forward struct {
	Name string `xml:"-"`
	GroupDuration
	GroupEditorialVoice
	GroupStaff
}

// Glissando Named source named complex type "glissando"
// Glissando and slide types both indicate rapidly moving from one pitch
// to the other so that individual notes are not discerned. A glissando sounds the
// distinct notes in between the two pitches and defaults to a wavy line. The optional
// text is printed alongside the line.
type Glissando struct {
	Name   string    `xml:"-"`
	Type   StartStop `xml:"type,attr,omitempty"`
	Number int       `xml:"number,attr,omitempty"`
	AttributeGroupLineType
	AttributeGroupDashedFormatting
	AttributeGroupPrintStyle
	AttributeGroupOptionalUniqueId

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Grace Named source named complex type "grace"
// The grace type indicates the presence of a grace note. The slash
// attribute for a grace note is yes for slashed grace notes. The steal-time-previous
// attribute indicates the percentage of time to steal from the previous note for the
// grace note. The steal-time-following attribute indicates the percentage of time to
// steal from the following note for the grace note, as for appoggiaturas. The
// make-time attribute indicates to make time, not steal time; the units are in
// real-time divisions for the grace note.
type Grace struct {
	Name               string `xml:"-"`
	StealTimePrevious  string `xml:"steal-time-previous,attr,omitempty"`
	StealTimeFollowing string `xml:"steal-time-following,attr,omitempty"`
	MakeTime           string `xml:"make-time,attr,omitempty"`
	Slash              YesNo  `xml:"slash,attr,omitempty"`
}

// HammerOnPullOff Named source named complex type "hammer-on-pull-off"
// The hammer-on and pull-off elements are used in guitar and fretted
// instrument notation. Since a single slur can be marked over many notes, the
// hammer-on and pull-off elements are separate so the individual pair of notes can be
// specified. The element content can be used to specify how the hammer-on or pull-off
// should be notated. An empty element leaves this choice up to the application.
type HammerOnPullOff struct {
	Name   string    `xml:"-"`
	Type   StartStop `xml:"type,attr,omitempty"`
	Number int       `xml:"number,attr,omitempty"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Handbell Named source named complex type "handbell"
// The handbell element represents notation for various techniques used
// in handbell and handchime music.
type Handbell struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// HarmonClosed Named source named complex type "harmon-closed"
// The harmon-closed type represents whether the harmon mute is closed,
// open, or half-open. The optional location attribute indicates which portion of the
// symbol is filled in when the element value is half.
type HarmonClosed struct {
	Name     string `xml:"-"`
	Location string `xml:"location,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// HarmonMute Named source named complex type "harmon-mute"
// The harmon-mute type represents the symbols used for harmon mutes in
// brass notation.
type HarmonMute struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement
	HarmonClosed *HarmonClosed `xml:"harmon-closed,omitempty"`
}

// Harmonic Named source named complex type "harmonic"
// The harmonic type indicates natural and artificial harmonics. Allowing
// the type of pitch to be specified, combined with controls for appearance/playback
// differences, allows both the notation and the sound to be represented. Artificial
// harmonics can add a notated touching pitch; artificial pinch harmonics will usually
// not notate a touching pitch. The attributes for the harmonic element refer to the
// use of the circular harmonic symbol, typically but not always used with natural
// harmonics.
type Harmonic struct {
	Name string `xml:"-"`
	AttributeGroupPrintObject
	AttributeGroupPrintStyle
	AttributeGroupPlacement
	Natural       string `xml:"natural,omitempty"`
	Artificial    string `xml:"artificial,omitempty"`
	BasePitch     string `xml:"base-pitch,omitempty"`
	TouchingPitch string `xml:"touching-pitch,omitempty"`
	SoundingPitch string `xml:"sounding-pitch,omitempty"`
}

// HeelToe Named source named complex type "heel-toe"
// The heel and toe elements are used with organ pedals. The substitution
// value is "no" if the attribute is not present.
type HeelToe struct {
	Name string `xml:"-"`
}

// Hole Named source named complex type "hole"
// The hole type represents the symbols used for woodwind and brass
// fingerings as well as other notations.
type Hole struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement
	HoleType   string      `xml:"hole-type,omitempty"`
	HoleClosed *HoleClosed `xml:"hole-closed,omitempty"`
	HoleShape  string      `xml:"hole-shape,omitempty"`
}

// HoleClosed Named source named complex type "hole-closed"
// The hole-closed type represents whether the hole is closed, open, or
// half-open. The optional location attribute indicates which portion of the hole is
// filled in when the element value is half.
type HoleClosed struct {
	Name     string `xml:"-"`
	Location string `xml:"location,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Instrument Named source named complex type "instrument"
// The instrument type distinguishes between score-instrument elements in
// a score-part. The id attribute is an IDREF back to the score-instrument ID. If
// multiple score-instruments are specified in a score-part, there should be an
// instrument element for each note in the part. Notes that are shared between multiple
// score-instruments can have more than one instrument element.
type Instrument struct {
	Name string `xml:"-"`
	Id   string `xml:"id,attr,omitempty"`
}

// Listen Named source named complex type "listen"
// The listen and listening types, new in Version 4.0, specify different
// ways that a score following or machine listening application can interact with a
// performer. The listen type handles interactions that are specific to a note. If
// multiple child elements of the same type are present, they should have distinct
// player and/or time-only attributes.
type Listen struct {
	Name        string            `xml:"-"`
	Assess      []*Assess         `xml:"assess,omitempty"`
	Wait        []*Wait           `xml:"wait,omitempty"`
	OtherListen []*OtherListening `xml:"other-listen,omitempty"`
}

// Lyric Named source named complex type "lyric"
// The lyric type represents text underlays for lyrics. Two text elements
// that are not separated by an elision element are part of the same syllable, but may
// have different text formatting. The MusicXML XSD is more strict than the DTD in
// enforcing this by disallowing a second syllabic element unless preceded by an
// elision element. The lyric number indicates multiple lines, though a name can be
// used as well. Common name examples are verse and chorus. Justification is center by
// default; placement is below by default. Vertical alignment is to the baseline of the
// text and horizontal alignment matches justification. The print-object attribute can
// override a note's print-lyric attribute in cases where only some lyrics on a note
// are printed, as when lyrics for later verses are printed in a block of text rather
// than with each note. The time-only attribute precisely specifies which lyrics are to
// be sung which time through a repeated section.
type Lyric struct {
	Name     string   `xml:"-"`
	Number   string   `xml:"number,attr,omitempty"`
	NameXSD  string   `xml:"name,attr,omitempty"`
	TimeOnly TimeOnly `xml:"time-only,attr,omitempty"`
	AttributeGroupJustify
	AttributeGroupPosition
	AttributeGroupPlacement
	AttributeGroupColor
	AttributeGroupPrintObject
	AttributeGroupOptionalUniqueId
	Elision      []*Elision         `xml:"elision,omitempty"`
	Syllabic     string             `xml:"syllabic,omitempty"`
	Text         []*TextElementData `xml:"text,omitempty"`
	Extend       *Extend            `xml:"extend,omitempty"`
	Laughing     string             `xml:"laughing,omitempty"`
	Humming      string             `xml:"humming,omitempty"`
	EndLine      string             `xml:"end-line,omitempty"`
	EndParagraph string             `xml:"end-paragraph,omitempty"`
	GroupEditorial
}

// Mordent Named source named complex type "mordent"
// The mordent type is used for both represents the mordent sign with the
// vertical line and the inverted-mordent sign without the line. The long attribute is
// "no" by default. The approach and departure attributes are used for compound
// ornaments, indicating how the beginning and ending of the ornament look relative to
// the main part of the mordent.
type Mordent struct {
	Name string `xml:"-"`
}

// NonArpeggiate Named source named complex type "non-arpeggiate"
// The non-arpeggiate type indicates that this note is at the top or
// bottom of a bracket indicating to not arpeggiate these notes. Since this does not
// involve playback, it is only used on the top or bottom notes, not on each note as
// for the arpeggiate type.
type NonArpeggiate struct {
	Name   string `xml:"-"`
	Type   string `xml:"type,attr,omitempty"`
	Number int    `xml:"number,attr,omitempty"`
	AttributeGroupPosition
	AttributeGroupPlacement
	AttributeGroupColor
	AttributeGroupOptionalUniqueId
}

// Notations Named source named complex type "notations"
// Notations refer to musical notations, not XML notations. Multiple
// notations are allowed in order to represent multiple editorial levels. The
// print-object attribute, added in Version 3.0, allows notations to represent details
// of performance technique, such as fingerings, without having them appear in the
// score.
type Notations struct {
	Name string `xml:"-"`
	AttributeGroupPrintObject
	AttributeGroupOptionalUniqueId
	GroupEditorial
	Tied           []*Tied           `xml:"tied,omitempty"`
	Slur           []*Slur           `xml:"slur,omitempty"`
	Tuplet         []*Tuplet         `xml:"tuplet,omitempty"`
	Glissando      []*Glissando      `xml:"glissando,omitempty"`
	Slide          []*Slide          `xml:"slide,omitempty"`
	Ornaments      []*Ornaments      `xml:"ornaments,omitempty"`
	Technical      []*Technical      `xml:"technical,omitempty"`
	Articulations  []*Articulations  `xml:"articulations,omitempty"`
	Dynamics       []*Dynamics       `xml:"dynamics,omitempty"`
	Fermata        []*Fermata        `xml:"fermata,omitempty"`
	Arpeggiate     []*Arpeggiate     `xml:"arpeggiate,omitempty"`
	NonArpeggiate  []*NonArpeggiate  `xml:"non-arpeggiate,omitempty"`
	AccidentalMark []*AccidentalMark `xml:"accidental-mark,omitempty"`
	OtherNotation  []*OtherNotation  `xml:"other-notation,omitempty"`
}

// Note Named source named complex type "note"
// Notes are the most common type of MusicXML data. The MusicXML format
// distinguishes between elements used for sound information and elements used for
// notation information (e.g., tie is used for sound, tied for notation). Thus grace
// notes do not have a duration element. Cue notes have a duration element, as do
// forward elements, but no tie elements. Having these two types of information
// available can make interchange easier, as some programs handle one type of
// information more readily than the other. The print-leger attribute is used to
// indicate whether leger lines are printed. Notes without leger lines are used to
// indicate indeterminate high and low notes. By default, it is set to yes. If
// print-object is set to no, print-leger is interpreted to also be set to no if not
// present. This attribute is ignored for rests. The dynamics and end-dynamics
// attributes correspond to MIDI 1.0's Note On and Note Off velocities, respectively.
// They are expressed in terms of percentages of the default forte value (90 for MIDI
// 1.0). The attack and release attributes are used to alter the starting and stopping
// time of the note from when it would otherwise occur based on the flow of durations -
// information that is specific to a performance. They are expressed in terms of
// divisions, either positive or negative. A note that starts a tie should not have a
// release attribute, and a note that stops a tie should not have an attack attribute.
// The attack and release attributes are independent of each other. The attack
// attribute only changes the starting time of a note, and the release attribute only
// changes the stopping time of a note. If a note is played only particular times
// through a repeat, the time-only attribute shows which times to play the note. The
// pizzicato attribute is used when just this note is sounded pizzicato, vs. the
// pizzicato element which changes overall playback between pizzicato and arco.
type Note struct {
	Name        string   `xml:"-"`
	PrintLeger  YesNo    `xml:"print-leger,attr,omitempty"`
	Dynamics    string   `xml:"dynamics,attr,omitempty"`
	EndDynamics string   `xml:"end-dynamics,attr,omitempty"`
	Attack      string   `xml:"attack,attr,omitempty"`
	Release     string   `xml:"release,attr,omitempty"`
	TimeOnly    TimeOnly `xml:"time-only,attr,omitempty"`
	Pizzicato   YesNo    `xml:"pizzicato,attr,omitempty"`
	AttributeGroupXPosition
	AttributeGroupFont
	AttributeGroupColor
	AttributeGroupPrintout
	AttributeGroupOptionalUniqueId
	Grace *Grace `xml:"grace,omitempty"`
	GroupFullNote
	Tie []*Tie `xml:"tie,omitempty"`
	Cue string `xml:"cue,omitempty"`
	GroupDuration
	Instrument []*Instrument `xml:"instrument,omitempty"`
	GroupEditorialVoice
	Type             *NoteType         `xml:"type,omitempty"`
	Dot              []*EmptyPlacement `xml:"dot,omitempty"`
	Accidental       *Accidental       `xml:"accidental,omitempty"`
	TimeModification *TimeModification `xml:"time-modification,omitempty"`
	Stem             *Stem             `xml:"stem,omitempty"`
	Notehead         *Notehead         `xml:"notehead,omitempty"`
	NoteheadText     *NoteheadText     `xml:"notehead-text,omitempty"`
	GroupStaff
	Beam      *Beam        `xml:"beam,omitempty"`
	Notations []*Notations `xml:"notations,omitempty"`
	Lyric     []*Lyric     `xml:"lyric,omitempty"`
	Play      *Play        `xml:"play,omitempty"`
	Listen    *Listen      `xml:"listen,omitempty"`
}

// NoteType Named source named complex type "note-type"
// The note-type type indicates the graphic note type. Values range from
// 1024th to maxima. The size attribute indicates full, cue, grace-cue, or large size.
// The default is full for regular notes, grace-cue for notes that contain both grace
// and cue elements, and cue for notes that contain either a cue or a grace element,
// but not both.
type NoteType struct {
	Name string     `xml:"-"`
	Size SymbolSize `xml:"size,attr,omitempty"`

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText NoteTypeValue `xml:",chardata"`
}

// Notehead Named source named complex type "notehead"
// The notehead type indicates shapes other than the open and closed
// ovals associated with note durations. The smufl attribute can be used to specify a
// particular notehead, allowing application interoperability without requiring every
// SMuFL glyph to have a MusicXML element equivalent. This attribute can be used either
// with the "other" value, or to refine a specific notehead value such as "cluster".
// Noteheads in the SMuFL Note name noteheads and Note name noteheads supplement ranges
// (U+E150–U+E1AF and U+EEE0–U+EEFF) should not use the smufl attribute or the "other"
// value, but instead use the notehead-text element. For the enclosed shapes, the
// default is to be hollow for half notes and longer, and filled otherwise. The filled
// attribute can be set to change this if needed. If the parentheses attribute is set
// to yes, the notehead is parenthesized. It is no by default.
type Notehead struct {
	Name        string `xml:"-"`
	Filled      YesNo  `xml:"filled,attr,omitempty"`
	Parentheses YesNo  `xml:"parentheses,attr,omitempty"`
	AttributeGroupFont
	AttributeGroupColor
	AttributeGroupSmufl

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// NoteheadText Named source named complex type "notehead-text"
// The notehead-text type represents text that is displayed inside a
// notehead, as is done in some educational music. It is not needed for the numbers
// used in tablature or jianpu notation. The presence of a TAB or jianpu clefs is
// sufficient to indicate that numbers are used. The display-text and accidental-text
// elements allow display of fully formatted text and accidentals.
type NoteheadText struct {
	Name           string            `xml:"-"`
	DisplayText    []*FormattedText  `xml:"display-text,omitempty"`
	AccidentalText []*AccidentalText `xml:"accidental-text,omitempty"`
}

// Ornaments Named source named complex type "ornaments"
// Ornaments can be any of several types, followed optionally by
// accidentals. The accidental-mark element's content is represented the same as an
// accidental element, but with a different name to reflect the different musical
// meaning.
type Ornaments struct {
	Name string `xml:"-"`
	AttributeGroupOptionalUniqueId
	TrillMark            []*EmptyTrillSound    `xml:"trill-mark,omitempty"`
	Turn                 []*HorizontalTurn     `xml:"turn,omitempty"`
	DelayedTurn          []*HorizontalTurn     `xml:"delayed-turn,omitempty"`
	InvertedTurn         []*HorizontalTurn     `xml:"inverted-turn,omitempty"`
	DelayedInvertedTurn  []*HorizontalTurn     `xml:"delayed-inverted-turn,omitempty"`
	VerticalTurn         []*EmptyTrillSound    `xml:"vertical-turn,omitempty"`
	InvertedVerticalTurn []*EmptyTrillSound    `xml:"inverted-vertical-turn,omitempty"`
	Shake                []*EmptyTrillSound    `xml:"shake,omitempty"`
	WavyLine             []*WavyLine           `xml:"wavy-line,omitempty"`
	Mordent              []*Mordent            `xml:"mordent,omitempty"`
	InvertedMordent      []*Mordent            `xml:"inverted-mordent,omitempty"`
	Schleifer            []*EmptyPlacement     `xml:"schleifer,omitempty"`
	Tremolo              []*Tremolo            `xml:"tremolo,omitempty"`
	Haydn                []*EmptyTrillSound    `xml:"haydn,omitempty"`
	OtherOrnament        []*OtherPlacementText `xml:"other-ornament,omitempty"`
	AccidentalMark       []*AccidentalMark     `xml:"accidental-mark,omitempty"`
}

// OtherNotation Named source named complex type "other-notation"
// The other-notation type is used to define any notations not yet in the
// MusicXML format. It handles notations where more specific extension elements such as
// other-dynamics and other-technical are not appropriate. The smufl attribute can be
// used to specify a particular notation, allowing application interoperability without
// requiring every SMuFL glyph to have a MusicXML element equivalent. Using the
// other-notation type without the smufl attribute allows for extended representation,
// though without application interoperability.
type OtherNotation struct {
	Name   string          `xml:"-"`
	Type   StartStopSingle `xml:"type,attr,omitempty"`
	Number int             `xml:"number,attr,omitempty"`
	AttributeGroupPrintObject
	AttributeGroupPrintStyle
	AttributeGroupPlacement
	AttributeGroupSmufl
	AttributeGroupOptionalUniqueId

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// OtherPlacementText Named source named complex type "other-placement-text"
// The other-placement-text type represents a text element with
// print-style, placement, and smufl attribute groups. This type is used by MusicXML
// notation extension elements to allow specification of specific SMuFL glyphs without
// needed to add every glyph as a MusicXML element.
type OtherPlacementText struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement
	AttributeGroupSmufl

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// OtherText Named source named complex type "other-text"
// The other-text type represents a text element with a smufl attribute
// group. This type is used by MusicXML direction extension elements to allow
// specification of specific SMuFL glyphs without needed to add every glyph as a
// MusicXML element.
type OtherText struct {
	Name string `xml:"-"`
	AttributeGroupSmufl

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Pitch Named source named complex type "pitch"
// Pitch is represented as a combination of the step of the diatonic
// scale, the chromatic alteration, and the octave.
type Pitch struct {
	Name   string `xml:"-"`
	Step   Step   `xml:"step,omitempty"`
	Alter  string `xml:"alter,omitempty"`
	Octave int    `xml:"octave,omitempty"`
}

// PlacementText Named source named complex type "placement-text"
// The placement-text type represents a text element with print-style and
// placement attribute groups.
type PlacementText struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Release Named source named complex type "release"
// The release type indicates that a bend is a release rather than a
// normal bend or pre-bend. The offset attribute specifies where the release starts in
// terms of divisions relative to the current note. The first-beat and last-beat
// attributes of the parent bend element are relative to the original note position,
// not this offset value.
type Release struct {
	Name string `xml:"-"`
}

// Rest Named source named complex type "rest"
// The rest element indicates notated rests or silences. Rest elements
// are usually empty, but placement on the staff can be specified using display-step
// and display-octave elements. If the measure attribute is set to yes, this indicates
// this is a complete measure rest.
type Rest struct {
	Name    string `xml:"-"`
	Measure YesNo  `xml:"measure,attr,omitempty"`
	GroupDisplayStepOctave
}

// Slide Named source named complex type "slide"
// Glissando and slide types both indicate rapidly moving from one pitch
// to the other so that individual notes are not discerned. A slide is continuous
// between the two pitches and defaults to a solid line. The optional text for a is
// printed alongside the line.
type Slide struct {
	Name   string    `xml:"-"`
	Type   StartStop `xml:"type,attr,omitempty"`
	Number int       `xml:"number,attr,omitempty"`
	AttributeGroupLineType
	AttributeGroupDashedFormatting
	AttributeGroupPrintStyle
	AttributeGroupBendSound
	AttributeGroupOptionalUniqueId

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Slur Named source named complex type "slur"
// Slur types are empty. Most slurs are represented with two elements:
// one with a start type, and one with a stop type. Slurs can add more elements using a
// continue type. This is typically used to specify the formatting of cross-system
// slurs, or to specify the shape of very complex slurs.
type Slur struct {
	Name   string            `xml:"-"`
	Type   StartStopContinue `xml:"type,attr,omitempty"`
	Number int               `xml:"number,attr,omitempty"`
	AttributeGroupLineType
	AttributeGroupDashedFormatting
	AttributeGroupPosition
	AttributeGroupPlacement
	AttributeGroupOrientation
	AttributeGroupBezier
	AttributeGroupColor
	AttributeGroupOptionalUniqueId
}

// Stem Named source named complex type "stem"
// Stems can be down, up, none, or double. For down and up stems, the
// position attributes can be used to specify stem length. The relative values specify
// the end of the stem relative to the program default. Default values specify an
// absolute end stem position. Negative values of relative-y that would flip a stem
// instead of shortening it are ignored. A stem element associated with a rest refers
// to a stemlet.
type Stem struct {
	Name string `xml:"-"`
	AttributeGroupYPosition
	AttributeGroupColor

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// StrongAccent Named source named complex type "strong-accent"
// The strong-accent type indicates a vertical accent mark. The type
// attribute indicates if the point of the accent is down or up.
type StrongAccent struct {
	Name string `xml:"-"`
}

// StyleText Named source named complex type "style-text"
// The style-text type represents a text element with a print-style
// attribute group.
type StyleText struct {
	Name string `xml:"-"`
	AttributeGroupPrintStyle

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Tap Named source named complex type "tap"
// The tap type indicates a tap on the fretboard. The text content allows
// specification of the notation; + and T are common choices. If the element is empty,
// the hand attribute is used to specify the symbol to use. The hand attribute is
// ignored if the tap glyph is already specified by the text content. If neither text
// content nor the hand attribute are present, the display is application-specific.
type Tap struct {
	Name string `xml:"-"`
	Hand string `xml:"hand,attr,omitempty"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Technical Named source named complex type "technical"
// Technical indications give performance information for individual
// instruments.
type Technical struct {
	Name string `xml:"-"`
	AttributeGroupOptionalUniqueId
	UpBow          []*EmptyPlacement      `xml:"up-bow,omitempty"`
	DownBow        []*EmptyPlacement      `xml:"down-bow,omitempty"`
	Harmonic       []*Harmonic            `xml:"harmonic,omitempty"`
	OpenString     []*EmptyPlacement      `xml:"open-string,omitempty"`
	ThumbPosition  []*EmptyPlacement      `xml:"thumb-position,omitempty"`
	Fingering      []*Fingering           `xml:"fingering,omitempty"`
	Pluck          []*PlacementText       `xml:"pluck,omitempty"`
	DoubleTongue   []*EmptyPlacement      `xml:"double-tongue,omitempty"`
	TripleTongue   []*EmptyPlacement      `xml:"triple-tongue,omitempty"`
	Stopped        []*EmptyPlacementSmufl `xml:"stopped,omitempty"`
	SnapPizzicato  []*EmptyPlacement      `xml:"snap-pizzicato,omitempty"`
	Fret           []*Fret                `xml:"fret,omitempty"`
	String         []*StringType          `xml:"string,omitempty"`
	HammerOn       []*HammerOnPullOff     `xml:"hammer-on,omitempty"`
	PullOff        []*HammerOnPullOff     `xml:"pull-off,omitempty"`
	Bend           []*Bend                `xml:"bend,omitempty"`
	Tap            []*Tap                 `xml:"tap,omitempty"`
	Heel           []*HeelToe             `xml:"heel,omitempty"`
	Toe            []*HeelToe             `xml:"toe,omitempty"`
	Fingernails    []*EmptyPlacement      `xml:"fingernails,omitempty"`
	Hole           []*Hole                `xml:"hole,omitempty"`
	Arrow          []*Arrow               `xml:"arrow,omitempty"`
	Handbell       []*Handbell            `xml:"handbell,omitempty"`
	BrassBend      []*EmptyPlacement      `xml:"brass-bend,omitempty"`
	Flip           []*EmptyPlacement      `xml:"flip,omitempty"`
	Smear          []*EmptyPlacement      `xml:"smear,omitempty"`
	Open           []*EmptyPlacementSmufl `xml:"open,omitempty"`
	HalfMuted      []*EmptyPlacementSmufl `xml:"half-muted,omitempty"`
	HarmonMute     []*HarmonMute          `xml:"harmon-mute,omitempty"`
	Golpe          []*EmptyPlacement      `xml:"golpe,omitempty"`
	OtherTechnical []*OtherPlacementText  `xml:"other-technical,omitempty"`
}

// TextElementData Named source named complex type "text-element-data"
// The text-element-data type represents a syllable or portion of a
// syllable for lyric text underlay. A hyphen in the string content should only be used
// for an actual hyphenated word. Language names for text elements come from ISO 639,
// with optional country subcodes from ISO 3166.
type TextElementData struct {
	Name string `xml:"-"`
	Lang string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
	AttributeGroupFont
	AttributeGroupColor
	AttributeGroupTextDecoration
	AttributeGroupTextRotation
	AttributeGroupLetterSpacing
	AttributeGroupTextDirection

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Tie Named source named complex type "tie"
// The tie element indicates that a tie begins or ends with this note. If
// the tie element applies only particular times through a repeat, the time-only
// attribute indicates which times to apply it. The tie element indicates sound; the
// tied element indicates notation.
type Tie struct {
	Name     string    `xml:"-"`
	Type     StartStop `xml:"type,attr,omitempty"`
	TimeOnly TimeOnly  `xml:"time-only,attr,omitempty"`
}

// Tied Named source named complex type "tied"
// The tied element represents the notated tie. The tie element
// represents the tie sound. The number attribute is rarely needed to disambiguate
// ties, since note pitches will usually suffice. The attribute is implied rather than
// defaulting to 1 as with most elements. It is available for use in more complex tied
// notation situations. Ties that join two notes of the same pitch together should be
// represented with a tied element on the first note with type="start" and a tied
// element on the second note with type="stop". This can also be done if the two notes
// being tied are enharmonically equivalent, but have different step values. It is not
// recommended to use tied elements to join two notes with enharmonically inequivalent
// pitches. Ties that indicate that an instrument should be undamped are specified with
// a single tied element with type="let-ring". Ties that are visually attached to only
// one note, other than undamped ties, should be specified with two tied elements on
// the same note, first type="start" then type="stop". This can be used to represent
// ties into or out of repeated sections or codas.
type Tied struct {
	Name   string `xml:"-"`
	Type   string `xml:"type,attr,omitempty"`
	Number int    `xml:"number,attr,omitempty"`
	AttributeGroupLineType
	AttributeGroupDashedFormatting
	AttributeGroupPosition
	AttributeGroupPlacement
	AttributeGroupOrientation
	AttributeGroupBezier
	AttributeGroupColor
	AttributeGroupOptionalUniqueId
}

// TimeModification Named source named complex type "time-modification"
// Time modification indicates tuplets, double-note tremolos, and other
// durational changes. A time-modification element shows how the cumulative, sounding
// effect of tuplets and double-note tremolos compare to the written note type
// represented by the type and dot elements. Nested tuplets and other notations that
// use more detailed information need both the time-modification and tuplet elements to
// be represented accurately.
type TimeModification struct {
	Name        string        `xml:"-"`
	ActualNotes int           `xml:"actual-notes,omitempty"`
	NormalNotes int           `xml:"normal-notes,omitempty"`
	NormalType  NoteTypeValue `xml:"normal-type,omitempty"`
	NormalDot   string        `xml:"normal-dot,omitempty"`
}

// Tremolo Named source named complex type "tremolo"
// The tremolo ornament can be used to indicate single-note, double-note,
// or unmeasured tremolos. Single-note tremolos use the single type, double-note
// tremolos use the start and stop types, and unmeasured tremolos use the unmeasured
// type. The default is "single" for compatibility with Version 1.1. The text of the
// element indicates the number of tremolo marks and is an integer from 0 to 8. Note
// that the number of attached beams is not included in this value, but is represented
// separately using the beam element. The value should be 0 for unmeasured tremolos.
// When using double-note tremolos, the duration of each note in the tremolo should
// correspond to half of the notated type value. A time-modification element should
// also be added with an actual-notes value of 2 and a normal-notes value of 1. If used
// within a tuplet, this 2/1 ratio should be multiplied by the existing tuplet ratio.
// The smufl attribute specifies the glyph to use from the SMuFL Tremolos range for an
// unmeasured tremolo. It is ignored for other tremolo types. The SMuFL buzzRoll glyph
// is used by default if the attribute is missing. Using repeater beams for indicating
// tremolos is deprecated as of MusicXML 3.0.
type Tremolo struct {
	Name string `xml:"-"`
	Type string `xml:"type,attr,omitempty"`
	AttributeGroupPrintStyle
	AttributeGroupPlacement
	AttributeGroupSmufl

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText int `xml:",chardata"`
}

// Tuplet Named source named complex type "tuplet"
// A tuplet element is present when a tuplet is to be displayed
// graphically, in addition to the sound data provided by the time-modification
// elements. The number attribute is used to distinguish nested tuplets. The bracket
// attribute is used to indicate the presence of a bracket. If unspecified, the results
// are implementation-dependent. The line-shape attribute is used to specify whether
// the bracket is straight or in the older curved or slurred style. It is straight by
// default. Whereas a time-modification element shows how the cumulative, sounding
// effect of tuplets and double-note tremolos compare to the written note type, the
// tuplet element describes how this is displayed. The tuplet element also provides
// more detailed representation information than the time-modification element, and is
// needed to represent nested tuplets and other complex tuplets accurately. The
// show-number attribute is used to display either the number of actual notes, the
// number of both actual and normal notes, or neither. It is actual by default. The
// show-type attribute is used to display either the actual type, both the actual and
// normal types, or neither. It is none by default.
type Tuplet struct {
	Name       string     `xml:"-"`
	Type       StartStop  `xml:"type,attr,omitempty"`
	Number     int        `xml:"number,attr,omitempty"`
	Bracket    YesNo      `xml:"bracket,attr,omitempty"`
	ShowNumber string     `xml:"show-number,attr,omitempty"`
	ShowType   ShowTuplet `xml:"show-type,attr,omitempty"`
	AttributeGroupLineShape
	AttributeGroupPosition
	AttributeGroupPlacement
	AttributeGroupOptionalUniqueId
	TupletActual *TupletPortion `xml:"tuplet-actual,omitempty"`
	TupletNormal *TupletPortion `xml:"tuplet-normal,omitempty"`
}

// TupletDot Named source named complex type "tuplet-dot"
// The tuplet-dot type is used to specify dotted tuplet types.
type TupletDot struct {
	Name string `xml:"-"`
	AttributeGroupFont
	AttributeGroupColor
}

// TupletNumber Named source named complex type "tuplet-number"
// The tuplet-number type indicates the number of notes for this portion
// of the tuplet.
type TupletNumber struct {
	Name string `xml:"-"`
	AttributeGroupFont
	AttributeGroupColor

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText int `xml:",chardata"`
}

// TupletPortion Named source named complex type "tuplet-portion"
// The tuplet-portion type provides optional full control over tuplet
// specifications. It allows the number and note type (including dots) to be set for
// the actual and normal portions of a single tuplet. If any of these elements are
// absent, their values are based on the time-modification element.
type TupletPortion struct {
	Name         string        `xml:"-"`
	TupletNumber *TupletNumber `xml:"tuplet-number,omitempty"`
	TupletType   *TupletType   `xml:"tuplet-type,omitempty"`
	TupletDot    []*TupletDot  `xml:"tuplet-dot,omitempty"`
}

// TupletType Named source named complex type "tuplet-type"
// The tuplet-type type indicates the graphical note type of the notes
// for this portion of the tuplet.
type TupletType struct {
	Name string `xml:"-"`
	AttributeGroupFont
	AttributeGroupColor

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText NoteTypeValue `xml:",chardata"`
}

// Unpitched Named source named complex type "unpitched"
// The unpitched type represents musical elements that are notated on the
// staff but lack definite pitch, such as unpitched percussion and speaking voice. If
// the child elements are not present, the note is placed on the middle line of the
// staff. This is generally used with a one-line staff. Notes in percussion clef should
// always use an unpitched element rather than a pitch element.
type Unpitched struct {
	Name string `xml:"-"`
	GroupDisplayStepOctave
}

// Wait Named source named complex type "wait"
// The wait type specifies a point where the accompaniment should wait
// for a performer event before continuing. This typically happens at the start of new
// sections or after a held note or indeterminate music. These waiting points cannot
// always be inferred reliably from the contents of the displayed score. The optional
// player and time-only attributes restrict the type to apply to a single player or set
// of times through a repeated section, respectively.
type Wait struct {
	Name     string   `xml:"-"`
	Player   string   `xml:"player,attr,omitempty"`
	TimeOnly TimeOnly `xml:"time-only,attr,omitempty"`
}

// Credit Named source named complex type "credit"
// The credit type represents the appearance of the title, composer,
// arranger, lyricist, copyright, dedication, and other text, symbols, and graphics
// that commonly appear on the first page of a score. The credit-words, credit-symbol,
// and credit-image elements are similar to the words, symbol, and image elements for
// directions. However, since the credit is not part of a measure, the default-x and
// default-y attributes adjust the origin relative to the bottom left-hand corner of
// the page. The enclosure for credit-words and credit-symbol is none by default. By
// default, a series of credit-words and credit-symbol elements within a single credit
// element follow one another in sequence visually. Non-positional formatting
// attributes are carried over from the previous element by default. The page attribute
// for the credit element specifies the page number where the credit should appear.
// This is an integer value that starts with 1 for the first page. Its value is 1 by
// default. Since credits occur before the music, these page numbers do not refer to
// the page numbering specified by the print element's page-number attribute. The
// credit-type element indicates the purpose behind a credit. Multiple types of data
// may be combined in a single credit, so multiple elements may be used. Standard
// values include page number, title, subtitle, composer, arranger, lyricist, rights,
// and part name.
type Credit struct {
	Name string `xml:"-"`
	Page int    `xml:"page,attr,omitempty"`
	AttributeGroupOptionalUniqueId
	CreditType   string               `xml:"credit-type,omitempty"`
	CreditImage  *Image               `xml:"credit-image,omitempty"`
	Link         []*Link              `xml:"link,omitempty"`
	Bookmark     []*Bookmark          `xml:"bookmark,omitempty"`
	CreditWords  []*FormattedTextId   `xml:"credit-words,omitempty"`
	CreditSymbol []*FormattedSymbolId `xml:"credit-symbol,omitempty"`
}

// Defaults Named source named complex type "defaults"
// The defaults type specifies score-wide defaults for scaling; whether
// or not the file is a concert score; layout; and default values for the music font,
// word font, lyric font, and lyric language. Except for the concert-score element, if
// any defaults are missing, the choice of what to use is determined by the
// application.
type Defaults struct {
	Name         string   `xml:"-"`
	Scaling      *Scaling `xml:"scaling,omitempty"`
	ConcertScore string   `xml:"concert-score,omitempty"`
	GroupLayout
	Appearance    *Appearance      `xml:"appearance,omitempty"`
	MusicFont     *EmptyFont       `xml:"music-font,omitempty"`
	WordFont      *EmptyFont       `xml:"word-font,omitempty"`
	LyricFont     []*LyricFont     `xml:"lyric-font,omitempty"`
	LyricLanguage []*LyricLanguage `xml:"lyric-language,omitempty"`
}

// EmptyFont Named source named complex type "empty-font"
// The empty-font type represents an empty element with font attributes.
type EmptyFont struct {
	Name string `xml:"-"`
	AttributeGroupFont
}

// GroupBarline Named source named complex type "group-barline"
// The group-barline type indicates if the group should have common
// barlines.
type GroupBarline struct {
	Name string `xml:"-"`
	AttributeGroupColor

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// GroupName Named source named complex type "group-name"
// The group-name type describes the name or abbreviation of a part-group
// element. Formatting attributes in the group-name type are deprecated in Version 2.0
// in favor of the new group-name-display and group-abbreviation-display elements.
type GroupName struct {
	Name string `xml:"-"`
	AttributeGroupGroupNameText

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// GroupSymbol Named source named complex type "group-symbol"
// The group-symbol type indicates how the symbol for a group is
// indicated in the score. It is none if not specified.
type GroupSymbol struct {
	Name string `xml:"-"`
	AttributeGroupPosition
	AttributeGroupColor

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// InstrumentLink Named source named complex type "instrument-link"
// Multiple part-link elements can link a condensed part within a score
// file to multiple MusicXML parts files. For example, a "Clarinet 1 and 2" part in a
// score file could link to separate "Clarinet 1" and "Clarinet 2" part files. The
// instrument-link type distinguish which of the score-instruments within a score-part
// are in which part file. The instrument-link id attribute refers to a
// score-instrument id attribute.
type InstrumentLink struct {
	Name string `xml:"-"`
	Id   string `xml:"id,attr,omitempty"`
}

// LyricFont Named source named complex type "lyric-font"
// The lyric-font type specifies the default font for a particular name
// and number of lyric.
type LyricFont struct {
	Name    string `xml:"-"`
	Number  string `xml:"number,attr,omitempty"`
	NameXSD string `xml:"name,attr,omitempty"`
	AttributeGroupFont
}

// LyricLanguage Named source named complex type "lyric-language"
// The lyric-language type specifies the default language for a
// particular name and number of lyric.
type LyricLanguage struct {
	Name    string `xml:"-"`
	Number  string `xml:"number,attr,omitempty"`
	NameXSD string `xml:"name,attr,omitempty"`
	Lang    string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
}

// Opus Named source named complex type "opus"
// The opus type represents a link to a MusicXML opus document that
// composes multiple MusicXML scores into a collection.
type Opus struct {
	Name string `xml:"-"`
	AttributeGroupLinkAttributes
}

// PartGroup Named source named complex type "part-group"
// The part-group element indicates groupings of parts in the score,
// usually indicated by braces and brackets. Braces that are used for multi-staff parts
// should be defined in the attributes element for that part. The part-group start
// element appears before the first score-part in the group. The part-group stop
// element appears after the last score-part in the group. The number attribute is used
// to distinguish overlapping and nested part-groups, not the sequence of groups. As
// with parts, groups can have a name and abbreviation. Values for the child elements
// are ignored at the stop of a group. A part-group element is not needed for a single
// multi-staff part. By default, multi-staff parts include a brace symbol and (if
// appropriate given the bar-style) common barlines. The symbol formatting for a
// multi-staff part can be more fully specified using the part-symbol element.
type PartGroup struct {
	Name                     string        `xml:"-"`
	Type                     StartStop     `xml:"type,attr,omitempty"`
	Number                   string        `xml:"number,attr,omitempty"`
	GroupName                *GroupName    `xml:"group-name,omitempty"`
	GroupNameDisplay         *NameDisplay  `xml:"group-name-display,omitempty"`
	GroupAbbreviation        *GroupName    `xml:"group-abbreviation,omitempty"`
	GroupAbbreviationDisplay *NameDisplay  `xml:"group-abbreviation-display,omitempty"`
	GroupSymbol              *GroupSymbol  `xml:"group-symbol,omitempty"`
	GroupBarline             *GroupBarline `xml:"group-barline,omitempty"`
	GroupTime                string        `xml:"group-time,omitempty"`
	GroupEditorial
}

// PartLink Named source named complex type "part-link"
// The part-link type allows MusicXML data for both score and parts to be
// contained within a single compressed MusicXML file. It links a score-part from a
// score document to MusicXML documents that contain parts data. In the case of a
// single compressed MusicXML file, the link href values are paths that are relative to
// the root folder of the zip file.
type PartLink struct {
	Name string `xml:"-"`
	AttributeGroupLinkAttributes
	InstrumentLink []*InstrumentLink `xml:"instrument-link,omitempty"`
	GroupLink      string            `xml:"group-link,omitempty"`
}

// PartList Named source named complex type "part-list"
// The part-list identifies the different musical parts in this document.
// Each part has an ID that is used later within the musical data. Since parts may be
// encoded separately and combined later, identification elements are present at both
// the score and score-part levels. There must be at least one score-part, combined as
// desired with part-group elements that indicate braces and brackets. Parts are
// ordered from top to bottom in a score based on the order in which they appear in the
// part-list.
type PartList struct {
	Name string `xml:"-"`
	GroupPartGroup
	GroupScorePart
}

// PartName Named source named complex type "part-name"
// The part-name type describes the name or abbreviation of a score-part
// element. Formatting attributes for the part-name element are deprecated in Version
// 2.0 in favor of the new part-name-display and part-abbreviation-display elements.
type PartName struct {
	Name string `xml:"-"`
	AttributeGroupPartNameText

	// in case the extension has base type xs:string, one has to had the chardata stuff
	EnclosedText string `xml:",chardata"`
}

// Player Named source named complex type "player"
// The player type allows for multiple players per score-part for use in
// listening applications. One player may play multiple instruments, while a single
// instrument may include multiple players in divisi sections.
type Player struct {
	Name       string `xml:"-"`
	Id         string `xml:"id,attr,omitempty"`
	PlayerName string `xml:"player-name,omitempty"`
}

// ScoreInstrument Named source named complex type "score-instrument"
// The score-instrument type represents a single instrument within a
// score-part. As with the score-part type, each score-instrument has a required ID
// attribute, a name, and an optional abbreviation. A score-instrument type is also
// required if the score specifies MIDI 1.0 channels, banks, or programs. An initial
// midi-instrument assignment can also be made here. MusicXML software should be able
// to automatically assign reasonable channels and instruments without these elements
// in simple cases, such as where part names match General MIDI instrument names. The
// score-instrument element can also distinguish multiple instruments of the same type
// that are on the same part, such as Clarinet 1 and Clarinet 2 instruments within a
// Clarinets 1 and 2 part.
type ScoreInstrument struct {
	Name                   string `xml:"-"`
	Id                     string `xml:"id,attr,omitempty"`
	InstrumentName         string `xml:"instrument-name,omitempty"`
	InstrumentAbbreviation string `xml:"instrument-abbreviation,omitempty"`
	GroupVirtualInstrumentData
}

// ScorePart Named source named complex type "score-part"
// The score-part type collects part-wide information for each part in a
// score. Often, each MusicXML part corresponds to a track in a Standard MIDI Format 1
// file. In this case, the midi-device element is used to make a MIDI device or port
// assignment for the given track or specific MIDI instruments. Initial midi-instrument
// assignments may be made here as well. The score-instrument elements are used when
// there are multiple instruments per track.
type ScorePart struct {
	Name                    string             `xml:"-"`
	Id                      string             `xml:"id,attr,omitempty"`
	Identification          *Identification    `xml:"identification,omitempty"`
	PartLink                []*PartLink        `xml:"part-link,omitempty"`
	PartName                *PartName          `xml:"part-name,omitempty"`
	PartNameDisplay         *NameDisplay       `xml:"part-name-display,omitempty"`
	PartAbbreviation        *PartName          `xml:"part-abbreviation,omitempty"`
	PartAbbreviationDisplay *NameDisplay       `xml:"part-abbreviation-display,omitempty"`
	Group                   string             `xml:"group,omitempty"`
	ScoreInstrument         []*ScoreInstrument `xml:"score-instrument,omitempty"`
	Player                  []*Player          `xml:"player,omitempty"`
	MidiDevice              []*MidiDevice      `xml:"midi-device,omitempty"`
	MidiInstrument          []*MidiInstrument  `xml:"midi-instrument,omitempty"`
}

// VirtualInstrument Named source named complex type "virtual-instrument"
// The virtual-instrument element defines a specific virtual instrument
// used for an instrument sound.
type VirtualInstrument struct {
	Name           string `xml:"-"`
	VirtualLibrary string `xml:"virtual-library,omitempty"`
	VirtualName    string `xml:"virtual-name,omitempty"`
}

// Work Named source named complex type "work"
// Works are optionally identified by number and title. The work type
// also may indicate a link to the opus document that composes multiple scores into a
// collection.
type Work struct {
	Name       string `xml:"-"`
	WorkNumber string `xml:"work-number,omitempty"`
	WorkTitle  string `xml:"work-title,omitempty"`
	Opus       *Opus  `xml:"opus,omitempty"`
}

// GroupEditorial UnNamed source named group "editorial"
type GroupEditorial struct {
	GroupFootnote
	GroupLevel
}

// GroupEditorialVoice UnNamed source named group "editorial-voice"
type GroupEditorialVoice struct {
	GroupFootnote
	GroupLevel
	GroupVoice
}

// GroupEditorialVoiceDirection UnNamed source named group "editorial-voice-direction"
type GroupEditorialVoiceDirection struct {
	GroupFootnote
	GroupLevel
	GroupVoice
}

// GroupFootnote UnNamed source named group "footnote"
type GroupFootnote struct {
	Footnote *FormattedText `xml:"footnote,omitempty"`
}

// GroupLevel UnNamed source named group "level"
type GroupLevel struct {
	Level *Level `xml:"level,omitempty"`
}

// GroupStaff UnNamed source named group "staff"
type GroupStaff struct {
	Staff int `xml:"staff,omitempty"`
}

// GroupTuning UnNamed source named group "tuning"
type GroupTuning struct {
	TuningStep   Step   `xml:"tuning-step,omitempty"`
	TuningAlter  string `xml:"tuning-alter,omitempty"`
	TuningOctave int    `xml:"tuning-octave,omitempty"`
}

// GroupVirtualInstrumentData UnNamed source named group "virtual-instrument-data"
type GroupVirtualInstrumentData struct {
	InstrumentSound   string             `xml:"instrument-sound,omitempty"`
	Solo              string             `xml:"solo,omitempty"`
	Ensemble          string             `xml:"ensemble,omitempty"`
	VirtualInstrument *VirtualInstrument `xml:"virtual-instrument,omitempty"`
}

// GroupVoice UnNamed source named group "voice"
type GroupVoice struct {
	Voice string `xml:"voice,omitempty"`
}

// GroupClef UnNamed source named group "clef"
type GroupClef struct {
	Sign             string `xml:"sign,omitempty"`
	Line             int    `xml:"line,omitempty"`
	ClefOctaveChange int    `xml:"clef-octave-change,omitempty"`
}

// GroupNonTraditionalKey UnNamed source named group "non-traditional-key"
type GroupNonTraditionalKey struct {
	KeyStep       Step           `xml:"key-step,omitempty"`
	KeyAlter      string         `xml:"key-alter,omitempty"`
	KeyAccidental *KeyAccidental `xml:"key-accidental,omitempty"`
}

// GroupSlash UnNamed source named group "slash"
type GroupSlash struct {
	SlashType   NoteTypeValue `xml:"slash-type,omitempty"`
	SlashDot    string        `xml:"slash-dot,omitempty"`
	ExceptVoice string        `xml:"except-voice,omitempty"`
}

// GroupTimeSignature UnNamed source named group "time-signature"
type GroupTimeSignature struct {
	Beats    string `xml:"beats,omitempty"`
	BeatType int `xml:"beat-type,omitempty"`
}

// GroupTraditionalKey UnNamed source named group "traditional-key"
type GroupTraditionalKey struct {
	Cancel *Cancel `xml:"cancel,omitempty"`
	Fifths int     `xml:"fifths,omitempty"`
	Mode   string  `xml:"mode,omitempty"`
}

// GroupTranspose UnNamed source named group "transpose"
type GroupTranspose struct {
	Diatonic     int     `xml:"diatonic,omitempty"`
	Chromatic    string  `xml:"chromatic,omitempty"`
	OctaveChange int     `xml:"octave-change,omitempty"`
	Double       float64 `xml:"double,omitempty"`
}

// GroupBeatUnit UnNamed source named group "beat-unit"
type GroupBeatUnit struct {
	BeatUnit    NoteTypeValue `xml:"beat-unit,omitempty"`
	BeatUnitDot string        `xml:"beat-unit-dot,omitempty"`
}

// GroupHarmonyChord UnNamed source named group "harmony-chord"
type GroupHarmonyChord struct {
	Root      *Root      `xml:"root,omitempty"`
	Numeral   *Numeral   `xml:"numeral,omitempty"`
	Function  *StyleText `xml:"function,omitempty"`
	Kind      *Kind      `xml:"kind,omitempty"`
	Inversion *Inversion `xml:"inversion,omitempty"`
	Bass      *Bass      `xml:"bass,omitempty"`
	Degree    []*Degree  `xml:"degree,omitempty"`
}

// GroupAllMargins UnNamed source named group "all-margins"
type GroupAllMargins struct {
	GroupLeftRightMargins
	TopMargin    string `xml:"top-margin,omitempty"`
	BottomMargin string `xml:"bottom-margin,omitempty"`
}

// GroupLayout UnNamed source named group "layout"
type GroupLayout struct {
	PageLayout   *PageLayout    `xml:"page-layout,omitempty"`
	SystemLayout *SystemLayout  `xml:"system-layout,omitempty"`
	StaffLayout  []*StaffLayout `xml:"staff-layout,omitempty"`
}

// GroupLeftRightMargins UnNamed source named group "left-right-margins"
type GroupLeftRightMargins struct {
	LeftMargin  string `xml:"left-margin,omitempty"`
	RightMargin string `xml:"right-margin,omitempty"`
}

// GroupDuration UnNamed source named group "duration"
type GroupDuration struct {
	Duration string `xml:"duration,omitempty"`
}

// GroupDisplayStepOctave UnNamed source named group "display-step-octave"
type GroupDisplayStepOctave struct {
	DisplayStep   Step `xml:"display-step,omitempty"`
	DisplayOctave int  `xml:"display-octave,omitempty"`
}

// GroupFullNote UnNamed source named group "full-note"
type GroupFullNote struct {
	Chord     string     `xml:"chord,omitempty"`
	Pitch     *Pitch     `xml:"pitch,omitempty"`
	Unpitched *Unpitched `xml:"unpitched,omitempty"`
	Rest      *Rest      `xml:"rest,omitempty"`
}

// GroupMusicData UnNamed source named group "music-data"
type GroupMusicData struct {
	Element []GroupMusicElement `xml:",any,omitempty"`
	/*
	Note        []*Note        `xml:"note,omitempty"`
	Backup      []*Backup      `xml:"backup,omitempty"`
	Forward     []*Forward     `xml:"forward,omitempty"`
	Direction   []*Direction   `xml:"direction,omitempty"`
	Attributes  []*Attributes  `xml:"attributes,omitempty"`
	Harmony     []*Harmony     `xml:"harmony,omitempty"`
	FiguredBass []*FiguredBass `xml:"figured-bass,omitempty"`
	Print       []*Print       `xml:"print,omitempty"`
	Sound       []*Sound       `xml:"sound,omitempty"`
	Listening   []*Listening   `xml:"listening,omitempty"`
	Barline     []*Barline     `xml:"barline,omitempty"`
	Grouping    []*Grouping    `xml:"grouping,omitempty"`
	Link        []*Link        `xml:"link,omitempty"`
	Bookmark    []*Bookmark    `xml:"bookmark,omitempty"`
	*/
}

type GroupMusicElement struct {
	Value any
	/*
	one of:

	Note        *Note        `xml:"note,omitempty"`
	Backup      *Backup      `xml:"backup,omitempty"`
	Forward     *Forward     `xml:"forward,omitempty"`
	Direction   *Direction   `xml:"direction,omitempty"`
	Attributes  *Attributes  `xml:"attributes,omitempty"`
	Harmony     *Harmony     `xml:"harmony,omitempty"`
	FiguredBass *FiguredBass `xml:"figured-bass,omitempty"`
	Print       *Print       `xml:"print,omitempty"`
	Sound       *Sound       `xml:"sound,omitempty"`
	Listening   *Listening   `xml:"listening,omitempty"`
	Barline     *Barline     `xml:"barline,omitempty"`
	Grouping    *Grouping    `xml:"grouping,omitempty"`
	Link        *Link        `xml:"link,omitempty"`
	Bookmark    *Bookmark    `xml:"bookmark,omitempty"`
	*/
}

// Implements encoding.xml.Unmarshaler interface
func (element *GroupMusicElement) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) (err error) {
	switch start.Name.Local {
	case "note":
		element.Value = new(Note)
		err = decoder.DecodeElement(element.Value, &start)
	case "backup":
		element.Value = new(Backup)
		err = decoder.DecodeElement(element.Value, &start)
	case "forward":
		element.Value = new(Forward)
		err = decoder.DecodeElement(element.Value, &start)
	case "direction":
		element.Value = new(Direction)
		err = decoder.DecodeElement(element.Value, &start)
	case "attributes":
		element.Value = new(Attributes)
		err = decoder.DecodeElement(element.Value, &start)
	case "harmony":
		element.Value = new(Harmony)
		err = decoder.DecodeElement(element.Value, &start)
	case "figured-bass":
		element.Value = new(FiguredBass)
		err = decoder.DecodeElement(element.Value, &start)
	case "print":
		element.Value = new(Print)
		err = decoder.DecodeElement(element.Value, &start)
	case "sound":
		element.Value = new(Sound)
		err = decoder.DecodeElement(element.Value, &start)
	case "listening":
		element.Value = new(Listening)
		err = decoder.DecodeElement(element.Value, &start)
	case "barline":
		element.Value = new(Barline)
		err = decoder.DecodeElement(element.Value, &start)
	case "grouping":
		element.Value = new(Grouping)
		err = decoder.DecodeElement(element.Value, &start)
	case "link":
		element.Value = new(Link)
		err = decoder.DecodeElement(element.Value, &start)
	case "bookmark":
		element.Value = new(Bookmark)
		err = decoder.DecodeElement(element.Value, &start)
	default:
		return fmt.Errorf("unknown tag %q", start.Name)
	}

	return err
}

// GroupPartGroup UnNamed source named group "part-group"
type GroupPartGroup struct {
	PartGroup *PartGroup `xml:"part-group,omitempty"`
}

// GroupScoreHeader UnNamed source named group "score-header"
type GroupScoreHeader struct {
	Work           *Work           `xml:"work,omitempty"`
	MovementNumber string          `xml:"movement-number,omitempty"`
	MovementTitle  string          `xml:"movement-title,omitempty"`
	Identification *Identification `xml:"identification,omitempty"`
	Defaults       *Defaults       `xml:"defaults,omitempty"`
	Credit         []*Credit       `xml:"credit,omitempty"`
	PartList       *PartList       `xml:"part-list,omitempty"`
}

// GroupScorePart UnNamed source named group "score-part"
type GroupScorePart struct {
	ScorePart []*ScorePart `xml:"score-part,omitempty"`
}

// ScorePartwise Named source element score-partwise within root schema
// The score-partwise element is the root element for a partwise MusicXML
// score. It includes a score-header group followed by a series of parts with measures
// inside. The document-attributes attribute group includes the version attribute.
type ScorePartwise struct {
	Name    string   `xml:"-"`
	XMLName xml.Name `xml:"score-partwise"`
	ScorePartwiseInfo
}

// ScorePartwiseInfo Named source within outer element "score-partwise"
type ScorePartwiseInfo struct {
	Name string `xml:"-"`
	AttributeGroupDocumentAttributes
	GroupScoreHeader
	Part []*Part `xml:"part,omitempty"`
}

// Part Named source within outer element "part"
type Part struct {
	Name string `xml:"-"`
	AttributeGroupPartAttributes
	Measure []*Measure `xml:"measure,omitempty"`
}

// Measure Named source within outer element "measure"
type Measure struct {
	Name string `xml:"-"`
	AttributeGroupMeasureAttributes
	GroupMusicData
}

// ScoreTimewise Named source element score-timewise within root schema
// The score-timewise element is the root element for a timewise MusicXML
// score. It includes a score-header group followed by a series of measures with parts
// inside. The document-attributes attribute group includes the version attribute.
type ScoreTimewise struct {
	Name    string   `xml:"-"`
	XMLName xml.Name `xml:"score-timewise"`
	ScoreTimewiseInfo
}

// ScoreTimewiseInfo Named source within outer element "score-timewise"
type ScoreTimewiseInfo struct {
	Name string `xml:"-"`
	AttributeGroupDocumentAttributes
	GroupScoreHeader
	Measure []*TimewiseMeasure `xml:"measure,omitempty"`
}

// TimewiseMeasure Named source within outer element "measure"
type TimewiseMeasure struct {
	Name string `xml:"-"`
	AttributeGroupMeasureAttributes
	Part []*TimewisePart `xml:"part,omitempty"`
}

// TimewisePart Named source within outer element "part"
type TimewisePart struct {
	Name string `xml:"-"`
	AttributeGroupPartAttributes
	GroupMusicData
}
