package music

import "strings"

// Optimization: pre-computed lookup maps
var (
	noteByNameMap = make(map[string]*Note)
	chordPatternMap = make(map[string]string) // pattern -> intervals
	noteConversionMap = make(map[string]string) // accidental -> sharp
)

var noteDB = []Note{
	Note{MidiValue: -12, NameSharp: "C-2", Frequency: 4.0879, NamesOther: []string{"c-2"}},
	Note{MidiValue: -11, NameSharp: "C#-2", Frequency: 4.331, NamesOther: []string{"cs-2", "db-2"}},
	Note{MidiValue: -10, NameSharp: "D-2", Frequency: 4.5885, NamesOther: []string{"d-2"}},
	Note{MidiValue: -9, NameSharp: "D#-2", Frequency: 4.8614, NamesOther: []string{"ds-2", "eb-2"}},
	Note{MidiValue: -8, NameSharp: "E-2", Frequency: 5.1504, NamesOther: []string{"e-2", "fb-2"}},
	Note{MidiValue: -7, NameSharp: "F-2", Frequency: 5.4567, NamesOther: []string{"f-2"}},
	Note{MidiValue: -6, NameSharp: "F#-2", Frequency: 5.7812, NamesOther: []string{"fs-2", "gb-2"}},
	Note{MidiValue: -5, NameSharp: "G-2", Frequency: 6.1249, NamesOther: []string{"g-2"}},
	Note{MidiValue: -4, NameSharp: "G#-2", Frequency: 6.4891, NamesOther: []string{"gs-2", "ab-2"}},
	Note{MidiValue: -3, NameSharp: "A-2", Frequency: 6.875, NamesOther: []string{"a-2"}},
	Note{MidiValue: -2, NameSharp: "A#-2", Frequency: 7.2838, NamesOther: []string{"as-2", "bb-2"}},
	Note{MidiValue: -1, NameSharp: "B-2", Frequency: 7.7169, NamesOther: []string{"b-2", "cb-2"}},
	Note{MidiValue: 0, NameSharp: "C-1", Frequency: 8.1758, NamesOther: []string{"c-1"}},
	Note{MidiValue: 1, NameSharp: "C#-1", Frequency: 8.662, NamesOther: []string{"cs-1", "db-1"}},
	Note{MidiValue: 2, NameSharp: "D-1", Frequency: 9.177, NamesOther: []string{"d-1"}},
	Note{MidiValue: 3, NameSharp: "D#-1", Frequency: 9.7227, NamesOther: []string{"ds-1", "eb-1"}},
	Note{MidiValue: 4, NameSharp: "E-1", Frequency: 10.3009, NamesOther: []string{"e-1", "fb-1"}},
	Note{MidiValue: 5, NameSharp: "F-1", Frequency: 10.9134, NamesOther: []string{"f-1"}},
	Note{MidiValue: 6, NameSharp: "F#-1", Frequency: 11.5623, NamesOther: []string{"fs-1", "gb-1"}},
	Note{MidiValue: 7, NameSharp: "G-1", Frequency: 12.2499, NamesOther: []string{"g-1"}},
	Note{MidiValue: 8, NameSharp: "G#-1", Frequency: 12.9783, NamesOther: []string{"gs-1", "ab-1"}},
	Note{MidiValue: 9, NameSharp: "A-1", Frequency: 13.75, NamesOther: []string{"a-1"}},
	Note{MidiValue: 10, NameSharp: "A#-1", Frequency: 14.5676, NamesOther: []string{"as-1", "bb-1"}},
	Note{MidiValue: 11, NameSharp: "B-1", Frequency: 15.4339, NamesOther: []string{"b-1", "cb-1"}},
	Note{MidiValue: 12, NameSharp: "C0", Frequency: 16.351, NamesOther: []string{"c0"}},
	Note{MidiValue: 13, NameSharp: "C#0", Frequency: 17.324, NamesOther: []string{"cs0", "db0"}},
	Note{MidiValue: 14, NameSharp: "D0", Frequency: 18.354, NamesOther: []string{"d0"}},
	Note{MidiValue: 15, NameSharp: "D#0", Frequency: 19.445, NamesOther: []string{"ds0", "eb0"}},
	Note{MidiValue: 16, NameSharp: "E0", Frequency: 20.601, NamesOther: []string{"e0", "fb0"}},
	Note{MidiValue: 17, NameSharp: "F0", Frequency: 21.827, NamesOther: []string{"f0"}},
	Note{MidiValue: 18, NameSharp: "F#0", Frequency: 23.124, NamesOther: []string{"fs0", "gb0"}},
	Note{MidiValue: 19, NameSharp: "G0", Frequency: 24.499, NamesOther: []string{"g0"}},
	Note{MidiValue: 20, NameSharp: "G#0", Frequency: 25.956, NamesOther: []string{"gs0", "ab0"}},
	Note{MidiValue: 21, NameSharp: "A0", Frequency: 27.5, NamesOther: []string{"a0"}},
	Note{MidiValue: 22, NameSharp: "A#0", Frequency: 29.135, NamesOther: []string{"as0", "bb0"}},
	Note{MidiValue: 23, NameSharp: "B0", Frequency: 30.868, NamesOther: []string{"b0", "cb0"}},
	Note{MidiValue: 24, NameSharp: "C1", Frequency: 32.703, NamesOther: []string{"c1"}},
	Note{MidiValue: 25, NameSharp: "C#1", Frequency: 34.648, NamesOther: []string{"cs1", "db1"}},
	Note{MidiValue: 26, NameSharp: "D1", Frequency: 36.708, NamesOther: []string{"d1"}},
	Note{MidiValue: 27, NameSharp: "D#1", Frequency: 38.891, NamesOther: []string{"ds1", "eb1"}},
	Note{MidiValue: 28, NameSharp: "E1", Frequency: 41.203, NamesOther: []string{"e1", "fb1"}},
	Note{MidiValue: 29, NameSharp: "F1", Frequency: 43.654, NamesOther: []string{"f1"}},
	Note{MidiValue: 30, NameSharp: "F#1", Frequency: 46.249, NamesOther: []string{"fs1", "gb1"}},
	Note{MidiValue: 31, NameSharp: "G1", Frequency: 48.999, NamesOther: []string{"g1"}},
	Note{MidiValue: 32, NameSharp: "G#1", Frequency: 51.913, NamesOther: []string{"gs1", "ab1"}},
	Note{MidiValue: 33, NameSharp: "A1", Frequency: 55, NamesOther: []string{"a1"}},
	Note{MidiValue: 34, NameSharp: "A#1", Frequency: 58.27, NamesOther: []string{"as1", "bb1"}},
	Note{MidiValue: 35, NameSharp: "B1", Frequency: 61.735, NamesOther: []string{"b1", "cb1"}},
	Note{MidiValue: 36, NameSharp: "C2", Frequency: 65.406, NamesOther: []string{"c2"}},
	Note{MidiValue: 37, NameSharp: "C#2", Frequency: 69.296, NamesOther: []string{"cs2", "db2"}},
	Note{MidiValue: 38, NameSharp: "D2", Frequency: 73.416, NamesOther: []string{"d2"}},
	Note{MidiValue: 39, NameSharp: "D#2", Frequency: 77.782, NamesOther: []string{"ds2", "eb2"}},
	Note{MidiValue: 40, NameSharp: "E2", Frequency: 82.407, NamesOther: []string{"e2", "fb2"}},
	Note{MidiValue: 41, NameSharp: "F2", Frequency: 87.307, NamesOther: []string{"f2"}},
	Note{MidiValue: 42, NameSharp: "F#2", Frequency: 92.499, NamesOther: []string{"fs2", "gb2"}},
	Note{MidiValue: 43, NameSharp: "G2", Frequency: 97.999, NamesOther: []string{"g2"}},
	Note{MidiValue: 44, NameSharp: "G#2", Frequency: 103.826, NamesOther: []string{"gs2", "ab2"}},
	Note{MidiValue: 45, NameSharp: "A2", Frequency: 110, NamesOther: []string{"a2"}},
	Note{MidiValue: 46, NameSharp: "A#2", Frequency: 116.541, NamesOther: []string{"as2", "bb2"}},
	Note{MidiValue: 47, NameSharp: "B2", Frequency: 123.471, NamesOther: []string{"b2", "cb2"}},
	Note{MidiValue: 48, NameSharp: "C3", Frequency: 130.813, NamesOther: []string{"c3"}},
	Note{MidiValue: 49, NameSharp: "C#3", Frequency: 138.591, NamesOther: []string{"cs3", "db3"}},
	Note{MidiValue: 50, NameSharp: "D3", Frequency: 146.832, NamesOther: []string{"d3"}},
	Note{MidiValue: 51, NameSharp: "D#3", Frequency: 155.563, NamesOther: []string{"ds3", "eb3"}},
	Note{MidiValue: 52, NameSharp: "E3", Frequency: 164.814, NamesOther: []string{"e3", "fb3"}},
	Note{MidiValue: 53, NameSharp: "F3", Frequency: 174.614, NamesOther: []string{"f3"}},
	Note{MidiValue: 54, NameSharp: "F#3", Frequency: 184.997, NamesOther: []string{"fs3", "gb3"}},
	Note{MidiValue: 55, NameSharp: "G3", Frequency: 195.998, NamesOther: []string{"g3"}},
	Note{MidiValue: 56, NameSharp: "G#3", Frequency: 207.652, NamesOther: []string{"gs3", "ab3"}},
	Note{MidiValue: 57, NameSharp: "A3", Frequency: 220, NamesOther: []string{"a3"}},
	Note{MidiValue: 58, NameSharp: "A#3", Frequency: 233.082, NamesOther: []string{"as3", "bb3"}},
	Note{MidiValue: 59, NameSharp: "B3", Frequency: 246.942, NamesOther: []string{"b3", "cb3"}},
	Note{MidiValue: 60, NameSharp: "C4", Frequency: 261.626, NamesOther: []string{"c4"}},
	Note{MidiValue: 61, NameSharp: "C#4", Frequency: 277.183, NamesOther: []string{"cs4", "db4"}},
	Note{MidiValue: 62, NameSharp: "D4", Frequency: 293.665, NamesOther: []string{"d4"}},
	Note{MidiValue: 63, NameSharp: "D#4", Frequency: 311.127, NamesOther: []string{"ds4", "eb4"}},
	Note{MidiValue: 64, NameSharp: "E4", Frequency: 329.628, NamesOther: []string{"e4", "fb4"}},
	Note{MidiValue: 65, NameSharp: "F4", Frequency: 349.228, NamesOther: []string{"f4"}},
	Note{MidiValue: 66, NameSharp: "F#4", Frequency: 369.994, NamesOther: []string{"fs4", "gb4"}},
	Note{MidiValue: 67, NameSharp: "G4", Frequency: 391.995, NamesOther: []string{"g4"}},
	Note{MidiValue: 68, NameSharp: "G#4", Frequency: 415.305, NamesOther: []string{"gs4", "ab4"}},
	Note{MidiValue: 69, NameSharp: "A4", Frequency: 440, NamesOther: []string{"a4"}},
	Note{MidiValue: 70, NameSharp: "A#4", Frequency: 466.164, NamesOther: []string{"as4", "bb4"}},
	Note{MidiValue: 71, NameSharp: "B4", Frequency: 493.883, NamesOther: []string{"b4", "cb4"}},
	Note{MidiValue: 72, NameSharp: "C5", Frequency: 523.251, NamesOther: []string{"c5"}},
	Note{MidiValue: 73, NameSharp: "C#5", Frequency: 554.365, NamesOther: []string{"cs5", "db5"}},
	Note{MidiValue: 74, NameSharp: "D5", Frequency: 587.33, NamesOther: []string{"d5"}},
	Note{MidiValue: 75, NameSharp: "D#5", Frequency: 622.254, NamesOther: []string{"ds5", "eb5"}},
	Note{MidiValue: 76, NameSharp: "E5", Frequency: 659.255, NamesOther: []string{"e5", "fb5"}},
	Note{MidiValue: 77, NameSharp: "F5", Frequency: 698.456, NamesOther: []string{"f5"}},
	Note{MidiValue: 78, NameSharp: "F#5", Frequency: 739.989, NamesOther: []string{"fs5", "gb5"}},
	Note{MidiValue: 79, NameSharp: "G5", Frequency: 783.991, NamesOther: []string{"g5"}},
	Note{MidiValue: 80, NameSharp: "G#5", Frequency: 830.609, NamesOther: []string{"gs5", "ab5"}},
	Note{MidiValue: 81, NameSharp: "A5", Frequency: 880, NamesOther: []string{"a5"}},
	Note{MidiValue: 82, NameSharp: "A#5", Frequency: 932.328, NamesOther: []string{"as5", "bb5"}},
	Note{MidiValue: 83, NameSharp: "B5", Frequency: 987.767, NamesOther: []string{"b5", "cb5"}},
	Note{MidiValue: 84, NameSharp: "C6", Frequency: 1046.502, NamesOther: []string{"c6"}},
	Note{MidiValue: 85, NameSharp: "C#6", Frequency: 1108.731, NamesOther: []string{"cs6", "db6"}},
	Note{MidiValue: 86, NameSharp: "D6", Frequency: 1174.659, NamesOther: []string{"d6"}},
	Note{MidiValue: 87, NameSharp: "D#6", Frequency: 1244.508, NamesOther: []string{"ds6", "eb6"}},
	Note{MidiValue: 88, NameSharp: "E6", Frequency: 1318.51, NamesOther: []string{"e6", "fb6"}},
	Note{MidiValue: 89, NameSharp: "F6", Frequency: 1396.913, NamesOther: []string{"f6"}},
	Note{MidiValue: 90, NameSharp: "F#6", Frequency: 1479.978, NamesOther: []string{"fs6", "gb6"}},
	Note{MidiValue: 91, NameSharp: "G6", Frequency: 1567.982, NamesOther: []string{"g6"}},
	Note{MidiValue: 92, NameSharp: "G#6", Frequency: 1661.219, NamesOther: []string{"gs6", "ab6"}},
	Note{MidiValue: 93, NameSharp: "A6", Frequency: 1760, NamesOther: []string{"a6"}},
	Note{MidiValue: 94, NameSharp: "A#6", Frequency: 1864.655, NamesOther: []string{"as6", "bb6"}},
	Note{MidiValue: 95, NameSharp: "B6", Frequency: 1975.533, NamesOther: []string{"b6", "cb6"}},
	Note{MidiValue: 96, NameSharp: "C7", Frequency: 2093.005, NamesOther: []string{"c7"}},
	Note{MidiValue: 97, NameSharp: "C#7", Frequency: 2217.461, NamesOther: []string{"cs7", "db7"}},
	Note{MidiValue: 98, NameSharp: "D7", Frequency: 2349.318, NamesOther: []string{"d7"}},
	Note{MidiValue: 99, NameSharp: "D#7", Frequency: 2489.016, NamesOther: []string{"ds7", "eb7"}},
	Note{MidiValue: 100, NameSharp: "E7", Frequency: 2637.021, NamesOther: []string{"e7", "fb7"}},
	Note{MidiValue: 101, NameSharp: "F7", Frequency: 2793.826, NamesOther: []string{"f7"}},
	Note{MidiValue: 102, NameSharp: "F#7", Frequency: 2959.955, NamesOther: []string{"fs7", "gb7"}},
	Note{MidiValue: 103, NameSharp: "G7", Frequency: 3135.964, NamesOther: []string{"g7"}},
	Note{MidiValue: 104, NameSharp: "G#7", Frequency: 3322.438, NamesOther: []string{"gs7", "ab7"}},
	Note{MidiValue: 105, NameSharp: "A7", Frequency: 3520, NamesOther: []string{"a7"}},
	Note{MidiValue: 106, NameSharp: "A#7", Frequency: 3729.31, NamesOther: []string{"as7", "bb7"}},
	Note{MidiValue: 107, NameSharp: "B7", Frequency: 3951.066, NamesOther: []string{"b7", "cb7"}},
	Note{MidiValue: 108, NameSharp: "C8", Frequency: 4186.009, NamesOther: []string{"c8"}},
	Note{MidiValue: 109, NameSharp: "C#8", Frequency: 4434.922, NamesOther: []string{"cs8", "db8"}},
	Note{MidiValue: 110, NameSharp: "D8", Frequency: 4698.636, NamesOther: []string{"d8"}},
	Note{MidiValue: 111, NameSharp: "D#8", Frequency: 4978.032, NamesOther: []string{"ds8", "eb8"}},
	Note{MidiValue: 112, NameSharp: "E8", Frequency: 5274.042, NamesOther: []string{"e8", "fb8"}},
	Note{MidiValue: 113, NameSharp: "F8", Frequency: 5587.652, NamesOther: []string{"f8"}},
	Note{MidiValue: 114, NameSharp: "F#8", Frequency: 5919.91, NamesOther: []string{"fs8", "gb8"}},
	Note{MidiValue: 115, NameSharp: "G8", Frequency: 6271.928, NamesOther: []string{"g8"}},
	Note{MidiValue: 116, NameSharp: "G#8", Frequency: 6644.876, NamesOther: []string{"gs8", "ab8"}},
	Note{MidiValue: 117, NameSharp: "A8", Frequency: 7040, NamesOther: []string{"a8"}},
	Note{MidiValue: 118, NameSharp: "A#8", Frequency: 7458.62, NamesOther: []string{"as8", "bb8"}},
	Note{MidiValue: 119, NameSharp: "B8", Frequency: 7902.132, NamesOther: []string{"b8", "cb8"}},
	Note{MidiValue: 120, NameSharp: "C9", Frequency: 8372.018, NamesOther: []string{"c9"}},
	Note{MidiValue: 121, NameSharp: "C#9", Frequency: 8869.844, NamesOther: []string{"cs9", "db9"}},
	Note{MidiValue: 122, NameSharp: "D9", Frequency: 9397.272, NamesOther: []string{"d9"}},
	Note{MidiValue: 123, NameSharp: "D#9", Frequency: 9956.064, NamesOther: []string{"ds9", "eb9"}},
	Note{MidiValue: 124, NameSharp: "E9", Frequency: 10548.084, NamesOther: []string{"e9", "fb9"}},
	Note{MidiValue: 125, NameSharp: "F9", Frequency: 11175.304, NamesOther: []string{"f9"}},
	Note{MidiValue: 126, NameSharp: "F#9", Frequency: 11839.82, NamesOther: []string{"fs9", "gb9"}},
	Note{MidiValue: 127, NameSharp: "G9", Frequency: 12543.856, NamesOther: []string{"g9"}},
	Note{MidiValue: 128, NameSharp: "G#9", Frequency: 13289.752, NamesOther: []string{"gs9", "ab9"}},
	Note{MidiValue: 129, NameSharp: "A9", Frequency: 14080, NamesOther: []string{"a9"}},
	Note{MidiValue: 130, NameSharp: "A#9", Frequency: 14917.24, NamesOther: []string{"as9", "bb9"}},
	Note{MidiValue: 131, NameSharp: "B9", Frequency: 15804.264, NamesOther: []string{"b9", "cb9"}},
}
var dbChords = [][]string{
	[]string{"1P 3M 5P", "major", "maj", "^", ""},
	[]string{"1P 3M 5P 7M", "major seventh", "maj7", "ma7", "Maj7", "^7"},
	[]string{"1P 3M 5P 7M 9M", "major ninth", "maj9", "^9"},
	[]string{"1P 3M 5P 7M 9M 13M", "major thirteenth", "maj13", "Maj13 ^13"},
	[]string{"1P 3M 5P 6M", "sixth", "6", "add6", "add13"},
	[]string{"1P 3M 5P 6M 9M", "sixth/ninth", "6/9", "69"},
	[]string{"1P 3M 6m 7M", "major seventh flat sixth", "maj7b6", "^7b6"},
	[]string{"1P 3M 5P 7M 11A", "major seventh sharp eleventh", "majs4", "^7#11", "maj7#11"},
	// ==Minor==
	// '''Normal'''
	[]string{"1P 3m 5P", "minor", "m", "min", "-"},
	[]string{"1P 3m 5P 7m", "minor seventh", "m7", "min7", "mi7", "-7"},
	[]string{"1P 3m 5P 7M", "minor/major seventh", "maj7", "majmaj7", "mM7", "mMaj7", "m/M7", "-^7"},
	[]string{"1P 3m 5P 6M", "minor sixth", "m6", "-6"},
	[]string{"1P 3m 5P 7m 9M", "minor ninth", "m9", "-9"},
	[]string{"1P 3m 5P 7M 9M", "minor/major ninth", "minmaj9", "mMaj9", "-^9"},
	[]string{"1P 3m 5P 7m 9M 11P", "minor eleventh", "m11", "-11"},
	[]string{"1P 3m 5P 7m 9M 13M", "minor thirteenth", "m13", "-13"},
	// '''Diminished'''
	[]string{"1P 3m 5d", "diminished", "dim", "°", "o"},
	[]string{"1P 3m 5d 7d", "diminished seventh", "dim7", "°7", "o7"},
	[]string{"1P 3m 5d 7m", "half-diminished", "m7b5", "ø", "-7b5", "h7", "h"},
	// ==Dominant/Seventh==
	// '''Normal'''
	[]string{"1P 3M 5P 7m", "dominant seventh", "7", "dom"},
	[]string{"1P 3M 5P 7m 9M", "dominant ninth", "9"},
	[]string{"1P 3M 5P 7m 9M 13M", "dominant thirteenth", "13"},
	[]string{"1P 3M 5P 7m 11A", "lydian dominant seventh", "7s11", "7#4"},
	// '''Altered'''
	[]string{"1P 3M 5P 7m 9m", "dominant flat ninth", "7b9"},
	[]string{"1P 3M 5P 7m 9A", "dominant sharp ninth", "7s9"},
	[]string{"1P 3M 7m 9m", "altered", "alt7"},
	// '''Suspended'''
	[]string{"1P 4P 5P", "suspended fourth", "sus4", "sus"},
	[]string{"1P 2M 5P", "suspended second", "sus2"},
	[]string{"1P 4P 5P 7m", "suspended fourth seventh", "7sus4", "7sus"},
	[]string{"1P 5P 7m 9M 11P", "eleventh", "11"},
	[]string{"1P 4P 5P 7m 9m", "suspended fourth flat ninth", "b9sus", "phryg", "7b9sus", "7b9sus4"},
	// ==Other==
	[]string{"1P 5P", "fifth", "5"},
	[]string{"1P 3M 5A", "augmented", "aug", "+", "+5", "^#5"},
	[]string{"1P 3m 5A", "minor augmented", "ms5", "-#5", "m+"},
	[]string{"1P 3M 5A 7M", "augmented seventh", "maj75", "maj7+5", "+maj7", "^7#5"},
	[]string{"1P 3M 5P 7M 9M 11A", "major sharp eleventh (lydian)", "maj9s11", "^9#11"},
	// ==Legacy==
	[]string{"1P 2M 4P 5P", "", "sus24", "sus4add9"},
	[]string{"1P 3M 5A 7M 9M", "", "maj9s5", "Maj9s5"},
	[]string{"1P 3M 5A 7m", "", "7s5", "+7", "7+", "7aug", "aug7"},
	[]string{"1P 3M 5A 7m 9A", "", "7s5s9", "7s9s5", "7alt"},
	[]string{"1P 3M 5A 7m 9M", "", "9s5", "9+"},
	[]string{"1P 3M 5A 7m 9M 11A", "", "9s5s11"},
	[]string{"1P 3M 5A 7m 9m", "", "7s5b9", "7b9s5"},
	[]string{"1P 3M 5A 7m 9m 11A", "", "7s5b9s11"},
	[]string{"1P 3M 5A 9A", "", "padds9"},
	[]string{"1P 3M 5A 9M", "", "ms5add9", "padd9"},
	[]string{"1P 3M 5P 6M 11A", "", "M6s11", "M6b5", "6s11", "6b5"},
	[]string{"1P 3M 5P 6M 7M 9M", "", "maj7add13"},
	[]string{"1P 3M 5P 6M 9M 11A", "", "69s11"},
	[]string{"1P 3m 5P 6M 9M", "", "m69", "-69"},
	[]string{"1P 3M 5P 6m 7m", "", "7b6"},
	[]string{"1P 3M 5P 7M 9A 11A", "", "maj7s9s11"},
	[]string{"1P 3M 5P 7M 9M 11A 13M", "", "M13s11", "maj13s11", "M13+4", "M13s4"},
	[]string{"1P 3M 5P 7M 9m", "", "maj7b9"},
	[]string{"1P 3M 5P 7m 11A 13m", "", "7s11b13", "7b5b13"},
	[]string{"1P 3M 5P 7m 13M", "", "7add6", "67", "7add13"},
	[]string{"1P 3M 5P 7m 9A 11A", "", "7s9s11", "7b5s9", "7s9b5"},
	[]string{"1P 3M 5P 7m 9A 11A 13M", "", "13s9s11"},
	[]string{"1P 3M 5P 7m 9A 11A 13m", "", "7s9s11b13"},
	[]string{"1P 3M 5P 7m 9A 13M", "", "13s9"},
	[]string{"1P 3M 5P 7m 9A 13m", "", "7s9b13"},
	[]string{"1P 3M 5P 7m 9M 11A", "", "9s11", "9+4", "9s4"},
	[]string{"1P 3M 5P 7m 9M 11A 13M", "", "13s11", "13+4", "13s4"},
	[]string{"1P 3M 5P 7m 9M 11A 13m", "", "9s11b13", "9b5b13"},
	[]string{"1P 3M 5P 7m 9m 11A", "", "7b9s11", "7b5b9", "7b9b5"},
	[]string{"1P 3M 5P 7m 9m 11A 13M", "", "13b9s11"},
	[]string{"1P 3M 5P 7m 9m 11A 13m", "", "7b9b13s11", "7b9s11b13", "7b5b9b13"},
	[]string{"1P 3M 5P 7m 9m 13M", "", "13b9"},
	[]string{"1P 3M 5P 7m 9m 13m", "", "7b9b13"},
	[]string{"1P 3M 5P 7m 9m 9A", "", "7b9s9"},
	[]string{"1P 3M 5P 9M", "", "Madd9", "2", "add9", "add2"},
	[]string{"1P 3M 5P 9m", "", "majaddb9"},
	[]string{"1P 3M 5d", "", "majb5"},
	[]string{"1P 3M 5d 6M 7m 9M", "", "13b5"},
	[]string{"1P 3M 5d 7M", "", "maj7b5"},
	[]string{"1P 3M 5d 7M 9M", "", "maj9b5"},
	[]string{"1P 3M 5d 7m", "", "7b5"},
	[]string{"1P 3M 5d 7m 9M", "", "9b5"},
	[]string{"1P 3M 7m", "", "7no5"},
	[]string{"1P 3M 7m 13m", "", "7b13"},
	[]string{"1P 3M 7m 9M", "", "9no5"},
	[]string{"1P 3M 7m 9M 13M", "", "13no5"},
	[]string{"1P 3M 7m 9M 13m", "", "9b13"},
	[]string{"1P 3m 4P 5P", "", "madd4"},
	[]string{"1P 3m 5P 6m 7M", "", "mmaj7b6"},
	[]string{"1P 3m 5P 6m 7M 9M", "", "mmaj9b6"},
	[]string{"1P 3m 5P 7m 11P", "", "m7add11", "m7add4"},
	[]string{"1P 3m 5P 9M", "", "madd9"},
	[]string{"1P 3m 5d 6M 7M", "", "o7maj7"},
	[]string{"1P 3m 5d 7M", "", "omaj7"},
	[]string{"1P 3m 6m 7M", "", "mb6maj7"},
	[]string{"1P 3m 6m 7m", "", "m7s5"},
	[]string{"1P 3m 6m 7m 9M", "", "m9s5"},
	[]string{"1P 3m 5A 7m 9M 11P", "", "m11A"},
	[]string{"1P 3m 6m 9m", "", "mb6b9"},
	[]string{"1P 2M 3m 5d 7m", "", "m9b5"},
	[]string{"1P 4P 5A 7M", "", "maj7s5sus4"},
	[]string{"1P 4P 5A 7M 9M", "", "maj9s5sus4"},
	[]string{"1P 4P 5A 7m", "", "7s5sus4"},
	[]string{"1P 4P 5P 7M", "", "maj7sus4"},
	[]string{"1P 4P 5P 7M 9M", "", "maj9sus4"},
	[]string{"1P 4P 5P 7m 9M", "", "9sus4", "9sus"},
	[]string{"1P 4P 5P 7m 9M 13M", "", "13sus4", "13sus"},
	[]string{"1P 4P 5P 7m 9m 13m", "", "7sus4b9b13", "7b9b13sus4"},
	[]string{"1P 4P 7m 10m", "", "4", "quartal"},
	[]string{"1P 5P 7m 9m 11P", "", "11b9"},
}
var notesWhite = []string{"C", "D", "E", "F", "G", "A", "B"}
var notesScaleSharp = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
var notesScaleAcc1 = []string{"B#", "Db", "D", "Eb", "Fb", "E#", "Gb", "G", "Ab", "A", "Bb", "Cb"}
var notesScaleAcc2 = []string{"C", "Cs", "D", "Ds", "E", "F", "Fs", "G", "Gs", "A", "As", "B"}
var notesScaleAcc3 = []string{"Bs", "Db", "D", "Eb", "Fb", "Es", "Gb", "G", "Ab", "A", "Bb", "Cb"}
var notesAdds = []string{"", "#", "b", "s"}
var notesAll = []string{}

