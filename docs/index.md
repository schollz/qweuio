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

+# part1*3 [part2 part3] * 2

# part1 
a3 b c d

# part2
Cmaj/G;3@u4d8

# part3
g _
[g ~] d4
```

### comments

Comments like on line 1:
```
// first example 
```

are prefixed with `//` and can be placed anywhere in the file. They are ignored by **museq**.

### specify player

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