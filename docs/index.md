# museq 

**museq** is a music sequencer. It works like a [tracker](https://en.wikipedia.org/wiki/Music_tracker), but the main difference is that it 
designed to be use with any text editor. 

## features

- **Text-based**: You can write music in any text editor, and it will work on any platform.
- **MIDI support**: It can play music using MIDI devices, so you can use it with any MIDI instrument or software.
- **SuperCollider support**: It can also play music using SuperCollider, so you can use it with any SuperCollider synth.
- **Patterns**: You can create patterns of music that can be reused and combined.
- **Chains**: You can create chains of patterns that can be run independently.
- **Arpeggios**: You can create arpeggios using a simple syntax.
- **Chords**: You can create chords using a simple syntax.
- **Asynchronous velocity/transpose**: You can change the velocity and transpose notes asynchronously.

## quickstart

First install **museq**:

```bash
curl https://museq.schollz.com/install.sh | bash
```

Now you can identify the name of midi devices:

```
museq -midi
```

When you call a midi device, use any part of the name (case insensitive). 

## syntax

The syntax for **museq** is called **tli** (text-limited interface). '

Let's start with a full-fledged example of a **museq** file, and go through it piece by piece which will help understand 90% of what **museq** can do. 

```bash
// first example 

midi thingy 1
// supercollider /synth1
bpm 180
transpose -12
gate 0.2

# part1 
a3 b c d

# part2
Cmaj/G;3@u4d8,u2,d2

# part3
g _
[g ~] d4,e5

// pattern chaining
+# part1*3 [part2 part3] * 2

```

### playing

To play the above example, you would save it to a file called `example.tli` and then run:

```bash
museq example.tli
```

### comments

Comments like on line 1:
```
// first example 
```

are prefixed with `//` and can be placed anywhere in the file. They are ignored by **museq**.

### players

There are two players: `midi` and `supercollider`. You can specify which player to use by using the `midi` or `supercollider` keyword followed by the name of the device. For example, on line 2:

```bash
midi thingy 1
```

This specifies that the `thingy` MIDI device should be used, and the `1` specifies the MIDI channel (1-16). If you want to use SuperCollider, you would write

```bash
supercollider /synth1
```

This specifies that the SuperCollider path `/synth1` should be used. Here's an examle SuperCollider synth definition:
```supercollider
(
s.waitForBoot({
	SynthDef(\simpleSynth, {
		arg freq = 440, amp = 0.5, gate = 1;
		var env, osc, out;
		env = EnvGen.kr(Env.adsr(0.01, 0.3, 0.5, 0.8), gate, doneAction: 2);
		osc = Mix.new(Saw.ar(freq * [1, 1.008, 0.993], amp / 3));
		osc = MoogFF.ar(osc,MouseX.kr(100,10000,1));
		out = osc * env;
		out = Pan2.ar(out,Rand(-0.1,0.1));
		Out.ar(0, out);
	}).add;
	~activeNotes = Dictionary.new;
	OSCdef(\noteOn, {
		arg msg, time, addr, recvPort;
		var note, velocity, freq, amp, synth;
		note = msg[1];
		velocity = msg[2];
		freq = note.midicps;
		amp = velocity / 127.0;
		synth = Synth(\simpleSynth, [\freq, freq, \amp, amp]);
		~activeNotes[note] = synth;
		("Note ON: " ++ note ++ " (freq: " ++ freq.round(0.1) ++ " Hz, vel: " ++ velocity ++ ")").postln;
	}, '/synth1/noteOn');
	OSCdef(\noteOff, {
		arg msg, time, addr, recvPort;
		var note, synth;
		note = msg[1];
		synth = ~activeNotes[note];
		if(synth.notNil, {
			synth.set(\gate, 0);
			~activeNotes.removeAt(note);
			("Note OFF: " ++ note).postln;
		}, {
			("Warning: Note OFF received for inactive note: " ++ note).postln;
		});
	}, '/synth1/noteOff');
	"OSC Synth setup complete!".postln;
	"Listening for:".postln;
	"  /synth1/noteOn [note, velocity]".postln;
	"  /synth1/noteOff [note]".postln;
	"Default OSC port: 57120".postln;
});
)
```

### globals

You can specify global settings that apply to the whole section. Global settings are specified with the `bpm`, `transpose`, and `gate` keywords. For example, on lines 3-5:

```bash
bpm 180
transpose -12
gate 0.2
```

- `bpm` specifies the beats per minute (default is 120).
- `transpose` specifies the number of semitones to transpose the notes (default is 0).
- `gate` specifies the duration of the notes in beats (default is 0.5).

## basic note patterns

Note patterns are specified by first specifying the pattern name, prefixed with `#`, followed by the notes. For example, on line 7:

```bash
# part1
a3 b c d
```

This specifies a pattern called `part1` that contains the notes `a3`, `b`, `c`, and `d`. The notes can be specified in any order, and can include octaves (e.g., `a3` is an A note in the 3rd octave). If it does not include an octave, it will default to the closest octave from the previous note. 

Each line is one measure, and each entity gets an equal proportion of the measure. For example, in the above example, each note gets 1/4 of the measure (i.e., a quarter note).

### chords and arpeggios

You can specify chords and arpeggios using a simple syntax. For example, on line 10:

```bash
# part2
Cmaj/G;3@u4d8,u2,d2
```

This specifies a chord `Cmaj` which is a C major chord. It is optionally transposed by including `/G` which means that the bass note is G. The `;3` optionally specifies the octave. Finally, the `@u4d8` specifies that the chord should be played as an arpeggio. If you do not have the arpeggio syntax, it will be played as a chord, polyphonically.

The arpeggio syntax has several options:

- `u` specifies the upward direction (default).
- `d` specifies the downward direction.
- The number specifies the number of steps to include.

Also notice that there is a comma `,` which separates different arpeggio patterns. When using a `,` to separate different patterns (arpeggios or chords or notes) it will move through which one to play in a round-robin fashion. For example, the first playthrough will play the first pattern, the second playthrough will play the second pattern, and so on.

### legato, rests, and subdivisions

In the third pattern we have examples of legato, rests, and subdivisions:

```bash
# part3
g _
[g ~] d4,e5
```

The `# part3` specifies a pattern called `part3`. 

The first line `g _` specifies that the note `g` should be played legato (i.e., held) until the next note. The underscore `_` indicates that the note should be held until the next note is played.

The second line has `g ~` which specifies that the note `g` is followed by a rest. It is enclosed in `[]` which indicates that it is a single entity. The last part, `d4,e5`, specifies the two notes `d4` and `e5` - where the first, `d4` is played the first time through and the second, `e5` is played the second time through. 

Both `[g ~]` and `d4,e5` are considered single entities, so each is proportioned 1/2 of the measure. The `[g ~]` is then further subdivided in time - the `g` gets 1/2 of its subdivision and the `~` gets the other half. So in the end, `g` gets 1/4 of the measure, and `d4`/`e5` each get 1/2 of the measure. This is a powerful way to create complex rhythms and patterns. However, it could also be written as:

```bash
g ~ d4,e5 _
```

I.e. The `g` is played for 1/4 measure, then rested, and then the `d4` or `e5` are played for the rest of the measure.

## pattern chaining

When you have multiple patterns, you can chain them together using the `+` operator. For example, on line 16:

```bash
+# part1*3 [part2 part3] * 2
```

This specifies that the `part1` pattern should be played 3 times, followed by the `part2` and `part3` patterns played after each other twice. The `*` operator specifies that the pattern should be repeated a certain number of times and the `[]` indicates that the patterns inside should be played considered a single entity (useful when you want to play multiple patterns together). This could be written in an expanded form as:

```bash	
+# part1 part1 part1 part2 part3 part2 part3
```

which is equivalent to the above.