var C4 = Note{MidiValue: 60, NameSharp: "C4", Frequency: 261.626, NamesOther: []string{"c4"}}
var D4 = Note{MidiValue: 62, NameSharp: "D4", Frequency: 293.665, NamesOther: []string{"d4"}}
var E4 = Note{MidiValue: 64, NameSharp: "E4", Frequency: 329.628, NamesOther: []string{"e4", "fb4"}}
var F4 = Note{MidiValue: 65, NameSharp: "F4", Frequency: 349.228, NamesOther: []string{"f4"}}
var G4 = Note{MidiValue: 67, NameSharp: "G4", Frequency: 391.995, NamesOther: []string{"g4"}}
var A4 = Note{MidiValue: 69, NameSharp: "A4", Frequency: 440, NamesOther: []string{"a4"}}
var B4 = Note{MidiValue: 71, NameSharp: "B4", Frequency: 493.883, NamesOther: []string{"b4", "cb4"}}
var C5 = Note{MidiValue: 72, NameSharp: "C5", Frequency: 523.251, NamesOther: []string{"c5"}}
var D5 = Note{MidiValue: 74, NameSharp: "D5", Frequency: 587.33, NamesOther: []string{"d5"}}
var E5 = Note{MidiValue: 76, NameSharp: "E5", Frequency: 659.255, NamesOther: []string{"e5", "fb5"}}
var F5 = Note{MidiValue: 77, NameSharp: "F5", Frequency: 698.456, NamesOther: []string{"f5"}}
var G5 = Note{MidiValue: 79, NameSharp: "G5", Frequency: 783.991, NamesOther: []string{"g5"}}
var A5 = Note{MidiValue: 81, NameSharp: "A5", Frequency: 880, NamesOther: []string{"a5"}}
var B5 = Note{MidiValue: 83, NameSharp: "B5", Frequency: 987.767, NamesOther: []string{"b5", "cb5"}}
var C6 = Note{MidiValue: 84, NameSharp: "C6", Frequency: 1046.502, NamesOther: []string{"c6"}}
var D6 = Note{MidiValue: 86, NameSharp: "D6", Frequency: 1174.659, NamesOther: []string{"d6"}}
var E6 = Note{MidiValue: 88, NameSharp: "E6", Frequency: 1318.51, NamesOther: []string{"e6", "fb6"}}

