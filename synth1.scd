(
s.waitForBoot({
	SynthDef(\simpleSynth, {
		arg freq = 440, amp = 1.0, gate = 1;
		var env, osc, out;
		
		env = EnvGen.kr(Env.adsr(0.01, 0.3, 0.5, 0.8), gate, doneAction: 2);
		
		osc = Mix.new(Saw.ar(freq * [1, 1.008, 0.993], amp / 3));
		
		osc = MoogFF.ar(osc,4000);
		
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

/*
NetAddr.localAddr.sendMsg('/synth1/noteOn', 60, 100);
NetAddr.localAddr.sendMsg('/synth1/noteOn', 64, 80);
NetAddr.localAddr.sendMsg('/synth1/noteOn', 67, 90);

NetAddr.localAddr.sendMsg('/synth1/noteOff', 60);
NetAddr.localAddr.sendMsg('/synth1/noteOff', 64);
NetAddr.localAddr.sendMsg('/synth1/noteOff', 67);
*/

/*
(
OSCdef(\noteOn).free;
OSCdef(\noteOff).free;
~activeNotes.do(_.free);
~activeNotes.clear;
"OSC responders cleaned up.".postln;
)
*/
