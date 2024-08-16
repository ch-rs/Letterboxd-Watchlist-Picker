export default function (velocity, frequency) {
    var context = new AudioContext();

    // Create an oscillator for the main tone
    var oscillator = context.createOscillator();
    oscillator.type = "sine";
    oscillator.frequency.value = frequency; // Frequency for the metal ding sound

    // Create a gain node to control the volume
    var gainNode = context.createGain();
    oscillator.connect(gainNode);
    gainNode.connect(context.destination);

    // Create a noise buffer
    var bufferSize = context.sampleRate * 0.5; // 0.5 seconds of noise
    var noiseBuffer = context.createBuffer(1, bufferSize, context.sampleRate);
    var output = noiseBuffer.getChannelData(0);
    for (var i = 0; i < bufferSize; i++) {
        output[i] = Math.random() * 2 - 1;
    }

    // Create a noise source
    var noise = context.createBufferSource();
    noise.buffer = noiseBuffer;
    var noiseGain = context.createGain();
    //noise.connect(noiseGain);
    noiseGain.connect(context.destination);

    // Apply a gain envelope to the oscillator and noise based on velocity
    gainNode.gain.setValueAtTime(velocity, context.currentTime);
    gainNode.gain.exponentialRampToValueAtTime(0.001, context.currentTime + 1); // Fade out over 0.5 seconds
    noiseGain.gain.setValueAtTime(velocity * 0.5, context.currentTime); // Adjust the noise volume based on velocity
    noiseGain.gain.exponentialRampToValueAtTime(0.001, context.currentTime + 1); // Fade out over 0.5 seconds

    // Start the oscillator and noise
    oscillator.start(context.currentTime);
    noise.start(context.currentTime);

    // Stop the oscillator and noise after 0.5 seconds
    oscillator.stop(context.currentTime + 1);
    noise.stop(context.currentTime + 1);
}