func init() {
	for _, n := range notesWhite {
		for _, a := range notesAdds {
			notesAll = append(notesAll, n+a)
		}
	}
	// convert everything to lowercase
	for i, n := range notesWhite {
		notesWhite[i] = strings.ToLower(n)
	}
	for i, n := range notesScaleSharp {
		notesScaleSharp[i] = strings.ToLower(n)
	}
	for i, n := range notesScaleAcc1 {
		notesScaleAcc1[i] = strings.ToLower(n)
	}
	for i, n := range notesScaleAcc2 {
		notesScaleAcc2[i] = strings.ToLower(n)
	}
	for i, n := range notesScaleAcc3 {
		notesScaleAcc3[i] = strings.ToLower(n)
	}
	for i, n := range notesAdds {
		notesAdds[i] = strings.ToLower(n)
	}
	for i, n := range notesAll {
		notesAll[i] = strings.ToLower(n)
	}
	for i, n := range noteDB {
		noteDB[i].NameSharp = strings.ToLower(n.NameSharp)
	}
	// add name sharp into nameother
	for i, n := range noteDB {
		noteDB[i].NamesOther = append(noteDB[i].NamesOther, n.NameSharp)
	}
	
	// Build optimization lookup maps
	buildLookupMaps()
}

func buildLookupMaps() {
	// Build note lookup map for faster access
	for i := range noteDB {
		note := &noteDB[i]
		noteByNameMap[note.NameSharp] = note
		for _, name := range note.NamesOther {
			noteByNameMap[name] = note
		}
	}
	
	// Build chord pattern lookup map
	for _, chordType := range dbChords {
		if len(chordType) < 2 {
			continue
		}
		intervals := chordType[0]
		for i := 2; i < len(chordType); i++ {
			pattern := strings.ToLower(chordType[i])
			if len(pattern) > len(chordPatternMap[pattern]) {
				chordPatternMap[pattern] = intervals
			}
		}
	}
	
	// Build note conversion map for accidentals
	for i, n := range notesScaleAcc1 {
		if i < len(notesScaleSharp) {
			noteConversionMap[n] = notesScaleSharp[i]
		}
	}
	for i, n := range notesScaleAcc2 {
		if i < len(notesScaleSharp) {
			noteConversionMap[n] = notesScaleSharp[i]
		}
	}
	for i, n := range notesScaleAcc3 {
		if i < len(notesScaleSharp) {
			noteConversionMap[n] = notesScaleSharp[i]
		}
	}
}
