const chatWindow = document.getElementById('chat-window');
const recordBtn = document.getElementById('record-btn');
let mediaRecorder;
let audioChunks = [];

recordBtn.addEventListener('click', () => {
    if (recordBtn.textContent === 'Record Audio') {
        startRecording();
    } else {
        stopRecording();
    }
});

function startRecording() {
    navigator.mediaDevices.getUserMedia({ audio: true })
        .then(stream => {
            mediaRecorder = new MediaRecorder(stream);
            mediaRecorder.start();
            recordBtn.textContent = 'Stop Recording';

            mediaRecorder.ondataavailable = event => {
                audioChunks.push(event.data);
            };

            mediaRecorder.onstop = () => {
                const audioBlob = new Blob(audioChunks, { type: 'audio/wav' });
                audioChunks = [];
                handleAudioBlob(audioBlob);
            };
        });
}

function stopRecording() {
    mediaRecorder.stop();
    recordBtn.textContent = 'Record Audio';
}

function handleAudioBlob(audioBlob) {
    const formData = new FormData();
    formData.append('file', audioBlob, 'audio.wav');

    // Send audio to speech-to-text API
    fetch('http://localhost:8080/chat/answer', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        const userMessage = data.prompt.text;
        appendMessage(userMessage, 'user');
        
        const replyMessage = data.answer.text;
        appendMessage(replyMessage, 'reply')

        const replyAudio = data.answer.audio;
        
        const audioData = atob(replyAudio);
        const audioBuffer = new ArrayBuffer(audioData.length);
        const audioView = new Uint8Array(audioBuffer);
        for (let i = 0; i < audioData.length; i++) {
            audioView[i] = audioData.charCodeAt(i);
        }
        const audioBlob = new Blob([audioBuffer], { type: 'audio/wav' });
        const audioUrl = URL.createObjectURL(audioBlob);

        const audio = new Audio(audioUrl);
        audio.play();
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

function appendMessage(message, type) {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type}`;
    messageDiv.textContent = message;
    chatWindow.appendChild(messageDiv);
    chatWindow.scrollTop = chatWindow.scrollHeight;
}
